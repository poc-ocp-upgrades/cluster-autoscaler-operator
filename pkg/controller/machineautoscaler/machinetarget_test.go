package machineautoscaler

import (
	"testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	TargetName	= "test-name"
	TargetNamespace	= "test-namespace"
)

type TargetOwner struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`
}

func NewTargetOwner(namespace, name string) *TargetOwner {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TargetOwner{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
}
func NewTarget() *MachineTarget {
	_logClusterCodePath()
	defer _logClusterCodePath()
	firstGVK := DefaultSupportedTargetGVKs()[0]
	u := unstructured.Unstructured{}
	u.SetGroupVersionKind(firstGVK)
	u.SetName(TargetName)
	u.SetNamespace(TargetNamespace)
	target, err := MachineTargetFromObject(u.DeepCopyObject())
	if err != nil {
		panic(err)
	}
	return target
}
func TestNeedsUpdate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	target.SetLimits(4, 6)
	if !target.NeedsUpdate(2, 4) {
		t.Fatal("target should need update")
	}
	if target.NeedsUpdate(4, 6) {
		t.Fatal("target should not need update")
	}
	target.SetAnnotations(map[string]string{minSizeAnnotation: "not-an-int", maxSizeAnnotation: "not-an-int"})
	if !target.NeedsUpdate(1, 2) {
		t.Fatal("target should need update")
	}
}
func TestSetLimits(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	expectedMin, expectedMax := 2, 4
	target.SetLimits(expectedMin, expectedMax)
	min, max, err := target.GetLimits()
	if err != nil {
		t.Fatalf("error getting limits: %v", err)
	}
	if min != expectedMin || max != expectedMax {
		t.Fatalf("got %d-%d, want %d-%d", min, max, expectedMin, expectedMax)
	}
}
func TestGetLimits(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	_, _, err := target.GetLimits()
	if err != ErrTargetMissingAnnotations {
		t.Fatal("expected missing annotations error")
	}
	target.SetAnnotations(map[string]string{minSizeAnnotation: "not-an-int", maxSizeAnnotation: "4"})
	_, _, err = target.GetLimits()
	if err == nil {
		t.Fatal("expected bad annotations error")
	}
	target.SetAnnotations(map[string]string{minSizeAnnotation: "2", maxSizeAnnotation: "not-an-int"})
	_, _, err = target.GetLimits()
	if err == nil {
		t.Fatal("expected bad annotation error")
	}
	expectedMin, expectedMax := 2, 4
	target.SetLimits(expectedMin, expectedMax)
	min, max, err := target.GetLimits()
	if err != nil {
		t.Fatal("error getting limits")
	}
	if min != 2 || max != 4 {
		t.Fatalf("got %d-%d, want %d-%d", min, max, expectedMin, expectedMax)
	}
}
func TestRemoveLimits(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	target.SetLimits(2, 4)
	target.RemoveLimits()
	annotations := target.GetAnnotations()
	_, minOK := annotations[minSizeAnnotation]
	_, maxOK := annotations[maxSizeAnnotation]
	if minOK || maxOK {
		t.Fatal("found annotations after removal")
	}
}
func TestSetOwner(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	owner := NewTargetOwner("owner", "owner")
	otherOwner := NewTargetOwner("other-owner", "other-owner")
	modified, err := target.SetOwner(owner)
	if err != nil {
		t.Fatalf("error setting owner: %v", err)
	}
	if !modified {
		t.Fatal("setting new owner did not report modifed")
	}
	modified, err = target.SetOwner(owner)
	if err != nil {
		t.Fatalf("error setting owner: %v", err)
	}
	if modified {
		t.Fatal("setting same owner reported modifed")
	}
	_, err = target.SetOwner(otherOwner)
	if err != ErrTargetAlreadyOwned {
		t.Fatal("changing owner did not report ErrTargetAlreadyOwned")
	}
}
func TestRemoveOwner(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	owner := NewTargetOwner("owner", "owner")
	if _, err := target.SetOwner(owner); err != nil {
		t.Fatalf("error setting owner: %v", err)
	}
	target.RemoveOwner()
	annotations := target.GetAnnotations()
	if _, ok := annotations[MachineTargetOwnerAnnotation]; ok {
		t.Fatal("found owner annotation after removal")
	}
}
func TestGetOwner(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	nn, err := target.GetOwner()
	if err != ErrTargetMissingOwner {
		t.Errorf("target with no owner did not report ErrTargetMissingOwner")
	}
	owner := NewTargetOwner("owner", "owner")
	if _, err := target.SetOwner(owner); err != nil {
		t.Fatalf("error setting owner: %v", err)
	}
	nn, err = target.GetOwner()
	if err != nil {
		t.Fatalf("failed to get owner: %v", err)
	}
	if nn.Name != "owner" || nn.Namespace != "owner" {
		t.Error("target returned unexpected owner")
	}
	target.SetAnnotations(map[string]string{MachineTargetOwnerAnnotation: "too/many/parts/here"})
	nn, err = target.GetOwner()
	if err != ErrTargetBadOwner {
		t.Errorf("target with bad owner did not report ErrTargetBadOwner")
	}
}
func TestFinalize(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	owner := NewTargetOwner("owner", "owner")
	if _, err := target.SetOwner(owner); err != nil {
		t.Fatalf("error setting owner: %v", err)
	}
	target.SetLimits(4, 6)
	modified := target.Finalize()
	annotations := target.GetAnnotations()
	_, minOK := annotations[minSizeAnnotation]
	_, maxOK := annotations[maxSizeAnnotation]
	_, ownerOk := annotations[MachineTargetOwnerAnnotation]
	if minOK || maxOK || ownerOk {
		t.Errorf("Annotations present after Finailze()")
	}
	if !modified {
		t.Errorf("Finailze() did not report modification")
	}
	modified = target.Finalize()
	if modified {
		t.Errorf("Finailze() reported modification unnecessarily")
	}
}
func TestNamespacedName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := NewTarget()
	nn := target.NamespacedName()
	if nn.Name != TargetName {
		t.Errorf("NamespacedName() returned bad name. Got: %s, Want: %s", nn.Name, TargetName)
	}
	if nn.Namespace != TargetNamespace {
		t.Errorf("NamespacedName() returned bad namespace. Got: %s, Want: %s", nn.Namespace, TargetNamespace)
	}
}
