package util

import (
	"math"
	"os"
	"time"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/config"

	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var startTime metav1.Time = metav1.Now()

func Default(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func EnsureNamespacedName(
	partialNamespacedName *integrationv1alpha2.NamespacedName,
	defaultNamespace string,
) types.NamespacedName {
	return types.NamespacedName{
		Name:      partialNamespacedName.Name,
		Namespace: Default(partialNamespacedName.Namespace, defaultNamespace),
	}
}

func GetOperatorNamespace() string {
	operatorNamespace := os.Getenv("POD_NAMESPACE")
	if operatorNamespace == "" {
		operatorNamespace = "kube-system"
	}
	return operatorNamespace
}

func CalculateExponentialRequireAfter(
	lastUpdate metav1.Time,
	factor int64,
) time.Duration {
	if factor == 0 {
		factor = 2
	}
	now := metav1.Now()
	if startTime.Unix() > lastUpdate.Unix() {
		return time.Duration(time.Second * 2)
	}
	retryInterval := time.Second
	if !lastUpdate.Time.IsZero() {
		retryInterval = now.Sub(lastUpdate.Time).Round(time.Second)
	}
	return time.Duration(math.Min(
		float64(retryInterval.Nanoseconds()*factor),
		float64(config.MaxRequeueDuration),
	))
}

func GetEndpoint(endpoint string) string {
	if endpoint == "" {
		return endpoint
	}
	if endpoint[0:8] != "https://" && endpoint[0:7] != "http://" {
		endpoint = "http://" + endpoint
	}
	if endpoint[len(endpoint)-1] == '/' {
		endpoint = string(endpoint[0 : len(endpoint)-2])
	}
	return endpoint
}
