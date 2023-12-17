//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigInterface) DeepCopyInto(out *ConfigInterface) {
	*out = *in
	if in.Plug != nil {
		in, out := &in.Plug, &out.Plug
		*out = make(map[string]*SchemaProperty, len(*in))
		for key, val := range *in {
			var outVal *SchemaProperty
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SchemaProperty)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.Socket != nil {
		in, out := &in.Socket, &out.Socket
		*out = make(map[string]*SchemaProperty, len(*in))
		for key, val := range *in {
			var outVal *SchemaProperty
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SchemaProperty)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigInterface.
func (in *ConfigInterface) DeepCopy() *ConfigInterface {
	if in == nil {
		return nil
	}
	out := new(ConfigInterface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CoupledPlug) DeepCopyInto(out *CoupledPlug) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CoupledPlug.
func (in *CoupledPlug) DeepCopy() *CoupledPlug {
	if in == nil {
		return nil
	}
	out := new(CoupledPlug)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CoupledResult) DeepCopyInto(out *CoupledResult) {
	*out = *in
	if in.Plug != nil {
		in, out := &in.Plug, &out.Plug
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Socket != nil {
		in, out := &in.Socket, &out.Socket
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CoupledResult.
func (in *CoupledResult) DeepCopy() *CoupledResult {
	if in == nil {
		return nil
	}
	out := new(CoupledResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CoupledResultStatus) DeepCopyInto(out *CoupledResultStatus) {
	*out = *in
	in.CoupledResult.DeepCopyInto(&out.CoupledResult)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CoupledResultStatus.
func (in *CoupledResultStatus) DeepCopy() *CoupledResultStatus {
	if in == nil {
		return nil
	}
	out := new(CoupledResultStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CoupledSocket) DeepCopyInto(out *CoupledSocket) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CoupledSocket.
func (in *CoupledSocket) DeepCopy() *CoupledSocket {
	if in == nil {
		return nil
	}
	out := new(CoupledSocket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeferResource) DeepCopyInto(out *DeferResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeferResource.
func (in *DeferResource) DeepCopy() *DeferResource {
	if in == nil {
		return nil
	}
	out := new(DeferResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeferResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeferResourceList) DeepCopyInto(out *DeferResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DeferResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeferResourceList.
func (in *DeferResourceList) DeepCopy() *DeferResourceList {
	if in == nil {
		return nil
	}
	out := new(DeferResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeferResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeferResourceSpec) DeepCopyInto(out *DeferResourceSpec) {
	*out = *in
	if in.WaitFor != nil {
		in, out := &in.WaitFor, &out.WaitFor
		*out = new([]*WaitForTarget)
		if **in != nil {
			in, out := *in, *out
			*out = make([]*WaitForTarget, len(*in))
			for i := range *in {
				if (*in)[i] != nil {
					in, out := &(*in)[i], &(*out)[i]
					*out = new(WaitForTarget)
					**out = **in
				}
			}
		}
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = new(v1.JSON)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeferResourceSpec.
func (in *DeferResourceSpec) DeepCopy() *DeferResourceSpec {
	if in == nil {
		return nil
	}
	out := new(DeferResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeferResourceStatus) DeepCopyInto(out *DeferResourceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.OwnerReference.DeepCopyInto(&out.OwnerReference)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeferResourceStatus.
func (in *DeferResourceStatus) DeepCopy() *DeferResourceStatus {
	if in == nil {
		return nil
	}
	out := new(DeferResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Interface) DeepCopyInto(out *Interface) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ConfigInterface)
		(*in).DeepCopyInto(*out)
	}
	if in.Result != nil {
		in, out := &in.Result, &out.Result
		*out = new(ResultInterface)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Interface.
func (in *Interface) DeepCopy() *Interface {
	if in == nil {
		return nil
	}
	out := new(Interface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Plug) DeepCopyInto(out *Plug) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Plug.
func (in *Plug) DeepCopy() *Plug {
	if in == nil {
		return nil
	}
	out := new(Plug)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Plug) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlugList) DeepCopyInto(out *PlugList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Plug, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlugList.
func (in *PlugList) DeepCopy() *PlugList {
	if in == nil {
		return nil
	}
	out := new(PlugList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlugList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlugSpec) DeepCopyInto(out *PlugSpec) {
	*out = *in
	out.Socket = in.Socket
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make([]*Var, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Var)
				**out = **in
			}
		}
	}
	if in.ResultVars != nil {
		in, out := &in.ResultVars, &out.ResultVars
		*out = make([]*Var, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Var)
				**out = **in
			}
		}
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ConfigTemplate != nil {
		in, out := &in.ConfigTemplate, &out.ConfigTemplate
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Result != nil {
		in, out := &in.Result, &out.Result
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ResultTemplate != nil {
		in, out := &in.ResultTemplate, &out.ResultTemplate
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Apparatus != nil {
		in, out := &in.Apparatus, &out.Apparatus
		*out = new(SpecApparatus)
		(*in).DeepCopyInto(*out)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]*Resource, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Resource)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.ResultResources != nil {
		in, out := &in.ResultResources, &out.ResultResources
		*out = make([]*ResourceAction, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ResourceAction)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlugSpec.
func (in *PlugSpec) DeepCopy() *PlugSpec {
	if in == nil {
		return nil
	}
	out := new(PlugSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlugStatus) DeepCopyInto(out *PlugStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CoupledSocket != nil {
		in, out := &in.CoupledSocket, &out.CoupledSocket
		*out = new(CoupledSocket)
		**out = **in
	}
	if in.CoupledResult != nil {
		in, out := &in.CoupledResult, &out.CoupledResult
		*out = new(CoupledResultStatus)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlugStatus.
func (in *PlugStatus) DeepCopy() *PlugStatus {
	if in == nil {
		return nil
	}
	out := new(PlugStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
	in.ResourceAction.DeepCopyInto(&out.ResourceAction)
	if in.When != nil {
		in, out := &in.When, &out.When
		*out = new([]When)
		if **in != nil {
			in, out := *in, *out
			*out = make([]When, len(*in))
			copy(*out, *in)
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceAction) DeepCopyInto(out *ResourceAction) {
	*out = *in
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(v1.JSON)
		(*in).DeepCopyInto(*out)
	}
	if in.Templates != nil {
		in, out := &in.Templates, &out.Templates
		*out = new([]*v1.JSON)
		if **in != nil {
			in, out := *in, *out
			*out = make([]*v1.JSON, len(*in))
			for i := range *in {
				if (*in)[i] != nil {
					in, out := &(*in)[i], &(*out)[i]
					*out = new(v1.JSON)
					(*in).DeepCopyInto(*out)
				}
			}
		}
	}
	if in.StringTemplates != nil {
		in, out := &in.StringTemplates, &out.StringTemplates
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceAction.
func (in *ResourceAction) DeepCopy() *ResourceAction {
	if in == nil {
		return nil
	}
	out := new(ResourceAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResultInterface) DeepCopyInto(out *ResultInterface) {
	*out = *in
	if in.Plug != nil {
		in, out := &in.Plug, &out.Plug
		*out = make(map[string]*SchemaProperty, len(*in))
		for key, val := range *in {
			var outVal *SchemaProperty
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SchemaProperty)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.Socket != nil {
		in, out := &in.Socket, &out.Socket
		*out = make(map[string]*SchemaProperty, len(*in))
		for key, val := range *in {
			var outVal *SchemaProperty
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SchemaProperty)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResultInterface.
func (in *ResultInterface) DeepCopy() *ResultInterface {
	if in == nil {
		return nil
	}
	out := new(ResultInterface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchemaProperty) DeepCopyInto(out *SchemaProperty) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchemaProperty.
func (in *SchemaProperty) DeepCopy() *SchemaProperty {
	if in == nil {
		return nil
	}
	out := new(SchemaProperty)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Socket) DeepCopyInto(out *Socket) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Socket.
func (in *Socket) DeepCopy() *Socket {
	if in == nil {
		return nil
	}
	out := new(Socket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Socket) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SocketList) DeepCopyInto(out *SocketList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Socket, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SocketList.
func (in *SocketList) DeepCopy() *SocketList {
	if in == nil {
		return nil
	}
	out := new(SocketList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SocketList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SocketSpec) DeepCopyInto(out *SocketSpec) {
	*out = *in
	if in.Interface != nil {
		in, out := &in.Interface, &out.Interface
		*out = new(Interface)
		(*in).DeepCopyInto(*out)
	}
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make([]*Var, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Var)
				**out = **in
			}
		}
	}
	if in.ResultVars != nil {
		in, out := &in.ResultVars, &out.ResultVars
		*out = make([]*Var, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Var)
				**out = **in
			}
		}
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ConfigTemplate != nil {
		in, out := &in.ConfigTemplate, &out.ConfigTemplate
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Result != nil {
		in, out := &in.Result, &out.Result
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ResultTemplate != nil {
		in, out := &in.ResultTemplate, &out.ResultTemplate
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Apparatus != nil {
		in, out := &in.Apparatus, &out.Apparatus
		*out = new(SpecApparatus)
		(*in).DeepCopyInto(*out)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]*Resource, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Resource)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.ResultResources != nil {
		in, out := &in.ResultResources, &out.ResultResources
		*out = make([]*ResourceAction, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ResourceAction)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Validation != nil {
		in, out := &in.Validation, &out.Validation
		*out = new(SocketSpecValidation)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SocketSpec.
func (in *SocketSpec) DeepCopy() *SocketSpec {
	if in == nil {
		return nil
	}
	out := new(SocketSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SocketSpecValidation) DeepCopyInto(out *SocketSpecValidation) {
	*out = *in
	if in.NamespaceWhitelist != nil {
		in, out := &in.NamespaceWhitelist, &out.NamespaceWhitelist
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NamespaceBlacklist != nil {
		in, out := &in.NamespaceBlacklist, &out.NamespaceBlacklist
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SocketSpecValidation.
func (in *SocketSpecValidation) DeepCopy() *SocketSpecValidation {
	if in == nil {
		return nil
	}
	out := new(SocketSpecValidation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SocketStatus) DeepCopyInto(out *SocketStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CoupledPlugs != nil {
		in, out := &in.CoupledPlugs, &out.CoupledPlugs
		*out = make([]*CoupledPlug, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(CoupledPlug)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SocketStatus.
func (in *SocketStatus) DeepCopy() *SocketStatus {
	if in == nil {
		return nil
	}
	out := new(SocketStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SpecApparatus) DeepCopyInto(out *SpecApparatus) {
	*out = *in
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = new([]corev1.Container)
		if **in != nil {
			in, out := *in, *out
			*out = make([]corev1.Container, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SpecApparatus.
func (in *SpecApparatus) DeepCopy() *SpecApparatus {
	if in == nil {
		return nil
	}
	out := new(SpecApparatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Target) DeepCopyInto(out *Target) {
	*out = *in
	out.Gvk = in.Gvk
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Target.
func (in *Target) DeepCopy() *Target {
	if in == nil {
		return nil
	}
	out := new(Target)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Var) DeepCopyInto(out *Var) {
	*out = *in
	out.ObjRef = in.ObjRef
	out.FieldRef = in.FieldRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Var.
func (in *Var) DeepCopy() *Var {
	if in == nil {
		return nil
	}
	out := new(Var)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WaitForTarget) DeepCopyInto(out *WaitForTarget) {
	*out = *in
	out.Gvk = in.Gvk
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WaitForTarget.
func (in *WaitForTarget) DeepCopy() *WaitForTarget {
	if in == nil {
		return nil
	}
	out := new(WaitForTarget)
	in.DeepCopyInto(out)
	return out
}
