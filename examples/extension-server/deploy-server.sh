#!/bin/bash

#1. make clean command to clean up the previous build
make clean

#default manifests command to generate combined CRD
make manifests-combined

# Build and push the extension-server to the docker hub
make build
make image

# upgrade the helm chart with latest iamge
helm upgrade -n envoy-gateway-system extension-server /home/avindu/Desktop/envoy-gateway/examples/extension-server/charts/extension-server \
  --set image.repository=avidocker692/extension-server \
  --set image.tag=latest1 \
  --set imagePullPolicy=Always
  
# Restart deployments
kubectl rollout restart deployment -n envoy-gateway-system