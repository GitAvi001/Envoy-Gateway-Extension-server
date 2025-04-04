#!/bin/bash

make clean

#default mainifests command to generate combined CRD
make manifests-combined

# Build and push the extension-server
make build
make image

# upgrade the helm chart
helm upgrade -n envoy-gateway-system extension-server /home/avindu/Desktop/envoy-gateway/examples/extension-server/charts/extension-server \
  --set image.repository=avidocker692/extension-server \
  --set image.tag=latest1 \
  --set imagePullPolicy=Always
  
# Restart deployments
kubectl rollout restart deployment -n envoy-gateway-system