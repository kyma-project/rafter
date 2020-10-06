package loader

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/client-go/dynamic"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/pkg/errors"
)

type loader struct {
	temporaryDir  string
	dynamicClient dynamic.Interface

	// for testing
	osRemoveAllFunc func(string) error
	osCreateFunc    func(name string) (*os.File, error)
	httpGetFunc     func(url string) (*http.Response, error)
	ioutilTempDir   func(dir, prefix string) (string, error)
}

//go:generate mockery -name=Loader -output=automock -outpkg=automock -case=underscore
type Loader interface {
	Load(src, assetName string, mode v1beta1.AssetMode, filter string) (string, []string, error)
	Clean(path string) error
}

func New(dynamicClient dynamic.Interface, temporaryDir string, verifySSL bool) Loader {
	if len(temporaryDir) == 0 {
		temporaryDir = os.TempDir()
	}

	if !verifySSL {
		cfg := &tls.Config{
			InsecureSkipVerify: true,
		}
		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: cfg,
		}
	}

	return &loader{
		temporaryDir:    temporaryDir,
		dynamicClient:   dynamicClient,
		osRemoveAllFunc: os.RemoveAll,
		osCreateFunc:    os.Create,
		httpGetFunc:     http.Get,
		ioutilTempDir:   ioutil.TempDir,
	}
}

func (l *loader) Load(src, assetName string, mode v1beta1.AssetMode, filter string) (string, []string, error) {
	switch mode {
	case v1beta1.AssetSingle:
		return l.loadSingle(src, assetName)
	case v1beta1.AssetPackage:
		return l.loadPackage(src, assetName, filter)
	case v1beta1.AssetConfigMap:
		return l.loadConfigMap(src, assetName, filter)
	}

	return "", nil, fmt.Errorf("not supported source mode %+v", mode)
}

func (l *loader) Clean(path string) error {
	return l.osRemoveAllFunc(path)
}

func (l *loader) download(destination, source string) error {
	file, err := l.osCreateFunc(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := l.httpGetFunc(source)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.New(response.Status)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (l *loader) fileName(source string) string {
	_, filename := path.Split(source)
	if len(filename) == 0 {
		filename = "asset"
	}

	// trim query params, fragments etc
	parsedUrl, err := url.Parse(filename)
	if err != nil {
		filename = "asset"
	} else {
		filename = parsedUrl.Path
	}

	return filename
}
