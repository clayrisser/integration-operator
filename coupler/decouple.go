package coupler

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

	if err := GlobalCoupler.Decoupled(plug, socket); err != nil {
		return plugUtil.Error(err)
	}

	if socket != nil {
		if err := socketUtil.UpdateStatusRemovePlug(plug); err != nil {
			return plugUtil.Error(err)
		}
	}

	controllerutil.RemoveFinalizer(plug, integrationv1alpha2.PlugFinalizer)
	if err := plugUtil.Update(plug); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
