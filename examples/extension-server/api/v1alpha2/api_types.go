// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha2

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/envoyproxy/gateway/examples/extension-server/api/v1alpha1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="API Name",type=string,JSONPath=`.spec.apiName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.apiVersion`
// +kubebuilder:printcolumn:name="BasePath",type=string,JSONPath=`.spec.basePath`
// +kubebuilder:printcolumn:name="Organization",type=string,JSONPath=`.spec.organization`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +groupName=dp.wso2.com
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

	// APIType denotes the type of the API. Possible values could be REST, GraphQL
	// +kubebuilder:validation:Enum=REST;GraphQL
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
	// +kubebuilder:validation:MaxItems=1
	// +nullable
	Production []EnvConfig `json:"production,omitempty"`

	// IsDefaultVersion indicates whether this API version should be used as a default API
	IsDefaultVersion bool `json:"isDefaultVersion"`

	// SystemAPI denotes if it is an internal system API.
	SystemAPI bool `json:"systemAPI"`

	// APIProperties denotes the custom properties of the API.
	// +nullable
	APIProperties []Property `json:"apiProperties,omitempty"`

	// DefinitionFileRef contains the definition of the API in a ConfigMap.
	DefinitionFileRef string `json:"definitionFileRef,omitempty"`

	// Sandbox contains a list of references to HttpRoutes of type HttpRoute.
	// +kubebuilder:validation:MaxItems=1
	// +nullable
	Sandbox []EnvConfig `json:"sandbox,omitempty"`

	// Environment denotes the environment of the API.
	// +nullable
	Environment string `json:"environment,omitempty"`
}

// EnvConfig contains the environment specific configuration
type EnvConfig struct {
	// RouteRefs denotes the environment of the API.
	// +kubebuilder:validation:Required
	RouteRefs []string `json:"routeRefs"`
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

	// Convert Production and Sandbox (rename RouteRefs to HTTPRouteRefs)
	if src.Spec.Production != nil {
		dst.Spec.Production = make([]v1alpha1.RouteRef, len(src.Spec.Production))
		for i, prod := range src.Spec.Production {
			dst.Spec.Production[i].HTTPRouteRefs = prod.RouteRefs
		}
	}
	if src.Spec.Sandbox != nil {
		dst.Spec.Sandbox = make([]v1alpha1.RouteRef, len(src.Spec.Sandbox))
		for i, sandbox := range src.Spec.Sandbox {
			dst.Spec.Sandbox[i].HTTPRouteRefs = sandbox.RouteRefs
		}
	}

	// v1alpha2 has an Environment field, which v1alpha1 does not. Ignore it during conversion to hub.

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

	// Convert Production and Sandbox (rename HTTPRouteRefs to RouteRefs)
	if src.Spec.Production != nil {
		dst.Spec.Production = make([]EnvConfig, len(src.Spec.Production))
		for i, prod := range src.Spec.Production {
			dst.Spec.Production[i].RouteRefs = prod.HTTPRouteRefs
		}
	}
	if src.Spec.Sandbox != nil {
		dst.Spec.Sandbox = make([]EnvConfig, len(src.Spec.Sandbox))
		for i, sandbox := range src.Spec.Sandbox {
			dst.Spec.Sandbox[i].RouteRefs = sandbox.HTTPRouteRefs
		}
	}

	// v1alpha1 does not have an Environment field, so leave it as default (empty string).

	// Convert Status
	dst.Status.DeploymentStatus.Accepted = src.Status.DeploymentStatus.Accepted
	dst.Status.DeploymentStatus.Events = src.Status.DeploymentStatus.Events
	dst.Status.DeploymentStatus.Message = src.Status.DeploymentStatus.Message
	dst.Status.DeploymentStatus.Status = src.Status.DeploymentStatus.Status
	dst.Status.DeploymentStatus.TransitionTime = src.Status.DeploymentStatus.TransitionTime

	return nil
}
