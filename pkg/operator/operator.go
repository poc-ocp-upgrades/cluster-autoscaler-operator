package operator

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis"
	"github.com/openshift/cluster-autoscaler-operator/pkg/controller/clusterautoscaler"
	"github.com/openshift/cluster-autoscaler-operator/pkg/controller/machineautoscaler"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const OperatorName = "cluster-autoscaler"

type Operator struct {
	config	*Config
	manager	manager.Manager
}

func New(cfg *Config) (*Operator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	operator := &Operator{config: cfg}
	clientConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	managerOptions := manager.Options{Namespace: cfg.WatchNamespace, LeaderElection: cfg.LeaderElection, LeaderElectionNamespace: cfg.LeaderElectionNamespace, LeaderElectionID: cfg.LeaderElectionID}
	operator.manager, err = manager.New(clientConfig, managerOptions)
	if err != nil {
		return nil, err
	}
	if err := apis.AddToScheme(operator.manager.GetScheme()); err != nil {
		return nil, err
	}
	if err := operator.AddControllers(); err != nil {
		return nil, err
	}
	statusConfig := &StatusReporterConfig{ClusterAutoscalerName: cfg.ClusterAutoscalerName, ClusterAutoscalerNamespace: cfg.ClusterAutoscalerNamespace, ReleaseVersion: cfg.ReleaseVersion, RelatedObjects: operator.RelatedObjects()}
	statusReporter, err := NewStatusReporter(operator.manager, statusConfig)
	if err != nil {
		return nil, err
	}
	operator.manager.Add(statusReporter)
	return operator, nil
}
func (o *Operator) RelatedObjects() []configv1.ObjectReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	relatedNamespaces := map[string]string{}
	relatedNamespaces[o.config.WatchNamespace] = ""
	relatedNamespaces[o.config.LeaderElectionNamespace] = ""
	relatedNamespaces[o.config.ClusterAutoscalerNamespace] = ""
	relatedObjects := []configv1.ObjectReference{}
	for namespace := range relatedNamespaces {
		relatedObjects = append(relatedObjects, configv1.ObjectReference{Resource: "namespaces", Name: namespace})
	}
	return relatedObjects
}
func (o *Operator) AddControllers() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ca := clusterautoscaler.NewReconciler(o.manager, &clusterautoscaler.Config{ReleaseVersion: o.config.ReleaseVersion, Name: o.config.ClusterAutoscalerName, Image: o.config.ClusterAutoscalerImage, Replicas: o.config.ClusterAutoscalerReplicas, Namespace: o.config.ClusterAutoscalerNamespace, CloudProvider: o.config.ClusterAutoscalerCloudProvider, Verbosity: o.config.ClusterAutoscalerVerbosity, ExtraArgs: o.config.ClusterAutoscalerExtraArgs})
	if err := ca.AddToManager(o.manager); err != nil {
		return err
	}
	ma := machineautoscaler.NewReconciler(o.manager, &machineautoscaler.Config{Namespace: o.config.ClusterAutoscalerNamespace, SupportedTargetGVKs: machineautoscaler.DefaultSupportedTargetGVKs()})
	if err := ma.AddToManager(o.manager); err != nil {
		return err
	}
	return nil
}
func (o *Operator) Start() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	stopCh := signals.SetupSignalHandler()
	return o.manager.Start(stopCh)
}
