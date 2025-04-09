// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha4

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/envoyproxy/gateway/examples/extension-server/api/v1alpha1"
	"github.com/envoyproxy/gateway/examples/extension-server/api/v1alpha3"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="API Name",type=string,JSONPath=`.spec.apiName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.apiVersion`
// +kubebuilder:printcolumn:name="BasePath",type=string,JSONPath=`.spec.basePath`
// +kubebuilder:printcolumn:name="Organization",type=string,JSONPath=`.spec.organization`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +groupName=dp.wso2.com
// +kubebuilder:storageversion
// API is the Schema for the apis API
type API struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APISpec   `json:"spec"`
	Status APIStatus `json:"status,omitempty"`
}

// APISpec defines the desired state of API
type APISpec struct {
	// APIName is the unique name of the API can be used to uniquely identify an API.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=60
	// +kubebuilder:validation:Pattern=`^[^~!@#;:%^*()+={}|\<>"'',&$\[\]\/]*$`
	APIName string `json:"apiName"`

	// APIType denotes the type of the API. Possible values could be REST, GraphQL, Async
	// +kubebuilder:validation:Enum=REST
	APIType string `json:"apiType"`

	// APIVersion is the version number of the API.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:Pattern=`^[^~!@#;:%^*()+={}|\<>"'',&/$\[\]\s+\/]+$`
	APIVersion string `json:"apiVersion"`

	// BasePath denotes the basepath of the API. e.g: /pet-store-api/1.0.6
	// +kubebuilder:validation:Pattern=`^[/][a-zA-Z0-9~/_.-]*$`
	BasePath string `json:"basePath"`

	// DefinitionPath contains the path to expose the API definition.
	// +kubebuilder:default=/api-definition
	// +kubebuilder:validation:MinLength=1
	DefinitionPath string `json:"definitionPath"`

	// Organization denotes the organization related to the API
	Organization string `json:"organization"`

	// Production contains a list of references to HttpRoutes of type HttpRoute.
	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
	// +kubebuilder:validation:MaxItems=1
	// +nullable
	Production []RouteRef `json:"production,omitempty"`

	// IsDefaultVersion indicates whether this API version should be used as a default API
	IsDefaultVersion bool `json:"isDefaultVersion"`

	// SystemAPI denotes if it is an internal system API.
	SystemAPI bool `json:"systemAPI"`

	// APIProperties denotes the custom properties of the API.
	// +nullable
	APIProperties []Property `json:"apiProperties,omitempty"`

	// DefinitionFileRef contains the OpenAPI 3 or Swagger definition of the API in a ConfigMap.
	DefinitionFileRef string `json:"definitionFileRef,omitempty"`

	// Sandbox contains a list of references to HttpRoutes of type HttpRoute.
	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
	// +kubebuilder:validation:MaxItems=1
	// +nullable
	Sandbox []RouteRef `json:"sandbox,omitempty"`
}

// RouteRef contains the environment specific configuration
type RouteRef struct {
	// HTTPRouteRefs denotes the environment of the API.
	// +kubebuilder:validation:Required
	HTTPRouteRefs []string `json:"httpRouteRefs"`
}

// Property holds key value pair of APIProperties
type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// APIStatus defines the observed state of API
type APIStatus struct {
	// DeploymentStatus denotes the deployment status of the API
	DeploymentStatus DeploymentStatus `json:"deploymentStatus,omitempty"`
}

// DeploymentStatus denotes the deployment status of the API
type DeploymentStatus struct {
	// Accepted represents whether the API is accepted or not.
	// +kubebuilder:validation:Required
	Accepted bool `json:"accepted"`

	// Events contains a list of events related to the API.
	Events []string `json:"events,omitempty"`

	// Message represents a user friendly message that explains the current state of the API.
	Message string `json:"message,omitempty"`

	// Status denotes the state of the API in its lifecycle.
	// Possible values could be Accepted, Invalid, Deploy etc.
	// +kubebuilder:validation:Required
	Status string `json:"status"`

	// TransitionTime represents the last known transition timestamp.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=date-time
	TransitionTime *metav1.Time `json:"transitionTime"`
}

// +kubebuilder:object:root=true
// APIList contains a list of API resources.
type APIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []API `json:"items"`
}

// Allowing Kubernetes API server to recognize the resources defined in the API CR yaml
func init() {
	SchemeBuilder.Register(&API{}, &APIList{})
}

