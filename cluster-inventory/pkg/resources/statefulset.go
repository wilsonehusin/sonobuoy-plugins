/*
Copyright the Sonobuoy contributors 2020

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

package resources

import (
	"fmt"

	"github.com/vmware-tanzu/sonobuoy-plugins/cluster-inventory/pkg/reports"

	appsv1 "k8s.io/api/apps/v1"
)

type StatefulSet struct {
	appsv1.StatefulSet
	Pods map[string]*Pod
}

func (s StatefulSet) statusMessage() string {
	return fmt.Sprintf("Desired: %d, Total: %d, Current: %d, Ready: %d",
		s.Spec.Replicas, s.Status.Replicas, s.Status.CurrentReplicas, s.Status.ReadyReplicas)
}

func (s StatefulSet) GenerateSonobuoyItem() reports.SonobuoyResultsItem {
	item := reports.SonobuoyResultsItem{
		Name:   s.Name,
		Status: s.statusMessage(),
		Metadata: map[string]string{
			"kind": "StatefulSet",
			"uid":  string(s.UID),
		},
		Details: map[string]interface{}{
			"status":              s.Status,
			"replicas":            s.Spec.Replicas,
			"podManagementPolicy": s.Spec.PodManagementPolicy,
			"updateStrategy":      s.Spec.UpdateStrategy,
			"serviceName":         s.Spec.ServiceName,
		},
	}

	if s.Spec.Selector != nil {
		item.Details["selector"] = s.Spec.Selector
	}

	if s.Spec.Template.Spec.NodeSelector != nil {
		item.Details["nodeSelector"] = s.Spec.Template.Spec.NodeSelector
	}

	if s.Spec.RevisionHistoryLimit != nil {
		item.Details["revisionHistoryLimit"] = s.Spec.RevisionHistoryLimit
	}

	if len(s.Labels) > 0 {
		item.Details["labels"] = s.Labels
	}

	for _, pod := range s.Pods {
		item.Items = append(item.Items, pod.GenerateSonobuoyItem())
	}

	return item
}
