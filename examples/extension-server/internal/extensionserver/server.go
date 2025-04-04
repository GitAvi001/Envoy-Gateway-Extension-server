/*Implementation - basePath both monitors and receives the response(v1alpha1 CRD)*/

// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package extensionserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	pb "github.com/envoyproxy/gateway/proto/extension"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	"github.com/envoyproxy/gateway/examples/extension-server/api/v1alpha1"
)

type Server struct {
	pb.UnimplementedEnvoyGatewayExtensionServer

	log *slog.Logger
}

func New(logger *slog.Logger) *Server {
	logger.Info("extension-server starting up")
	return &Server{
		log: logger,
	}
}

func (s *Server) PostRouteModify(ctx context.Context, req *pb.PostRouteModifyRequest) (*pb.PostRouteModifyResponse, error) {
	s.log.Info("postRouteModify callback was invoked")

	// Log the incoming route for debugging
	s.log.Info("Incoming route", slog.Any("route", req.Route))

	basePath := ""
	httpRouteName := ""

	// Extract basePath and HTTPRoute name from the API CR
	for _, ext := range req.PostRouteContext.ExtensionResources {
		var api v1alpha1.API
		s.log.Info("extension resource", slog.String("extension", string(ext.GetUnstructuredBytes())))
		if err := json.Unmarshal(ext.GetUnstructuredBytes(), &api); err != nil {
			s.log.Error("failed to unmarshal the extension", slog.String("error", err.Error()))
			continue
		}
		// Extract basePath and HTTPRoute name from the API CR
		basePath = api.Spec.BasePath
		if len(api.Spec.Production) > 0 && len(api.Spec.Production[0].HTTPRouteRefs) > 0 {
			httpRouteName = api.Spec.Production[0].HTTPRouteRefs[0]
		}
		s.log.Info(fmt.Sprintf("Extracted basePath: %s and HTTPRoute name: %s from API: %s", basePath, httpRouteName, api.ObjectMeta.Name))
		break // Assuming we only need the first matching API CR
	}

	if basePath == "" || httpRouteName == "" {
		s.log.Info("Missing basePath or HTTPRoute name, skipping modification")
		return &pb.PostRouteModifyResponse{
			Route: req.Route,
		}, nil
	}

	r := req.Route
	if r.Match != nil {
		switch r.Match.PathSpecifier.(type) {
		case *routev3.RouteMatch_Prefix:
			s.log.Info("Path specifier is a prefix")
			originalPrefix := r.Match.PathSpecifier.(*routev3.RouteMatch_Prefix).Prefix
			// Construct the new prefix with basePath and HTTPRoute name
			newPrefix := fmt.Sprintf("%s/%s%s", basePath, httpRouteName, originalPrefix)
			r.Match.PathSpecifier = &routev3.RouteMatch_Prefix{
				Prefix: newPrefix,
			}
			// Add rewrite to strip "/<basePath>/<http-route-name>" and send the original prefix to the backend
			if routeAction, ok := r.Action.(*routev3.Route_Route); ok {
				routeAction.Route.PrefixRewrite = originalPrefix
			} else {
				r.Action = &routev3.Route_Route{
					Route: &routev3.RouteAction{
						PrefixRewrite: originalPrefix, // Strips "/my-api/apk-http-route" to "/"
					},
				}
			}
			s.log.Info(fmt.Sprintf("Matched HTTPRoute %s, applying prefix %s with rewrite to %s", httpRouteName, newPrefix, originalPrefix))
		default:
			s.log.Info("Path specifier is not handled", slog.Any("pathSpecifier", r.Match.PathSpecifier))
		}
	} else {
		s.log.Info("Route match is nil, skipping modification")
	}

	// Log the modified route for debugging
	s.log.Info("Modified route", slog.Any("route", r))

	return &pb.PostRouteModifyResponse{
		Route: r,
	}, nil
}

func (s *Server) PostVirtualHostModify(ctx context.Context, req *pb.PostVirtualHostModifyRequest) (*pb.PostVirtualHostModifyResponse, error) {
	s.log.Info("PostVirtualHostModify callback was invoked")
	// Log the virtual host for debugging
	s.log.Info("VirtualHost", slog.Any("virtualHost", req.VirtualHost))
	return &pb.PostVirtualHostModifyResponse{
		VirtualHost: req.VirtualHost,
	}, nil
}

