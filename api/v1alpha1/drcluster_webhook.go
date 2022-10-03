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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var drclusterlog = logf.Log.WithName("drcluster-webhook")

func (r *DRCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//nolint
//+kubebuilder:webhook:path=/validate-ramendr-openshift-io-v1alpha1-drcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=ramendr.openshift.io,resources=drclusters,verbs=create;update,versions=v1alpha1,name=vdrcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DRCluster{}

// ValidateCreate checks if region has a value or not while creating
func (r *DRCluster) ValidateCreate() error {
	drclusterlog.Info("validate create", "name", r.Name)

	return r.ValidateDRCluster()
}

// ValidateUpdate implements immutability for Region and S3ProfileName
func (r *DRCluster) ValidateUpdate(old runtime.Object) error {
	drclusterlog.Info("validate update", "name", r.Name)

	oldDRCluster, ok := old.(*DRCluster)
	if !ok {
		return fmt.Errorf("error casting old DRCluster")
	}

	if r.Spec.Region != oldDRCluster.Spec.Region {
		return fmt.Errorf("Region cannot be changed")
	}

	if r.Spec.S3ProfileName != oldDRCluster.Spec.S3ProfileName {
		return fmt.Errorf("S3ProfileName cannot be changed")
	}

	return r.ValidateDRCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DRCluster) ValidateDelete() error {
	drclusterlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *DRCluster) ValidateDRCluster() error {
	if r.Spec.Region == "" {
		return fmt.Errorf("Region cannot be empty")
	}

	if r.Spec.S3ProfileName == "" {
		return fmt.Errorf("S3ProfileName cannot be empty")
	}

	// TODO: We can add other validations like validation of CIDRs format

	return nil
}
