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

package controllers

import (
	"context"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	quotav1alpha1 "github.com/snapp-cab/quota-operator/api/v1alpha1"
)

// QuotaReconciler reconciles a Quota object
type QuotaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=quota.snappcloud.io,resources=quotas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=quota.snappcloud.io,resources=quotas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=quota.snappcloud.io,resources=quotas/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Quota object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *QuotaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Lookup the quota instance for this reconcile request
	quota := &quotav1alpha1.Quota{}
	err := r.Get(ctx, req.NamespacedName, quota)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get quota")
		return ctrl.Result{}, err
	}

	var found *corev1.ResourceQuota
	err = r.Get(ctx, types.NamespacedName{Name: quota.Name, Namespace: quota.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new quota
		corequota, err := r.corequotaFromQuota(quota)
		if err != nil {
			log.Error(err, "Error converting quota.snappcloud.io/v1 to quota/v1")
			return ctrl.Result{}, err
		}
		log.Info("Creating a new quota", "coreQuota.Namespace", corequota.Namespace, "coreQuota.Name", corequota.Name)
		err = r.Create(ctx, corequota)
		if err != nil {
			log.Error(err, "Failed to create new corequota", "corequota.Namespace", corequota.Namespace, "corequota.Name", corequota.Name)
			return ctrl.Result{}, err
		}
		// corequota created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get corequota")
		return ctrl.Result{}, err
	}

	// If Gslbcontents already exist, check if it is deeply equal with desrired state
	if !reflect.DeepEqual(quota.Spec, found.Spec) {
		log.Info("Updating corequota", "quota.Namespace", found.Namespace, "quota.Name", found.Name, "quota.Name", quota.Name)
		found, err = r.corequotaFromQuota(quota)
		if err != nil {
			log.Error(err, "Error converting quota.snappcloud.io/v1 to quota/v1")
			return ctrl.Result{}, err
		}
		err := r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update corequota", "quota.Namespace", found.Namespace, "Gslb.Name", found.Name)
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *QuotaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&quotav1alpha1.Quota{}).
		Owns(&corev1.ResourceQuota{}).
		Complete(r)
}

func (r *QuotaReconciler) corequotaFromQuota(q *quotav1alpha1.Quota) (*corev1.ResourceQuota, error) {
	corequota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      q.Name,
			Namespace: q.Namespace,
			Labels:    q.Labels,
		},
		Spec: q.Spec,
	}
	// Set Route instance as the owner and controller
	ctrl.SetControllerReference(q, corequota, r.Scheme)
	return corequota, nil
}
