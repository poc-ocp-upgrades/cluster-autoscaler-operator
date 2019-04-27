package clusterautoscaler

import (
	"context"
	"fmt"
	autoscalingv1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName		= "cluster-autoscaler-controller"
	caServiceAccount	= "cluster-autoscaler"
	caPriorityClassName	= "system-cluster-critical"
)

func NewReconciler(mgr manager.Manager, cfg *Config) *Reconciler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Reconciler{client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetRecorder(controllerName), config: cfg}
}

type Config struct {
	ReleaseVersion	string
	Name		string
	Namespace	string
	Image		string
	Replicas	int32
	CloudProvider	string
	Verbosity	int
	ExtraArgs	string
}

var _ reconcile.Reconciler = &Reconciler{}

type Reconciler struct {
	client		client.Client
	recorder	record.EventRecorder
	scheme		*runtime.Scheme
	config		*Config
}

func (r *Reconciler) AddToManager(mgr manager.Manager) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	p := predicate.Funcs{CreateFunc: func(e event.CreateEvent) bool {
		return r.NamePredicate(e.Meta)
	}, UpdateFunc: func(e event.UpdateEvent) bool {
		return r.NamePredicate(e.MetaNew)
	}, DeleteFunc: func(e event.DeleteEvent) bool {
		return r.NamePredicate(e.Meta)
	}, GenericFunc: func(e event.GenericEvent) bool {
		return r.NamePredicate(e.Meta)
	}}
	err = c.Watch(&source.Kind{Type: &autoscalingv1.ClusterAutoscaler{}}, &handler.EnqueueRequestForObject{}, p)
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{IsController: true, OwnerType: &autoscalingv1.ClusterAutoscaler{}})
	if err != nil {
		return err
	}
	return nil
}
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Infof("Reconciling ClusterAutoscaler %s\n", request.Name)
	ca := &autoscalingv1.ClusterAutoscaler{}
	err := r.client.Get(context.TODO(), request.NamespacedName, ca)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("ClusterAutoscaler %s not found, will not reconcile", request.Name)
			return reconcile.Result{}, nil
		}
		klog.Errorf("Error reading ClusterAutoscaler: %v", err)
		return reconcile.Result{}, err
	}
	caRef := r.objectReference(ca)
	_, err = r.GetAutoscaler(ca)
	if err != nil && !errors.IsNotFound(err) {
		errMsg := fmt.Sprintf("Error getting cluster-autoscaler deployment: %v", err)
		r.recorder.Event(caRef, corev1.EventTypeWarning, "FailedGetDeployment", errMsg)
		klog.Error(errMsg)
		return reconcile.Result{}, err
	}
	if errors.IsNotFound(err) {
		if err := r.CreateAutoscaler(ca); err != nil {
			errMsg := fmt.Sprintf("Error creating ClusterAutoscaler deployment: %v", err)
			r.recorder.Event(caRef, corev1.EventTypeWarning, "FailedCreate", errMsg)
			klog.Error(errMsg)
			return reconcile.Result{}, err
		}
		msg := fmt.Sprintf("Created ClusterAutoscaler deployment: %s", r.AutoscalerName(ca))
		r.recorder.Eventf(caRef, corev1.EventTypeNormal, "SuccessfulCreate", msg)
		klog.Info(msg)
		return reconcile.Result{}, nil
	}
	if err := r.UpdateAutoscaler(ca); err != nil {
		errMsg := fmt.Sprintf("Error updating cluster-autoscaler deployment: %v", err)
		r.recorder.Event(caRef, corev1.EventTypeWarning, "FailedUpdate", errMsg)
		klog.Error(errMsg)
		return reconcile.Result{}, err
	}
	msg := fmt.Sprintf("Updated ClusterAutoscaler deployment: %s", r.AutoscalerName(ca))
	r.recorder.Eventf(caRef, corev1.EventTypeNormal, "SuccessfulUpdate", msg)
	klog.Info(msg)
	return reconcile.Result{}, nil
}
func (r *Reconciler) SetConfig(cfg *Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.config = cfg
}
func (r *Reconciler) NamePredicate(meta metav1.Object) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if meta.GetName() != r.config.Name {
		klog.Warningf("Not processing ClusterAutoscaler %s", meta.GetName())
		return false
	}
	return true
}
func (r *Reconciler) CreateAutoscaler(ca *autoscalingv1.ClusterAutoscaler) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Infof("Creating ClusterAutoscaler deployment: %s\n", r.AutoscalerName(ca))
	deployment := r.AutoscalerDeployment(ca)
	if err := controllerutil.SetControllerReference(ca, deployment, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), deployment)
}
func (r *Reconciler) UpdateAutoscaler(ca *autoscalingv1.ClusterAutoscaler) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	existingDeployment, err := r.GetAutoscaler(ca)
	if err != nil {
		return err
	}
	existingSpec := existingDeployment.Spec.Template.Spec
	expectedSpec := r.AutoscalerPodSpec(ca)
	if equality.Semantic.DeepEqual(existingSpec, expectedSpec) && util.ReleaseVersionMatches(ca, r.config.ReleaseVersion) {
		return nil
	}
	existingDeployment.Spec.Template.Spec = *expectedSpec
	r.UpdateAnnotations(existingDeployment)
	r.UpdateAnnotations(&existingDeployment.Spec.Template)
	return r.client.Update(context.TODO(), existingDeployment)
}
func (r *Reconciler) GetAutoscaler(ca *autoscalingv1.ClusterAutoscaler) (*appsv1.Deployment, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	deployment := &appsv1.Deployment{}
	nn := r.AutoscalerName(ca)
	if err := r.client.Get(context.TODO(), nn, deployment); err != nil {
		return nil, err
	}
	return deployment, nil
}
func (r *Reconciler) AutoscalerName(ca *autoscalingv1.ClusterAutoscaler) types.NamespacedName {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return types.NamespacedName{Name: fmt.Sprintf("cluster-autoscaler-%s", ca.Name), Namespace: r.config.Namespace}
}
func (r *Reconciler) UpdateAnnotations(obj metav1.Object) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[util.CriticalPodAnnotation] = ""
	annotations[util.ReleaseVersionAnnotation] = r.config.ReleaseVersion
	obj.SetAnnotations(annotations)
}
func (r *Reconciler) AutoscalerDeployment(ca *autoscalingv1.ClusterAutoscaler) *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	namespacedName := r.AutoscalerName(ca)
	labels := map[string]string{"cluster-autoscaler": ca.Name, "app": "cluster-autoscaler"}
	annotations := map[string]string{util.CriticalPodAnnotation: "", util.ReleaseVersionAnnotation: r.config.ReleaseVersion}
	podSpec := r.AutoscalerPodSpec(ca)
	deployment := &appsv1.Deployment{TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"}, ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace, Annotations: annotations}, Spec: appsv1.DeploymentSpec{Replicas: &r.config.Replicas, Selector: &metav1.LabelSelector{MatchLabels: labels}, Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: labels, Annotations: annotations}, Spec: *podSpec}}}
	return deployment
}
func (r *Reconciler) AutoscalerPodSpec(ca *autoscalingv1.ClusterAutoscaler) *corev1.PodSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := AutoscalerArgs(ca, r.config)
	if r.config.ExtraArgs != "" {
		args = append(args, r.config.ExtraArgs)
	}
	spec := &corev1.PodSpec{ServiceAccountName: caServiceAccount, PriorityClassName: caPriorityClassName, NodeSelector: map[string]string{"node-role.kubernetes.io/master": "", "beta.kubernetes.io/os": "linux"}, Containers: []corev1.Container{{Name: "cluster-autoscaler", Image: r.config.Image, Command: []string{"cluster-autoscaler"}, Args: args}}, Tolerations: []corev1.Toleration{{Key: "CriticalAddonsOnly", Operator: corev1.TolerationOpExists}, {Key: "node-role.kubernetes.io/master", Effect: corev1.TaintEffectNoSchedule, Operator: corev1.TolerationOpExists}}}
	return spec
}
func (r *Reconciler) objectReference(obj runtime.Object) *corev1.ObjectReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref, err := reference.GetReference(r.scheme, obj)
	if err != nil {
		klog.Errorf("Error creating object reference: %v", err)
		return nil
	}
	if ref != nil && ref.Namespace == "" {
		ref.Namespace = r.config.Namespace
	}
	return ref
}
