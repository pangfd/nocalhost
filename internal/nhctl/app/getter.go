/*
Copyright 2020 The Nocalhost Authors.
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

package app

import "nocalhost/pkg/nhctl/clientgoutils"

func (a *Application) GetNamespace() string {
	return a.AppProfile.Namespace
}

func (a *Application) GetType() AppType {
	return a.AppProfile.AppType
}

func (a *Application) GetKubeconfig() string {
	return a.AppProfile.Kubeconfig
}

func (a *Application) IsHelm() bool {
	return a.AppProfile.AppType == Helm || a.AppProfile.AppType == HelmRepo
}

func (a *Application) IsManifest() bool {
	return a.AppProfile.AppType == Manifest
}

func (a *Application) GetClient() *clientgoutils.ClientGoUtils {
	return a.client
}
