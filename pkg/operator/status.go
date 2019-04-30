package operator

import (
	"context"
	"fmt"
	"time"
	configv1 "github.com/openshift/api/config/v1"
	osconfig "github.com/openshift/client-go/config/clientset/versioned"
	autoscalingv1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	cvorm "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	ReasonEmpty		= ""
	ReasonMissingDependency	= "MissingDependency"
	ReasonSyncing		= "SyncingResources"
	ReasonCheckAutoscaler	= "UnableToCheckAutoscalers"
)

type StatusReporter struct {
	client		client.Client
	configClient	osconfig.Interface
	config		*StatusReporterConfig
}
type StatusReporterConfig struct {
	ClusterAutoscalerName		string
	ClusterAutoscalerNamespace	string
	ReleaseVersion			string
	RelatedObjects			[]configv1.ObjectReference
}

func NewStatusReporter(mgr manager.Manager, cfg *StatusReporterConfig) (*StatusReporter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	reporter := &StatusReporter{client: mgr.GetClient(), config: cfg}
	reporter.configClient, err = osconfig.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}
	return reporter, nil
}
func (r *StatusReporter) SetReleaseVersion(version string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.config.ReleaseVersion = version
}
func (r *StatusReporter) SetRelatedObjects(objs []configv1.ObjectReference) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.config.RelatedObjects = objs
}
func (r *StatusReporter) AddRelatedObjects(objs []configv1.ObjectReference) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, obj := range objs {
		r.config.RelatedObjects = append(r.config.RelatedObjects, obj)
	}
}
func (r *StatusReporter) GetClusterOperator() (*configv1.ClusterOperator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.configClient.ConfigV1().ClusterOperators().Get(OperatorName, metav1.GetOptions{})
}
func (r *StatusReporter) GetOrCreateClusterOperator() (*configv1.ClusterOperator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	clusterOperator := &configv1.ClusterOperator{ObjectMeta: metav1.ObjectMeta{Name: OperatorName}}
	existing, err := r.GetClusterOperator()
	if errors.IsNotFound(err) {
		return r.configClient.ConfigV1().ClusterOperators().Create(clusterOperator)
	}
	return existing, err
}
func (r *StatusReporter) ApplyStatus(status configv1.ClusterOperatorStatus) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var modified bool
	co, err := r.GetOrCreateClusterOperator()
	if err != nil {
		return err
	}
	status.RelatedObjects = r.config.RelatedObjects
	if status.Versions == nil {
		status.Versions = co.Status.Versions
	}
	for i := range status.Conditions {
		condType := status.Conditions[i].Type
		timestamp := metav1.NewTime(time.Now())
		c := cvorm.FindOperatorStatusCondition(co.Status.Conditions, condType)
		if c != nil && c.Status != status.Conditions[i].Status {
			status.Conditions[i].LastTransitionTime = timestamp
		}
		if c != nil && c.Status == status.Conditions[i].Status {
			status.Conditions[i].LastTransitionTime = c.LastTransitionTime
		}
		if status.Conditions[i].LastTransitionTime.IsZero() {
			status.Conditions[i].LastTransitionTime = timestamp
		}
	}
	if !equality.Semantic.DeepEqual(status.Versions, co.Status.Versions) {
		util.ResetProgressingTime(&status.Conditions)
	}
	requiredCO := &configv1.ClusterOperator{}
	co.DeepCopyInto(requiredCO)
	requiredCO.Status = status
	cvorm.EnsureClusterOperatorStatus(&modified, co, *requiredCO)
	if modified {
		_, err := r.configClient.ConfigV1().ClusterOperators().UpdateStatus(co)
		return err
	}
	return nil
}
func (r *StatusReporter) available(reason, message string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	status := configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, Reason: reason, Message: message}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorDegraded, Status: configv1.ConditionFalse}}, Versions: []configv1.OperandVersion{{Name: "operator", Version: r.config.ReleaseVersion}}}
	klog.Infof("Operator status available: %s", message)
	return r.ApplyStatus(status)
}
func (r *StatusReporter) degraded(reason, message string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	status := configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorDegraded, Status: configv1.ConditionTrue, Reason: reason, Message: message}}}
	klog.Warningf("Operator status degraded: %s", message)
	return r.ApplyStatus(status)
}
func (r *StatusReporter) progressing(reason, message string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	status := configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionTrue, Reason: reason, Message: message}, {Type: configv1.OperatorDegraded, Status: configv1.ConditionFalse}}}
	klog.Infof("Operator status progressing: %s", message)
	return r.ApplyStatus(status)
}
func (r *StatusReporter) Start(stopCh <-chan struct{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	interval := 15 * time.Second
	pollFunc := func() (bool, error) {
		return r.ReportStatus()
	}
	err := wait.PollImmediateUntil(interval, pollFunc, stopCh)
	<-stopCh
	return err
}
func (r *StatusReporter) ReportStatus() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ok, err := r.CheckMachineAPI()
	if err != nil {
		msg := fmt.Sprintf("error checking machine-api status: %v", err)
		r.degraded(ReasonMissingDependency, msg)
		return false, nil
	}
	if !ok {
		r.degraded(ReasonMissingDependency, "machine-api not ready")
		return false, nil
	}
	ok, err = r.CheckClusterAutoscaler()
	if err != nil {
		msg := fmt.Sprintf("error checking autoscaler status: %v", err)
		r.degraded(ReasonCheckAutoscaler, msg)
		return false, nil
	}
	if !ok {
		msg := fmt.Sprintf("updating to %s", r.config.ReleaseVersion)
		r.progressing(ReasonSyncing, msg)
		return false, nil
	}
	msg := fmt.Sprintf("at version %s", r.config.ReleaseVersion)
	r.available(ReasonEmpty, msg)
	return true, nil
}
func (r *StatusReporter) CheckMachineAPI() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mao, err := r.configClient.ConfigV1().ClusterOperators().Get("machine-api", metav1.GetOptions{})
	if err != nil {
		klog.Errorf("failed to get dependency machine-api status: %v", err)
		return false, err
	}
	conds := mao.Status.Conditions
	if cvorm.IsOperatorStatusConditionTrue(conds, configv1.OperatorAvailable) && (cvorm.IsOperatorStatusConditionFalse(conds, configv1.OperatorFailing) || cvorm.IsOperatorStatusConditionFalse(conds, configv1.OperatorDegraded)) {
		return true, nil
	}
	klog.Infof("machine-api-operator not ready yet")
	return false, nil
}
func (r *StatusReporter) CheckClusterAutoscaler() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ca := &autoscalingv1.ClusterAutoscaler{}
	caName := client.ObjectKey{Name: r.config.ClusterAutoscalerName}
	if err := r.client.Get(context.TODO(), caName, ca); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("No ClusterAutoscaler. Reporting available.")
			return true, nil
		}
		klog.Errorf("Error getting ClusterAutoscaler: %v", err)
		return false, err
	}
	deployment := &appsv1.Deployment{}
	deploymentName := client.ObjectKey{Name: fmt.Sprintf("%s-%s", OperatorName, r.config.ClusterAutoscalerName), Namespace: r.config.ClusterAutoscalerNamespace}
	if err := r.client.Get(context.TODO(), deploymentName, deployment); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("No ClusterAutoscaler deployment. Reporting unavailable.")
			return false, nil
		}
		klog.Errorf("Error getting ClusterAutoscaler deployment: %v", err)
		return false, err
	}
	if !util.ReleaseVersionMatches(deployment, r.config.ReleaseVersion) {
		klog.Info("ClusterAutoscaler deployment version not current.")
		return false, nil
	}
	if !util.DeploymentUpdated(deployment) {
		klog.Info("ClusterAutoscaler deployment updating.")
		return false, nil
	}
	klog.Info("ClusterAutoscaler deployment is available and updated.")
	return true, nil
}
