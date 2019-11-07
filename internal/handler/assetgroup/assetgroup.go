package assetgroup

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/kyma-project/rafter/internal/webhookconfig"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
)

type CommonAsset struct {
	v1.ObjectMeta
	Spec   v1beta1.CommonAssetSpec
	Status v1beta1.CommonAssetStatus
}

const (
	assetGroupLabel          = "rafter.kyma-project.io/asset-group"
	accessLabel              = "rafter.kyma-project.io/access"
	assetShortNameAnnotation = "rafter.kyma-project.io/asset-short-name"
	typeLabel                = "rafter.kyma-project.io/type"
)

var (
	errDuplicatedAssetName = errors.New("duplicated asset name")
)

//go:generate mockery -name=AssetService -output=automock -outpkg=automock -case=underscore
type AssetService interface {
	List(ctx context.Context, namespace string, labels map[string]string) ([]CommonAsset, error)
	Create(ctx context.Context, assetGroup v1.Object, commonAsset CommonAsset) error
	Update(ctx context.Context, commonAsset CommonAsset) error
	Delete(ctx context.Context, commonAsset CommonAsset) error
}

//go:generate mockery -name=BucketService -output=automock -outpkg=automock -case=underscore
type BucketService interface {
	List(ctx context.Context, namespace string, labels map[string]string) ([]string, error)
	Create(ctx context.Context, namespacedName types.NamespacedName, private bool, labels map[string]string) error
}

type Handler interface {
	Handle(ctx context.Context, instance ObjectMetaAccessor, spec v1beta1.CommonAssetGroupSpec, status v1beta1.CommonAssetGroupStatus) (*v1beta1.CommonAssetGroupStatus, error)
}

type ObjectMetaAccessor interface {
	v1.Object
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() runtime.Object
}

type assetgroupHandler struct {
	log              logr.Logger
	recorder         record.EventRecorder
	assetSvc         AssetService
	bucketSvc        BucketService
	webhookConfigSvc webhookconfig.AssetWebhookConfigService
}

func New(log logr.Logger, recorder record.EventRecorder, assetSvc AssetService, bucketSvc BucketService, webhookConfigSvc webhookconfig.AssetWebhookConfigService) Handler {
	return &assetgroupHandler{
		log:              log,
		recorder:         recorder,
		assetSvc:         assetSvc,
		bucketSvc:        bucketSvc,
		webhookConfigSvc: webhookConfigSvc,
	}
}

func (h *assetgroupHandler) Handle(ctx context.Context, instance ObjectMetaAccessor, spec v1beta1.CommonAssetGroupSpec, status v1beta1.CommonAssetGroupStatus) (*v1beta1.CommonAssetGroupStatus, error) {
	h.logInfof("Start common AssetGroup handling")
	defer h.logInfof("Finish common AssetGroup handling")

	err := h.validateSpec(spec)
	if err != nil {
		h.recordWarningEventf(instance, v1beta1.AssetGroupAssetsSpecValidationFailed, err.Error())
		return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupAssetsSpecValidationFailed, err.Error()), status), err
	}

	bucketName := spec.BucketRef.Name
	if bucketName == "" {
		bucketName, err = h.ensureBucketExits(ctx, instance.GetNamespace())

		if err != nil {
			h.recordWarningEventf(instance, v1beta1.AssetGroupBucketError, err.Error())
			return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupBucketError, err.Error()), status), err
		}
	}

	commonAssets, err := h.assetSvc.List(ctx, instance.GetNamespace(), h.buildLabels(instance.GetName(), ""))
	if err != nil {
		h.recordWarningEventf(instance, v1beta1.AssetGroupAssetsListingFailed, err.Error())
		return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupAssetsListingFailed, err.Error()), status), err
	}

	commonAssetsMap := h.convertToAssetMap(commonAssets)

	webhookCfg, err := h.webhookConfigSvc.Get(ctx)
	if err != nil {
		h.recordWarningEventf(instance, v1beta1.AssetGroupAssetsWebhookGetFailed, err.Error())
		return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupAssetsWebhookGetFailed, err.Error()), status), err
	}
	h.logInfof("Webhook configuration loaded")

	switch {
	case h.isOnChange(commonAssetsMap, spec, bucketName, webhookCfg):
		return h.onChange(ctx, instance, spec, status, commonAssetsMap, bucketName, webhookCfg)
	case h.isOnPhaseChange(commonAssetsMap, status):
		return h.onPhaseChange(instance, status, commonAssetsMap)
	default:
		h.logInfof("Instance is up-to-date, action not taken")
		return nil, nil
	}
}

