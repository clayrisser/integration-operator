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
	previousLastUpdate metav1.Time,
	succeeded bool,
	lastUpdate metav1.Time,
	factor int64,
) time.Duration {
	if factor == 0 {
		factor = 2
	}
	if lastUpdate.IsZero() {
		lastUpdate = metav1.Now()
	}
	retryInterval := time.Second
	if !previousLastUpdate.Time.IsZero() && !succeeded {
		retryInterval = lastUpdate.Sub(previousLastUpdate.Time).Round(time.Second)
	}
	return time.Duration(math.Min(
		float64(retryInterval.Nanoseconds()*factor),
		float64(config.MaxRequeueDuration),
	))
}
