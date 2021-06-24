package util

import integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

type DataUtil struct{}

func NewDataUtil() *DataUtil {
	return &DataUtil{}
}

func (u *DataUtil) GetPlugData(plug *integrationv1alpha2.Plug) (map[string]string, error) {
	return plug.Spec.Data, nil
}

func (u *DataUtil) GetSocketData(socket *integrationv1alpha2.Socket) (map[string]string, error) {
	return socket.Spec.Data, nil
}