func (h *assetgroupHandler) validateSpec(spec v1beta1.CommonAssetGroupSpec) error {
	h.logInfof("validating CommonAssetGroupSpec")
	names := map[v1beta1.AssetGroupSourceName]map[v1beta1.AssetGroupSourceType]struct{}{}
	for _, src := range spec.Sources {
		if nameTypes, exists := names[src.Name]; exists {
			if _, exists := nameTypes[src.Type]; exists {
				return errDuplicatedAssetName
			}
			names[src.Name][src.Type] = struct{}{}
			continue
		}
		names[src.Name] = map[v1beta1.AssetGroupSourceType]struct{}{}
		names[src.Name][src.Type] = struct{}{}
	}
	h.logInfof("CommonAssetGroupSpec validated")

	return nil
}

func (h *assetgroupHandler) ensureBucketExits(ctx context.Context, namespace string) (string, error) {
	h.logInfof("Listing buckets")
	labels := map[string]string{accessLabel: "public"}
	names, err := h.bucketSvc.List(ctx, namespace, labels)
	if err != nil {
		return "", err
	}

	bucketCount := len(names)
	if bucketCount > 1 {
		return "", fmt.Errorf("too many buckets with labels: %+v", labels)
	}
	if bucketCount == 1 {
		h.logInfof("Bucket %s already exits", names[0])
		return names[0], nil
	}

	name := h.generateBucketName(false)
	h.logInfof("Creating bucket %s", name)
	if err := h.bucketSvc.Create(ctx, types.NamespacedName{Name: name, Namespace: namespace}, false, labels); err != nil {
		return "", err
	}
	h.logInfof("Bucket created %s", name)

	return name, nil
}

func (h *assetgroupHandler) isOnChange(existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec, bucketName string, config webhookconfig.AssetWebhookConfigMap) bool {
	return h.shouldCreateAssets(existing, spec) || h.shouldDeleteAssets(existing, spec) || h.shouldUpdateAssets(existing, spec, bucketName, config)
}

func (h *assetgroupHandler) isOnPhaseChange(existing map[v1beta1.AssetGroupSourceName]CommonAsset, status v1beta1.CommonAssetGroupStatus) bool {
	return status.Phase != h.calculateAssetPhase(existing)
}

func (h *assetgroupHandler) shouldCreateAssets(existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec) bool {
	for _, source := range spec.Sources {
		if _, exists := existing[source.Name]; !exists {
			return true
		}
	}

	return false
}

func (h *assetgroupHandler) shouldUpdateAssets(existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec, bucketName string, config webhookconfig.AssetWebhookConfigMap) bool {
	for key, existingAsset := range existing {
		expectedSpec := findSource(spec.Sources, key, v1beta1.AssetGroupSourceType(existingAsset.Labels[typeLabel]))
		if expectedSpec == nil {
			continue
		}

		assetWhsMap := config[expectedSpec.Type]
		expected := h.convertToCommonAssetSpec(*expectedSpec, bucketName, assetWhsMap)
		if !reflect.DeepEqual(expected, existingAsset.Spec) {
			return true
		}
	}

	return false
}

func (h *assetgroupHandler) shouldDeleteAssets(existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec) bool {
	for key, existingAsset := range existing {
		if findSource(spec.Sources, key, v1beta1.AssetGroupSourceType(existingAsset.Labels[typeLabel])) == nil {
			return true
		}
	}

	return false
}

