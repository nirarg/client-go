/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2017 Red Hat, Inc.
 *
 */

package v1

import (
	"encoding/json"
	"fmt"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/klog"
	"kubevirt.io/client-go/precond"
	"sigs.k8s.io/yaml"
)

// This is meant for testing
func NewMinimalVMI(name string) *VirtualMachineInstance {
	return NewMinimalVMIWithNS(k8sv1.NamespaceDefault, name)
}

// This is meant for testing
func NewMinimalVMIWithNS(namespace, name string) *VirtualMachineInstance {
	precond.CheckNotEmpty(name)
	vmi := NewVMIReferenceFromNameWithNS(namespace, name)
	vmi.Spec = VirtualMachineInstanceSpec{Domain: DomainSpec{}}
	vmi.Spec.Domain.Resources.Requests = k8sv1.ResourceList{
		k8sv1.ResourceMemory: resource.MustParse("8192Ki"),
	}
	vmi.TypeMeta = k8smetav1.TypeMeta{
		APIVersion: GroupVersion.String(),
		Kind:       "VirtualMachineInstance",
	}
	return vmi
}

// RawExtensionFromProviderSpec marshals the machine provider spec.
func RawExtensionFromProviderSpec(spec *VirtualMachineInstanceSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error
	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, fmt.Errorf("error marshalling providerSpec: %v", err)
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// RawExtensionFromProviderStatus marshals the machine provider status
func RawExtensionFromProviderStatus(status *VirtualMachineInstanceStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, fmt.Errorf("error marshalling providerStatus: %v", err)
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// ProviderSpecFromRawExtension unmarshals a raw extension into an AWSMachineProviderSpec type
func ProviderSpecFromRawExtension(rawExtension *runtime.RawExtension) (*VirtualMachineInstanceSpec, error) {
	if rawExtension == nil {
		return &VirtualMachineInstanceSpec{}, nil
	}

	spec := new(VirtualMachineInstanceSpec)
	if err := yaml.Unmarshal(rawExtension.Raw, &spec); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerSpec: %v", err)
	}

	klog.V(5).Infof("Got provider Spec from raw extension: %+v", spec)
	return spec, nil
}

// ProviderStatusFromRawExtension unmarshals a raw extension into an AWSMachineProviderStatus type
func ProviderStatusFromRawExtension(rawExtension *runtime.RawExtension) (*VirtualMachineInstanceStatus, error) {
	if rawExtension == nil {
		return &VirtualMachineInstanceStatus{}, nil
	}

	providerStatus := new(VirtualMachineInstanceStatus)
	if err := yaml.Unmarshal(rawExtension.Raw, providerStatus); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerStatus: %v", err)
	}

	klog.V(5).Infof("Got provider Status from raw extension: %+v", providerStatus)
	return providerStatus, nil
}

// VirtualMachineInstancePhaseToString returns a pointer to the string value of the VirtualMachineInstancePhase passed.
func VirtualMachineInstancePhaseToString(src VirtualMachineInstancePhase) *string {
	value := string(src)
	return &value
}
