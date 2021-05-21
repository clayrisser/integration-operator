module github.com/silicon-hills/integration-operator

go 1.15

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.3.0
	github.com/go-resty/resty/v2 v2.6.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/prometheus/common v0.10.0
	github.com/tdewolff/minify v2.3.6+incompatible
	github.com/tdewolff/minify/v2 v2.9.16
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tidwall/gjson v1.7.5
	k8s.io/apimachinery v0.20.2
	k8s.io/apiserver v0.20.1
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/kustomize/api v0.8.9
	sigs.k8s.io/yaml v1.2.0
)
