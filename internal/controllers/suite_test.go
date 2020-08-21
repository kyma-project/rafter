package controllers

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	"github.com/kyma-project/rafter/internal/webhookconfig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	assethook "github.com/kyma-project/rafter/internal/assethook/automock"
	loader "github.com/kyma-project/rafter/internal/loader/automock"
	store "github.com/kyma-project/rafter/internal/store/automock"
	assetstorev1beta1 "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var webhookConfigSvc webhookconfig.AssetWebhookConfigService
var fixParameters = runtime.RawExtension{Raw: []byte(`{"json":"true","complex":{"data":"true"}}`)}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager, g *GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		g.Expect(mgr.Start(stop)).NotTo(HaveOccurred())
		wg.Done()
	}()
	return stop, wg
}

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = assetstorev1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	webhookConfigSvc = initWebhookConfigService(webhookconfig.Config{CfgMapName: "test", CfgMapNamespace: "test"}, cfg)

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func initWebhookConfigService(webhookCfg webhookconfig.Config, config *rest.Config) webhookconfig.AssetWebhookConfigService {
	dc, err := dynamic.NewForConfig(config)
	Expect(err).To(Succeed())

	configmapsResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
	resourceGetter := dc.Resource(configmapsResource)
	webhookCfgService := webhookconfig.New(resourceGetter, webhookCfg.CfgMapName, webhookCfg.CfgMapNamespace)

	return webhookCfgService
}

type MockContainer struct {
	Store     *store.Store
	Extractor *assethook.MetadataExtractor
	Mutator   *assethook.Mutator
	Validator *assethook.Validator
	Loader    *loader.Loader
}

func NewMockContainer() *MockContainer {
	return &MockContainer{
		Store:     new(store.Store),
		Extractor: new(assethook.MetadataExtractor),
		Mutator:   new(assethook.Mutator),
		Validator: new(assethook.Validator),
		Loader:    new(loader.Loader),
	}
}

func (c *MockContainer) AssertExpetactions(t GinkgoTInterface) {
	c.Store.AssertExpectations(t)
	c.Extractor.AssertExpectations(t)
	c.Mutator.AssertExpectations(t)
	c.Validator.AssertExpectations(t)
	c.Loader.AssertExpectations(t)
}

func validateReconcilation(err error, result controllerruntime.Result) {
	Expect(err).ToNot(HaveOccurred())
	Expect(result.Requeue).To(BeFalse())
	Expect(result.RequeueAfter).To(Equal(60 * time.Hour))
}

func validateAsset(status assetstorev1beta1.CommonAssetStatus, meta controllerruntime.ObjectMeta, expectedBaseURL string, files []string, expectedPhase assetstorev1beta1.AssetPhase, expectedReason assetstorev1beta1.AssetReason) {
	Expect(status.Phase).To(Equal(expectedPhase))
	Expect(status.Reason).To(Equal(expectedReason))
	Expect(status.AssetRef.BaseURL).To(Equal(expectedBaseURL))
	Expect(status.AssetRef.Files).To(HaveLen(len(files)))
	Expect(meta.Finalizers).To(ContainElement("test"))
}

func validateAssetGroup(status assetstorev1beta1.CommonAssetGroupStatus, meta controllerruntime.ObjectMeta, expectedPhase assetstorev1beta1.AssetGroupPhase, expectedReason assetstorev1beta1.AssetGroupReason) {
	Expect(status.Phase).To(Equal(expectedPhase))
	Expect(status.Reason).To(Equal(expectedReason))
}
