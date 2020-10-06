package loader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/onsi/gomega"
)

func TestLoader_Clean(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		loader := &loader{
			temporaryDir:    "/tmp",
			osRemoveAllFunc: removeAll,
		}

		// When
		err := loader.Clean("test")

		// Then
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Fail", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		loader := &loader{
			temporaryDir:    "/tmp",
			osRemoveAllFunc: removeAll,
		}

		// When
		err := loader.Clean("error1")

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
	})
}

func TestLoader_Load_NotSupported(t *testing.T) {
	// Given
	g := gomega.NewGomegaWithT(t)
	loader := &loader{
		temporaryDir:    "/tmp",
		osRemoveAllFunc: os.RemoveAll,
		osCreateFunc:    os.Create,
		httpGetFunc:     get,
		ioutilTempDir:   ioutil.TempDir,
	}

	// When
	_, files, err := loader.Load("test", "asset", "other", "")

	// Then
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(files).To(gomega.HaveLen(0))
}

func TestLoader_filename(t *testing.T) {
	// Given
	g := gomega.NewGomegaWithT(t)
	loader := &loader{
		temporaryDir:    "/tmp",
		osRemoveAllFunc: removeAll,
	}

	t.Run("Empty filename", func(t *testing.T) {
		// When
		filename := loader.fileName("")
		// Then
		g.Expect(filename).To(gomega.Equal("asset"))
	})

	t.Run("Filename with directory", func(t *testing.T) {
		// When
		filename := loader.fileName("directory/foobar")
		// Then
		g.Expect(filename).To(gomega.Equal("foobar"))
	})

	t.Run("Filename as url", func(t *testing.T) {
		// When
		filename := loader.fileName("directory/foobar?query=foo&query2=bar#fragment")
		// Then
		g.Expect(filename).To(gomega.Equal("foobar"))
	})
}

func removeAll(s string) error {
	if s == "error1" {
		return fmt.Errorf("nope")
	}

	return nil
}

func createError(name string) (*os.File, error) {
	return nil, fmt.Errorf("nope")

}

func get(url string) (*http.Response, error) {
	if url == "error3" {
		return nil, fmt.Errorf("nope")
	}

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("ala ma kota"))),
	}

	return response, nil
}

func tempDirError(dir, prefix string) (string, error) {
	return "", fmt.Errorf("nope")
}
