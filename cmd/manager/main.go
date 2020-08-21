package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/minio/minio-go"
	"github.com/vrischmann/envconfig"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports

	"github.com/kyma-project/rafter/internal/assethook"
	"github.com/kyma-project/rafter/internal/controllers"
	"github.com/kyma-project/rafter/internal/loader"
	"github.com/kyma-project/rafter/internal/store"
	"github.com/kyma-project/rafter/internal/webhookconfig"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = v1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

type Config struct {
	Store               store.Config
	Loader              loader.Config
	Webhook             assethook.Config
	Asset               controllers.AssetConfig
	ClusterAsset        controllers.ClusterAssetConfig
	Bucket              controllers.BucketConfig
	ClusterBucket       controllers.ClusterBucketConfig
	AssetGroup          controllers.AssetGroupConfig
	ClusterAssetGroup   controllers.ClusterAssetGroupConfig
	WebhookConfigMap    webhookconfig.Config
	BucketRegion        string `envconfig:"optional"`
	ClusterBucketRegion string `envconfig:"optional"`
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	cfg, err := loadConfig("APP")
	if err != nil {
		setupLog.Error(err, "unable to load config")
		os.Exit(1)
	}

	httpClient := &http.Client{}
	minioClient, err := minio.New(cfg.Store.Endpoint, cfg.Store.AccessKey, cfg.Store.SecretKey, cfg.Store.UseSSL)
	if err != nil {
		setupLog.Error(err, "unable initialize Minio client")
		os.Exit(1)
	}

	restConfig := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "rafter-controller-leader-election-helper",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		setupLog.Error(err, "unable to initialize dynamic client")
		os.Exit(1)
	}

	container := &controllers.Container{
		Manager:   mgr,
		Store:     store.New(minioClient, cfg.Store.UploadWorkersCount),
		Loader:    loader.New(dynamicClient, cfg.Loader.TemporaryDirectory, cfg.Loader.VerifySSL),
		Validator: assethook.NewValidator(httpClient, cfg.Webhook.ValidationTimeout, cfg.Webhook.ValidationWorkersCount),
		Mutator:   assethook.NewMutator(httpClient, cfg.Webhook.MutationTimeout, cfg.Webhook.MutationWorkersCount),
		Extractor: assethook.NewMetadataExtractor(httpClient, cfg.Webhook.MetadataExtractionTimeout),
	}

	webhookSvc := initWebhookConfigService(cfg.WebhookConfigMap, dynamicClient)

	if err = controllers.NewClusterAsset(cfg.ClusterAsset, ctrl.Log.WithName("controllers").WithName("ClusterAsset"), container).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterAsset")
		os.Exit(1)
	}
	if err = controllers.NewClusterBucket(cfg.ClusterBucket, ctrl.Log.WithName("controllers").WithName("ClusterBucket"), container).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterBucket")
		os.Exit(1)
	}
	if err = controllers.NewAsset(cfg.Asset, ctrl.Log.WithName("controllers").WithName("Asset"), container).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Asset")
		os.Exit(1)
	}
	if err = controllers.NewBucket(cfg.Bucket, ctrl.Log.WithName("controllers").WithName("Bucket"), container).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Bucket")
		os.Exit(1)
	}
	if err = controllers.NewClusterAssetGroup(cfg.ClusterAssetGroup, ctrl.Log.WithName("controllers").WithName("ClusterAssetGroup"), mgr, webhookSvc).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterAssetGroup")
		os.Exit(1)
	}
	if err = controllers.NewAssetGroup(cfg.AssetGroup, ctrl.Log.WithName("controllers").WithName("AssetGroup"), mgr, webhookSvc).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AssetGroup")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func loadConfig(prefix string) (Config, error) {
	cfg := Config{}
	err := envconfig.InitWithPrefix(&cfg, prefix)
	if err != nil {
		return cfg, err
	}
	cfg.Bucket.ExternalEndpoint = cfg.Store.ExternalEndpoint
	cfg.ClusterBucket.ExternalEndpoint = cfg.Store.ExternalEndpoint
	cfg.ClusterAssetGroup.BucketRegion = cfg.ClusterBucketRegion
	cfg.AssetGroup.BucketRegion = cfg.BucketRegion
	return cfg, nil
}

func initWebhookConfigService(webhookCfg webhookconfig.Config, dc dynamic.Interface) webhookconfig.AssetWebhookConfigService {
	configmapsResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
	resourceGetter := dc.Resource(configmapsResource).Namespace(webhookCfg.CfgMapNamespace)

	webhookCfgService := webhookconfig.New(resourceGetter, webhookCfg.CfgMapName, webhookCfg.CfgMapNamespace)
	return webhookCfgService
}
