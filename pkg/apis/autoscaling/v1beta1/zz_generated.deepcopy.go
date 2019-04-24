package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *CrossVersionObjectReference) DeepCopyInto(out *CrossVersionObjectReference) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *CrossVersionObjectReference) DeepCopy() *CrossVersionObjectReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(CrossVersionObjectReference)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineAutoscaler) DeepCopyInto(out *MachineAutoscaler) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}
func (in *MachineAutoscaler) DeepCopy() *MachineAutoscaler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineAutoscaler)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineAutoscaler) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *MachineAutoscalerList) DeepCopyInto(out *MachineAutoscalerList) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MachineAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *MachineAutoscalerList) DeepCopy() *MachineAutoscalerList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineAutoscalerList)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineAutoscalerList) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *MachineAutoscalerSpec) DeepCopyInto(out *MachineAutoscalerSpec) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.ScaleTargetRef = in.ScaleTargetRef
	return
}
func (in *MachineAutoscalerSpec) DeepCopy() *MachineAutoscalerSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineAutoscalerStatus) DeepCopyInto(out *MachineAutoscalerStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	if in.LastTargetRef != nil {
		in, out := &in.LastTargetRef, &out.LastTargetRef
		*out = new(CrossVersionObjectReference)
		**out = **in
	}
	return
}
func (in *MachineAutoscalerStatus) DeepCopy() *MachineAutoscalerStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}
