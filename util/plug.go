package util

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PlugUtil struct {
	client         *client.Client
	ctx            *context.Context
	log            *logr.Logger
	mutex          *sync.Mutex
	namespacedName types.NamespacedName
	req            *ctrl.Request
}

func NewPlugUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	namespacedName *integrationv1alpha2.NamespacedName,
	mutex *sync.Mutex,
) *PlugUtil {
	operatorNamespace := GetOperatorNamespace()
	if mutex == nil {
		mutex = &sync.Mutex{}
	}
	return &PlugUtil{
		client:         client,
		ctx:            ctx,
		log:            log,
		mutex:          mutex,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
	}
}

func (u *PlugUtil) Get() (*integrationv1alpha2.Plug, error) {
	client := *u.client
	ctx := *u.ctx
	plug := &integrationv1alpha2.Plug{}
	if err := client.Get(ctx, u.namespacedName, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (u *PlugUtil) Update(plug *integrationv1alpha2.Plug) error {
	client := *u.client
	ctx := *u.ctx
	u.mutex.Lock()
	if err := client.Update(ctx, plug); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *PlugUtil) UpdateStatus(plug *integrationv1alpha2.Plug) error {
	client := *u.client
	ctx := *u.ctx
	u.mutex.Lock()
	if err := client.Status().Update(ctx, plug); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *PlugUtil) CommonUpdateStatus(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
	message string,
	reset bool,
) error {
	plug, err := u.Get()
	if err != nil {
		return err
	}
	if reset {
		plug.Status = integrationv1alpha2.PlugStatus{}
	}
	if joinedStatusCondition != "" {
		joinedStatus := false
		if message == "" {
			if joinedStatusCondition == PlugCreatedStatusCondition {
				message = "plug created"
			} else if joinedStatusCondition == SocketNotCreatedStatusCondition {
				message = "waiting for socket to be created"
			} else if joinedStatusCondition == SocketNotReadyStatusCondition {
				message = "waiting for socket to be ready"
			} else if joinedStatusCondition == CouplingInProcessStatusCondition {
				message = "coupling to socket"
			} else if joinedStatusCondition == CouplingSucceededStatusCondition {
				message = "coupling succeeded"
			} else if joinedStatusCondition == ErrorStatusCondition {
				message = "unknown error"
			}
		}
		if joinedStatusCondition == CouplingSucceededStatusCondition {
			joinedStatus = true
		}
		c := metav1.Condition{
			Message:            message,
			ObservedGeneration: plug.Generation,
			Reason:             string(joinedStatusCondition),
			Status:             "False",
			Type:               "Joined",
		}
		if joinedStatus {
			c.Status = "True"
		}
		meta.SetStatusCondition(&plug.Status.Conditions, c)
	}
	if socket != nil {
		coupledSocket := integrationv1alpha2.CoupledSocket{
			APIVersion: socket.APIVersion,
			Kind:       socket.Kind,
			Name:       socket.Name,
			Namespace:  socket.Namespace,
			UID:        socket.UID,
		}
		plug.Status.CoupledSocket = coupledSocket
	}
	if phase != "" {
		plug.Status.Phase = phase
	}
	if err := u.UpdateStatus(plug); err != nil {
		return err
	}
	return nil
}

func (u *PlugUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
) error {
	return u.CommonUpdateStatus(phase, joinedStatusCondition, socket, "", false)
}

func (u *PlugUtil) UpdateStatusSocket(
	socket *integrationv1alpha2.Socket,
) error {
	return u.CommonUpdateStatus("", "", socket, "", false)
}

func (u *PlugUtil) UpdateStatusPhase(
	phase integrationv1alpha2.Phase,
) error {
	return u.CommonUpdateStatus(phase, "", nil, "", false)
}

func (u *PlugUtil) UpdateStatusJoinedCondition(
	joinedStatusCondition StatusCondition,
	message string,
) error {
	return u.CommonUpdateStatus("", joinedStatusCondition, nil, message, false)
}

func (u *PlugUtil) UpdateStatusJoinedConditionError(
	err error,
) error {
	return u.CommonUpdateStatus(integrationv1alpha2.FailedPhase, ErrorStatusCondition, nil, err.Error(), false)
}

func (u *PlugUtil) GetJoinedCondition() (*metav1.Condition, error) {
	plug, err := u.Get()
	if err != nil {
		return nil, err
	}
	joinedCondition := meta.FindStatusCondition(plug.Status.Conditions, "Joined")
	return joinedCondition, nil
}

func (u *PlugUtil) Error(err error) error {
	log := *u.log
	log.Error(err, err.Error())
	return u.UpdateStatusJoinedConditionError(err)
}

var GlobalPlugMutex *sync.Mutex = &sync.Mutex{}
