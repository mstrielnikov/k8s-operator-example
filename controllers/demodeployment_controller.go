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

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog"

	scalev1 "mstrielnikov/k8s-operator-sample/api/v1"
)

const (
	MaxReplicasNum int32 = 2
)

var (
	ScaleError error = fmt.Errorf("Unable to scale Demo Deployment replicas greater then %v", MaxReplicasNum)
)

// DemoDeploymentReconciler reconciles a DemoDeployment object
type DemoDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=scale.mstrielnikov,resources=demodeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scale.mstrielnikov,resources=demodeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scale.mstrielnikov,resources=demodeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DemoDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *DemoDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Operator logic here
	var demoDeployment scalev1.DemoDeployment

	if err := r.Get(ctx, req.NamespacedName, &demoDeployment); err != nil {
		if apierrors.IsNotFound(err) {
			klog.Error(err, "The scalev1/DemoDeployment object %s is not found", demoDeployment.DemoDeploymentName)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		} else {
			klog.Error(err, "unable to fetch scalev1/demoDeployment")
			return ctrl.Result{}, err
		}
	}

	klog.Infof("Found scalev1/DemoDeployment object %v", demoDeployment)

	// Handle scale error
	if *demoDeployment.Spec.Replicas > MaxReplicasNum {
		klog.Error("Unable scalev1/DemoDeployment larger the ", MaxReplicasNum)
		return ctrl.Result{}, ScaleError
	}

	if err := r.Update(ctx, &demoDeployment); err != nil {
		if apierrors.IsConflict(err) {
			// The DemoDeployment has been updated since we read it.
			// Requeue the Pod to try to reconciliate again.
			return ctrl.Result{Requeue: true}, nil
		}
		if apierrors.IsNotFound(err) {
			klog.Error(err, "unable to fetch scalev1/DemoDeployment object %s", demoDeployment.DemoDeploymentName)
			return ctrl.Result{Requeue: true}, nil
		} else {
			klog.Error(err, "unable to update scalev1/DemoDeployment %s", demoDeployment.DemoDeploymentName)
			return ctrl.Result{}, err
		}
	}

	klog.Infof("Successfully updated DemoDeployment %s", demoDeployment.DemoDeploymentName)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalev1.DemoDeployment{}).
		Complete(r)
}
