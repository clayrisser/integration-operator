package util

import (
	"context"
	"encoding/json"

	"github.com/tidwall/gjson"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type VarUtil struct {
	client       *kubernetes.Clientset
	ctx          *context.Context
	resourceUtil *ResourceUtil
	kubectlUtil  *KubectlUtil
}

func NewVarUtil(ctx *context.Context) *VarUtil {
	return &VarUtil{
		client:       kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		kubectlUtil:  NewKubectlUtil(ctx, &rest.Config{}),
		resourceUtil: NewResourceUtil(ctx),
	}
}

func (u *VarUtil) GetVars(vars []kustomizeTypes.Var) (map[string]string, error) {
	resultMap := make(map[string]string)
	for _, v := range vars {
		varResult, err := u.GetVar(v)
		if err != nil {
			return nil, err
		}
		resultMap[v.Name] = varResult
	}
	return resultMap, nil
}

func (u *VarUtil) GetVar(v kustomizeTypes.Var) (string, error) {
	resource, err := u.resourceUtil.GetResource(v.ObjRef)
	if err != nil {
		return "", err
	}
	bResource, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	return gjson.Parse(string(bResource)).Get(v.FieldRef.FieldPath).String(), nil
}
