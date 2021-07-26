/*
Copyright The Helm Authors.

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

package diffutil

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/yaml"

	"helm.sh/helm/v3/pkg/kube"
)

// WriteKubeResult writes contents from a kube.Result to differ.
func WriteKubeResult(differ *Differ, result *kube.Result) error {
	infoName := func(info *resource.Info) string {
		return fmt.Sprintf(
			// so that resources within the same namespace appear near
			"%s.%s.%s.%s.%s.yaml",
			info.Namespace,
			info.Name,
			info.Mapping.GroupVersionKind.Group,
			info.Mapping.GroupVersionKind.Version,
			info.Mapping.GroupVersionKind.Kind,
		)
	}

	for _, info := range result.Created {
		b, err := yaml.Marshal(info.Object)
		if err != nil {
			return err
		}
		if err := differ.WriteNew(infoName(info), b); err != nil {
			return err
		}
	}
	for _, info := range result.Updated {
		if result.LiveBeforeUpdate[info] == nil {
			// shouldn't happen
			return fmt.Errorf("no live before update version of %s", infoName(info))
		}
		bOld, err := yaml.Marshal(result.LiveBeforeUpdate[info])
		if err != nil {
			return err
		}
		bNew, err := yaml.Marshal(info.Object)
		if err != nil {
			return err
		}

		if err := differ.WriteNew(infoName(info), bNew); err != nil {
			return err
		}
		if err := differ.WriteOld(infoName(info), bOld); err != nil {
			return err
		}
	}
	for _, info := range result.Deleted {
		b, err := yaml.Marshal(info.Object)
		if err != nil {
			return err
		}
		if err := differ.WriteOld(infoName(info), b); err != nil {
			return err
		}
	}

	return nil
}
