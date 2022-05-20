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

package validator

import (
	"context"
	"fmt"
	"os"

	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/admission/common"
	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux"
	gardenlinuxinstall "github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux/install"
	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux/v1alpha1"
	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux/validation"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"
)

var decoder runtime.Decoder

func init() {
	scheme := runtime.NewScheme()
	if err := gardenlinuxinstall.AddToScheme(scheme); err != nil {
		runtimelog.Log.Error(err, "Could not update scheme")
		os.Exit(1)
	}
	decoder = serializer.NewCodecFactory(scheme).UniversalDecoder()
}

// NewShootValidator returns a new instance of a shoot validator.
func NewShootValidator() extensionswebhook.Validator {
	return &shoot{}
}

// shoot validates shoots
type shoot struct {
	common.ShootAdmissionHandler
}

// gardenlinuxOSConfig encapsulates an OperatingSystemConfiguration and its fieldPath
type gardenlinuxOSConfig struct {
	osc   *v1alpha1.OperatingSystemConfiguration
	field *field.Path
}

// Validate implements extensionswebhook.Validator.Validate
func (s *shoot) Validate(ctx context.Context, new, old client.Object) error {
	shoot, ok := new.(*core.Shoot)
	if !ok {
		return fmt.Errorf("wrong object type %T", new)
	}

	var oldShoot *core.Shoot
	if old != nil {
		var ok bool
		oldShoot, ok = old.(*core.Shoot)
		if !ok {
			return fmt.Errorf("wrong object type %T for old object", old)
		}
	}

	return s.validateShoot(ctx, oldShoot, shoot)
}

// validateShoot validated a Shoot rescource with focus on Garden Linux operating system configuration
func (s *shoot) validateShoot(_ context.Context, _, shoot *core.Shoot) error {
	osConfigurations, err := s.extractGardenLinuxOperatingSystemConfigurations(shoot)
	if err != nil {
		return fmt.Errorf("failed to extract Garden Linux relevant OperatingSystemConfigurations from Shoot: %e", err)
	}

	if len(osConfigurations) == 0 {
		return nil
	}

	allErrs := field.ErrorList{}
	for _, osc := range osConfigurations {
		allErrs = append(allErrs, validation.ValidateOperatingSystemConfiguration(osc.osc, osc.field)...)
	}

	return allErrs.ToAggregate()
}

// extractGardenLinuxOperatingSystemConfigurations extracts all Garden Linux relevant provider configurations from the Shoot spec
func (s *shoot) extractGardenLinuxOperatingSystemConfigurations(shoot *core.Shoot) ([]gardenlinuxOSConfig, error) {
	gardenLinuxOSConfigs := []gardenlinuxOSConfig{}

	for _, worker := range shoot.Spec.Provider.Workers {
		imageConfig := worker.Machine.Image
		if imageConfig.Name != gardenlinux.OSTypeGardenLinux {
			continue
		}

		if imageConfig.ProviderConfig == nil {
			continue
		}

		obj := &v1alpha1.OperatingSystemConfiguration{}
		if _, _, err := decoder.Decode(imageConfig.ProviderConfig.Raw, nil, obj); err != nil {
			return nil, fmt.Errorf("failed to decode provider config: %+v", err)
		}

		gardenLinuxOSConfig := gardenlinuxOSConfig{
			osc:   obj,
			field: field.NewPath(".Spec.Provider.Workers[" + worker.Name + "].Machine.Image.ProviderConfig"),
		}
		gardenLinuxOSConfigs = append(gardenLinuxOSConfigs, gardenLinuxOSConfig)
	}
	return gardenLinuxOSConfigs, nil
}
