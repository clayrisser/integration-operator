package coupler

import (
	"context"
	"errors"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
)

func (c *Coupler) Couple(
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

	coupledCondition, err := plugUtil.GetCoupledCondition()
	if err != nil {
		return plugUtil.Error(err)
	}
	if coupledCondition == nil {
		if err := GlobalCoupler.CreatedPlug(plug); err != nil {
			return plugUtil.Error(err)
		}
		return plugUtil.UpdateStatusSimple(
			integrationv1alpha2.PendingPhase,
			util.PlugCreatedStatusCondition,
			nil,
		)
	}

	plugInterfaceUtil := util.NewInterfaceUtil(client, ctx, req, log, &plug.Spec.Interface)
	plugInterface, err := plugInterfaceUtil.Get()
	if err != nil {
		return plugUtil.Error(err)
	}

	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	socket, err := socketUtil.Get()
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotCreatedStatusCondition, nil)
		}
		return plugUtil.Error(err)
	}
	if !socket.Status.Ready {
		return plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotReadyStatusCondition, nil)
	}

	socketInterfaceUtil := util.NewInterfaceUtil(client, ctx, req, log, &socket.Spec.Interface)
	socketInterface, err := socketInterfaceUtil.Get()
	if err != nil {
		return plugUtil.Error(err)
	}

	if plugInterface.UID != socketInterface.UID {
		return plugUtil.Error(errors.New("plug and socket interface do not match"))
	}

	coupledCondition, _ = plugUtil.GetCoupledCondition()
	isCoupled := coupledCondition != nil && coupledCondition.Status != "True"

	var plugConfig []byte
	if plug.Spec.IntegrationEndpoint != "" {
		plugConfig, err = GlobalCoupler.GetConfig(plug.Spec.IntegrationEndpoint, plug, nil)
		if err != nil {
			return plugUtil.Error(err)
		}
	}
	var socketConfig []byte
	if socket.Spec.IntegrationEndpoint != "" {
		socketConfig, err = GlobalCoupler.GetConfig(socket.Spec.IntegrationEndpoint, nil, socket)
		if err != nil {
			return plugUtil.Error(err)
		}
	}

	if isCoupled {
		err = GlobalCoupler.CoupledPlug(plug, socket, plugConfig, socketConfig)
		if err != nil {
			return plugUtil.Error(err)
		}
		err = GlobalCoupler.CoupledSocket(plug, socket, plugConfig, socketConfig)
		if err != nil {
			result, err := plugUtil.Error(err)
			if _, err := socketUtil.UpdateErrorStatus(err); err != nil {
				return plugUtil.Error(err)
			}
			return result, err
		}
	} else {
		err = GlobalCoupler.UpdatedPlug(plug, socket, plugConfig, socketConfig)
		if err != nil {
			return plugUtil.Error(err)
		}
		err = GlobalCoupler.UpdatedSocket(plug, socket, plugConfig, socketConfig)
		if err != nil {
			result, err := plugUtil.Error(err)
			if _, err := socketUtil.UpdateErrorStatus(err); err != nil {
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	coupledCondition, err = plugUtil.GetCoupledCondition()
	if err != nil {
		return plugUtil.Error(err)
	}
	if !socketUtil.CoupledPlugExits(&socket.Status.CoupledPlugs, plug) {
		if _, err := socketUtil.UpdateStatusAppendPlug(plug); err != nil {
			return plugUtil.Error(err)
		}
		return ctrl.Result{}, nil
	}
	if plug.Status.Phase != integrationv1alpha2.SucceededPhase || coupledCondition.Reason != string(util.CouplingSucceededStatusCondition) {
		return plugUtil.UpdateStatusSimple(integrationv1alpha2.SucceededPhase, util.CouplingSucceededStatusCondition, socket)
	}
	return ctrl.Result{}, nil
}
