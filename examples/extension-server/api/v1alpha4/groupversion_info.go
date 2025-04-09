// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha4

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const GroupName = "dp.wso2.com"

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme

	// Add additional group versions for conversion
	V1Alpha2GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha2"}
	V1Alpha3GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha3"}
)

func init() {
	SchemeBuilder.Register(&API{}, &APIList{})
	SchemeBuilder.GroupVersion = V1Alpha2GroupVersion
	SchemeBuilder.Register(&API{}, &APIList{})
	SchemeBuilder.GroupVersion = V1Alpha3GroupVersion
	SchemeBuilder.Register(&API{}, &APIList{})
}
