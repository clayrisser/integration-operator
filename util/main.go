package util

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

type Condition struct {
	Message string
	Reason  string
	Status  bool
	Type    string
}

type StatusCondition string

const (
	CouplingInProcessStatusCondition StatusCondition = "CouplingInProcess"
	CouplingSucceededStatusCondition StatusCondition = "CouplingSucceeded"
	ErrorStatusCondition             StatusCondition = "Error"
	PlugCreatedStatusCondition       StatusCondition = "PlugCreated"
	SocketCoupledStatusCondition     StatusCondition = "SocketCoupled"
	SocketCreatedStatusCondition     StatusCondition = "SocketCreated"
	SocketEmptyStatusCondition       StatusCondition = "SocketEmpty"
	SocketNotCreatedStatusCondition  StatusCondition = "SocketNotCreated"
	SocketNotReadyStatusCondition    StatusCondition = "SocketNotReady"
)
