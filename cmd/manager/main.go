package main

import (
	"flag"
	"bytes"
	"net/http"
	"fmt"
	"runtime"
	"github.com/openshift/cluster-autoscaler-operator/pkg/operator"
	"github.com/openshift/cluster-autoscaler-operator/pkg/version"
	"k8s.io/klog"
)

func printVersion() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Infof("Go Version: %s", runtime.Version())
	klog.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	klog.Infof("Version: %s", version.String)
}
func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("alsologtostderr", "true")
	flag.Parse()
	printVersion()
	config := operator.ConfigFromEnvironment()
	operator, err := operator.New(config)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info("Starting cluster-autoscaler-operator")
	if err := operator.Start(); err != nil {
		klog.Fatal(err)
	}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}
