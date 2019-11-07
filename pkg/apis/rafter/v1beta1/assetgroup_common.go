package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type CommonAssetGroupSpec struct {
	DisplayName string         `json:"displayName,omitempty"`
	Description string         `json:"description,omitempty"`
	BucketRef   AssetBucketRef `json:"bucketRef,omitempty"`
	// +kubebuilder:validation:MinItems=1
	Sources []Source `json:"sources"`
}

type AssetGroupBucketRef struct {
	Name string `json:"name"`
}

// +kubebuilder:validation:Enum=single;package;index
type AssetGroupSourceMode string

const (
	AssetGroupSingle  AssetGroupSourceMode = "single"
	AssetGroupPackage AssetGroupSourceMode = "package"
	AssetGroupIndex   AssetGroupSourceMode = "index"
)

// +kubebuilder:validation:Pattern=^[a-z][a-zA-Z0-9-]*[a-zA-Z0-9]$
type AssetGroupSourceName string

// +kubebuilder:validation:Pattern=^[a-z][a-zA-Z0-9\._-]*[a-zA-Z0-9]$
type AssetGroupSourceType string

type Source struct {
	Name   AssetGroupSourceName `json:"name"`
	Type   AssetGroupSourceType `json:"type"`
	URL    string               `json:"url"`
	Mode   AssetGroupSourceMode `json:"mode"`
	Filter string               `json:"filter,omitempty"`
	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Ready;Failed
type AssetGroupPhase string

const (
	AssetGroupPending AssetGroupPhase = "Pending"
	AssetGroupReady   AssetGroupPhase = "Ready"
	AssetGroupFailed  AssetGroupPhase = "Failed"
)

type CommonAssetGroupStatus struct {
	Phase             AssetGroupPhase  `json:"phase"`
	Reason            AssetGroupReason `json:"reason,omitempty"`
	Message           string           `json:"message,omitempty"`
	LastHeartbeatTime metav1.Time      `json:"lastHeartbeatTime"`
}

type AssetGroupReason string

const (
	AssetGroupAssetCreated               AssetGroupReason = "AssetCreated"
	AssetGroupAssetCreationFailed        AssetGroupReason = "AssetCreationFailed"
	AssetGroupAssetsCreationFailed       AssetGroupReason = "AssetsCreationFailed"
	AssetGroupAssetsListingFailed        AssetGroupReason = "AssetsListingFailed"
	AssetGroupAssetDeleted               AssetGroupReason = "AssetDeleted"
	AssetGroupAssetDeletionFailed        AssetGroupReason = "AssetDeletionFailed"
	AssetGroupAssetsDeletionFailed       AssetGroupReason = "AssetsDeletionFailed"
	AssetGroupAssetUpdated               AssetGroupReason = "AssetUpdated"
	AssetGroupAssetUpdateFailed          AssetGroupReason = "AssetUpdateFailed"
	AssetGroupAssetsUpdateFailed         AssetGroupReason = "AssetsUpdateFailed"
	AssetGroupAssetsReady                AssetGroupReason = "AssetsReady"
	AssetGroupWaitingForAssets           AssetGroupReason = "WaitingForAssets"
	AssetGroupBucketError                AssetGroupReason = "BucketError"
	AssetGroupAssetsWebhookGetFailed     AssetGroupReason = "AssetsWebhookGetFailed"
	AssetGroupAssetsSpecValidationFailed AssetGroupReason = "AssetsSpecValidationFailed"
)

func (r AssetGroupReason) String() string {
	return string(r)
}

func (r AssetGroupReason) Message() string {
	switch r {
	case AssetGroupAssetCreated:
		return "Asset %s has been created"
	case AssetGroupAssetCreationFailed:
		return "Asset %s couldn't be created due to error %s"
	case AssetGroupAssetsCreationFailed:
		return "Assets couldn't be created due to error %s"
	case AssetGroupAssetsListingFailed:
		return "Assets couldn't be listed due to error %s"
	case AssetGroupAssetDeleted:
		return "Assets %s has been deleted"
	case AssetGroupAssetDeletionFailed:
		return "Assets %s couldn't be deleted due to error %s"
	case AssetGroupAssetsDeletionFailed:
		return "Assets couldn't be deleted due to error %s"
	case AssetGroupAssetUpdated:
		return "Asset %s has been updated"
	case AssetGroupAssetUpdateFailed:
		return "Asset %s couldn't be updated due to error %s"
	case AssetGroupAssetsUpdateFailed:
		return "Assets couldn't be updated due to error %s"
	case AssetGroupAssetsReady:
		return "Assets are ready to use"
	case AssetGroupWaitingForAssets:
		return "Waiting for assets to be in Ready phase"
	case AssetGroupBucketError:
		return "Couldn't ensure if bucket exist due to error %s"
	case AssetGroupAssetsWebhookGetFailed:
		return "Unable to get webhook configuration %s"
	case AssetGroupAssetsSpecValidationFailed:
		return "Invalid asset specification, %s"
	default:
		return ""
	}
}
