/**
 * File: /decouple.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 14-08-2022 14:34:43
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Risser Labs LLC (c) Copyright 2021
 */

package coupler

import (
	"context"

	"github.com/go-logr/logr"
	integrationv1alpha2 "gitlab.com/risserlabs/internal/integration-operator/api/v1alpha2"
	"gitlab.com/risserlabs/internal/integration-operator/util"
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
	plug *integrationv1alpha2.Plug,
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
	socket, err := socketUtil.Get()
	if err != nil {
		if !errors.IsNotFound(err) {
			return plugUtil.Error(err)
		}
	}

	plugConfig, err := configUtil.GetPlugConfig(plug, nil)
	if err != nil {
		return plugUtil.Error(err)
	}
	socketConfig := map[string]string{}
	if socket != nil {
		socketConfig, err = configUtil.GetSocketConfig(socket, nil)
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

	socket, err = socketUtil.Get()
	if err != nil {
		if !errors.IsNotFound(err) {
			return plugUtil.Error(err)
		}
	}
	if socket != nil {
		if _, err := socketUtil.UpdateStatusRemovePlug(plug.UID, false); err != nil {
			if !errors.IsNotFound(err) {
				return plugUtil.Error(err)
			}
		}
	}
	return ctrl.Result{}, nil
}