func (s *Server) PostTranslateModify(ctx context.Context, req *pb.PostTranslateModifyRequest) (*pb.PostTranslateModifyResponse, error) {
	s.log.Info("PostTranslateModify callback was invoked")
	return &pb.PostTranslateModifyResponse{
		Clusters: req.Clusters,
		Secrets:  req.Secrets,
	}, nil
}

func (s *Server) PostHTTPListenerModify(ctx context.Context, req *pb.PostHTTPListenerModifyRequest) (*pb.PostHTTPListenerModifyResponse, error) {
	s.log.Info("postHTTPListenerModify callback was invoked")
	// Log the listener for debugging
	s.log.Info("Listener", slog.Any("listener", req.Listener))
	return &pb.PostHTTPListenerModifyResponse{
		Listener: req.Listener,
	}, nil
}

/*extract basePath and receive the response (combined CRD apply)*/
// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// package extensionserver

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"

// 	pb "github.com/envoyproxy/gateway/proto/extension"
// 	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

// 	"github.com/envoyproxy/gateway/examples/extension-server/api/v1alpha1"
// )

// type Server struct {
// 	pb.UnimplementedEnvoyGatewayExtensionServer

// 	log *slog.Logger
// }

// func New(logger *slog.Logger) *Server {
// 	logger.Info("extension-server starting up")
// 	return &Server{
// 		log: logger,
// 	}
// }

// func (s *Server) PostRouteModify(ctx context.Context, req *pb.PostRouteModifyRequest) (*pb.PostRouteModifyResponse, error) {
// 	s.log.Info("postRouteModify callback was invoked")

// 	// Log the entire request for debugging
// 	s.log.Info("PostRouteModify request", slog.Any("request", req))

// 	// Log the extension resources
// 	s.log.Info("Extension resources", slog.Any("resources", req.PostRouteContext.ExtensionResources))

// 	basePath := ""
// 	httpRouteName := ""

// 	// Extract basePath and HTTPRoute name from the API CR
// 	for i, ext := range req.PostRouteContext.ExtensionResources {
// 		var api v1alpha1.API
// 		s.log.Info("Processing extension resource", slog.Int("index", i), slog.String("extension", string(ext.GetUnstructuredBytes())))
// 		if err := json.Unmarshal(ext.GetUnstructuredBytes(), &api); err != nil {
// 			s.log.Error("failed to unmarshal the extension", slog.String("error", err.Error()))
// 			continue
// 		}

// 		// Log the unmarshaled API object
// 		s.log.Info("Unmarshaled API", slog.Any("api", api))

// 		// Log the Production field specifically
// 		s.log.Info("API Production field", slog.Any("production", api.Spec.Production))

// 		// Extract basePath and HTTPRoute name from the API CR
// 		basePath = api.Spec.BasePath
// 		if len(api.Spec.Production) > 0 && len(api.Spec.Production[0].HTTPRouteRefs) > 0 {
// 			httpRouteName = api.Spec.Production[0].HTTPRouteRefs[0]
// 		} else {
// 			s.log.Warn("No HTTPRouteRefs found in Production field")
// 		}
// 		s.log.Info(fmt.Sprintf("Extracted basePath: %s and HTTPRoute name: %s from API: %s", basePath, httpRouteName, api.ObjectMeta.Name))
// 		break // Assuming we only need the first matching API CR
// 	}

// 	if basePath == "" || httpRouteName == "" {
// 		s.log.Info("Missing basePath or HTTPRoute name, skipping modification")
// 		return &pb.PostRouteModifyResponse{
// 			Route: req.Route,
// 		}, nil
// 	}

