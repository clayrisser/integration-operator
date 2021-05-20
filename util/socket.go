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

func (u *SocketUtil) CoupledPlugExits(coupledPlugs *[]integrationv1alpha2.CoupledPlug, plug *integrationv1alpha2.Plug) bool {
	coupledPlugExits := false
	for _, coupledPlug := range *coupledPlugs {
		if coupledPlug.UID == plug.UID {
			coupledPlugExits = true
		}
	}
	return coupledPlugExits
}

func (u *SocketUtil) CommonUpdateStatus(
	phase integrationv1alpha2.Phase,
	coupledStatusCondition StatusCondition,
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
		if !u.CoupledPlugExits(&socket.Status.CoupledPlugs, appendPlug) {
			socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, integrationv1alpha2.CoupledPlug{
				APIVersion: appendPlug.APIVersion,
				Kind:       appendPlug.Kind,
				Name:       appendPlug.Name,
				Namespace:  appendPlug.Namespace,
				UID:        appendPlug.UID,
			})
		}
		coupledCondition, err := u.GetCoupledCondition()
		if err != nil {
			return err
		}
		if (*coupledCondition).Reason == string(SocketCoupledStatusCondition) {
			coupledStatusCondition = SocketCoupledStatusCondition
		}
	}
	if coupledStatusCondition != "" {
		socket.Status.Ready = false
		coupledStatus := false
		coupledPlugsCount := len(socket.Status.CoupledPlugs)
		if message == "" {
			if coupledStatusCondition == SocketCreatedStatusCondition {
				message = "socket created"
			} else if coupledStatusCondition == ErrorStatusCondition {
				message = "unknown error"
			} else if coupledStatusCondition == SocketCoupledStatusCondition {
				message = "socket ready with " + fmt.Sprint(coupledPlugsCount) + " plugs coupled"
			} else if coupledStatusCondition == SocketEmptyStatusCondition {
				message = "socket ready with 0 plugs coupled"
			}
		}
		if coupledStatusCondition == SocketCoupledStatusCondition {
			if coupledPlugsCount > 0 {
				coupledStatus = true
			} else {
				coupledStatusCondition = SocketEmptyStatusCondition
			}
		}
		c := metav1.Condition{
			Message:            message,
			ObservedGeneration: socket.Generation,
			Reason:             string(coupledStatusCondition),
			Status:             "False",
			Type:               "Coupled",
		}
		if coupledStatus {
			c.Status = "True"
		}
		meta.SetStatusCondition(&socket.Status.Conditions, c)
		if coupledStatusCondition == SocketCoupledStatusCondition || coupledStatusCondition == SocketEmptyStatusCondition {
			socket.Status.Ready = true
		}
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
	coupledStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
) error {
	return u.CommonUpdateStatus(phase, coupledStatusCondition, appendPlug, "", false)
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

func (u *SocketUtil) UpdateStatusCoupledCondition(
	coupledStatusCondition StatusCondition,
	message string,
) error {
	return u.CommonUpdateStatus("", coupledStatusCondition, nil, message, false)
}

func (u *SocketUtil) UpdateStatusCoupledConditionError(
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
	coupledCondition, err := u.GetCoupledCondition()
	if (*coupledCondition).Reason == string(SocketCoupledStatusCondition) {
		condition := metav1.Condition{
			Message:            "socket ready with " + fmt.Sprint(len(coupledPlugs)) + " plugs coupled",
			ObservedGeneration: socket.Generation,
			Reason:             "SocketReady",
			Status:             "False",
			Type:               "Coupled",
		}
		if len(coupledPlugs) > 0 {
			condition.Reason = string(SocketEmptyStatusCondition)
			condition.Status = "True"
		}
	}
	plug.Status.LastUpdate = metav1.Now()
	u.UpdateStatus(socket)
	return nil
}

func (u *SocketUtil) GetCoupledCondition() (*metav1.Condition, error) {
	socket, err := u.Get()
	if err != nil {
		return nil, err
	}
	coupledCondition := meta.FindStatusCondition(socket.Status.Conditions, "Coupled")
	return coupledCondition, nil
}

func (u *SocketUtil) Error(err error) error {
	log := *u.log
	log.Error(err, err.Error())
	return u.UpdateStatusCoupledConditionError(err)
}

var GlobalSocketMutex *sync.Mutex = &sync.Mutex{}
