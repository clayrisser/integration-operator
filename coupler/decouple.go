package coupler

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
)

func (c *Coupler) Decouple(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	plugNamespacedName *integrationv1alpha2.NamespacedName,
) (ctrl.Result, error) {
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

	var plugConfig []byte
	if plug.Spec.IntegrationEndpoint != "" {
		plugConfig, err = GlobalCoupler.GetConfig(plug.Spec.IntegrationEndpoint, plug, nil)
		if err != nil {
			return plugUtil.Error(err)
		}
	}
	var socketConfig []byte
	if socket != nil && socket.Spec.IntegrationEndpoint != "" {
		socketConfig, err = GlobalCoupler.GetConfig(socket.Spec.IntegrationEndpoint, nil, socket)
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
