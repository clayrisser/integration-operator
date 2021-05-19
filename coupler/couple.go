package coupler

import (
	"context"
	"errors"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
)

func (c *Coupler) Couple(client client.Client, ctx context.Context, req ctrl.Request, result *ctrl.Result, plugNamespacedName *integrationv1alpha2.NamespacedName) error {
	plugUtil := util.NewPlugUtil(&client, &ctx, &req, plugNamespacedName)
	plug, err := plugUtil.Get()
	if err != nil {
		return err
	}

	joinedCondition, _ := plugUtil.GetJoinedCondition(plug)
	if joinedCondition == nil {
		if err := plugUtil.UpdateStatusSimple(
			integrationv1alpha2.PendingPhase,
			util.PlugCreatedStatusCondition,
			nil,
		); err != nil {
			return err
		}
		err = GlobalCoupler.CreatedPlug(plug)
		if err != nil {
			return nil
		}
	}

	plugInterfaceUtil := util.NewInterfaceUtil(&client, &ctx, &req, &plug.Spec.Interface)
	plugInterface, err := plugInterfaceUtil.Get()
	if err != nil {
		plugUtil.UpdateStatusJoinedConditionError(err)
		return nil
	}

	socketUtil := util.NewSocketUtil(&client, &ctx, &req, &plug.Spec.Socket)
	socket, err := socketUtil.Get()
	if err != nil {
		if k8serrors.IsNotFound(err) {
			plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotCreatedStatusCondition, nil)
		} else {
			plugUtil.UpdateStatusJoinedConditionError(err)
		}
		return nil
	}
	if !socket.Status.Ready {
		plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotReadyStatusCondition, nil)
		return nil
	}

	socketInterfaceUtil := util.NewInterfaceUtil(&client, &ctx, &req, &socket.Spec.Interface)
	socketInterface, err := socketInterfaceUtil.Get()
	if err != nil {
		plugUtil.UpdateStatusJoinedConditionError(err)
		return nil
	}

	if plugInterface.UID != socketInterface.UID {
		plugUtil.UpdateStatusJoinedConditionError(errors.New("plug and socket interface do not match"))
		return nil
	}

	joinedCondition, _ = plugUtil.GetJoinedCondition(plug)
	isJoined := joinedCondition != nil && joinedCondition.Status != "True"
	plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.CouplingInProcessStatusCondition, nil)

	var plugConfig []byte
	if plug.Spec.ConfigEndpoint != "" {
		plugConfig, err = GlobalCoupler.GetConfig(plug.Spec.ConfigEndpoint)
		if err != nil {
			plugUtil.UpdateStatusJoinedConditionError(err)
			return nil
		}
	}
	var socketConfig []byte
	if socket.Spec.ConfigEndpoint != "" {
		socketConfig, err = GlobalCoupler.GetConfig(socket.Spec.ConfigEndpoint)
		if err != nil {
			plugUtil.UpdateStatusJoinedConditionError(err)
			return nil
		}
	}

	if isJoined {
		err = GlobalCoupler.JoinedPlug(plug, socket, socketConfig)
		if err != nil {
			plugUtil.UpdateStatusJoinedConditionError(err)
			return nil
		}
		err = GlobalCoupler.JoinedSocket(plug, socket, plugConfig)
		if err != nil {
			plugUtil.UpdateStatusJoinedConditionError(err)
			socketUtil.UpdateStatusJoinedConditionError(err)
			return nil
		}
		plugUtil.UpdateStatusSocket(socket)
		socketUtil.UpdateStatusAppendPlug(plug)
	} else {
		err = GlobalCoupler.ChangedPlug(plug, socket, socketConfig)
		if err != nil {
			plugUtil.UpdateStatusJoinedConditionError(err)
			return nil
		}
		err = GlobalCoupler.ChangedSocket(plug, socket, socketConfig)
		if err != nil {
			plugUtil.UpdateStatusJoinedConditionError(err)
			socketUtil.UpdateStatusJoinedConditionError(err)
			return nil
		}
	}
	plugUtil.UpdateStatusSimple(integrationv1alpha2.SucceededPhase, util.CouplingSucceededStatusCondition, nil)
	return nil
}
