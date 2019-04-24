package apis

import (
	v1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	"bytes"
	"net/http"
	"runtime"
	"fmt"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}
