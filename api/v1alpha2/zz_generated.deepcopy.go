// +build !ignore_autogenerated

/*
Copyright 2021.

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

package v1alpha2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kustomize/api/types"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Interface) DeepCopyInto(out *Interface) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
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

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Interface) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceList) DeepCopyInto(out *InterfaceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Interface, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceList.
func (in *InterfaceList) DeepCopy() *InterfaceList {
	if in == nil {
		return nil
	}
	out := new(InterfaceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InterfaceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceSpec) DeepCopyInto(out *InterfaceSpec) {
	*out = *in
	if in.Schemas != nil {
		in, out := &in.Schemas, &out.Schemas
		*out = make([]*InterfaceSpecSchema, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(InterfaceSpecSchema)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceSpec.
func (in *InterfaceSpec) DeepCopy() *InterfaceSpec {
	if in == nil {
		return nil
	}
	out := new(InterfaceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceSpecSchema) DeepCopyInto(out *InterfaceSpecSchema) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceSpecSchema.
func (in *InterfaceSpecSchema) DeepCopy() *InterfaceSpecSchema {
	if in == nil {
		return nil
	}
	out := new(InterfaceSpecSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceStatus) DeepCopyInto(out *InterfaceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceStatus.
func (in *InterfaceStatus) DeepCopy() *InterfaceStatus {
	if in == nil {
		return nil
	}
	out := new(InterfaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Plug) DeepCopyInto(out *Plug) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
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
	out.Interface = in.Interface
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make([]types.Var, len(*in))
		copy(*out, *in)
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
func (in *PlugSpecInterface) DeepCopyInto(out *PlugSpecInterface) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlugSpecInterface.
func (in *PlugSpecInterface) DeepCopy() *PlugSpecInterface {
	if in == nil {
		return nil
	}
	out := new(PlugSpecInterface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlugSpecSocket) DeepCopyInto(out *PlugSpecSocket) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlugSpecSocket.
func (in *PlugSpecSocket) DeepCopy() *PlugSpecSocket {
	if in == nil {
		return nil
	}
	out := new(PlugSpecSocket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlugStatus) DeepCopyInto(out *PlugStatus) {
	*out = *in
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
func (in *Socket) DeepCopyInto(out *Socket) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
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
	out.Interface = in.Interface
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make([]types.Var, len(*in))
		copy(*out, *in)
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
func (in *SocketStatus) DeepCopyInto(out *SocketStatus) {
	*out = *in
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