func (h *assetgroupHandler) onPhaseChange(instance ObjectMetaAccessor, status v1beta1.CommonAssetGroupStatus, existing map[v1beta1.AssetGroupSourceName]CommonAsset) (*v1beta1.CommonAssetGroupStatus, error) {
	phase := h.calculateAssetPhase(existing)
	h.logInfof("Updating phase to %s", phase)

	if phase == v1beta1.AssetGroupPending {
		h.recordNormalEventf(instance, v1beta1.AssetGroupWaitingForAssets)
		return h.buildStatus(phase, v1beta1.AssetGroupWaitingForAssets), nil
	}

	h.recordNormalEventf(instance, v1beta1.AssetGroupAssetsReady)
	return h.buildStatus(phase, v1beta1.AssetGroupAssetsReady), nil
}

func (h *assetgroupHandler) onChange(ctx context.Context, instance ObjectMetaAccessor, spec v1beta1.CommonAssetGroupSpec, status v1beta1.CommonAssetGroupStatus, existing map[v1beta1.AssetGroupSourceName]CommonAsset, bucketName string, cfg webhookconfig.AssetWebhookConfigMap) (*v1beta1.CommonAssetGroupStatus, error) {
	if err := h.createMissingAssets(ctx, instance, existing, spec, bucketName, cfg); err != nil {
		return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupAssetsCreationFailed, err.Error()), status), err
	}

	if err := h.updateOutdatedAssets(ctx, instance, existing, spec, bucketName, cfg); err != nil {
		return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupAssetsUpdateFailed, err.Error()), status), err
	}

	if err := h.deleteNotExisting(ctx, instance, existing, spec); err != nil {
		return h.onFailedStatus(h.buildStatus(v1beta1.AssetGroupFailed, v1beta1.AssetGroupAssetsDeletionFailed, err.Error()), status), err
	}

	h.recordNormalEventf(instance, v1beta1.AssetGroupWaitingForAssets)
	return h.buildStatus(v1beta1.AssetGroupPending, v1beta1.AssetGroupWaitingForAssets), nil
}

func (h *assetgroupHandler) createMissingAssets(ctx context.Context, instance ObjectMetaAccessor, existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec, bucketName string, cfg webhookconfig.AssetWebhookConfigMap) error {
	for _, spec := range spec.Sources {
		name := spec.Name
		if _, exists := existing[name]; exists {
			continue
		}

		if err := h.createAsset(ctx, instance, spec, bucketName, cfg[spec.Type]); err != nil {
			return err
		}
	}

	return nil
}

func (h *assetgroupHandler) createAsset(ctx context.Context, instance ObjectMetaAccessor, assetSpec v1beta1.Source, bucketName string, cfg webhookconfig.AssetWebhookConfig) error {
	commonAsset := CommonAsset{
		ObjectMeta: v1.ObjectMeta{
			Name:        h.generateFullAssetName(instance.GetName(), assetSpec.Name, assetSpec.Type),
			Namespace:   instance.GetNamespace(),
			Labels:      h.buildLabels(instance.GetName(), assetSpec.Type),
			Annotations: h.buildAnnotations(assetSpec.Name),
		},
		Spec: h.convertToCommonAssetSpec(assetSpec, bucketName, cfg),
	}

	h.logInfof("Creating asset %s", commonAsset.Name)
	if err := h.assetSvc.Create(ctx, instance, commonAsset); err != nil {
		h.recordWarningEventf(instance, v1beta1.AssetGroupAssetCreationFailed, commonAsset.Name, err.Error())
		return err
	}
	h.logInfof("Asset %s created", commonAsset.Name)
	h.recordNormalEventf(instance, v1beta1.AssetGroupAssetCreated, commonAsset.Name)

	return nil
}

