package util

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/config"

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
	socket.Status.LastUpdate = metav1.Now()
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

func (u *SocketUtil) GetCoupledCondition() (*metav1.Condition, error) {
	socket, err := u.Get()
	if err != nil {
		return nil, err
	}
	coupledCondition := meta.FindStatusCondition(socket.Status.Conditions, "Coupled")
	return coupledCondition, nil
}

func (u *SocketUtil) Error(err error) (ctrl.Result, error) {
	stashedErr := err
	log := *u.log
	log.Error(err, err.Error())
	plug, err := u.Get()
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: config.MaxRequeueDuration,
		}, err
	}
	requeueAfter := CalculateExponentialRequireAfter(
		plug.Status.LastUpdate,
		plug.Status.Phase == integrationv1alpha2.SucceededPhase,
		metav1.Now(),
		999,
	)
	if strings.Index(stashedErr.Error(), "the object has been modified; please apply your changes to the latest version and try again") <= -1 {
		if _, err := u.UpdateErrorStatus(stashedErr); err != nil {
			if strings.Index(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") > -1 {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: requeueAfter,
			}, err
		}
	}
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: requeueAfter,
	}, nil
}

func (u *SocketUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	coupledStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
) (ctrl.Result, error) {
	socket, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if phase != "" {
		u.setPhaseStatus(socket, phase)
	}
	if coupledStatusCondition != "" {
		u.setCoupledStatusCondition(socket, coupledStatusCondition, "")
	}
	if appendPlug != nil {
		if err := u.appendCoupledPlugStatus(socket, appendPlug); err != nil {
			return u.Error(err)
		}
	}
	if err := u.UpdateStatus(socket); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) UpdateErrorStatus(err error) (ctrl.Result, error) {
	stashedErr := err
	socket, err := u.Get()
	if err != nil {
		return ctrl.Result{}, err
	}
	u.setErrorStatus(socket, stashedErr)
	if err := u.UpdateStatus(socket); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) UpdateStatusRemovePlug(
	plug *integrationv1alpha2.Plug,
) (ctrl.Result, error) {
	socket, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if err := u.removeCoupledPlugStatus(socket, plug); err != nil {
		return u.Error(err)
	}
	if err := u.UpdateStatus(socket); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) UpdateStatusAppendPlug(
	plug *integrationv1alpha2.Plug,
) (ctrl.Result, error) {
	socket, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if err := u.appendCoupledPlugStatus(socket, plug); err != nil {
		return u.Error(err)
	}
	if err := u.UpdateStatus(socket); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) setPhaseStatus(
	socket *integrationv1alpha2.Socket,
	phase integrationv1alpha2.Phase,

) {
	if phase != integrationv1alpha2.FailedPhase {
		socket.Status.Message = ""
	}
	socket.Status.Phase = phase
}

func (u *SocketUtil) setCoupledStatusCondition(
	socket *integrationv1alpha2.Socket,
	coupledStatusCondition StatusCondition,
	message string,
) {
	u.setReadyStatus(socket, false)
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
		u.setReadyStatus(socket, true)
	}
}

func (u *SocketUtil) setReadyStatus(socket *integrationv1alpha2.Socket, ready bool) {
	socket.Status.Ready = ready
}

func (u *SocketUtil) setErrorStatus(socket *integrationv1alpha2.Socket, err error) {
	message := err.Error()
	u.setCoupledStatusCondition(socket, ErrorStatusCondition, message)
	socket.Status.Phase = integrationv1alpha2.FailedPhase
	socket.Status.Message = message
}

func (u *SocketUtil) appendCoupledPlugStatus(
	socket *integrationv1alpha2.Socket,
	plug *integrationv1alpha2.Plug,
) error {
	if !u.CoupledPlugExits(&socket.Status.CoupledPlugs, plug) {
		socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, integrationv1alpha2.CoupledPlug{
			APIVersion: plug.APIVersion,
			Kind:       plug.Kind,
			Name:       plug.Name,
			Namespace:  plug.Namespace,
			UID:        plug.UID,
		})
	}
	coupledCondition, err := u.GetCoupledCondition()
	if err != nil {
		return err
	}
	if (*coupledCondition).Reason == string(SocketCoupledStatusCondition) {
		u.setCoupledStatusCondition(socket, SocketCoupledStatusCondition, "")
	}
	return nil
}

func (u *SocketUtil) removeCoupledPlugStatus(
	socket *integrationv1alpha2.Socket,
	plug *integrationv1alpha2.Plug,
) error {
	coupledPlugs := []integrationv1alpha2.CoupledPlug{}
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		if coupledPlug.UID != plug.UID {
			coupledPlugs = append(coupledPlugs, coupledPlug)
		}
	}
	socket.Status.CoupledPlugs = coupledPlugs
	coupledCondition, err := u.GetCoupledCondition()
	if err != nil {
		return err
	}
	if (*coupledCondition).Reason == string(SocketCoupledStatusCondition) {
		u.setCoupledStatusCondition(socket, SocketCoupledStatusCondition, "")
	}
	return nil
}

var GlobalSocketMutex *sync.Mutex = &sync.Mutex{}
