package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/kyma-project/rafter/internal/handler/assetgroup"
	"github.com/kyma-project/rafter/internal/webhookconfig"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	cmsv1alpha1 "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/pkg/errors"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AssetGroupReconciler reconciles a AssetGroup object
type AssetGroupReconciler struct {
	client.Client
	Log logr.Logger

	relistInterval   time.Duration
	recorder         record.EventRecorder
	assetSvc         assetgroup.AssetService
	bucketSvc        assetgroup.BucketService
	webhookConfigSvc webhookconfig.AssetWebhookConfigService
}

type AssetGroupConfig struct {
	RelistInterval time.Duration `envconfig:"default=5m"`
	BucketRegion   string        `envconfig:"-"`
}

func NewAssetGroup(config AssetGroupConfig, log logr.Logger, mgr ctrl.Manager, webhookConfigSvc webhookconfig.AssetWebhookConfigService) *AssetGroupReconciler {
	assetService := newAssetService(mgr.GetClient(), mgr.GetScheme())
	bucketService := newBucketService(mgr.GetClient(), mgr.GetScheme(), config.BucketRegion)

	return &AssetGroupReconciler{
		Client:           mgr.GetClient(),
		Log:              log,
		relistInterval:   config.RelistInterval,
		recorder:         mgr.GetEventRecorderFor("assetgroup-controller"),
		assetSvc:         assetService,
		bucketSvc:        bucketService,
		webhookConfigSvc: webhookConfigSvc,
	}
}

// Reconcile reads that state of the cluster for a AssetGroup object and makes changes based on the state read
// Automatically generate RBAC rules to allow the Controller to read and write AssetGroups, Assets, and Buckets
// +kubebuilder:rbac:groups=cms.kyma-project.io,resources=assetgroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cms.kyma-project.io,resources=assetgroups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rafter.kyma-project.io,resources=assets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rafter.kyma-project.io,resources=assets/status,verbs=get;list
// +kubebuilder:rbac:groups=rafter.kyma-project.io,resources=buckets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=rafter.kyma-project.io,resources=buckets/status,verbs=get;list
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;watch

func (r *AssetGroupReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	instance := &cmsv1alpha1.AssetGroup{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	assetGroupLogger := r.Log.WithValues("kind", instance.GetObjectKind().GroupVersionKind().Kind, "namespace", instance.GetNamespace(), "name", instance.GetName())
	commonHandler := assetgroup.New(assetGroupLogger, r.recorder, r.assetSvc, r.bucketSvc, r.webhookConfigSvc)
	commonStatus, err := commonHandler.Handle(ctx, instance, instance.Spec.CommonAssetGroupSpec, instance.Status.CommonAssetGroupStatus)
	if updateErr := r.updateStatus(ctx, instance, commonStatus); updateErr != nil {
		finalErr := updateErr
		if err != nil {
			finalErr = errors.Wrapf(err, "along with update error %s", updateErr.Error())
		}
		return ctrl.Result{}, finalErr
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{
		RequeueAfter: r.relistInterval,
	}, nil
}

func (r *AssetGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cmsv1alpha1.AssetGroup{}).
		Owns(&v1beta1.Asset{}).
		Complete(r)
}

func (r *AssetGroupReconciler) updateStatus(ctx context.Context, instance *cmsv1alpha1.AssetGroup, commonStatus *cmsv1alpha1.CommonAssetGroupStatus) error {
	if commonStatus == nil {
		return nil
	}

	copy := instance.DeepCopy()
	copy.Status.CommonAssetGroupStatus = *commonStatus

	return r.Status().Update(ctx, copy)
}
