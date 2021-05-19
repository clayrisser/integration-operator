package util

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	ctrl "sigs.k8s.io/controller-runtime"
)

type SocketUtil struct {
	client         *client.Client
	ctx            *context.Context
	log            *logr.Logger
	mutex          *sync.Mutex
	namespacedName types.NamespacedName
	req            *ctrl.Request
}

func NewSocketUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	namespacedName *integrationv1alpha2.NamespacedName,
	mutex *sync.Mutex,
) *SocketUtil {
	operatorNamespace := GetOperatorNamespace()
	if mutex == nil {
		mutex = &sync.Mutex{}
	}
	return &SocketUtil{
		client:         client,
		ctx:            ctx,
		log:            log,
		mutex:          mutex,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
	}
}

func (u *SocketUtil) Get() (*integrationv1alpha2.Socket, error) {
	client := *u.client
	ctx := *u.ctx
	socket := &integrationv1alpha2.Socket{}
	if err := client.Get(ctx, u.namespacedName, socket); err != nil {
		return nil, err
	}
	return socket.DeepCopy(), nil
}

func (u *SocketUtil) Update(socket *integrationv1alpha2.Socket) error {
	client := *u.client
	ctx := *u.ctx
	u.mutex.Lock()
	if err := client.Update(ctx, socket); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *SocketUtil) UpdateStatus(socket *integrationv1alpha2.Socket) error {
	client := *u.client
	ctx := *u.ctx
	u.mutex.Lock()
	if err := client.Status().Update(ctx, socket); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *SocketUtil) CommonUpdateStatus(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
	message string,
	reset bool,
) error {
	socket, err := u.Get()
	if err != nil {
		return err
	}
	if reset {
		socket.Status = integrationv1alpha2.SocketStatus{}
	}
	if appendPlug != nil {
		socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, integrationv1alpha2.CoupledPlug{
			APIVersion: appendPlug.APIVersion,
			Kind:       appendPlug.Kind,
			Name:       appendPlug.Name,
			Namespace:  appendPlug.Namespace,
			UID:        appendPlug.UID,
		})
		joinedCondition, err := u.GetJoinedCondition()
		if err != nil {
			return err
		}
		if (*joinedCondition).Reason == string(SocketReadyStatusCondition) {
			joinedStatusCondition = SocketReadyStatusCondition
		}
	}
	if joinedStatusCondition != "" {
		joinedStatus := false
		coupledPlugsCount := len(socket.Status.CoupledPlugs)
		if message == "" {
			if joinedStatusCondition == SocketCreatedStatusCondition {
				message = "socket created"
			} else if joinedStatusCondition == ErrorStatusCondition {
				message = "unknown error"
			} else if joinedStatusCondition == SocketReadyStatusCondition {
				message = "socket ready with " + fmt.Sprint(coupledPlugsCount) + " plugs coupled"
			}
		}
		if joinedStatusCondition == SocketReadyStatusCondition && coupledPlugsCount > 0 {
			joinedStatus = true
		}
		c := metav1.Condition{
			Message:            message,
			ObservedGeneration: socket.Generation,
			Reason:             string(joinedStatusCondition),
			Status:             "False",
			Type:               "Joined",
		}
		if joinedStatus {
			c.Status = "True"
		}
		meta.SetStatusCondition(&socket.Status.Conditions, c)
	}
	if phase != "" {
		socket.Status.Phase = phase
	}
	if err := u.UpdateStatus(socket); err != nil {
		return err
	}
	return nil
}

func (u *SocketUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
) error {
	return u.CommonUpdateStatus(phase, joinedStatusCondition, appendPlug, "", false)
}

func (u *SocketUtil) UpdateStatusAppendPlug(
	appendPlug *integrationv1alpha2.Plug,
) error {
	return u.CommonUpdateStatus("", "", appendPlug, "", false)
}

func (u *SocketUtil) UpdateStatusPhase(
	phase integrationv1alpha2.Phase,
) error {
	return u.CommonUpdateStatus(phase, "", nil, "", false)
}

func (u *SocketUtil) UpdateStatusJoinedCondition(
	joinedStatusCondition StatusCondition,
	message string,
) error {
	return u.CommonUpdateStatus("", joinedStatusCondition, nil, message, false)
}

func (u *SocketUtil) UpdateStatusJoinedConditionError(
	err error,
) error {
	return u.CommonUpdateStatus(integrationv1alpha2.FailedPhase, ErrorStatusCondition, nil, err.Error(), false)
}

func (u *SocketUtil) UpdateStatusRemovePlug(
	plug *integrationv1alpha2.Plug,
) error {
	socket, err := u.Get()
	if err != nil {
		return err
	}
	coupledPlugs := []integrationv1alpha2.CoupledPlug{}
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		if coupledPlug.UID != plug.UID {
			coupledPlugs = append(coupledPlugs, coupledPlug)
		}
	}
	socket.Status.CoupledPlugs = coupledPlugs
	joinedCondition, err := u.GetJoinedCondition()
	if (*joinedCondition).Reason == string(SocketReadyStatusCondition) {
		condition := metav1.Condition{
			Message:            "socket ready with " + fmt.Sprint(len(coupledPlugs)) + " plugs coupled",
			ObservedGeneration: socket.Generation,
			Reason:             "SocketReady",
			Status:             "False",
			Type:               "Joined",
		}
		if len(coupledPlugs) > 0 {
			condition.Status = "True"
		}
	}
	u.UpdateStatus(socket)
	return nil
}

func (u *SocketUtil) GetJoinedCondition() (*metav1.Condition, error) {
	socket, err := u.Get()
	if err != nil {
		return nil, err
	}
	joinedCondition := meta.FindStatusCondition(socket.Status.Conditions, "Joined")
	return joinedCondition, nil
}

func (u *SocketUtil) Error(err error) error {
	log := *u.log
	log.Error(err, err.Error())
	return u.UpdateStatusJoinedConditionError(err)
}

var GlobalSocketMutex *sync.Mutex = &sync.Mutex{}