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

	scalev1 "mstrielnikov/k8s-operator-sample/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	MaxReplicasNum           int32 = 2
	DemoDeploymentController       = "DemoDeploymentController"
)

var (
	ErrScale = fmt.Errorf("unable to scale demo deployment replicas greater then %v", MaxReplicasNum)
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
	log := log.FromContext(ctx)

	// Operator logic here
	if err := r.handleCreate(ctx, req); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.handleUpdate(ctx, req); err != nil {
		if apierrors.IsConflict(err) {
			// The DemoDeployment has been updated since we read it.
			// Requeue the DemoDeployment to try to reconciliate again.
			log.Info("DemoDeployment has been updated since we read")
			return ctrl.Result{Requeue: true}, nil
		} else if apierrors.IsNotFound(err) {
			// The DemoDeployment has been deleted since we read it.
			// Requeue the DemoDeployment to try to reconciliate again.
			log.Info("object scalev1/DemoDeployment is not found")
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	if err := r.handleList(ctx, req); err != nil {
		if apierrors.IsNotFound(err) {
			// The DemoDeployment has been deleted since we read it.
			// Requeue the DemoDeployment to try to reconciliate again.
			log.Info("object scalev1/DemoDeployment is not found")
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalev1.DemoDeployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func (r *DemoDeploymentReconciler) handleList(ctx context.Context, req ctrl.Request) error {
	var demoDeployment scalev1.DemoDeployment
	if err := r.Get(ctx, req.NamespacedName, &demoDeployment); err != nil {
		klog.Error(err, "unable to fetch scalev1/demoDeployment %v", demoDeployment)
		return err
	}
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(demoDeployment.Namespace),
		client.MatchingLabels(demoDeployment.GetLabels()),
	}

	if err := r.List(ctx, podList, listOpts...); err != nil {
		klog.Error(err, "Falied to list pods", "DemoDeployment.Namespace", demoDeployment.Namespace, "DemoDeployment.Name", demoDeployment.Name)
		return err
	}
	klog.Infof("Successfully listed DemoDeployment %s", demoDeployment.Name)
	return nil
}

func (r *DemoDeploymentReconciler) handleUpdate(ctx context.Context, req ctrl.Request) error {
	var demoDeployment scalev1.DemoDeployment
	if err := r.Get(ctx, req.NamespacedName, &demoDeployment); err != nil {
		klog.Error(err, "unable to fetch scalev1/demoDeployment %v", demoDeployment)
		return err
	}
	// // Handle scale error
	if *demoDeployment.Spec.Replicas > MaxReplicasNum {
		klog.Error(ErrScale, "Unable to create scalev1/DemoDeployment with replicas num larger then", MaxReplicasNum)
		return ErrScale
	}
	if err := r.Update(ctx, &demoDeployment); err != nil {
		klog.Error(err, "unable to update scalev1/DemoDeployment %s", demoDeployment.Name)
		return err
	}
	klog.Infof("Successfully updated DemoDeployment %s", demoDeployment.Name)
	return nil
}

func (r *DemoDeploymentReconciler) handleCreate(ctx context.Context, req ctrl.Request) error {
	var demoDeployment scalev1.DemoDeployment
	if err := r.Get(ctx, req.NamespacedName, &demoDeployment); err != nil {
		klog.Error(err, "unable to fetch scalev1/demoDeployment %v", demoDeployment)
		return err
	}
	// Handle scale error
	if *demoDeployment.Spec.Replicas > MaxReplicasNum {
		klog.Error(ErrScale, "unable to create scalev1/DemoDeployment with replicas num larger then", MaxReplicasNum)
		return ErrScale
	}
	deployment := newDemoDeployment(demoDeployment)
	err := r.Create(ctx, &deployment)
	if apierrors.IsAlreadyExists(err) {
		klog.Info(err, "DemoDeployment in namespace ", deployment.Namespace, "Deployment.Name ", deployment.Name, "exists already")
		return nil
	}
	if err != nil && !apierrors.IsAlreadyExists(err) {
		klog.Error(err, "failed to create new Deployment ", "Deployment.Namespace ", deployment.Namespace, "Deployment.Name ", deployment.Name)
		return err
	}
	return nil
}

func newDemoDeployment(demoDeployment scalev1.DemoDeployment) appsv1.Deployment {
	labels := map[string]string{
		"app":        demoDeployment.Name,
		"namespace":  demoDeployment.Namespace,
		"controller": DemoDeploymentController,
	}
	return appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      demoDeployment.Name,
			Labels:    labels,
			Namespace: demoDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&demoDeployment, schema.GroupVersionKind{
					Group:   scalev1.GroupVersion.Group,
					Version: scalev1.GroupVersion.Version,
					Kind:    demoDeployment.Kind,
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: demoDeployment.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      demoDeployment.Name,
					Namespace: demoDeployment.Namespace,
					Labels:    labels,
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(&demoDeployment, schema.GroupVersionKind{
							Group:   scalev1.GroupVersion.Group,
							Version: scalev1.GroupVersion.Version,
							Kind:    demoDeployment.Kind,
						}),
					},
				},
				Spec: corev1.PodSpec{
					Containers: demoDeployment.Spec.Template.Spec.Containers,
				},
			},
			Strategy: appsv1.DeploymentStrategy{},
		},
		Status: appsv1.DeploymentStatus{},
	}
}
