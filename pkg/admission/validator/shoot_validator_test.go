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
	"encoding/json"
	"testing"

	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux"
	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux/v1alpha1"
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	gardencore "github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx            context.Context = context.Background()
	sh             *gardencore.Shoot
	worker         gardencore.Worker
	shootValidator extensionswebhook.Validator = NewShootValidator()
)

var _ = Describe("Garden Linux OS Validator Test", func() {

	BeforeEach(func() {
		providerConfig := encode(&v1alpha1.OperatingSystemConfiguration{})

		worker = gardencore.Worker{
			Machine: gardencore.Machine{
				Image: &gardencore.ShootMachineImage{
					Name: gardenlinux.OSTypeGardenLinux,
					ProviderConfig: &runtime.RawExtension{
						Raw: providerConfig,
					},
				},
			},
		}

		sh = &gardencore.Shoot{
			Spec: gardencore.ShootSpec{
				Provider: gardencore.Provider{
					Workers: []gardencore.Worker{
						worker,
					},
				},
			},
		}
	})

	It("should extract the Garden Linux operating system configs from the Shoot manifest", func() {
		w1 := worker.DeepCopy()
		w2 := w1.DeepCopy()

		sh.Spec.Provider.Workers = append(sh.Spec.Provider.Workers, *w1, *w2)
		validator := shoot{}

		oscs, err := validator.extractGardenLinuxOperatingSystemConfigurations(sh)
		Expect(err).ToNot(HaveOccurred())
		Expect(oscs).To(HaveLen(3))
	})

	It("should omit operating system configs that are not for Garden Linux", func() {
		w1 := worker.DeepCopy()
		w1.Machine.Image.Name = "notgardenlinux"

		sh.Spec.Provider.Workers = append(sh.Spec.Provider.Workers, *w1)
		validator := shoot{}

		oscs, err := validator.extractGardenLinuxOperatingSystemConfigurations(sh)
		Expect(err).ToNot(HaveOccurred())
		Expect(oscs).To(HaveLen(1))
	})

	It("should accept a Shoot manifest with provider config with default values", func() {
		err := shootValidator.Validate(ctx, sh, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should deny a Shoot manifest with provider config invalid values", func() {
		providerConfig := &v1alpha1.OperatingSystemConfiguration{
			LinuxSecurityModule: "foo",
			NetFilterBackend:    "bar",
			CgroupVersion:       "v1337",
		}

		sh.Spec.Provider.Workers[0].Machine.Image.ProviderConfig = &runtime.RawExtension{
			Raw: encode(providerConfig),
		}

		err := shootValidator.Validate(ctx, sh, nil)
		Expect(err).To(HaveOccurred())
	})
})

func encode(obj runtime.Object) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validator Suite")
}
