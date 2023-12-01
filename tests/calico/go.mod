module github.com/projectcalico/calico

go 1.20

replace (
	github.com/Microsoft/hcsshim => github.com/projectcalico/hcsshim v0.8.9-calico

	// Replace the envoy data-plane-api dependency with the projectcalico fork that includes the generated
	// go bindings for the API. Upstream only includes the protobuf definitions, so we need to fork in order to
	// supply the go code.
	github.com/envoyproxy/data-plane-api => github.com/projectcalico/data-plane-api v0.0.0-20210121211707-a620ff3c8f7e

	k8s.io/api => k8s.io/api v0.26.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.8
	k8s.io/apiserver => k8s.io/apiserver v0.26.8
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.26.8
	k8s.io/client-go => k8s.io/client-go v0.26.8
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.26.8
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.8
	k8s.io/code-generator => k8s.io/code-generator v0.26.8
	k8s.io/component-base => k8s.io/component-base v0.26.8
	k8s.io/component-helpers => k8s.io/component-helpers v0.26.8
	k8s.io/controller-manager => k8s.io/controller-manager v0.26.8
	k8s.io/cri-api => k8s.io/cri-api v0.26.8
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.26.8
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.8
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.26.8
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.26.8
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.26.8
	k8s.io/kubectl => k8s.io/kubectl v0.26.8
	k8s.io/kubelet => k8s.io/kubelet v0.26.8
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.26.8
	k8s.io/metrics => k8s.io/metrics v0.26.8
	k8s.io/mount-utils => k8s.io/mount-utils v0.26.8
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.26.8
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.26.8
)

replace k8s.io/kubernetes => k8s.io/kubernetes v1.26.0