// ConvertTo converts this API to the Hub version (v1alpha1).
func (src *API) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha1.API)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.API, got %T", dstRaw)
	}

	// Copy TypeMeta
	dst.TypeMeta = src.TypeMeta

	// Copy ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	dst.Spec.APIName = src.Spec.APIName
	dst.Spec.APIType = src.Spec.APIType
	dst.Spec.APIVersion = src.Spec.APIVersion
	dst.Spec.BasePath = src.Spec.BasePath
	dst.Spec.DefinitionPath = src.Spec.DefinitionPath
	dst.Spec.Organization = src.Spec.Organization
	dst.Spec.IsDefaultVersion = src.Spec.IsDefaultVersion
	dst.Spec.SystemAPI = src.Spec.SystemAPI
	dst.Spec.APIProperties = src.Spec.APIProperties
	dst.Spec.DefinitionFileRef = src.Spec.DefinitionFileRef

	// Copy Production and Sandbox (already using HTTPRouteRefs in v1alpha4)
	dst.Spec.Production = src.Spec.Production
	dst.Spec.Sandbox = src.Spec.Sandbox

	// Convert Status
	dst.Status.DeploymentStatus.Accepted = src.Status.DeploymentStatus.Accepted
	dst.Status.DeploymentStatus.Events = src.Status.DeploymentStatus.Events
	dst.Status.DeploymentStatus.Message = src.Status.DeploymentStatus.Message
	dst.Status.DeploymentStatus.Status = src.Status.DeploymentStatus.Status
	dst.Status.DeploymentStatus.TransitionTime = src.Status.DeploymentStatus.TransitionTime

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha1) to this API.
func (dst *API) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha1.API)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.API, got %T", srcRaw)
	}

	// Copy TypeMeta
	dst.TypeMeta = src.TypeMeta

	// Copy ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	dst.Spec.APIName = src.Spec.APIName
	dst.Spec.APIType = src.Spec.APIType
	dst.Spec.APIVersion = src.Spec.APIVersion
	dst.Spec.BasePath = src.Spec.BasePath
	dst.Spec.DefinitionPath = src.Spec.DefinitionPath
	dst.Spec.Organization = src.Spec.Organization
	dst.Spec.IsDefaultVersion = src.Spec.IsDefaultVersion
	dst.Spec.SystemAPI = src.Spec.SystemAPI
	dst.Spec.APIProperties = src.Spec.APIProperties
	dst.Spec.DefinitionFileRef = src.Spec.DefinitionFileRef

	// Copy Production and Sandbox (already using HTTPRouteRefs in v1alpha1)
	dst.Spec.Production = src.Spec.Production
	dst.Spec.Sandbox = src.Spec.Sandbox

	// Convert Status
	dst.Status.DeploymentStatus.Accepted = src.Status.DeploymentStatus.Accepted
	dst.Status.DeploymentStatus.Events = src.Status.DeploymentStatus.Events
	dst.Status.DeploymentStatus.Message = src.Status.DeploymentStatus.Message
	dst.Status.DeploymentStatus.Status = src.Status.DeploymentStatus.Status
	dst.Status.DeploymentStatus.TransitionTime = src.Status.DeploymentStatus.TransitionTime

	return nil
}