func (h *assetgroupHandler) updateOutdatedAssets(ctx context.Context, instance ObjectMetaAccessor, existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec, bucketName string, cfg webhookconfig.AssetWebhookConfigMap) error {
	for key, existingAsset := range existing {
		expectedSpec := findSource(spec.Sources, key, v1beta1.AssetGroupSourceType(existingAsset.Labels[typeLabel]))
		if expectedSpec == nil {
			continue
		}

		h.logInfof("Updating asset %s", existingAsset.Name)
		expected := h.convertToCommonAssetSpec(*expectedSpec, bucketName, cfg[expectedSpec.Type])
		if reflect.DeepEqual(expected, existingAsset.Spec) {
			h.logInfof("Asset %s is up-to-date", existingAsset.Name)
			continue
		}

		existingAsset.Spec = expected
		if err := h.assetSvc.Update(ctx, existingAsset); err != nil {
			h.recordWarningEventf(instance, v1beta1.AssetGroupAssetUpdateFailed, existingAsset.Name, err.Error())
			return err
		}
		h.logInfof("Asset %s updated", existingAsset.Name)
		h.recordNormalEventf(instance, v1beta1.AssetGroupAssetUpdated, existingAsset.Name)
	}

	return nil
}

func findSource(slice []v1beta1.Source, sourceName v1beta1.AssetGroupSourceName, sourceType v1beta1.AssetGroupSourceType) *v1beta1.Source {
	for _, source := range slice {
		if source.Name == sourceName && source.Type == sourceType {
			return &source
		}
	}
	return nil
}

func (h *assetgroupHandler) deleteNotExisting(ctx context.Context, instance ObjectMetaAccessor, existing map[v1beta1.AssetGroupSourceName]CommonAsset, spec v1beta1.CommonAssetGroupSpec) error {
	for key, commonAsset := range existing {
		if findSource(spec.Sources, key, v1beta1.AssetGroupSourceType(commonAsset.Labels[typeLabel])) != nil {
			continue
		}

		h.logInfof("Deleting asset %s", commonAsset.Name)
		if err := h.assetSvc.Delete(ctx, commonAsset); err != nil {
			h.recordWarningEventf(instance, v1beta1.AssetGroupAssetDeletionFailed, commonAsset.Name, err.Error())
			return err
		}
		h.logInfof("Asset %s deleted", commonAsset.Name)
		h.recordNormalEventf(instance, v1beta1.AssetGroupAssetDeleted, commonAsset.Name)
	}

	return nil
}

func (h *assetgroupHandler) convertToAssetMap(assets []CommonAsset) map[v1beta1.AssetGroupSourceName]CommonAsset {
	result := make(map[v1beta1.AssetGroupSourceName]CommonAsset)

	for _, asset := range assets {
		assetShortName := asset.Annotations[assetShortNameAnnotation]
		result[v1beta1.AssetGroupSourceName(assetShortName)] = asset
	}

	return result
}

func (h *assetgroupHandler) convertToCommonAssetSpec(spec v1beta1.Source, bucketName string, cfg webhookconfig.AssetWebhookConfig) v1beta1.CommonAssetSpec {
	return v1beta1.CommonAssetSpec{
		Source: v1beta1.AssetSource{
			Mode:                     h.convertToAssetMode(spec.Mode),
			URL:                      spec.URL,
			Filter:                   spec.Filter,
			ValidationWebhookService: convertToAssetWebhookServices(cfg.Validations),
			MutationWebhookService:   convertToAssetWebhookServices(cfg.Mutations),
			MetadataWebhookService:   convertToWebhookService(cfg.MetadataExtractors),
		},
		BucketRef: v1beta1.AssetBucketRef{
			Name: bucketName,
		},
		Parameters: spec.Parameters,
	}
}

func convertToWebhookService(services []webhookconfig.WebhookService) []v1beta1.WebhookService {
	servicesLen := len(services)
	if servicesLen < 1 {
		return nil
	}
	result := make([]v1beta1.WebhookService, 0, servicesLen)
	for _, service := range services {
		result = append(result, v1beta1.WebhookService{
			Name:      service.Name,
			Namespace: service.Namespace,
			Endpoint:  service.Endpoint,
			Filter:    service.Filter,
		})
	}
	return result
}

