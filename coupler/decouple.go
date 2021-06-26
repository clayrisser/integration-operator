package coupler

import (
	"context"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *Coupler) Decouple(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	plugNamespacedName *integrationv1alpha2.NamespacedName,
) (ctrl.Result, error) {
	configUtil := util.NewConfigUtil(ctx)

	plugUtil := util.NewPlugUtil(client, ctx, req, log, plugNamespacedName, util.GlobalPlugMutex)
	plug, err := plugUtil.Get()
	if err != nil {
		return ctrl.Result{}, err
	}

	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	socket, err := socketUtil.Get()
	if err != nil {
		if !errors.IsNotFound(err) {
			return plugUtil.Error(err)
		}
	}

	var plugConfig map[string]string
	if plug.Spec.Apparatus.Endpoint != "" {
		plugConfig, err = configUtil.GetPlugConfig(plug)
		if err != nil {
			return plugUtil.Error(err)
		}
	}
	var socketConfig map[string]string
	if socket != nil && socket.Spec.Apparatus.Endpoint != "" {
		socketConfig, err = configUtil.GetSocketConfig(socket)
		if err != nil {
			return plugUtil.Error(err)
		}
	}

	if err := GlobalCoupler.DecoupledPlug(plug, socket, plugConfig, socketConfig); err != nil {
		return plugUtil.Error(err)
	}
	if err := GlobalCoupler.DecoupledSocket(plug, socket, plugConfig, socketConfig); err != nil {
		return plugUtil.Error(err)
	}

	if socket != nil {
		if _, err := socketUtil.UpdateStatusRemovePlug(plug); err != nil {
			return plugUtil.Error(err)
		}
	}

	return ctrl.Result{}, nil
}
