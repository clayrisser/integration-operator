/*
 * File: /util/helpers.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 10:53:30
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Silicon Hills LLC (c) Copyright 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/config"
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
