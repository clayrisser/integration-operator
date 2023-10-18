/**
 * File: /util/var.go
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
	"context"
	"encoding/json"

	"github.com/tidwall/gjson"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type VarUtil struct {
	client       *kubernetes.Clientset
	resourceUtil *ResourceUtil
}

func NewVarUtil(ctx *context.Context) *VarUtil {
	return &VarUtil{
		client:       kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		resourceUtil: NewResourceUtil(ctx),
	}
}

func (u *VarUtil) GetVars(namespace string, vars []*kustomizeTypes.Var, kubectlUtil *KubectlUtil) (map[string]string, error) {
	resultMap := make(map[string]string)
	for _, v := range vars {
		varResult, err := u.GetVar(namespace, v, kubectlUtil)
		if err != nil {
			return nil, err
		}
		resultMap[v.Name] = varResult
	}
	return resultMap, nil
}

func (u *VarUtil) GetVar(namespace string, v *kustomizeTypes.Var, kubectlUtil *KubectlUtil) (string, error) {
	resource, err := u.resourceUtil.GetResource(namespace, v.ObjRef, kubectlUtil)
	if err != nil {
		return "", err
	}
	bResource, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	return gjson.Parse(string(bResource)).Get(v.FieldRef.FieldPath).String(), nil
}