// 	r := req.Route
// 	if r.Match != nil {
// 		switch r.Match.PathSpecifier.(type) {
// 		case *routev3.RouteMatch_Prefix:
// 			s.log.Info("Path specifier is a prefix")
// 			originalPrefix := r.Match.PathSpecifier.(*routev3.RouteMatch_Prefix).Prefix
// 			// Construct the new prefix with basePath and HTTPRoute name
// 			newPrefix := fmt.Sprintf("%s/%s%s", basePath, httpRouteName, originalPrefix)
// 			r.Match.PathSpecifier = &routev3.RouteMatch_Prefix{
// 				Prefix: newPrefix,
// 			}
// 			// Add rewrite to strip "/<basePath>/<http-route-name>" and send the original prefix to the backend
// 			if routeAction, ok := r.Action.(*routev3.Route_Route); ok {
// 				routeAction.Route.PrefixRewrite = originalPrefix
// 			} else {
// 				r.Action = &routev3.Route_Route{
// 					Route: &routev3.RouteAction{
// 						PrefixRewrite: originalPrefix, // Strips "/my-api/apk-http-route" to "/"
// 					},
// 				}
// 			}
// 			s.log.Info(fmt.Sprintf("Matched HTTPRoute %s, applying prefix %s with rewrite to %s", httpRouteName, newPrefix, originalPrefix))
// 		default:
// 			s.log.Info("Path specifier is not handled", slog.Any("pathSpecifier", r.Match.PathSpecifier))
// 		}
// 	} else {
// 		s.log.Info("Route match is nil, skipping modification")
// 	}

// 	// Log the modified route
// 	s.log.Info("Modified route", slog.Any("route", r))

// 	return &pb.PostRouteModifyResponse{
// 		Route: r,
// 	}, nil
// }

// func (s *Server) PostVirtualHostModify(ctx context.Context, req *pb.PostVirtualHostModifyRequest) (*pb.PostVirtualHostModifyResponse, error) {
// 	s.log.Info("PostVirtualHostModify callback was invoked")
// 	// Log the virtual host for debugging
// 	s.log.Info("VirtualHost", slog.Any("virtualHost", req.VirtualHost))
// 	return &pb.PostVirtualHostModifyResponse{
// 		VirtualHost: req.VirtualHost,
// 	}, nil
// }

// func (s *Server) PostTranslateModify(ctx context.Context, req *pb.PostTranslateModifyRequest) (*pb.PostTranslateModifyResponse, error) {
// 	s.log.Info("PostTranslateModify callback was invoked")
// 	return &pb.PostTranslateModifyResponse{
// 		Clusters: req.Clusters,
// 		Secrets:  req.Secrets,
// 	}, nil
// }

// func (s *Server) PostHTTPListenerModify(ctx context.Context, req *pb.PostHTTPListenerModifyRequest) (*pb.PostHTTPListenerModifyResponse, error) {
// 	s.log.Info("postHTTPListenerModify callback was invoked")
// 	// Log the listener for debugging
// 	s.log.Info("Listener", slog.Any("listener", req.Listener))
// 	return &pb.PostHTTPListenerModifyResponse{
// 		Listener: req.Listener,
// 	}, nil
// }

/*Implementation - extract context path and receives the response (example.extensions.io_apis cr apply) */

// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// package extensionserver

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"

// 	pb "github.com/envoyproxy/gateway/proto/extension"
// 	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
// 	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
// 	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
// 	bav3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/basic_auth/v3"
// 	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"

// 	// v32 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
// 	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
// 	"google.golang.org/protobuf/types/known/anypb"

// 	"github.com/exampleorg/envoygateway-extension/api/v1alpha1"
// )

// type Server struct {
// 	pb.UnimplementedEnvoyGatewayExtensionServer

// 	log *slog.Logger
// }

// func New(logger *slog.Logger) *Server {
// 	return &Server{
// 		log: logger,
// 	}
// }

// func (s *Server) PostRouteModify(ctx context.Context, req *pb.PostRouteModifyRequest) (*pb.PostRouteModifyResponse, error) {
// 	s.log.Info("postRouteModify callback was invoked")
// 	contextPath := ""
// 	httpRouteName := ""

// 	// Extract context and HTTPRoute name from the API CR
// 	for _, ext := range req.PostRouteContext.ExtensionResources {
// 		var api v1alpha1.API
// 		s.log.Info("extension resource", slog.String("extension", string(ext.GetUnstructuredBytes())))
// 		if err := json.Unmarshal(ext.GetUnstructuredBytes(), &api); err != nil {
// 			s.log.Error("failed to unmarshal the extension", slog.String("error", err.Error()))
// 			continue
// 		} else {
// 			contextPath = api.Spec.Context
// 			// Find the HTTPRoute in targetRefs
// 			for _, ref := range api.Spec.TargetRefs {
// 				if ref.Kind == "HTTPRoute" && ref.Group == "gateway.networking.k8s.io" {
// 					httpRouteName = string(ref.Name)
// 					break
// 				}
// 			}
// 			s.log.Info(fmt.Sprintf("Extracted context path: %s and HTTPRoute name: %s from API: %s", contextPath, httpRouteName, api.Name))
// 			break
// 		}
// 	}

