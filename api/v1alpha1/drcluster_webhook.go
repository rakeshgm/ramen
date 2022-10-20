/*
Copyright 2022 The RamenDR authors.

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var drclusterlog = logf.Log.WithName("drcluster-resource")

func (r *DRCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-ramendr-openshift-io-v1alpha1-drcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=ramendr.openshift.io,resources=drclusters,verbs=create;update,versions=v1alpha1,name=mdrcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &DRCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *DRCluster) Default() {
	drclusterlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-ramendr-openshift-io-v1alpha1-drcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=ramendr.openshift.io,resources=drclusters,verbs=create;update,versions=v1alpha1,name=vdrcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DRCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DRCluster) ValidateCreate() error {
	drclusterlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DRCluster) ValidateUpdate(old runtime.Object) error {
	drclusterlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DRCluster) ValidateDelete() error {
	drclusterlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
