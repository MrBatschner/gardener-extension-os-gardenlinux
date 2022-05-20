// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"fmt"

	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux/v1alpha1"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateOperatingSystemConfiguration validates an OperatingSystemConfiguration object.
func ValidateOperatingSystemConfiguration(osc *v1alpha1.OperatingSystemConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	lsmPath := fldPath.Child("linuxSecurityModule")
	allErrs = append(allErrs, validateLinuxSecurityModule(osc.LinuxSecurityModule, lsmPath)...)

	nfFrontendPath := fldPath.Child("netfilterFrontend")
	allErrs = append(allErrs, validateNetfilterFrontend(osc.NetFilterBackend, nfFrontendPath)...)

	cgroupPath := fldPath.Child("cgroupVersion")
	allErrs = append(allErrs, validateCgroupVersion(osc.CgroupVersion, cgroupPath)...)

	return allErrs
}

// validateLinuxSecurityModule validates the Linux Security Module
func validateLinuxSecurityModule(lsm v1alpha1.LinuxSecurityModule, fldPath *field.Path) field.ErrorList {
	var allErrs = field.ErrorList{}

	if lsm != v1alpha1.LsmAppArmor && lsm != v1alpha1.LsmSeLinux {
		allErrs = append(allErrs, field.Invalid(fldPath, lsm, fmt.Sprintf("must be either %s or %s", v1alpha1.LsmAppArmor, v1alpha1.LsmSeLinux)))
	}

	return allErrs
}

// validateNetfilterFrontend validates the netfilter frontend
func validateNetfilterFrontend(nfFrontend v1alpha1.NetFilterBackend, fldPath *field.Path) field.ErrorList {
	var allErrs = field.ErrorList{}

	if nfFrontend != v1alpha1.NetFilterIpTables && nfFrontend != v1alpha1.NetFilterNfTables {
		allErrs = append(allErrs, field.Invalid(fldPath, nfFrontend, fmt.Sprintf("must be either %s or %s", v1alpha1.NetFilterIpTables, v1alpha1.NetFilterNfTables)))
	}

	return allErrs
}

// validateCgroupVersion validates the cgroup version
func validateCgroupVersion(cgroupVersion v1alpha1.CgroupVersion, fldPath *field.Path) field.ErrorList {
	var allErrs = field.ErrorList{}

	if cgroupVersion != v1alpha1.CgroupVersionV1 && cgroupVersion != v1alpha1.CgroupVersionV2 {
		allErrs = append(allErrs, field.Invalid(fldPath, cgroupVersion, fmt.Sprintf("must be either %s or %s", v1alpha1.CgroupVersionV1, v1alpha1.CgroupVersionV2)))
	}

	return allErrs
}
