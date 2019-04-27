package apis

import (
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1beta1"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	AddToSchemes = append(AddToSchemes, v1beta1.SchemeBuilder.AddToScheme)
}
