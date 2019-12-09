package testsuite

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/kyma-project/rafter/tests/asset-store/pkg/mockice"
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
	Namespace             string        `envconfig:"default=default"`
	BucketName            string        `envconfig:"default=test-bucket"`
	ClusterBucketName     string        `envconfig:"default=test-cluster-bucket"`
	AssetGroupName        string        `envconfig:"default=test-asset-group"`
	ClusterAssetGroupName string        `envconfig:"default=test-cluster-asset-group"`
	CommonAssetPrefix     string        `envconfig:"default=test"`
	UploadServiceUrl      string        `envconfig:"default=http://localhost:3000/v1/upload"`
	MockiceName           string        `envconfig:"default=rafter-test-svc"`
	WaitTimeout           time.Duration `envconfig:"default=2m"`
	Minio                 MinioConfig
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

	t *testing.T
	g *gomega.GomegaWithT

	assetDetails []assetData
	uploadResult *upload.Response

	systemBucketName string
	minioCli         *minio.Client
	cfg              Config

	testId      string
	stopMockice func()
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

	// parametrize this shit
	ag := newAssetGroup(dynamicCli, cfg.AssetGroupName, cfg.Namespace, cfg.BucketName, cfg.WaitTimeout, t.Logf)
	cag := newClusterAssetGroup(dynamicCli, cfg.ClusterAssetGroupName, cfg.ClusterBucketName, cfg.WaitTimeout, t.Logf)
	b := newBucket(dynamicCli, cfg.BucketName, cfg.Namespace, cfg.WaitTimeout, t.Logf)
	cb := newClusterBucket(dynamicCli, cfg.ClusterBucketName, cfg.WaitTimeout, t.Logf)
	a := newAsset(dynamicCli, cfg.Namespace, cfg.BucketName, cfg.WaitTimeout, t.Logf)
	ca := newClusterAsset(dynamicCli, cfg.ClusterBucketName, cfg.WaitTimeout, t.Logf)

	stopMockice := func() {
		mockice.Stop(dynamicCli, cfg.Namespace, cfg.MockiceName)
	}

	host, err := mockice.Start(dynamicCli, cfg.Namespace, cfg.MockiceName)
	if err != nil {
		return nil, errors.Wrap(err, "while creating Mockice client")
	}

	as := []assetData{{
		Name: "first-test-asset",
		URL:  mockice.ResourceURL(host),
		Mode: v1beta1.AssetSingle,
		Type: "markdown",
	}, {
		Name: "second-test-asset",
		URL:  fmt.Sprintf("http://rafter-rafter-controller-manager.%s.svc.cluster.local:8080/metrics", cfg.Namespace),
		Mode: v1beta1.AssetSingle,
		Type: "metrics",
	}}

	// 	- type: openapi
	// name: swagger
	// mode: single
	// url: https://petstore.swagger.io/v2/swagger.json

	return &TestSuite{
		namespace:         ns,
		bucket:            b,
		clusterBucket:     cb,
		fileUpload:        newTestData(cfg.UploadServiceUrl),
		asset:             a,
		clusterAsset:      ca,
		assetGroup:        ag,
		clusterAssetGroup: cag,
		assetDetails:      as,
		t:                 t,
		g:                 g,
		minioCli:          minioCli,
		testId:            "singularity",
		cfg:               cfg,
		stopMockice:       stopMockice,
	}, nil
}

