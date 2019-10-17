package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AssetGroupSpec defines the desired state of AssetGroup
type AssetGroupSpec struct {
	CommonAssetGroupSpec `json:",inline"`
}

// AssetGroupStatus defines the observed state of AssetGroup
type AssetGroupStatus struct {
	CommonAssetGroupStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// AssetGroup is the Schema for the assetgroups API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type AssetGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AssetGroupSpec   `json:"spec,omitempty"`
	Status AssetGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AssetGroupList contains a list of AssetGroup
type AssetGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AssetGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AssetGroup{}, &AssetGroupList{})
}