// ConvertToV1alpha3 converts this API to the v1alpha3 version.
func (src *API) ConvertToV1alpha3(dstRaw *v1alpha3.API) error {
	// Copy TypeMeta
	dstRaw.TypeMeta = src.TypeMeta

	// Copy ObjectMeta
	dstRaw.ObjectMeta = src.ObjectMeta

	// Convert Spec
	dstRaw.Spec.APIName = src.Spec.APIName
	dstRaw.Spec.APIType = src.Spec.APIType
	dstRaw.Spec.APIVersion = src.Spec.APIVersion
	dstRaw.Spec.BasePath = src.Spec.BasePath
	dstRaw.Spec.DefinitionPath = src.Spec.DefinitionPath
	dstRaw.Spec.Organization = src.Spec.Organization
	dstRaw.Spec.IsDefaultVersion = src.Spec.IsDefaultVersion
	dstRaw.Spec.SystemAPI = src.Spec.SystemAPI
	dstRaw.Spec.APIProperties = src.Spec.APIProperties
	dstRaw.Spec.DefinitionFileRef = src.Spec.DefinitionFileRef

	// Convert Production and Sandbox (rename HTTPRouteRefs to RouteRefs)
	if src.Spec.Production != nil {
		dstRaw.Spec.Production = make([]v1alpha3.EnvConfig, len(src.Spec.Production))
		for i, prod := range src.Spec.Production {
			dstRaw.Spec.Production[i].RouteRefs = prod.HTTPRouteRefs
		}
	}
	if src.Spec.Sandbox != nil {
		dstRaw.Spec.Sandbox = make([]v1alpha3.EnvConfig, len(src.Spec.Sandbox))
		for i, sandbox := range src.Spec.Sandbox {
			dstRaw.Spec.Sandbox[i].RouteRefs = sandbox.HTTPRouteRefs
		}
	}

	// v1alpha3 has an Environment field, which v1alpha4 does not. Leave it as default (empty string).

	// Convert Status
	dstRaw.Status.DeploymentStatus.Accepted = src.Status.DeploymentStatus.Accepted
	dstRaw.Status.DeploymentStatus.Events = src.Status.DeploymentStatus.Events
	dstRaw.Status.DeploymentStatus.Message = src.Status.DeploymentStatus.Message
	dstRaw.Status.DeploymentStatus.Status = src.Status.DeploymentStatus.Status
	dstRaw.Status.DeploymentStatus.TransitionTime = src.Status.DeploymentStatus.TransitionTime

	return nil
}

// ConvertFromV1alpha3 converts from the v1alpha3 version to this API.
func (dst *API) ConvertFromV1alpha3(srcRaw *v1alpha3.API) error {
	// Copy TypeMeta
	dst.TypeMeta = srcRaw.TypeMeta

	// Copy ObjectMeta
	dst.ObjectMeta = srcRaw.ObjectMeta

	// Convert Spec
	dst.Spec.APIName = srcRaw.Spec.APIName
	dst.Spec.APIType = srcRaw.Spec.APIType
	dst.Spec.APIVersion = srcRaw.Spec.APIVersion
	dst.Spec.BasePath = srcRaw.Spec.BasePath
	dst.Spec.DefinitionPath = srcRaw.Spec.DefinitionPath
	dst.Spec.Organization = srcRaw.Spec.Organization
	dst.Spec.IsDefaultVersion = srcRaw.Spec.IsDefaultVersion
	dst.Spec.SystemAPI = srcRaw.Spec.SystemAPI
	dst.Spec.APIProperties = srcRaw.Spec.APIProperties
	dst.Spec.DefinitionFileRef = srcRaw.Spec.DefinitionFileRef

	// Convert Production and Sandbox (rename RouteRefs to HTTPRouteRefs)
	if srcRaw.Spec.Production != nil {
		dst.Spec.Production = make([]RouteRef, len(srcRaw.Spec.Production))
		for i, prod := range srcRaw.Spec.Production {
			dst.Spec.Production[i].HTTPRouteRefs = prod.RouteRefs
		}
	}
	if srcRaw.Spec.Sandbox != nil {
		dst.Spec.Sandbox = make([]RouteRef, len(srcRaw.Spec.Sandbox))
		for i, sandbox := range srcRaw.Spec.Sandbox {
			dst.Spec.Sandbox[i].HTTPRouteRefs = sandbox.RouteRefs
		}
	}

	// Convert Status
	dst.Status.DeploymentStatus.Accepted = srcRaw.Status.DeploymentStatus.Accepted
	dst.Status.DeploymentStatus.Events = srcRaw.Status.DeploymentStatus.Events
	dst.Status.DeploymentStatus.Message = srcRaw.Status.DeploymentStatus.Message
	dst.Status.DeploymentStatus.Status = srcRaw.Status.DeploymentStatus.Status
	dst.Status.DeploymentStatus.TransitionTime = srcRaw.Status.DeploymentStatus.TransitionTime

	return nil
}

// ConvertToV1alpha2 converts this API to the v1alpha2 version.
func (src *API) ConvertToV1alpha2(dstRaw *v1alpha3.API) error {
	// Since v1alpha2 uses the same types as v1alpha3, we can reuse the v1alpha3 conversion logic
	return src.ConvertToV1alpha3(dstRaw)
}

// ConvertFromV1alpha2 converts from the v1alpha2 version to this API.
func (dst *API) ConvertFromV1alpha2(srcRaw *v1alpha3.API) error {
	// Since v1alpha2 uses the same types as v1alpha3, we can reuse the v1alpha3 conversion logic
	return dst.ConvertFromV1alpha3(srcRaw)
}
