/*
Copyright 2024.

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

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretReconciler reconciles a Secret object
type SecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.example.com,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=secrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Secret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	//get instance
	secret_instance := &cachev1alpha1.Secret{}

	//update secret based on fetched data
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret_instance.Name,
			Namespace: secret_instance.Namespace,
		},
		StringData: map[string]string{
			"yourKey": secret_instance.Spec.SecretData,
		},
	}

	if err := r.CreateOrUpdateSecret(ctx, secret); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Secret{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func (r *SecretReconciler) CreateOrUpdateSecret(ctx context.Context, secret *corev1.Secret) error {
	foundSecret := &corev1.Secret{}
	//get current state of secret
	err := r.Get(ctx,types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, foundSecret)

	//secret does not exist, create it
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, secret) //create
	} else if err != nil {
		return err
	}

	//secret exists so we update it
	foundSecret.StringData = secret.StringData
	return r.Update(ctx, foundSecret) //update
}
