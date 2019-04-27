package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *ClusterAutoscaler) DeepCopyInto(out *ClusterAutoscaler) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}
func (in *ClusterAutoscaler) DeepCopy() *ClusterAutoscaler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ClusterAutoscaler)
	in.DeepCopyInto(out)
	return out
}
func (in *ClusterAutoscaler) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *ClusterAutoscalerList) DeepCopyInto(out *ClusterAutoscalerList) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *ClusterAutoscalerList) DeepCopy() *ClusterAutoscalerList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ClusterAutoscalerList)
	in.DeepCopyInto(out)
	return out
}
func (in *ClusterAutoscalerList) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *ClusterAutoscalerSpec) DeepCopyInto(out *ClusterAutoscalerSpec) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	if in.ResourceLimits != nil {
		in, out := &in.ResourceLimits, &out.ResourceLimits
		*out = new(ResourceLimits)
		(*in).DeepCopyInto(*out)
	}
	if in.ScaleDown != nil {
		in, out := &in.ScaleDown, &out.ScaleDown
		*out = new(ScaleDownConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.MaxPodGracePeriod != nil {
		in, out := &in.MaxPodGracePeriod, &out.MaxPodGracePeriod
		*out = new(int32)
		**out = **in
	}
	if in.PodPriorityThreshold != nil {
		in, out := &in.PodPriorityThreshold, &out.PodPriorityThreshold
		*out = new(int32)
		**out = **in
	}
	return
}
func (in *ClusterAutoscalerSpec) DeepCopy() *ClusterAutoscalerSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ClusterAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}
func (in *ClusterAutoscalerStatus) DeepCopyInto(out *ClusterAutoscalerStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ClusterAutoscalerStatus) DeepCopy() *ClusterAutoscalerStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ClusterAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}
func (in *GPULimit) DeepCopyInto(out *GPULimit) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *GPULimit) DeepCopy() *GPULimit {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(GPULimit)
	in.DeepCopyInto(out)
	return out
}
func (in *ResourceLimits) DeepCopyInto(out *ResourceLimits) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	if in.MaxNodesTotal != nil {
		in, out := &in.MaxNodesTotal, &out.MaxNodesTotal
		*out = new(int32)
		**out = **in
	}
	if in.Cores != nil {
		in, out := &in.Cores, &out.Cores
		*out = new(ResourceRange)
		**out = **in
	}
	if in.Memory != nil {
		in, out := &in.Memory, &out.Memory
		*out = new(ResourceRange)
		**out = **in
	}
	if in.GPUS != nil {
		in, out := &in.GPUS, &out.GPUS
		*out = make([]GPULimit, len(*in))
		copy(*out, *in)
	}
	return
}
func (in *ResourceLimits) DeepCopy() *ResourceLimits {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ResourceLimits)
	in.DeepCopyInto(out)
	return out
}
func (in *ResourceRange) DeepCopyInto(out *ResourceRange) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ResourceRange) DeepCopy() *ResourceRange {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ResourceRange)
	in.DeepCopyInto(out)
	return out
}
func (in *ScaleDownConfig) DeepCopyInto(out *ScaleDownConfig) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	if in.DelayAfterAdd != nil {
		in, out := &in.DelayAfterAdd, &out.DelayAfterAdd
		*out = new(string)
		**out = **in
	}
	if in.DelayAfterDelete != nil {
		in, out := &in.DelayAfterDelete, &out.DelayAfterDelete
		*out = new(string)
		**out = **in
	}
	if in.DelayAfterFailure != nil {
		in, out := &in.DelayAfterFailure, &out.DelayAfterFailure
		*out = new(string)
		**out = **in
	}
	if in.UnneededTime != nil {
		in, out := &in.UnneededTime, &out.UnneededTime
		*out = new(string)
		**out = **in
	}
	return
}
func (in *ScaleDownConfig) DeepCopy() *ScaleDownConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ScaleDownConfig)
	in.DeepCopyInto(out)
	return out
}