func (t *TestSuite) Run() {

	// clean up leftovers from previous tests
	t.t.Log("Deleting old asset groups...")
	err := t.assetGroup.DeleteLeftovers(t.testId)
	failOnError(t.g, err)

	t.t.Log("Deleting old cluster asset groups...")
	err = t.clusterAssetGroup.DeleteLeftovers(t.testId)
	failOnError(t.g, err)

	t.t.Log("Deleting old cluster bucket...")
	err = t.clusterBucket.Delete(t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Deleting old bucket...")
	err = t.bucket.Delete(t.t.Log)
	failOnError(t.g, err)

	// setup environment
	// t.t.Log("Creating namespace...")
	// err = t.namespace.Create(t.t.Log)
	// failOnError(t.g, err)

	t.t.Log("Creating cluster bucket...")
	var resourceVersion string
	resourceVersion, err = t.clusterBucket.Create(t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Waiting for cluster bucket to have ready phase...")
	err = t.clusterBucket.WaitForStatusReady(resourceVersion, t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Creating bucket...")
	resourceVersion, err = t.bucket.Create(t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Waiting for bucket to have ready phase...")
	err = t.bucket.WaitForStatusReady(resourceVersion, t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Creating assetgroup...")
	resourceVersion, err = t.assetGroup.Create(t.assetDetails, t.testId, t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Waiting for assetgroup to have ready phase...")
	err = t.assetGroup.WaitForStatusReady(resourceVersion, t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Creating cluster asset group...")
	resourceVersion, err = t.clusterAssetGroup.Create(t.assetDetails, t.testId, t.t.Log)
	failOnError(t.g, err)

	t.t.Log("Waiting for cluster asset group to have ready phase...")
	err = t.clusterAssetGroup.WaitForStatusReady(resourceVersion, t.t.Log)
	failOnError(t.g, err)

	// t.t.Log("Uploading test files...")
	// uploadResult, err := t.uploadTestFiles()
	// failOnError(t.g, err)
	//
	// t.t.Log("Uploaded files:\n", uploadResult.UploadedFiles)
	//
	// t.uploadResult = uploadResult
	// t.systemBucketName = uploadResult.UploadedFiles[0].Bucket
	//
	// t.t.Log("Preparing metadata...")
	// t.assetDetails = convertToAssetResourceDetails(uploadResult, t.cfg.CommonAssetPrefix)
	//
	// t.t.Log("Creating assets...")
	// resourceVersion, err = t.asset.CreateMany(t.assetDetails, t.testId, t.t.Log)
	// failOnError(t.g, err)
	// t.t.Log("Waiting for assets to have ready phase...")
	// err = t.asset.WaitForStatusesReady(t.assetDetails, resourceVersion, t.t.Log)
	// failOnError(t.g, err)
	//
	// t.t.Log("Creating cluster assets...")
	// resourceVersion, err = t.clusterAsset.CreateMany(t.assetDetails, t.testId, t.t.Log)
	// failOnError(t.g, err)
	// t.t.Log("Waiting for cluster assets to have ready phase...")
	// err = t.clusterAsset.WaitForStatusesReady(t.assetDetails, resourceVersion, t.t.Log)
	// failOnError(t.g, err)
	//
	// t.t.Log(fmt.Sprintf("asset details:\n%v", t.assetDetails))
	// files, err := t.populateUploadedFiles(t.t.Log)
	// failOnError(t.g, err)
	//
	// t.t.Log("Verifying uploaded files...")
	// err = t.verifyUploadedFiles(files)
	// failOnError(t.g, err)
	//
	// t.t.Log("Removing assets...")
	// err = t.asset.DeleteLeftovers(t.testId, t.t.Log)
	// failOnError(t.g, err)
	//
	// t.t.Log("Removing cluster assets...")
	// err = t.clusterAsset.DeleteLeftovers(t.testId, t.t.Log)
	// failOnError(t.g, err)
	//
	// err = t.verifyDeletedFiles(files)
	// failOnError(t.g, err)
}

func (t *TestSuite) Cleanup() {
	t.t.Log("Cleaning up...")

	err := t.clusterBucket.Delete(t.t.Log)
	failOnError(t.g, err)

	err = t.bucket.Delete(t.t.Log)
	failOnError(t.g, err)

	err = t.clusterAssetGroup.Delete(t.t.Log)
	failOnError(t.g, err)

	err = t.assetGroup.Delete(t.t.Log)
	failOnError(t.g, err)

	t.stopMockice()
	// err = t.namespace.Delete(t.t.Log)
	// failOnError(t.g, err)
	//
	// err = deleteFiles(t.minioCli, t.uploadResult, t.t.Logf)
	// failOnError(t.g, err)
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

func (t *TestSuite) populateUploadedFiles(callbacks ...func(...interface{})) ([]uploadedFile, error) {
	var allFiles []uploadedFile
	assetFiles, err := t.asset.PopulateUploadFiles(t.assetDetails, callbacks...)
	if err != nil {
		return nil, err
	}

	t.g.Expect(assetFiles).NotTo(gomega.HaveLen(0))

	allFiles = append(allFiles, assetFiles...)

	clusterAssetFiles, err := t.clusterAsset.PopulateUploadFiles(t.assetDetails)
	if err != nil {
		return nil, err
	}

	t.g.Expect(clusterAssetFiles).NotTo(gomega.HaveLen(0))

	allFiles = append(allFiles, clusterAssetFiles...)

	return allFiles, nil
}

func (t *TestSuite) verifyUploadedFiles(files []uploadedFile) error {
	err := verifyUploadedAssets(files, t.t.Logf)
	if err != nil {
		return errors.Wrap(err, "while verifying uploaded files")
	}
	return nil
}

func (t *TestSuite) verifyDeletedFiles(files []uploadedFile) error {
	err := verifyDeletedAssets(files, t.t.Logf)
	if err != nil {
		return errors.Wrap(err, "while verifying deleted files")
	}
	return nil
}

func failOnError(g *gomega.GomegaWithT, err error) {
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
