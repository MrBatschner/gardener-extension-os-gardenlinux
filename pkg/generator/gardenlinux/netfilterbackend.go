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

package gardenlinux

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/gardener/gardener-extension-os-gardenlinux/pkg/apis/gardenlinux/v1alpha1"
	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/generator"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// defaultNetFilterBackend is the netfilter backend to fall back to
var defaultNetFilterBackend = v1alpha1.NetFilterBackend(v1alpha1.NetFilterIpTables)

var nfBackendUnitContent string = `[Unit]
Description=Configure netfilter backend for Gardener
After=cloud-config-downloader.service
Before=gardener-restart-system.service kubelet.service

[Install]
WantedBy=multi-user.target

[Service]
Type=oneshot
ExecStart=/opt/gardener/bin/configure_netfilter_backend.sh
RemainAfterExit=true
StandardOutput=journal
`

func ConfigureNetFilterBackend(osc *extensionsv1alpha1.OperatingSystemConfig, decoder runtime.Decoder) (*generator.File, *generator.Unit, error) {
	providerConfig := osc.Spec.ProviderConfig
	nfBackend := defaultNetFilterBackend

	if providerConfig != nil {
		obj := &v1alpha1.OperatingSystemConfiguration{}

		if _, _, err := decoder.Decode(providerConfig.Raw, nil, obj); err != nil {
			return nil, nil, fmt.Errorf("failed to decode provider config: %+v", err)
		}

		if len(obj.NetFilterBackend) != 0 {
			nfBackend = obj.NetFilterBackend
		}
	}

	config := map[string]interface{}{
		"netFilterBackend": nfBackend,
	}

	var buff bytes.Buffer
	nfBackendScriptTemplate, err := templates.ReadFile(filepath.Join("templates", "configure_netfilter_backend.sh.tpl"))
	if err != nil {
		return nil, nil, err
	}
	t, err := template.New("nfBackendScript").Parse(string(nfBackendScriptTemplate))
	if err != nil {
		return nil, nil, err
	}
	if err = t.Execute(&buff, config); err != nil {
		return nil, nil, err
	}

	nfBackendScript := &generator.File{
		Path:        filepath.Join(scriptLocation, "configure_netfilter_backend.sh"),
		Content:     buff.Bytes(),
		Permissions: &scriptPermissions,
	}

	nfBackendUnit := &generator.Unit{
		Name:    "gardener-configure-nfbackend.service",
		Content: []byte(nfBackendUnitContent),
	}

	return nfBackendScript, nfBackendUnit, nil
}
