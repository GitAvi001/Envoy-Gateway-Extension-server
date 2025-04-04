// // Copyright Envoy Gateway Authors
// // SPDX-License-Identifier: Apache-2.0
// // The full text of the Apache license is available in the LICENSE file at
// // the root of the repo.

// package v1alpha1

// import (
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// // +kubebuilder:object:root=true
// // +kubebuilder:subresource:status
// // +kubebuilder:printcolumn:name="API Name",type=string,JSONPath=`.spec.apiName`
// // +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.apiVersion`
// // +kubebuilder:printcolumn:name="BasePath",type=string,JSONPath=`.spec.basePath`
// // +kubebuilder:printcolumn:name="Organization",type=string,JSONPath=`.spec.organization`
// // +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// // +groupName=dp.wso2.com
// // +kubebuilder:storageversion
// // API is the Schema for the apis API
// type API struct {
// 	metav1.TypeMeta   `json:",inline"`
// 	metav1.ObjectMeta `json:"metadata,omitempty"`

// 	Spec   APISpec   `json:"spec"`
// 	Status APIStatus `json:"status,omitempty"`
// }

// // APISpec defines the desired state of API
// type APISpec struct {
// 	// APIName is the unique name of the API can be used to uniquely identify an API.
// 	// +kubebuilder:validation:MinLength=1
// 	// +kubebuilder:validation:MaxLength=60
// 	// +kubebuilder:validation:Pattern=`^[^~!@#;:%^*()+={}|\<>"'',&$\[\]\/]*$`
// 	APIName string `json:"apiName"`

// 	// APIType denotes the type of the API. Possible values could be REST, GraphQL, Async
// 	// +kubebuilder:validation:Enum=REST
// 	APIType string `json:"apiType"`

// 	// APIVersion is the version number of the API.
// 	// +kubebuilder:validation:MinLength=1
// 	// +kubebuilder:validation:MaxLength=30
// 	// +kubebuilder:validation:Pattern=`^[^~!@#;:%^*()+={}|\<>"'',&/$\[\]\s+\/]+$`
// 	APIVersion string `json:"apiVersion"`

// 	// BasePath denotes the basepath of the API. e.g: /pet-store-api/1.0.6
// 	// +kubebuilder:validation:Pattern=`^[/][a-zA-Z0-9~/_.-]*$`
// 	BasePath string `json:"basePath"`

// 	// DefinitionPath contains the path to expose the API definition.
// 	// +kubebuilder:default=/api-definition
// 	// +kubebuilder:validation:MinLength=1
// 	DefinitionPath string `json:"definitionPath"`

// 	// Organization denotes the organization related to the API
// 	Organization string `json:"organization"`

// 	// Production contains a list of references to HttpRoutes of type HttpRoute.
// 	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
// 	// +kubebuilder:validation:MaxItems=1
// 	// +nullable
// 	Production []RouteRef `json:"production,omitempty"`

// 	// IsDefaultVersion indicates whether this API version should be used as a default API
// 	IsDefaultVersion bool `json:"isDefaultVersion"`

// 	// SystemAPI denotes if it is an internal system API.
// 	SystemAPI bool `json:"systemAPI"`

// 	// APIProperties denotes the custom properties of the API.
// 	// +nullable
// 	APIProperties []Property `json:"apiProperties,omitempty"`

// 	// DefinitionFileRef contains the definition of the API in a ConfigMap.
// 	DefinitionFileRef string `json:"definitionFileRef,omitempty"`

// 	// Sandbox contains a list of references to HttpRoutes of type HttpRoute.
// 	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
// 	// +kubebuilder:validation:MaxItems=1
// 	// +nullable
// 	Sandbox []RouteRef `json:"sandbox,omitempty"`
// }

// // RouteRef contains the environment specific configuration
// type RouteRef struct {
// 	// HTTPRouteRefs denotes the environment of the API.
// 	// +kubebuilder:validation:Required
// 	HTTPRouteRefs []string `json:"httpRouteRefs"`
// }

// // Property holds key value pair of APIProperties
// type Property struct {
// 	Name  string `json:"name"`
// 	Value string `json:"value"`
// }

// // APIStatus defines the observed state of API
// type APIStatus struct {
// 	// DeploymentStatus denotes the deployment status of the API
// 	DeploymentStatus DeploymentStatus `json:"deploymentStatus,omitempty"`
// }

// // DeploymentStatus denotes the deployment status of the API
// type DeploymentStatus struct {
// 	// Accepted represents whether the API is accepted or not.
// 	// +kubebuilder:validation:Required
// 	Accepted bool `json:"accepted"`

// 	// Events contains a list of events related to the API.
// 	Events []string `json:"events,omitempty"`

// 	// Message represents a user friendly message that explains the current state of the API.
// 	Message string `json:"message,omitempty"`

// 	// Status denotes the state of the API in its lifecycle.
// 	// Possible values could be Accepted, Invalid, Deploy etc.
// 	// +kubebuilder:validation:Required
// 	Status string `json:"status"`

// 	// TransitionTime represents the last known transition timestamp.
// 	// +kubebuilder:validation:Required
// 	// +kubebuilder:validation:Format=date-time
// 	TransitionTime *metav1.Time `json:"transitionTime"`
// }

// // +kubebuilder:object:root=true
// // APIList contains a list of API resources.
// type APIList struct {
// 	metav1.TypeMeta `json:",inline"`
// 	metav1.ListMeta `json:"metadata,omitempty"`
// 	Items           []API `json:"items"`
// }

// // Allowing Kubernetes API server to recognize the resources defined in the API CR yaml
// func init() {
// 	SchemeBuilder.Register(&API{}, &APIList{})
// }

// // Hub marks this version as the hub for conversion
// func (in *API) Hub() {}

// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Hub marks this type as a conversion hub.
func (*API) Hub() {}
