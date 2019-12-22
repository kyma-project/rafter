package loader

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
	"regexp"
	"strings"
)

func (l *loader) loadConfigMap(src string, name string, filter string) (string, []string, error) {
	basePath, err := ioutil.TempDir(l.temporaryDir, name)
	if err != nil {
		return "", nil, err
	}

	filterRegexp, err := regexp.Compile(filter)
	if err != nil {
		return "", nil, errors.Wrapf(err, "while compiling filter")
	}

	srcs := strings.Split(src, "/")
	if len(srcs) != 2 {
		return "", nil, fmt.Errorf("%s: illegal source", src)
	}

	configMap, err := l.getConfigMap(srcs[1], srcs[0])
	if err != nil {
		return "", nil, errors.Wrapf(err, "while getting configmap")
	}

	var fileList []string
	for key, value := range configMap.Data {
		if fileList, err = l.copyBytesToFile([]byte(value), key, basePath, filterRegexp, fileList); err != nil {
			return "", nil, errors.Wrapf(err, "while copying data to file")
		}
	}

	for key, value := range configMap.BinaryData {
		if fileList, err = l.copyBytesToFile(value, key, basePath, filterRegexp, fileList); err != nil {
			return "", nil, errors.Wrapf(err, "while copying binary data to file")
		}
	}

	return basePath, fileList, nil
}

func (l *loader) getConfigMap(name string, namespace string) (*corev1.ConfigMap, error) {
	configmapsResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}

	item, err := l.dynamicClient.Resource(configmapsResource).Namespace(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var configMap corev1.ConfigMap
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.UnstructuredContent(), &configMap)
	if err != nil {
		return nil, err
	}

	return &configMap, nil
}

func (l *loader) copyBytesToFile(value []byte, name string, path string, regexp *regexp.Regexp, fileList []string) ([]string, error) {
	if regexp.MatchString(name) {
		destination := filepath.Join(path, name)
		file, err := l.osCreateFunc(destination)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(file, bytes.NewReader(value))
		if err != nil {
			return nil, err
		}

		fileList = append(fileList, name)
	}

	return fileList, nil
}
