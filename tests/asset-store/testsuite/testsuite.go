package testsuite

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/kyma-project/rafter/tests/asset-store/mockice"
	"github.com/kyma-project/rafter/tests/asset-store/pkg/upload"
	"github.com/minio/minio-go"

	"github.com/kyma-project/rafter/tests/asset-store/pkg/namespace"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"k8s.io/client-go/dynamic"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type Config struct {
	Namespace         string        `envconfig:"default=default"`
	BucketName        string        `envconfig:"default=test-bucket"`
	ClusterBucketName string        `envconfig:"default=test-cluster-bucket"`
	CommonAssetPrefix string        `envconfig:"default=test"`
	UploadServiceUrl  string        `envconfig:"default=http://localhost:3000/v1/upload"`
	WaitTimeout       time.Duration `envconfig:"default=2m"`
	Minio             MinioConfig
}

type TestSuite struct {
	namespace         *namespace.Namespace
	bucket            *bucket
	clusterBucket     *clusterBucket
	fileUpload        *testData
	asset             *asset
	clusterAsset      *clusterAsset
	assetGroup        *assetGroup
	clusterAssetGroup *clusterAssetGroup
	t                 *testing.T
	g                 *gomega.GomegaWithT

	assetDetails []assetData
	uploadResult *upload.Response

	systemBucketName string
	minioCli         *minio.Client
	cfg              Config
	stopMockice      func()
}

func New(restConfig *rest.Config, cfg Config, t *testing.T, g *gomega.GomegaWithT) (*TestSuite, error) {
	coreCli, err := corev1.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while creating K8s Core client")
	}

	dynamicCli, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while creating K8s Dynamic client")
	}

	minioCli, err := minio.New(cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.UseSSL)
	if err != nil {
		return nil, errors.Wrap(err, "while creating Minio client")
	}

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	minioCli.SetCustomTransport(transCfg)

	ns := namespace.New(coreCli, cfg.Namespace)
	ag := newAssetGroup(dynamicCli, "example-asset-group", cfg.Namespace, cfg.BucketName, cfg.WaitTimeout, t.Logf)
	cag := newClusterAssetGroup(dynamicCli, "cluster-asset-group", cfg.ClusterBucketName, cfg.WaitTimeout, t.Logf)
	b := newBucket(dynamicCli, cfg.BucketName, cfg.Namespace, cfg.WaitTimeout, t.Logf)
	cb := newClusterBucket(dynamicCli, cfg.ClusterBucketName, cfg.WaitTimeout, t.Logf)
	a := newAsset(dynamicCli, cfg.Namespace, cfg.BucketName, cfg.WaitTimeout, t.Logf)
	ca := newClusterAsset(dynamicCli, cfg.ClusterBucketName, cfg.WaitTimeout, t.Logf)

	stopMockice := func() {
		mockice.Stop(dynamicCli, cfg.Namespace, mockice.SvcName)
	}
	// stopMockice()

	// time.Sleep(5 * time.Second)

	host, err := mockice.Start(dynamicCli, cfg.Namespace, mockice.SvcName)
	if err != nil {
		return nil, errors.Wrap(err, "while creating Mockice client")
	}

	// stopMockice := func() {
	// 	mockice.Stop(dynamicCli, "default", "mockice-test-svc")
	// }

	as := []assetData{{
		Name: "markdownOne",
		URL:  mockice.ResourceURL(host),
		Mode: v1beta1.AssetSingle,
		Type: "markdown",
	}, {
		Name: "markdownTwo",
		URL:  mockice.ResourceURL(host),
		Mode: v1beta1.AssetSingle,
		Type: "markdown",
	}}

	return &TestSuite{
		assetGroup:        ag,
		namespace:         ns,
		bucket:            b,
		clusterBucket:     cb,
		fileUpload:        newTestData(cfg.UploadServiceUrl),
		asset:             a,
		clusterAsset:      ca,
		clusterAssetGroup: cag,
		t:                 t,
		g:                 g,
		minioCli:          minioCli,
		assetDetails:      as,
		cfg:               cfg,
		stopMockice:       stopMockice,
	}, nil
}

func (t *TestSuite) Run() {
	err := t.namespace.Create()
	failOnError(t.g, err)

	t.t.Log("Creating Buckets...")
	err = t.createBuckets()
	failOnError(t.g, err)

	t.t.Log("Waiting for ready Buckets...")
	err = t.waitForBucketsReady()
	failOnError(t.g, err)

	t.t.Log("Creating ClusterBuckets...")
	err = t.createClusterBuckets()
	failOnError(t.g, err)

	t.t.Log("Waiting for ready ClusterBuckets...")
	err = t.waitForClusterBucketsReady()
	failOnError(t.g, err)

	t.t.Log("Creating AssetGroup...")
	err = t.createAssetGroup()
	failOnError(t.g, err)

	t.t.Log("Waiting for ready AssetGroup...")
	err = t.waitForAssetGroupReady()
	failOnError(t.g, err)

	t.t.Log("Creating ClusterAssetGroup")
	err = t.createClusterAssetGroup()
	failOnError(t.g, err)

	t.t.Log("Waiting for ready ClusterAssetGroup...")
	err = t.waitForClusterAssetGroupReady()
	failOnError(t.g, err)

	// t.t.Log("Uploading test files...")
	// uploadResult, err := t.uploadTestFiles()
	// failOnError(t.g, err)
	//
	// t.uploadResult = uploadResult
	//
	// t.systemBucketName = t.systemBucketNameFromUploadResult(uploadResult)
	//
	// t.t.Log("Creating assets...")
	// err = t.createAssets(uploadResult)
	// failOnError(t.g, err)
	//
	// t.t.Log("Waiting for ready assets...")
	// err = t.waitForAssetsReady()
	// failOnError(t.g, err)
	//
	// files, err := t.populateUploadedFiles()
	// failOnError(t.g, err)
	//
	// t.t.Log("Verifying uploaded files...")
	// err = t.verifyUploadedFiles(files)
	// failOnError(t.g, err)
}

