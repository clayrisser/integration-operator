package util

import integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

type ConfigUtil struct {
	apparartusUtil *ApparatusUtil
}

func NewConfigUtil() *ConfigUtil {
	return &ConfigUtil{
		apparartusUtil: NewApparatusUtil(),
	}
}

func (u *ConfigUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) ([]byte, error) {
	return u.apparartusUtil.GetPlugConfig(plug)
}

func (u *ConfigUtil) GetSocketConfig(
	socket *integrationv1alpha2.Socket,
) ([]byte, error) {
	return u.apparartusUtil.GetSocketConfig(socket)
}
