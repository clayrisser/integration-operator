package util

import (
	"context"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type LookupUtil struct {
	ctx      *context.Context
	dataUtil *DataUtil
	varUtil  *VarUtil
}

func NewLookupUtil(ctx *context.Context) *LookupUtil {
	return &LookupUtil{
		dataUtil: NewDataUtil(ctx),
		varUtil:  NewVarUtil(ctx),
	}
}

func (u *LookupUtil) PlugLookup(plug *integrationv1alpha2.Plug, path string) (string, error) {
	plugLookup, err := u.BuildPlugLookup(plug)
	if err != nil {
		return "", err
	}
	return plugLookup.Get(path).String(), nil
}

func (u *LookupUtil) SocketLookup(socket *integrationv1alpha2.Socket, path string) (string, error) {
	socketLookup, err := u.BuildSocketLookup(socket)
	if err != nil {
		return "", err
	}
	return socketLookup.Get(path).String(), nil
}

func (u *LookupUtil) BuildPlugLookup(plug *integrationv1alpha2.Plug) (gjson.Result, error) {
	result := gjson.Parse("{}")

	resultStr, err := sjson.Set(result.String(), "resource", plug)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	dataMap, err := u.dataUtil.GetPlugData(plug)
	if err != nil {
		return result, err
	}
	resultStr, err = sjson.Set(result.String(), "data", dataMap)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	if plug.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(plug.Spec.Vars)
		if err != nil {
			return result, err
		}
		resultStr, err = sjson.Set(result.String(), "var", varsMap)
		if err != nil {
			return result, err
		}
		result = gjson.Parse(resultStr)
	}

	return result, nil
}

func (u *LookupUtil) BuildSocketLookup(socket *integrationv1alpha2.Socket) (gjson.Result, error) {
	result := gjson.Parse("{}")

	resultStr, err := sjson.Set(result.String(), "resource", socket)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	dataMap, err := u.dataUtil.GetSocketData(socket)
	if err != nil {
		return result, err
	}
	resultStr, err = sjson.Set(result.String(), "data", dataMap)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	if socket.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(socket.Spec.Vars)
		if err != nil {
			return result, err
		}
		resultStr, err = sjson.Set(result.String(), "var", varsMap)
		if err != nil {
			return result, err
		}
		result = gjson.Parse(resultStr)
	}

	return result, nil
}
