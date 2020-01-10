package configmap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/kyma-project/rafter/tests/pkg/retry"
)

type Configmap struct {
	configMapCli      corev1.ConfigMapInterface
	waitTimeout       time.Duration
	namespace         string
	createdConfigMaps []string
}

func New(coreCli *corev1.CoreV1Client, namespace string, waitTimeout time.Duration) *Configmap {
	return &Configmap{configMapCli: coreCli.ConfigMaps(namespace),
		namespace:   namespace,
		waitTimeout: waitTimeout,
	}
}

func (s *Configmap) Create(name string, files []*os.File, callbacks ...func(...interface{})) (string, error) {
	data := make(map[string]string)
	binaryData := make(map[string][]byte)

	for _, file := range files {
		if file == nil {
			return "", fmt.Errorf("file can't be nil")
		}
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			return "", errors.Wrapf(err, "while reading file %s", file.Name())
		}

		filename := filepath.Base(file.Name())

		s.log(string(fileData), callbacks...)

		if filepath.Ext(filename) == ".json" {
			data[filename] = string(fileData)
		}
		if filepath.Ext(filename) == ".png" {
			binaryData[filename] = fileData
		}
	}
	s.log(fmt.Sprintf("[CREATE] configmap: %s", name), callbacks...)
	configmap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.namespace,
		},
		Data:       data,
		BinaryData: binaryData,
	}

	confMap, err := s.configMapCli.Create(configmap)
	if err != nil {
		return "", errors.Wrapf(err, "while creating ConfigMap %s", name)
	}

	s.createdConfigMaps = append(s.createdConfigMaps, name)

	return confMap.ResourceVersion, err
}

func (s *Configmap) DeleteAll(callbacks ...func(...interface{})) error {
	for _, name := range s.createdConfigMaps {
		err := retry.WithIgnoreOnNotFound(retry.DefaultBackoff, func() error {
			s.log(fmt.Sprintf("DELETE: configmap: %s", name), callbacks...)
			return s.configMapCli.Delete(name, &metav1.DeleteOptions{})
		}, callbacks...)
		if err != nil {
			return errors.Wrapf(err, "While deleting configmap: %s", name)
		}
	}
	return nil
}

func (s *Configmap) log(message string, callbacks ...func(...interface{})) {
	for _, callback := range callbacks {
		callback(message)
	}
}
