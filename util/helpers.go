/**
 * File: /util/helpers.go
 * Project: new
 * File Created: 17-10-2023 13:49:54
 * Author: Clay Risser
 * -----
 * BitSpur (c) Copyright 2021 - 2023
 *
 * Licensed under the GNU Affero General Public License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.gnu.org/licenses/agpl-3.0.en.html
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving this software without disclosing
 * the source code of your own applications.
 */

package util

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

func Default(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func EnsureServiceAccount(serviceAccountName string) string {
	return Default(serviceAccountName, "default")
}

func EnsureNamespacedName(
	partialNamespacedName *integrationv1beta1.NamespacedName,
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

func WhenInWhenSlice(when integrationv1beta1.When, whenSlice *[]integrationv1beta1.When) bool {
	for _, whenItem := range *whenSlice {
		if when == whenItem {
			return true
		}
	}
	return false
}

func Validate(plug *integrationv1beta1.Plug, socket *integrationv1beta1.Socket) error {
	if socket.Spec.Validation == nil {
		return nil
	}
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
