/*
Copyright 2022 The Katanomi Authors.

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

package accessmode

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	StorageosProvisioner      = "kubernetes.io/storageos"
	QuobyteProvisioner        = "kubernetes.io/quobyte"
	HostPathProvisioner       = "kubernetes.io/host-path"
	AwsEbsProvisioner         = "kubernetes.io/aws-ebs"
	AzureFileProvisioner      = "kubernetes.io/azure-file"
	AzureDiskProvisioner      = "kubernetes.io/azure-disk"
	CephfsProvisioner         = "kubernetes.io/cephfs"
	CinderProvisioner         = "kubernetes.io/cinder"
	FCProvisioner             = "kubernetes.io/fc"
	FlockerProvisioner        = "kubernetes.io/flocker"
	GcePdProvisioner          = "kubernetes.io/gce-pd"
	GlusterfsProvisioner      = "kubernetes.io/glusterfs"
	IscsiProvisioner          = "kubernetes.io/iscsi"
	NFSProvisioner            = "kubernetes.io/nfs"
	RdbProvisioner            = "kubernetes.io/rbd"
	VsphereVolumeProvisioner  = "kubernetes.io/vsphere-volume"
	PortworxVolumeProvisioner = "kubernetes.io/portworx-volume"
)

var dftAccessModeRules = func() map[string][]corev1.PersistentVolumeAccessMode {
	return map[string][]corev1.PersistentVolumeAccessMode{StorageosProvisioner: {corev1.ReadWriteOnce, corev1.ReadOnlyMany},
		QuobyteProvisioner:        {corev1.ReadWriteOnce, corev1.ReadWriteMany, corev1.ReadOnlyMany},
		HostPathProvisioner:       {corev1.ReadWriteOnce},
		AwsEbsProvisioner:         {corev1.ReadWriteOnce},
		AzureFileProvisioner:      {corev1.ReadWriteOnce, corev1.ReadWriteMany, corev1.ReadOnlyMany},
		AzureDiskProvisioner:      {corev1.ReadWriteOnce},
		CephfsProvisioner:         {corev1.ReadWriteOnce, corev1.ReadWriteMany, corev1.ReadOnlyMany},
		CinderProvisioner:         {corev1.ReadWriteOnce},
		FCProvisioner:             {corev1.ReadWriteOnce, corev1.ReadOnlyMany},
		FlockerProvisioner:        {corev1.ReadWriteOnce},
		GcePdProvisioner:          {corev1.ReadWriteOnce, corev1.ReadOnlyMany},
		GlusterfsProvisioner:      {corev1.ReadWriteOnce, corev1.ReadWriteMany, corev1.ReadOnlyMany},
		IscsiProvisioner:          {corev1.ReadWriteOnce, corev1.ReadOnlyMany},
		NFSProvisioner:            {corev1.ReadWriteOnce, corev1.ReadWriteMany, corev1.ReadOnlyMany},
		RdbProvisioner:            {corev1.ReadWriteOnce, corev1.ReadOnlyMany},
		VsphereVolumeProvisioner:  {corev1.ReadWriteOnce},
		PortworxVolumeProvisioner: {corev1.ReadWriteOnce, corev1.ReadWriteMany},
	}
}