// 	if contextPath == "" || httpRouteName == "" {
// 		s.log.Info("Missing context or HTTPRoute name, skipping modification")
// 		return &pb.PostRouteModifyResponse{
// 			Route: req.Route,
// 		}, nil
// 	}

// 	r := req.Route
// 	if r.Match != nil {
// 		switch r.Match.PathSpecifier.(type) {
// 		case *routev3.RouteMatch_Prefix:
// 			s.log.Info("Path specifier is a prefix")
// 			originalPrefix := r.Match.PathSpecifier.(*routev3.RouteMatch_Prefix).Prefix
// 			// Construct the new prefix with context and HTTPRoute name
// 			newPrefix := fmt.Sprintf("/%s/%s%s", contextPath, httpRouteName, originalPrefix)
// 			r.Match.PathSpecifier = &routev3.RouteMatch_Prefix{
// 				Prefix: newPrefix,
// 			}
// 			// Add rewrite to strip "/<context>/<http-route-name>"
// 			if routeAction, ok := r.Action.(*routev3.Route_Route); ok {
// 				routeAction.Route.PrefixRewrite = originalPrefix
// 			} else {
// 				r.Action = &routev3.Route_Route{
// 					Route: &routev3.RouteAction{
// 						PrefixRewrite: originalPrefix, // Strips "/example-context-1/<http-route-name>" to "/"
// 					},
// 				}
// 			}
// 		default:
// 			s.log.Info("Path specifier is not handled")
// 		}
// 	}
// 	return &pb.PostRouteModifyResponse{
// 		Route: r,
// 	}, nil
// }
// func (s *Server) PostVirtualHostModify(ctx context.Context, req *pb.PostVirtualHostModifyRequest) (*pb.PostVirtualHostModifyResponse, error) {
// 	s.log.Info("PostVirtualHostModify callback was invoked")

// 	return &pb.PostVirtualHostModifyResponse{
// 		VirtualHost: req.VirtualHost,
// 	}, nil
// }

// func (s *Server) PostTranslateModify(ctx context.Context, req *pb.PostTranslateModifyRequest) (*pb.PostTranslateModifyResponse, error) {
// 	s.log.Info("PostVirtualHostModify callback was invoked")

// 	return &pb.PostTranslateModifyResponse{
// 		Clusters: req.Clusters,
// 		Secrets:  req.Secrets,
// 	}, nil

// }

// // PostHTTPListenerModify is called after Envoy Gateway is done generating a
// // Listener xDS configuration and before that configuration is passed on to
// // Envoy Proxy.
// // This example adds Basic Authentication on the Listener level as an example.
// // Note: This implementation is not secure, and should not be used to protect
// // anything important.
// func (s *Server) PostHTTPListenerModify(ctx context.Context, req *pb.PostHTTPListenerModifyRequest) (*pb.PostHTTPListenerModifyResponse, error) {
// 	s.log.Info("postHTTPListenerModify callback was invoked")
// 	// Collect all of the required username/password combinations from the
// 	// provided contexts that were attached to the gateway.
// 	passwords := NewHtpasswd()
// 	flag := false
// 	for _, ext := range req.PostListenerContext.ExtensionResources {
// 		var listenerContext v1alpha1.ListenerContextExample
// 		s.log.Info("extension resource", slog.String("extension", string(ext.GetUnstructuredBytes())))
// 		if err := json.Unmarshal(ext.GetUnstructuredBytes(), &listenerContext); err != nil {
// 			s.log.Error("failed to unmarshal the extension", slog.String("error", err.Error()))
// 			continue
// 		}
// 		s.log.Info("processing an extension context", slog.String("username", listenerContext.Spec.Username))
// 		passwords.AddUser(listenerContext.Spec.Username, listenerContext.Spec.Password)
// 		flag = true
// 	}
// 	if !flag {
// 		return &pb.PostHTTPListenerModifyResponse{
// 			Listener: req.Listener,
// 		}, nil
// 	}

