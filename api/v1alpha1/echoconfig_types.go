package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type EchoConfigSpec struct {
	// Message is mirrored into the ConfigMap's data.message
	// +kubebuilder:validation:MinLength=1
	Message string `json:"message"`
}

type EchoConfigStatus struct {
	// Name of the managed ConfigMap
	ConfigMapName string `json:"configMapName,omitempty"`
	// For knowing if status matches current spec generation
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=echoconfigs,scope=Namespaced,shortName=echo
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.spec.message`
// +kubebuilder:printcolumn:name="ConfigMap",type=string,JSONPath=`.status.configMapName`
type EchoConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EchoConfigSpec   `json:"spec,omitempty"`
	Status            EchoConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type EchoConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EchoConfig `json:"items"`
}
