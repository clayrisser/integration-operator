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
	result *ctrl.Result,
	log *logr.Logger,
	plugNamespacedName *integrationv1alpha2.NamespacedName,
) error {
	plugUtil := util.NewPlugUtil(client, ctx, req, log, plugNamespacedName, util.GlobalPlugMutex)
	plug, err := plugUtil.Get()
	if err != nil {
		return err
	}

	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	socket, err := socketUtil.Get()
	if err != nil {
		if !errors.IsNotFound(err) {
			if err := plugUtil.Error(err); err != nil {
				return err
			}
			return nil
		}
	}

	if err := GlobalCoupler.Departed(plug, socket); err != nil {
		if err := plugUtil.Error(err); err != nil {
			return err
		}
		return nil
	}

	if socket != nil {
		if err := socketUtil.UpdateStatusRemovePlug(plug); err != nil {
			if err := plugUtil.Error(err); err != nil {
				return err
			}
			return nil
		}
	}

	controllerutil.RemoveFinalizer(plug, integrationv1alpha2.PlugFinalizer)
	if err := plugUtil.Update(plug); err != nil {
		return err
	}
	return nil
}
