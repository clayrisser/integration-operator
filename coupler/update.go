/**
 * File: /update.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:49:19
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package coupler

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
	"gitlab.com/bitspur/rock8s/integration-operator/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *Coupler) Update(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	plugNamespacedName *integrationv1alpha2.NamespacedName,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	socketInterface *integrationv1alpha2.Interface,
) (ctrl.Result, error) {
	configUtil := util.NewConfigUtil(ctx)

	plugUtil := util.NewPlugUtil(client, ctx, req, log, plugNamespacedName, util.GlobalPlugMutex)
	if plug == nil {
		var err error
		plug, err = plugUtil.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	if socket == nil {
		var err error
		socket, err = socketUtil.Get()
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotCreatedStatusCondition, nil, false)
			}
			return plugUtil.Error(err)
		}
	}
	if !socket.Status.Ready {
		return plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotReadyStatusCondition, nil, false)
	}

	plugIsCoupled, err := plugUtil.IsCoupled(plug, nil)
	if err != nil {
		return plugUtil.Error(err)
	}
	if !socketUtil.CoupledPlugExists(socket.Status.CoupledPlugs, plug.UID) || !plugIsCoupled {
		return ctrl.Result{}, nil
	}

	plugInterfaceUtil := util.NewInterfaceUtil(client, ctx, &plug.Spec.Interface)
	plugInterface, err := plugInterfaceUtil.Get()
	if err != nil {
		return plugUtil.Error(err)
	}
	socketInterfaceUtil := util.NewInterfaceUtil(client, ctx, &socket.Spec.Interface)
	if socketInterface == nil {
		socketInterface, err = socketInterfaceUtil.Get()
		if err != nil {
			return plugUtil.Error(err)
		}
	}
	if plugInterface.UID != socketInterface.UID {
		return plugUtil.Error(errors.New("plug and socket interface do not match"))
	}

	plugConfig, err := configUtil.GetPlugConfig(plug, plugInterface, socket)
	if err != nil {
		return plugUtil.Error(err)
	}
	socketConfig, err := configUtil.GetSocketConfig(socket, socketInterface, plug)
	if err != nil {
		return socketUtil.Error(err)
	}

	err = GlobalCoupler.UpdatedPlug(plug, socket, plugConfig, socketConfig)
	if err != nil {
		return plugUtil.Error(err)
	}
	err = GlobalCoupler.UpdatedSocket(plug, socket, plugConfig, socketConfig)
	if err != nil {
		if err := plugUtil.SocketError(err); err != nil {
			return plugUtil.Error(err)
		}
		return socketUtil.Error(err)
	}

	return ctrl.Result{}, nil
}
