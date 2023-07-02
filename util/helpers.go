/**
 * File: /helpers.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 12:04:55
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package util

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"time"

	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
	"gitlab.com/bitspur/rock8s/integration-operator/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func JsonToHashMap(body []byte) (map[string]string, error) {
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

func WhenInWhenSlice(when integrationv1alpha2.When, whenSlice *[]integrationv1alpha2.When) bool {
	for _, whenItem := range *whenSlice {
		if when == whenItem {
			return true
		}
	}
	return false
}

func Validate(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket) error {
	if socket.Spec.Validation.NamespaceBlacklist != nil {
		for _, namespace := range socket.Spec.Validation.NamespaceBlacklist {
			match, _ := regexp.MatchString(namespace, plug.Namespace)
			if match {
				return fmt.Errorf("namespace %s is blacklisted", plug.Namespace)
			}
		}
	}
	if socket.Spec.Validation.NamespaceWhitelist != nil {
		for _, namespace := range socket.Spec.Validation.NamespaceWhitelist {
			match, _ := regexp.MatchString(namespace, plug.Namespace)
			if match {
				return nil
			}
		}
		return fmt.Errorf("namespace %s is not whitelisted", plug.Namespace)
	}
	return nil
}
