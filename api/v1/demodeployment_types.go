/*
Copyright 2022.

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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DemoDeployment defines the desired state of DemoDeployment
type DemoDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DemoDeploymentImage is an example field of DemoDeployment. Edit demodeployment_types.go to remove/update
	DemoDeploymentImage string `json:"deploymentName"`
	Replicas            *int32 `json:"replicas"`
}

// DemoDeploymentStatus defines the observed state of DemoDeployment
type DemoDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DemoDeployment is the Schema for the demodeployments API
type DemoDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec               DemoDeploymentSpec   `json:"spec,omitempty"`
	Status             DemoDeploymentStatus `json:"status,omitempty"`
	DemoDeploymentName string               `json:"deploymentName"`
}

//+kubebuilder:object:root=true

// DemoDeploymentList contains a list of DemoDeployment
type DemoDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DemoDeployment `json:"items"`
}

// NewDemoDeployment creates a new Deployment for a DemoDeployment resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the DemoDeployment resource that 'owns' it.
func NewDemoDeployment(demoDeployment *DemoDeployment) *appsv1.Deployment {
	labels := map[string]string{
		"app":        demoDeployment.Spec.DemoDeploymentImage,
		"controller": demoDeployment.DemoDeploymentName,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      demoDeployment.DemoDeploymentName,
			Namespace: demoDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(demoDeployment, schema.GroupVersionKind{
					Group:   GroupVersion.Group,
					Version: GroupVersion.Version,
					Kind:    "DemoDeployment",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: demoDeployment.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  demoDeployment.DemoDeploymentName,
							Image: demoDeployment.Spec.DemoDeploymentImage,
						},
					},
				},
			},
		},
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DemoDeployment) DeepCopyInto(out *DemoDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DemoDeployment.
func (in *DemoDeployment) DeepCopy() *DemoDeployment {
	if in == nil {
		return nil
	}
	out := new(DemoDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DemoDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DemoDeploymentList) DeepCopyInto(out *DemoDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DemoDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DemoDeploymentList.
func (in *DemoDeploymentList) DeepCopy() *DemoDeploymentList {
	if in == nil {
		return nil
	}
	out := new(DemoDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DemoDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DemoDeploymentSpec) DeepCopyInto(out *DemoDeploymentSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DemoDeploymentSpec.
func (in *DemoDeploymentSpec) DeepCopy() *DemoDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(DemoDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DemoDeploymentStatus) DeepCopyInto(out *DemoDeploymentStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DemoDeploymentStatus.
func (in *DemoDeploymentStatus) DeepCopy() *DemoDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(DemoDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

func init() {
	SchemeBuilder.Register(&DemoDeployment{}, &DemoDeploymentList{})
}