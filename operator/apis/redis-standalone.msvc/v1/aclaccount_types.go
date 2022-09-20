package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"operators.kloudlite.io/lib/constants"
	rApi "operators.kloudlite.io/lib/operator"
)

// ACLAccountSpec defines the desired state of ACLAccount
type ACLAccountSpec struct {
	KeyPrefix      string `json:"keyPrefix"`
	ManagedSvcName string `json:"managedSvcName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.isReady",name=Ready,type=boolean
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date

// ACLAccount is the Schema for the aclaccounts API
type ACLAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ACLAccountSpec `json:"spec,omitempty"`
	Status rApi.Status    `json:"status,omitempty"`
}

func (ac *ACLAccount) GetStatus() *rApi.Status {
	return &ac.Status
}

func (ac *ACLAccount) GetEnsuredLabels() map[string]string {
	return map[string]string{
		"kloudlite.io/msvc.name": ac.Spec.ManagedSvcName,
		"kloudlite.io/mres.name": ac.Name,
	}
}
func (m *ACLAccount) GetEnsuredAnnotations() map[string]string {
	return map[string]string{
		constants.AnnotationKeys.GroupVersionKind: GroupVersion.WithKind("ACLAccount").String(),
	}
}

// +kubebuilder:object:root=true

// ACLAccountList contains a list of ACLAccount
type ACLAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ACLAccount `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ACLAccount{}, &ACLAccountList{})
}
