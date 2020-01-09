package loader

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicFake "k8s.io/client-go/dynamic/fake"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
)

var (
	exampleDataBody = `{"type": "service_account","project_id": "wookiee","private_key_id": "12312312312312323123123123123123"}`
	examplePNGBody  = `iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAABg2lDQ1BJQ0MgcHJvZmlsZQAAKJF9kT1Iw0AcxV9TpVIqgmYQcchQnSyIijhqFYpQIdQKrTqYXPoFTRqSFBdHwbXg4Mdi1cHFWVcHV0EQ/ABxcXVSdJES/9cUWsR4cNyPd/ced+8AoV5mut01DuiGY6UScSmTXZVCr4hARBhAv8Jsc06Wk/AdX/cI8PUuxrP8z/05erWczYCARDzLTMsh3iCe3nRMzvvEIisqGvE58ZhFFyR+5Lrq8RvnQpMFnila6dQ8sUgsFTpY7WBWtHTiKeKophuUL2Q81jhvcdbLVda6J39hJGesLHOd5jASWMQSZEhQUUUJZTiI0WqQYiNF+3Ef/1DTL5NLJVcJjBwLqECH0vSD/8Hvbu385ISXFIkD3S+u+zEChHaBRs11v49dt3ECBJ+BK6Ptr9SBmU/Sa20tegT0bQMX121N3QMud4DBJ1OxlKYUpCnk88D7GX1TFhi4BcJrXm+tfZw+AGnqKnkDHBwCowXKXvd5d09nb/+eafX3A/3VcngnrbzqAAAACXBIWXMAAC4jAAAuIwF4pT92AAAAB3RJTUUH4wwUCgEsYNs2MQAAABl0RVh0Q29tbWVudABDcmVhdGVkIHdpdGggR0lNUFeBDhcAAAAPSURBVAjXY/jPAAP/GRgAD/4B/yIectQAAAAASUVORK5CYII=`
)

func TestLoader_Load_ConfigMap(t *testing.T) {
	fakedc, err := newFakeDynamicClient(
		fixConfigMap("text", "default", map[string]string{
			"example.json": exampleDataBody,
			"example2.txt": exampleDataBody,
		}, nil),
		fixConfigMap("binary", "default", nil, map[string][]byte{
			"example.png":  []byte(examplePNGBody),
			"example2.txt": []byte(examplePNGBody),
		}),
		fixConfigMap("mix", "default", map[string]string{
			"example.txt": exampleDataBody,
		}, map[string][]byte{
			"example2.png": []byte(examplePNGBody),
		}),
	)
	if err != nil {
		return
	}
	loader := &loader{
		temporaryDir:    "/tmp",
		dynamicClient:   fakedc,
		osRemoveAllFunc: os.RemoveAll,
		osCreateFunc:    os.Create,
		httpGetFunc:     get,
		ioutilTempDir:   ioutil.TempDir,
	}

	for testName, testData := range map[string]struct {
		src        string
		name       string
		mode       v1beta1.AssetMode
		filter     string
		files      int
		errMatcher types.GomegaMatcher
	}{
		"ConfigMap with text files": {
			src:        "default/text",
			name:       "asset-text",
			mode:       v1beta1.AssetConfigMap,
			filter:     "",
			files:      2,
			errMatcher: gomega.BeNil(),
		},
		"Config-Map with binary files": {
			src:        "default/binary",
			name:       "asset-text",
			mode:       v1beta1.AssetConfigMap,
			filter:     "",
			files:      2,
			errMatcher: gomega.BeNil(),
		},
		"Config-Map with mixed files": {
			src:        "default/mix",
			name:       "asset-text",
			mode:       v1beta1.AssetConfigMap,
			filter:     "",
			files:      2,
			errMatcher: gomega.BeNil(),
		},
		"Filtered Config-Map": {
			src:        "default/mix",
			name:       "asset-text",
			mode:       v1beta1.AssetConfigMap,
			filter:     "example.txt",
			files:      1,
			errMatcher: gomega.BeNil(),
		},
		"Config-Map not found": {
			src:        "default/notFound",
			name:       "asset-text",
			mode:       v1beta1.AssetConfigMap,
			filter:     "example.txt",
			files:      0,
			errMatcher: gomega.HaveOccurred(),
		},
		"Bad src": {
			src:        "default:notExist",
			name:       "asset-text",
			mode:       v1beta1.AssetConfigMap,
			filter:     "",
			files:      0,
			errMatcher: gomega.HaveOccurred(),
		},
	} {
		t.Run(testName, func(t *testing.T) {
			// Given
			g := gomega.NewGomegaWithT(t)

			// When
			_, files, err := loader.Load(testData.src, testData.name, testData.mode, testData.filter)

			// Then
			g.Expect(err).To(testData.errMatcher)
			g.Expect(files).To(gomega.HaveLen(testData.files))
		})
	}
}

func fixConfigMap(name string, namespace string, data map[string]string, binaryData map[string][]byte) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: v1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data:       data,
		BinaryData: binaryData,
	}
}

func newFakeDynamicClient(objects ...runtime.Object) (*dynamicFake.FakeDynamicClient, error) {
	scheme := runtime.NewScheme()
	err := corev1.AddToScheme(scheme)
	if err != nil {
		return &dynamicFake.FakeDynamicClient{}, err
	}

	result := make([]runtime.Object, len(objects))
	for i, obj := range objects {
		converted, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		result[i] = &unstructured.Unstructured{Object: converted}
	}
	return dynamicFake.NewSimpleDynamicClient(scheme, result...), nil
}