func convertToAssetWebhookServices(services []webhookconfig.AssetWebhookService) []v1beta1.AssetWebhookService {
	servicesLen := len(services)
	if servicesLen < 1 {
		return nil
	}
	result := make([]v1beta1.AssetWebhookService, 0, servicesLen)
	for _, s := range services {
		result = append(result, v1beta1.AssetWebhookService{
			WebhookService: v1beta1.WebhookService{
				Name:      s.Name,
				Namespace: s.Namespace,
				Endpoint:  s.Endpoint,
				Filter:    s.Filter,
			},
			Parameters: s.Parameters,
		})
	}
	return result
}

func (h *assetgroupHandler) buildLabels(topicName string, assetType v1beta1.AssetGroupSourceType) map[string]string {
	labels := make(map[string]string)

	labels[assetGroupLabel] = topicName
	if assetType != "" {
		labels[typeLabel] = string(assetType)
	}

	return labels

}

func (h *assetgroupHandler) buildAnnotations(assetShortName v1beta1.AssetGroupSourceName) map[string]string {
	return map[string]string{
		assetShortNameAnnotation: string(assetShortName),
	}
}

func (h *assetgroupHandler) convertToAssetMode(mode v1beta1.AssetGroupSourceMode) v1beta1.AssetMode {
	switch mode {
	case v1beta1.AssetGroupIndex:
		return v1beta1.AssetIndex
	case v1beta1.AssetGroupPackage:
		return v1beta1.AssetPackage
	default:
		return v1beta1.AssetSingle
	}
}

func (h *assetgroupHandler) calculateAssetPhase(existing map[v1beta1.AssetGroupSourceName]CommonAsset) v1beta1.AssetGroupPhase {
	for _, asset := range existing {
		if asset.Status.Phase != v1beta1.AssetReady {
			return v1beta1.AssetGroupPending
		}
	}

	return v1beta1.AssetGroupReady
}

func (h *assetgroupHandler) buildStatus(phase v1beta1.AssetGroupPhase, reason v1beta1.AssetGroupReason, args ...interface{}) *v1beta1.CommonAssetGroupStatus {
	return &v1beta1.CommonAssetGroupStatus{
		Phase:             phase,
		Reason:            reason,
		Message:           fmt.Sprintf(reason.Message(), args...),
		LastHeartbeatTime: v1.Now(),
	}
}

func (h *assetgroupHandler) onFailedStatus(newStatus *v1beta1.CommonAssetGroupStatus, oldStatus v1beta1.CommonAssetGroupStatus) *v1beta1.CommonAssetGroupStatus {
	if newStatus.Phase == oldStatus.Phase && newStatus.Reason == oldStatus.Reason {
		return nil
	}

	return newStatus
}

func (h *assetgroupHandler) logInfof(message string, args ...interface{}) {
	h.log.Info(fmt.Sprintf(message, args...))
}

func (h *assetgroupHandler) recordNormalEventf(object ObjectMetaAccessor, reason v1beta1.AssetGroupReason, args ...interface{}) {
	h.recordEventf(object, "Normal", reason, args...)
}

func (h *assetgroupHandler) recordWarningEventf(object ObjectMetaAccessor, reason v1beta1.AssetGroupReason, args ...interface{}) {
	h.recordEventf(object, "Warning", reason, args...)
}

func (h *assetgroupHandler) recordEventf(object ObjectMetaAccessor, eventType string, reason v1beta1.AssetGroupReason, args ...interface{}) {
	h.recorder.Eventf(object, eventType, reason.String(), reason.Message(), args...)
}

func (h *assetgroupHandler) generateFullAssetName(assetGroupName string, assetShortName v1beta1.AssetGroupSourceName, assetType v1beta1.AssetGroupSourceType) string {
	assetTypeLower := strings.ToLower(string(assetType))
	assetShortNameLower := strings.ToLower(string(assetShortName))
	return h.appendSuffix(fmt.Sprintf("%s-%s-%s", assetGroupName, assetShortNameLower, assetTypeLower))
}

func (h *assetgroupHandler) generateBucketName(private bool) string {
	access := "public"
	if private {
		access = "private"
	}

	return h.appendSuffix(fmt.Sprintf("rafter-%s", access))
}

func (h *assetgroupHandler) appendSuffix(name string) string {
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 32)

	return fmt.Sprintf("%s-%s", name, suffix)
}
