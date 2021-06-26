package util

import (
	"context"
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type LookupUtil struct {
	ctx     *context.Context
	varUtil *VarUtil
}

func NewLookupUtil(ctx *context.Context) *LookupUtil {
	return &LookupUtil{
		varUtil: NewVarUtil(ctx),
	}
}

func (u *LookupUtil) Lookup(resource gjson.Result, path string) (string, error) {
	lookup, err := u.BuildLookup(resource)
	if err != nil {
		return "", err
	}
	return lookup.Get(path).String(), nil
}

func (u *LookupUtil) BuildLookup(resource gjson.Result) (gjson.Result, error) {
	result := gjson.Parse("{}")

	resourceObj := unstructured.Unstructured{}
	_, _, err := decUnstructured.Decode([]byte(resource.String()), nil, &resourceObj)
	if err != nil {
		return result, err
	}
	resultStr, err := sjson.Set(result.String(), "resource", resourceObj)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	var vars []kustomizeTypes.Var
	err = json.Unmarshal([]byte(resource.Get("spec.vars").String()), &vars)
	if err != nil {
		return result, err
	}
	varsMap, err := u.varUtil.GetVars(vars)
	if err != nil {
		return result, err
	}
	resultStr, err = sjson.Set(result.String(), "var", varsMap)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	return result, nil
}
