/*
Copyright 2025.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	skjaifv1alpha1 "github.com/hauks1/skjaiferator/api/v1alpha1"
	skjaifv1beta1 "github.com/hauks1/skjaiferator/api/v1beta1"
)

// SvartSkjaifReconciler reconciles a SvartSkjaif object
type SvartSkjaifReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=skjaif.skjaiferator.no,resources=svartskjaifs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skjaif.skjaiferator.no,resources=svartskjaifs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=skjaif.skjaiferator.no,resources=svartskjaifs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SvartSkjaif object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *SvartSkjaifReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)
	svartSkjaif := &skjaifv1beta1.SvartSkjaif{}
	if err := r.Get(ctx, req.NamespacedName, svartSkjaif); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("svartSkjaif", "svartSkjaif", &svartSkjaif)
	// Set svartSkjaifAlpha1 values in the container if they are not kopp: mummi, vann: varmt og kaffe:svart
	spec := &svartSkjaif.Spec
	if spec.Kaffe != "svart" {
		logger.Info("handled kaffe not svart, setting to svart", "kaffe", spec.Kaffe)
		spec.Kaffe = "svart"
	}
	if spec.Kopp != "mummi" {
		logger.Info("handled kopp not mummi, setting to mummi", "kopp", spec.Kopp)
		spec.Kopp = "mummi"
	}
	if spec.Vann != "varmt" {
		logger.Info("handled vann not varmt, setting to varmt", "vann", spec.Vann)
		spec.Vann = "varmt"
	}
	logger.Info("Final spec state",
		"kaffe", spec.Kaffe,
		"kopp", spec.Kopp,
		"vann", spec.Vann)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SvartSkjaifReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skjaifv1alpha1.SvartSkjaif{}).
		Named("svartskjaif").
		Complete(r)
}
