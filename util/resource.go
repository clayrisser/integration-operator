package util

import (
	"bytes"
	"context"
	"text/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type ResourceUtil struct {
	client      *kubernetes.Clientset
	ctx         *context.Context
	kubectlUtil *KubectlUtil
}

func NewResourceUtil(ctx *context.Context) *ResourceUtil {
	return &ResourceUtil{
		client:      kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		kubectlUtil: NewKubectlUtil(ctx, &rest.Config{}),
	}
}

func (u *ResourceUtil) CreateResource(resource string) {}

func (u *ResourceUtil) UpdateResource(resource string) {}

func (u *ResourceUtil) GetResource(objRef kustomizeTypes.Target) (*unstructured.Unstructured, error) {
	const tpl = `
apiVersion: {{ .APIVersion }}
kind: {{ .Kind }}
meta:
  name: {{ .Name }}
  namespace: {{ .Namespace }}`
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return nil, err
	}
	if objRef.Group != "" && objRef.Version != "" {
		objRef.APIVersion = objRef.Group + objRef.Version
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, objRef)
	if err != nil {
		return nil, err
	}
	body := []byte(buff.String())
	return u.kubectlUtil.Get(body)
}
