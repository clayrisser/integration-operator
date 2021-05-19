package util

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PlugService struct {
	client *client.Client
	ctx    *context.Context
	plug   *integrationv1alpha2.Plug
	req    *ctrl.Request
	update *Update
}

func NewPlugUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	plug *integrationv1alpha2.Plug,
) *PlugService {
	return &PlugService{
		client: client,
		ctx:    ctx,
		plug:   plug,
		req:    req,
		update: NewUpdate(99),
	}
}

func (s *PlugService) Get() (*integrationv1alpha2.Plug, error) {
	client := *s.client
	ctx := *s.ctx
	plug := &integrationv1alpha2.Plug{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      s.plug.Name,
		Namespace: s.plug.Namespace,
	}, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (s *PlugService) Update(plug *integrationv1alpha2.Plug) {
	s.update.SchedulePlugUpdate(s.client, s.ctx, nil, plug)
}

func (s *PlugService) UpdateStatus(plug *integrationv1alpha2.Plug) {
	s.update.SchedulePlugUpdateStatus(s.client, s.ctx, nil, plug)
}

func (s *PlugService) SimpleUpdateStatus(
	message string,
	phase integrationv1alpha2.Phase,
	reason string,
	status bool,
	t string,
) error {
	plug, err := s.Get()
	if err != nil {
		return err
	}
	plug.Status.Phase = phase
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: plug.Generation,
		Reason:             reason,
		Status:             "False",
		Type:               t,
	}
	if status {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&plug.Status.Conditions, condition)
	s.UpdateStatus(plug)
	return nil
}
