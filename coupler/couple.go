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
		if err := plugUtil.Error(err); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if coupledCondition == nil {
		if err := plugUtil.UpdateStatusSimple(
			integrationv1alpha2.PendingPhase,
			util.PlugCreatedStatusCondition,
			nil,
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := GlobalCoupler.CreatedPlug(plug); err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{Requeue: true}, nil
	}

	plugInterfaceUtil := util.NewInterfaceUtil(client, ctx, req, log, &plug.Spec.Interface)
	plugInterface, err := plugInterfaceUtil.Get()
	if err != nil {
		if err := plugUtil.Error(err); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	socket, err := socketUtil.Get()
	if err != nil {
		if k8serrors.IsNotFound(err) {
			plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotCreatedStatusCondition, nil)
		} else {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if !socket.Status.Ready {
		plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotReadyStatusCondition, nil)
		return ctrl.Result{Requeue: true}, nil
	}

	socketInterfaceUtil := util.NewInterfaceUtil(client, ctx, req, log, &socket.Spec.Interface)
	socketInterface, err := socketInterfaceUtil.Get()
	if err != nil {
		if err := plugUtil.Error(err); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	if plugInterface.UID != socketInterface.UID {
		if err := plugUtil.Error(errors.New("plug and socket interface do not match")); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	coupledCondition, _ = plugUtil.GetCoupledCondition()
	isCoupled := coupledCondition != nil && coupledCondition.Status != "True"
	if coupledCondition.Reason != string(util.CouplingInProcessStatusCondition) && coupledCondition.Reason != string(util.CouplingSucceededStatusCondition) {
		plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.CouplingInProcessStatusCondition, nil)
		return ctrl.Result{Requeue: true}, nil
	}

	var plugConfig []byte
	if plug.Spec.ConfigEndpoint != "" {
		plugConfig, err = GlobalCoupler.GetConfig(plug.Spec.ConfigEndpoint)
		if err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}
	var socketConfig []byte
	if socket.Spec.ConfigEndpoint != "" {
		socketConfig, err = GlobalCoupler.GetConfig(socket.Spec.ConfigEndpoint)
		if err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	if isCoupled {
		err = GlobalCoupler.CoupledPlug(plug, socket, socketConfig)
		if err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		err = GlobalCoupler.CoupledSocket(plug, socket, plugConfig)
		if err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
			if err := socketUtil.UpdateStatusCoupledConditionError(err); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	} else {
		err = GlobalCoupler.UpdatedPlug(plug, socket, socketConfig)
		if err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		err = GlobalCoupler.UpdatedSocket(plug, socket, socketConfig)
		if err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
			if err := socketUtil.UpdateStatusCoupledConditionError(err); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	coupledCondition, err = plugUtil.GetCoupledCondition()
	if err != nil {
		if err := plugUtil.Error(err); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if !socketUtil.CoupledPlugExits(&socket.Status.CoupledPlugs, plug) {
		if err := socketUtil.UpdateStatusAppendPlug(plug); err != nil {
			return ctrl.Result{}, err
		}
	}
	if plug.Status.Phase != integrationv1alpha2.SucceededPhase || coupledCondition.Reason != string(util.CouplingSucceededStatusCondition) {
		if err := plugUtil.UpdateStatusSimple(integrationv1alpha2.SucceededPhase, util.CouplingSucceededStatusCondition, socket); err != nil {
			if err := plugUtil.Error(err); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}
