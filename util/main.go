package util

import (
	"encoding/json"
	"fmt"

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

func jsonToHashMap(body []byte) (map[string]string, error) {
	hashMap := make(map[string]string)
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	for key, value := range obj {
		hashMap[key] = fmt.Sprintf("%v", value)
	}
	return hashMap, nil
}