func (t *TestSuite) Cleanup() {
	t.t.Log("Cleaning up...")

	// err := t.deleteAssets()
	// failOnError(t.g, err)

	// err = t.waitForAssetsDeleted()
	// failOnError(t.g, err)

	// err = t.deleteClusterAssets()
	// t.t.Log("Verifying if files have been deleted...")
	// err = t.verifyDeletedFiles(files)
	// failOnError(t.g, err)

	err := t.deleteClusterAssetGroups()
	failOnError(t.g, err)

	err = t.deleteAssetGroups()
	failOnError(t.g, err)

	err = t.deleteBuckets()
	failOnError(t.g, err)

	err = t.deleteClusterBuckets()
	failOnError(t.g, err)

	// err = t.namespace.Delete()
	// failOnError(t.g, err)

	t.stopMockice()
}

func (t *TestSuite) createBuckets() error {
	return t.bucket.Create()
}

func (t *TestSuite) createClusterBuckets() error {
	return t.clusterBucket.Create()
}

func (t *TestSuite) createAssetGroup() error {
	return t.assetGroup.Create(t.assetDetails)
}

func (t *TestSuite) createClusterAssetGroup() error {
	return t.clusterAssetGroup.Create(t.assetDetails)
}

func (t *TestSuite) systemBucketNameFromUploadResult(result *upload.Response) string {
	return result.UploadedFiles[0].Bucket
}

func (t *TestSuite) uploadTestFiles() (*upload.Response, error) {
	uploadResult, err := t.fileUpload.Upload()
	if err != nil {
		return nil, err
	}

	if len(uploadResult.Errors) > 0 {
		return nil, fmt.Errorf("during file upload: %+v", uploadResult.Errors)
	}

	return uploadResult, nil
}

func (t *TestSuite) createAssets(uploadResult *upload.Response) error {
	t.assetDetails = convertToAssetResourceDetails(uploadResult, t.cfg.CommonAssetPrefix)

	err := t.asset.CreateMany(t.assetDetails)
	if err != nil {
		return err
	}

	err = t.clusterAsset.CreateMany(t.assetDetails)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestSuite) waitForAssetsReady() error {
	return t.asset.WaitForStatusesReady(t.assetDetails)
}

func (t *TestSuite) waitForClusterAssetsReady() error {
	return t.clusterAsset.WaitForStatusesReady(t.assetDetails)
}

// func (t *TestSuite) waitForAssetsDeleted() error {
// 	err := t.asset.WaitForDeletedResources(t.assetDetails)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = t.clusterAsset.WaitForDeletedResources(t.assetDetails)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func (t *TestSuite) populateUploadedFiles() ([]uploadedFile, error) {
// 	var allFiles []uploadedFile
// 	assetFiles, err := t.asset.PopulateUploadFiles(t.assetDetails)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	t.g.Expect(assetFiles).NotTo(gomega.HaveLen(0))
//
// 	allFiles = append(allFiles, assetFiles...)
//
// 	clusterAssetFiles, err := t.clusterAsset.PopulateUploadFiles(t.assetDetails)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	t.g.Expect(clusterAssetFiles).NotTo(gomega.HaveLen(0))
//
// 	allFiles = append(allFiles, clusterAssetFiles...)
//
// 	return allFiles, nil
// }

// func (t *TestSuite) verifyUploadedFiles(files []uploadedFile) error {
// 	err := verifyUploadedAssets(files, t.t.Logf)
// 	if err != nil {
// 		return errors.Wrap(err, "while verifying uploaded files")
// 	}
//
// 	return nil
// }
//
// func (t *TestSuite) verifyDeletedFiles(files []uploadedFile) error {
// 	err := verifyDeletedAssets(files, t.t.Logf)
// 	if err != nil {
// 		return errors.Wrap(err, "while verifying deleted files")
// 	}
//
// 	return nil
// }

func (t *TestSuite) waitForBucketsReady() error {
	return t.bucket.WaitForStatusReady()
}

func (t *TestSuite) waitForClusterBucketsReady() error {
	return t.clusterBucket.WaitForStatusReady()
}

func (t *TestSuite) waitForAssetGroupReady() error {
	return t.assetGroup.WaitForStatusReady()
}

func (t *TestSuite) waitForClusterAssetGroupReady() error {
	return t.clusterAssetGroup.WaitForStatusReady()
}

// func (t *TestSuite) deleteAssets() error {
// 	err := t.asset.DeleteMany(t.assetDetails)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = t.clusterAsset.DeleteMany(t.assetDetails)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (t *TestSuite) deleteBuckets() error {
	return t.bucket.Delete()
}

func (t *TestSuite) deleteClusterBuckets() error {
	return t.clusterBucket.Delete()
}

func (t *TestSuite) deleteAssetGroups() error {
	return t.assetGroup.Delete()
}

func (t *TestSuite) deleteClusterAssetGroups() error {
	return t.clusterAssetGroup.Delete()
}

func failOnError(g *gomega.GomegaWithT, err error) {
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
