package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return map[string]common.OpenAPIDefinition{"github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1.ClusterAutoscaler": schema_pkg_apis_autoscaling_v1_ClusterAutoscaler(ref)}
}
func schema_pkg_apis_autoscaling_v1_ClusterAutoscaler(ref common.ReferenceCallback) common.OpenAPIDefinition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return common.OpenAPIDefinition{Schema: spec.Schema{SchemaProps: spec.SchemaProps{Description: "ClusterAutoscaler is the Schema for the clusterautoscalers API", Type: []string{"object"}, Properties: map[string]spec.Schema{"kind": {SchemaProps: spec.SchemaProps{Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds", Type: []string{"string"}, Format: ""}}, "apiVersion": {SchemaProps: spec.SchemaProps{Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources", Type: []string{"string"}, Format: ""}}, "metadata": {SchemaProps: spec.SchemaProps{Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta")}}, "spec": {SchemaProps: spec.SchemaProps{Ref: ref("github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1.ClusterAutoscalerSpec")}}, "status": {SchemaProps: spec.SchemaProps{Ref: ref("github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1.ClusterAutoscalerStatus")}}}}}, Dependencies: []string{"github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1.ClusterAutoscalerSpec", "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1.ClusterAutoscalerStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"}}
}
