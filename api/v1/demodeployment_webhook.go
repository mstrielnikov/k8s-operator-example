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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const MaxReplicasNum int32 = 2

var ErrScale error = fmt.Errorf("unable to scale demo deployment replicas greater then %v", MaxReplicasNum)

// log is for logging in this package.
var demodeploymentlog = logf.Log.WithName("demodeployment-resource")

func (r *DemoDeployment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-scale-mstrielnikov-v1-demodeployment,mutating=true,failurePolicy=fail,sideEffects=None,groups=scale.mstrielnikov,resources=demodeployments,verbs=create;update,versions=v1,name=mdemodeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &DemoDeployment{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *DemoDeployment) Default() {
	demodeploymentlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-scale-mstrielnikov-v1-demodeployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=scale.mstrielnikov,resources=demodeployments,verbs=create;update,versions=v1,name=vdemodeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DemoDeployment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DemoDeployment) ValidateCreate() error {
	demodeploymentlog.Info("validate create", "name", r.Name)
	// Handle scale error
	if *r.Spec.Replicas > MaxReplicasNum {
		klog.Error(ErrScale, "unable to create scalev1/DemoDeployment ", r.Name, " with replicas num larger then", MaxReplicasNum)
		return ErrScale
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DemoDeployment) ValidateUpdate(old runtime.Object) error {
	demodeploymentlog.Info("validate update", "name", r.Name)
	// Handle scale error
	if *r.Spec.Replicas > MaxReplicasNum {
		klog.Error(ErrScale, "Unable to create scalev1/DemoDeployment ", r.Name, " with replicas num larger then", MaxReplicasNum)
		return ErrScale
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DemoDeployment) ValidateDelete() error {
	demodeploymentlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
