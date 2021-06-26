package util

import (
	"bytes"
	"context"
	"text/template"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
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

func (u *ResourceUtil) PlugCreated(plug *integrationv1alpha2.Plug) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.CreatedWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.CoupledWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.UpdatedWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.DecoupledWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDeleted(
	plug *integrationv1alpha2.Plug,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.DeletedWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugBroken(
	plug *integrationv1alpha2.Plug,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.BrokenWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCreated(socket *integrationv1alpha2.Socket) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.CreatedWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.CoupledWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.UpdatedWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.DecoupledWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDeleted(
	socket *integrationv1alpha2.Socket,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.DeletedWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketBroken(
	socket *integrationv1alpha2.Socket,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.BrokenWhen)
	if err := u.ProcessResources(resources); err != nil {
		return err
	}
	return nil
}

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

func (u *ResourceUtil) GetResources(
	resources []*integrationv1alpha2.Resource,
	when integrationv1alpha2.When,
) []*integrationv1alpha2.Resource {
	filteredResources := []*integrationv1alpha2.Resource{}
	if resources == nil {
		return resources
	}
	for _, resource := range resources {
		if resource.When == when {
			filteredResources = append(filteredResources, resource)
		}
	}
	return filteredResources
}

func (u *ResourceUtil) ProcessResources(resources []*integrationv1alpha2.Resource) error {
	for _, resource := range resources {
		if resource.Do == integrationv1alpha2.ApplyDo {
			if err := u.kubectlUtil.Apply([]byte(resource.Resource)); err != nil {
				return err
			}
		} else if resource.Do == integrationv1alpha2.DeleteDo {
			if err := u.kubectlUtil.Delete([]byte(resource.Resource)); err != nil {
				return err
			}
		}
	}
	return nil
}
