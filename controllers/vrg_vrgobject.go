// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/controllers/util"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (v *VRGInstance) vrgObjectProtect(result *ctrl.Result, s3StoreAccessors []s3StoreAccessor) {
	vrg := v.instance
	eventReporter := v.reconciler.eventRecorder
	log := v.log

	for _, s3StoreAccessor := range s3StoreAccessors {
		if err := VrgObjectProtect(s3StoreAccessor.ObjectStorer, *vrg); err != nil {
			util.ReportIfNotPresent(
				eventReporter, vrg, corev1.EventTypeWarning, util.EventReasonVrgUploadFailed, err.Error(),
			)

			const message = "VRG Kube object protect error"

			log.Error(err, message, "profile", s3StoreAccessor.profileName)

			v.vrgObjectProtected = newVRGClusterDataUnprotectedCondition(vrg.Generation, message)
			result.Requeue = true

			return
		}

		log.Info("VRG Kube object protected", "profile", s3StoreAccessor.profileName)

		v.vrgObjectProtected = newVRGClusterDataProtectedCondition(vrg.Generation, clusterDataProtectedTrueMessage)
	}
}

const vrgS3ObjectNameSuffix = "a"

func VrgObjectProtect(objectStorer ObjectStorer, vrg ramen.VolumeReplicationGroup) error {
	return uploadTypedObject(objectStorer, s3PathNamePrefix(vrg.Namespace, vrg.Name), vrgS3ObjectNameSuffix, vrg)
}

func VrgObjectUnprotect(objectStorer ObjectStorer, vrg ramen.VolumeReplicationGroup) error {
	return DeleteTypedObjects(objectStorer, s3PathNamePrefix(vrg.Namespace, vrg.Name), vrgS3ObjectNameSuffix, vrg)
}

func vrgObjectDownload(objectStorer ObjectStorer, pathName string, vrg *ramen.VolumeReplicationGroup) error {
	return downloadTypedObject(objectStorer, pathName, vrgS3ObjectNameSuffix, vrg)
}