package machineautoscaler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var (
	ErrTargetMissingAnnotations	= errors.New("missing min or max annotation")
	ErrTargetAlreadyOwned		= errors.New("already owned by another MachineAutoscaler")
	ErrTargetMissingOwner		= errors.New("missing owner annotation")
	ErrTargetBadOwner		= errors.New("incorrectly formatted owner annotation")
)

func MachineTargetFromObject(obj runtime.Object) (*MachineTarget, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	target := &MachineTarget{Unstructured: unstructured.Unstructured{Object: u}}
	return target, nil
}

type MachineTarget struct{ unstructured.Unstructured }

func (mt *MachineTarget) NeedsUpdate(min, max int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	currentMin, currentMax, err := mt.GetLimits()
	if err != nil {
		return true
	}
	minDiff := min != currentMin
	maxDiff := max != currentMax
	return minDiff || maxDiff
}
func (mt *MachineTarget) SetLimits(min, max int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := mt.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[minSizeAnnotation] = strconv.Itoa(min)
	annotations[maxSizeAnnotation] = strconv.Itoa(max)
	mt.SetAnnotations(annotations)
}
func (mt *MachineTarget) RemoveLimits() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := []string{minSizeAnnotation, maxSizeAnnotation}
	return mt.RemoveAnnotations(annotations)
}
func (mt *MachineTarget) GetLimits() (min, max int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := mt.GetAnnotations()
	minString, minOK := annotations[minSizeAnnotation]
	maxString, maxOK := annotations[maxSizeAnnotation]
	if !minOK || !maxOK {
		return 0, 0, ErrTargetMissingAnnotations
	}
	min, err = strconv.Atoi(minString)
	if err != nil {
		return 0, 0, fmt.Errorf("bad min annotation: %s", minString)
	}
	max, err = strconv.Atoi(maxString)
	if err != nil {
		return 0, 0, fmt.Errorf("bad max annotation: %s", maxString)
	}
	return min, max, nil
}
func (mt *MachineTarget) SetOwner(owner metav1.Object) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := mt.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	ownerRef := types.NamespacedName{Namespace: owner.GetNamespace(), Name: owner.GetName()}
	if a, ok := annotations[MachineTargetOwnerAnnotation]; ok {
		if a != ownerRef.String() {
			return false, ErrTargetAlreadyOwned
		}
		return false, nil
	}
	annotations[MachineTargetOwnerAnnotation] = ownerRef.String()
	mt.SetAnnotations(annotations)
	return true, nil
}
func (mt *MachineTarget) RemoveOwner() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := []string{MachineTargetOwnerAnnotation}
	return mt.RemoveAnnotations(annotations)
}
func (mt *MachineTarget) GetOwner() (types.NamespacedName, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nn := types.NamespacedName{}
	annotations := mt.GetAnnotations()
	if annotations == nil {
		return nn, ErrTargetMissingOwner
	}
	owner, found := annotations[MachineTargetOwnerAnnotation]
	if !found {
		return nn, ErrTargetMissingOwner
	}
	parts := strings.Split(owner, string(types.Separator))
	if len(parts) != 2 {
		return nn, ErrTargetBadOwner
	}
	nn.Namespace, nn.Name = parts[0], parts[1]
	return nn, nil
}
func (mt *MachineTarget) RemoveAnnotations(keys []string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := mt.GetAnnotations()
	modified := false
	for _, key := range keys {
		if _, found := annotations[key]; found {
			delete(annotations, key)
			modified = true
		}
	}
	mt.SetAnnotations(annotations)
	return modified
}
func (mt *MachineTarget) Finalize() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	limitsModified := mt.RemoveLimits()
	ownerModified := mt.RemoveOwner()
	return limitsModified || ownerModified
}
func (mt *MachineTarget) NamespacedName() types.NamespacedName {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return types.NamespacedName{Name: mt.GetName(), Namespace: mt.GetNamespace()}
}
