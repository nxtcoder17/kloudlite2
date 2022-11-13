//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	"operators.kloudlite.io/lib/harbor"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborProject) DeepCopyInto(out *HarborProject) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborProject.
func (in *HarborProject) DeepCopy() *HarborProject {
	if in == nil {
		return nil
	}
	out := new(HarborProject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarborProject) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborProjectList) DeepCopyInto(out *HarborProjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HarborProject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborProjectList.
func (in *HarborProjectList) DeepCopy() *HarborProjectList {
	if in == nil {
		return nil
	}
	out := new(HarborProjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarborProjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborProjectSpec) DeepCopyInto(out *HarborProjectSpec) {
	*out = *in
	if in.Project != nil {
		in, out := &in.Project, &out.Project
		*out = new(harbor.Project)
		**out = **in
	}
	if in.Webhook != nil {
		in, out := &in.Webhook, &out.Webhook
		*out = new(harbor.Webhook)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborProjectSpec.
func (in *HarborProjectSpec) DeepCopy() *HarborProjectSpec {
	if in == nil {
		return nil
	}
	out := new(HarborProjectSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborUserAccount) DeepCopyInto(out *HarborUserAccount) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborUserAccount.
func (in *HarborUserAccount) DeepCopy() *HarborUserAccount {
	if in == nil {
		return nil
	}
	out := new(HarborUserAccount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarborUserAccount) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborUserAccountList) DeepCopyInto(out *HarborUserAccountList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HarborUserAccount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborUserAccountList.
func (in *HarborUserAccountList) DeepCopy() *HarborUserAccountList {
	if in == nil {
		return nil
	}
	out := new(HarborUserAccountList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarborUserAccountList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborUserAccountSpec) DeepCopyInto(out *HarborUserAccountSpec) {
	*out = *in
	in.OperatorProps.DeepCopyInto(&out.OperatorProps)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborUserAccountSpec.
func (in *HarborUserAccountSpec) DeepCopy() *HarborUserAccountSpec {
	if in == nil {
		return nil
	}
	out := new(HarborUserAccountSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorProps) DeepCopyInto(out *OperatorProps) {
	*out = *in
	if in.HarborUser != nil {
		in, out := &in.HarborUser, &out.HarborUser
		*out = new(harbor.User)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorProps.
func (in *OperatorProps) DeepCopy() *OperatorProps {
	if in == nil {
		return nil
	}
	out := new(OperatorProps)
	in.DeepCopyInto(out)
	return out
}