// 	// First, get the filter chains from the listener
// 	filterChains := req.Listener.GetFilterChains()
// 	defaultFC := req.Listener.DefaultFilterChain
// 	if defaultFC != nil {
// 		filterChains = append(filterChains, defaultFC)
// 	}
// 	// Go over all of the chains, and add the basic authentication http filter
// 	for _, currChain := range filterChains {
// 		httpConManager, hcmIndex, err := findHCM(currChain)
// 		if err != nil {
// 			s.log.Error("failed to find an HCM in the current chain", slog.Any("error", err))
// 			continue
// 		}
// 		// If a basic authentication filter already exists, update it. Otherwise, create it.
// 		basicAuth, baIndex, err := findBasicAuthFilter(httpConManager.HttpFilters)
// 		if err != nil {
// 			s.log.Error("failed to unmarshal the existing basicAuth filter", slog.Any("error", err))
// 			continue
// 		}
// 		if baIndex == -1 {
// 			// Create a new basic auth filter
// 			basicAuth = &bav3.BasicAuth{
// 				Users: &corev3.DataSource{
// 					Specifier: &corev3.DataSource_InlineString{
// 						InlineString: passwords.String(),
// 					},
// 				},
// 				ForwardUsernameHeader: "X-Example-Ext",
// 			}
// 		} else {
// 			// Update the basic auth filter
// 			basicAuth.Users.Specifier = &corev3.DataSource_InlineString{
// 				InlineString: passwords.String(),
// 			}
// 		}
// 		// Add or update the Basic Authentication filter in the HCM
// 		anyBAFilter, _ := anypb.New(basicAuth)
// 		if baIndex > -1 {
// 			httpConManager.HttpFilters[baIndex].ConfigType = &hcm.HttpFilter_TypedConfig{
// 				TypedConfig: anyBAFilter,
// 			}
// 		} else {
// 			filters := []*hcm.HttpFilter{
// 				{
// 					Name: "envoy.filters.http.basic_auth",
// 					ConfigType: &hcm.HttpFilter_TypedConfig{
// 						TypedConfig: anyBAFilter,
// 					},
// 				},
// 			}
// 			filters = append(filters, httpConManager.HttpFilters...)
// 			httpConManager.HttpFilters = filters
// 		}

// 		// Write the updated HCM back to the filter chain
// 		anyConnectionMgr, _ := anypb.New(httpConManager)
// 		currChain.Filters[hcmIndex].ConfigType = &listenerv3.Filter_TypedConfig{
// 			TypedConfig: anyConnectionMgr,
// 		}
// 	}

// 	return &pb.PostHTTPListenerModifyResponse{
// 		Listener: req.Listener,
// 	}, nil
// }

// // Tries to find an HTTP connection manager in the provided filter chain.
// func findHCM(filterChain *listenerv3.FilterChain) (*hcm.HttpConnectionManager, int, error) {
// 	for filterIndex, filter := range filterChain.Filters {
// 		if filter.Name == wellknown.HTTPConnectionManager {
// 			hcm := new(hcm.HttpConnectionManager)
// 			if err := filter.GetTypedConfig().UnmarshalTo(hcm); err != nil {
// 				return nil, -1, err
// 			}
// 			return hcm, filterIndex, nil
// 		}
// 	}
// 	return nil, -1, fmt.Errorf("unable to find HTTPConnectionManager in FilterChain: %s", filterChain.Name)
// }

// // Tries to find the Basic Authentication HTTP filter in the provided chain
// func findBasicAuthFilter(chain []*hcm.HttpFilter) (*bav3.BasicAuth, int, error) {
// 	for i, filter := range chain {
// 		if filter.Name == "envoy.filters.http.basic_auth" {
// 			ba := new(bav3.BasicAuth)
// 			if err := filter.GetTypedConfig().UnmarshalTo(ba); err != nil {
// 				return nil, -1, err
// 			}
// 			return ba, i, nil
// 		}
// 	}
// 	return nil, -1, nil
// }
