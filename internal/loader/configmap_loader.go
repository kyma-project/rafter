package loader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (l *loader) loadConfigMap(src string, name string, filter string) (string, []string, error) {
	basePath, err := ioutil.TempDir(l.temporaryDir, name)
	if err != nil {
		return "", nil, err
	}

	filterRegexp, err := regexp.Compile(filter)
	if err != nil {
		return "", nil, errors.Wrap(err, "while compiling filter")
	}

	srcs := strings.Split(src, "/")
	if len(srcs) != 2 {
		return "", nil, fmt.Errorf("%s: invalid source format", src)
	}

	configMap, err := l.getConfigMap(srcs[0], srcs[1])
	if err != nil {
		return "", nil, err
	}

	var fileList []string
	for key, value := range configMap.Data {
		if fileList, err = l.copyBytesToFile([]byte(value), key, basePath, filterRegexp, fileList); err != nil {
			return "", nil, errors.Wrap(err, "while copying data to file")
		}
	}

	for key, value := range configMap.BinaryData {
		if fileList, err = l.copyBytesToFile(value, key, basePath, filterRegexp, fileList); err != nil {
			return "", nil, errors.Wrap(err, "while copying binary data to file")
		}
	}

	return basePath, fileList, nil
}

func (l *loader) getConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
	configmapsResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}

	item, err := l.dynamicClient.Resource(configmapsResource).Namespace(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "while getting ConfigMap %s from %s namespace", name, namespace)

	}

	var configMap corev1.ConfigMap
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.UnstructuredContent(), &configMap)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting Unstructured to ConfigMap %s from %s", name, namespace)
	}

	return &configMap, nil
}

func (l *loader) copyBytesToFile(value []byte, name string, path string, regexp *regexp.Regexp, fileList []string) ([]string, error) {
	if !regexp.MatchString(name) {
		return fileList, nil
	}

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

	return fileList, nil
}
