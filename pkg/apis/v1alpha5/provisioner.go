/*
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

package v1alpha5

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/karpenter-core/pkg/apis/provisioning/v1alpha5"
	"github.com/aws/karpenter/pkg/apis/v1alpha1"
)

// Provisioner is an alias type for additional validation
// +kubebuilder:object:root=true
type Provisioner v1alpha5.Provisioner

func (p *Provisioner) Validate(ctx context.Context) (errs *apis.FieldError) {
	if p.Spec.Provider == nil {
		return nil
	}
	provider, err := v1alpha1.Deserialize(p.Spec.Provider)
	if err != nil {
		return apis.ErrGeneric(err.Error())
	}
	return provider.Validate()
}

func (p *Provisioner) SetDefaults(ctx context.Context) {
	for key, value := range map[string]string{
		v1alpha5.LabelCapacityType: ec2.DefaultTargetCapacityTypeOnDemand,
		v1.LabelArchStable:         v1alpha5.ArchitectureAmd64,
	} {
		hasLabel := false
		if _, ok := p.Spec.Labels[key]; ok {
			hasLabel = true
		}
		for _, requirement := range p.Spec.Requirements {
			if requirement.Key == key {
				hasLabel = true
			}
		}
		if !hasLabel {
			p.Spec.Requirements = append(p.Spec.Requirements, v1.NodeSelectorRequirement{
				Key: key, Operator: v1.NodeSelectorOpIn, Values: []string{value},
			})
		}
	}
}