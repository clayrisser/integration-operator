package util

import (
	"context"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	"github.com/silicon-hills/integration-operator/config"

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
	plug.Status.LastUpdate = metav1.Now()
	u.mutex.Lock()
	if err := client.Status().Update(ctx, plug); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *PlugUtil) GetCoupledCondition() (*metav1.Condition, error) {
	plug, err := u.Get()
	if err != nil {
		return nil, err
	}
	coupledCondition := meta.FindStatusCondition(plug.Status.Conditions, "coupled")
	return coupledCondition, nil
}

func (u *PlugUtil) Error(err error) (ctrl.Result, error) {
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

func (u *PlugUtil) UpdateErrorStatus(err error) (ctrl.Result, error) {
	stashedErr := err
	plug, err := u.Get()
	if err != nil {
		return ctrl.Result{}, err
	}
	u.setErrorStatus(plug, stashedErr)
	if err := u.UpdateStatus(plug); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	coupledStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
) (ctrl.Result, error) {
	plug, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if coupledStatusCondition != "" {
		u.setCoupledStatusCondition(plug, coupledStatusCondition, "")
	}
	if socket != nil {
		u.setCoupledSocketStatus(plug, socket)
	}
	if phase != "" {
		u.setPhaseStatus(plug, phase)
	}
	if err := u.UpdateStatus(plug); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) setPhaseStatus(
	plug *integrationv1alpha2.Plug,
	phase integrationv1alpha2.Phase,
) {
	if phase != integrationv1alpha2.FailedPhase {
		plug.Status.Message = ""
	}
	plug.Status.Phase = phase
}

func (u *PlugUtil) setCoupledStatusCondition(
	plug *integrationv1alpha2.Plug,
	coupledStatusCondition StatusCondition,
	message string,
) {
	coupledStatus := false
	if message == "" {
		if coupledStatusCondition == PlugCreatedStatusCondition {
			message = "plug created"
		} else if coupledStatusCondition == SocketNotCreatedStatusCondition {
			message = "waiting for socket to be created"
		} else if coupledStatusCondition == SocketNotReadyStatusCondition {
			message = "waiting for socket to be ready"
		} else if coupledStatusCondition == CouplingInProcessStatusCondition {
			message = "coupling to socket"
		} else if coupledStatusCondition == CouplingSucceededStatusCondition {
			message = "coupling succeeded"
		} else if coupledStatusCondition == ErrorStatusCondition {
			message = "unknown error"
		}
	}
	if coupledStatusCondition == CouplingSucceededStatusCondition {
		coupledStatus = true
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: plug.Generation,
		Reason:             string(coupledStatusCondition),
		Status:             "False",
		Type:               "coupled",
	}
	if coupledStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&plug.Status.Conditions, condition)
}

func (u *PlugUtil) setErrorStatus(plug *integrationv1alpha2.Plug, err error) {
	message := err.Error()
	u.setCoupledStatusCondition(plug, ErrorStatusCondition, message)
	plug.Status.Phase = integrationv1alpha2.FailedPhase
	plug.Status.Message = message
}

func (u *PlugUtil) setCoupledSocketStatus(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
) {
	plug.Status.CoupledSocket = integrationv1alpha2.CoupledSocket{
		APIVersion: socket.APIVersion,
		Kind:       socket.Kind,
		Name:       socket.Name,
		Namespace:  socket.Namespace,
		UID:        socket.UID,
	}
}

var GlobalPlugMutex *sync.Mutex = &sync.Mutex{}
