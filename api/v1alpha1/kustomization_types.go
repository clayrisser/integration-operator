/**
 * Copyright 2020 Silicon Hills LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
        kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

// KustomizationSpec defines the desired state of Kustomization
type KustomizationSpec struct {
        // keep retrying when fails until timeout in milliseconds expires
        RetryTimeout uint `json:"retryTimeout,omitempty" yaml:"retryTimeout,omitempty"`

        // kustomization config
        Configuration TransformerConfig `json:"configuration,omitempty" yaml:"configuration,omitempty"`

	// CommonAnnotations to add to all objects.
	CommonAnnotations map[string]string `json:"commonAnnotations,omitempty" yaml:"commonAnnotations,omitempty"`

	// CommonLabels to add to all objects and selectors.
	CommonLabels map[string]string `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`

	// ConfigMapGenerator is a list of configmaps to generate from
	// local data (one configMap per list item).
	// The resulting resource is a normal operand, subject to
	// name prefixing, patching, etc.  By default, the name of
	// the map will have a suffix hash generated from its contents.
	ConfigMapGenerator []kustomizeTypes.ConfigMapArgs `json:"configMapGenerator,omitempty" yaml:"configMapGenerator,omitempty"`

	// Crds specifies relative paths to Custom Resource Definition files.
	// This allows custom resources to be recognized as operands, making
	// it possible to add them to the Resources list.
	// CRDs themselves are not modified.
	Crds []string `json:"crds,omitempty" yaml:"crds,omitempty"`

	// GeneratorOptions modify behavior of all ConfigMap and Secret generators.
	GeneratorOptions *kustomizeTypes.GeneratorOptions `json:"generatorOptions,omitempty" yaml:"generatorOptions,omitempty"`

	// Images is a list of (image name, new name, new tag or digest)
	// for changing image names, tags or digests. This can also be achieved with a
	// patch, but this operator is simpler to specify.
	Images []kustomizeTypes.Image `json:"images,omitempty" yaml:"images,omitempty"`

	// NamePrefix will prefix the names of all resources mentioned in the kustomization
	// file including generated configmaps and secrets.
	NamePrefix string `json:"namePrefix,omitempty" yaml:"namePrefix,omitempty"`

	// Namespace to add to all objects.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// NameSuffix will suffix the names of all resources mentioned in the kustomization
	// file including generated configmaps and secrets.
	NameSuffix string `json:"nameSuffix,omitempty" yaml:"nameSuffix,omitempty"`

	// Patches is a list of patches, where each one can be either a
	// Strategic Merge Patch or a JSON patch.
	// Each patch can be applied to multiple target objects.
	Patches []kustomizeTypes.Patch `json:"patches,omitempty" yaml:"patches,omitempty"`

	// JSONPatches is a list of JSONPatch for applying JSON patch.
	// Format documented at https://tools.ietf.org/html/rfc6902
	// and http://jsonpatch.com
	PatchesJson6902 []kustomizeTypes.PatchJson6902 `json:"patchesJson6902,omitempty" yaml:"patchesJson6902,omitempty"`

	// PatchesStrategicMerge specifies the relative path to a file
	// containing a strategic merge patch.  Format documented at
	// https://github.com/kubernetes/community/blob/master/contributors/devel/strategic-merge-patch.md
	// URLs and globs are not supported.
	PatchesStrategicMerge []kustomizeTypes.PatchStrategicMerge `json:"patchesStrategicMerge,omitempty" yaml:"patchesStrategicMerge,omitempty"`

	// Replicas is a list of {resourcename, count} that allows for simpler replica
	// specification. This can also be done with a patch.
	Replicas []kustomizeTypes.Replica `json:"replicas,omitempty" yaml:"replicas,omitempty"`

	// Resources refers to kubernetes resources subject to
	// kustomization.
	Resources []*kustomizeTypes.Selector `json:"resources,omitempty" yaml:"resources,omitempty"`

	// SecretGenerator is a list of secrets to generate from
	// local data (one secret per list item).
	// The resulting resource is a normal operand, subject to
	// name prefixing, patching, etc.  By default, the name of
	// the map will have a suffix hash generated from its contents.
	SecretGenerator []kustomizeTypes.SecretArgs `json:"secretGenerator,omitempty" yaml:"secretGenerator,omitempty"`

	// Vars allow things modified by kustomize to be injected into a
	// kubernetes object specification. A var is a name (e.g. FOO) associated
	// with a field in a specific resource instance.  The field must
	// contain a value of type string/bool/int/float, and defaults to the name field
	// of the instance.  Any appearance of "$(FOO)" in the object
	// spec will be replaced at kustomize build time, after the final
	// value of the specified field has been determined.
	Vars []kustomizeTypes.Var `json:"vars,omitempty" yaml:"vars,omitempty"`
}

type TransformerConfig struct {
	// NameReference     nbrSlice      `json:"nameReference,omitempty" yaml:"nameReference,omitempty"`
	CommonAnnotations kustomizeTypes.FsSlice `json:"commonAnnotations,omitempty" yaml:"commonAnnotations,omitempty"`
	CommonLabels      kustomizeTypes.FsSlice `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
	Images            kustomizeTypes.FsSlice `json:"images,omitempty" yaml:"images,omitempty"`
	NamePrefix        kustomizeTypes.FsSlice `json:"namePrefix,omitempty" yaml:"namePrefix,omitempty"`
	NameSpace         kustomizeTypes.FsSlice `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	NameSuffix        kustomizeTypes.FsSlice `json:"nameSuffix,omitempty" yaml:"nameSuffix,omitempty"`
	Replicas          kustomizeTypes.FsSlice `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	VarReference      kustomizeTypes.FsSlice `json:"varReference,omitempty" yaml:"varReference,omitempty"`
}
