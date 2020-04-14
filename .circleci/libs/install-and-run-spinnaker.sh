#!/bin/bash -xe
# Install Halyard
curl -O https://raw.githubusercontent.com/spinnaker/halyard/master/install/debian/InstallHalyard.sh
USERNAME=`whoami`
sudo bash InstallHalyard.sh --user $USERNAME -y
hal -v

# Create Kind cluster
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
        authorization-mode: "AlwaysAllow"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
kubectl config use-context kind-kind
kubectl cluster-info

# Configure Spinnaker installation
## Generate random password
GATE_PASS=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 32 ; echo '')
## Configure default provider
hal -q config provider kubernetes enable
CONTEXT=$(kubectl config current-context)
hal -q config provider kubernetes account add my-k8s-v2-account --provider-version v2 --context $CONTEXT
hal -q config features edit --artifacts true
hal -q config deploy edit --type distributed --account-name my-k8s-v2-account
## Prepare cluster
kubectl create namespace spinnaker
kubectl -n spinnaker create -f https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/kubernetes/minio-standalone-pvc.yaml
kubectl -n spinnaker create -f https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/kubernetes/minio-standalone-deployment.yaml
curl -O https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/kubernetes/minio-standalone-service.yaml
sed -i 's/LoadBalancer/ClusterIP/g' minio-standalone-service.yaml
kubectl -n spinnaker create -f minio-standalone-service.yaml
mkdir -p ~/.hal/default/profiles
echo "spinnaker.s3.versioning: false" >> ~/.hal/default/profiles/front50-local.yml
## Configure Spinnaker
export MINIO_ACCESS_KEY=minio
export MINIO_SECRET_KEY=minio123
echo $MINIO_SECRET_KEY | hal -q config storage s3 edit --path-style-access=true --endpoint "http://minio-service:9000" --access-key-id $MINIO_ACCESS_KEY --secret-access-key
hal -q config storage edit --type s3
hal -q config features edit --pipeline-templates true
export NEED_SPINNAKER_VERSION=${NEED_SPINNAKER_VERSION:-1.19}
hal -q version list
hal -q config version edit --version $(hal -q version list  | grep "^ -" | awk '{ print $2 }' | grep ${NEED_SPINNAKER_VERSION})
hal config security ui edit --override-base-url='http://spinnaker'
hal config security api edit --override-base-url='http://spinnaker/api/v1'
### Extra parameters for Gate and Deck
echo 'window.spinnakerSettings.authEnabled = true;' > ~/.hal/default/profiles/settings-local.js
echo 'healthEndpoint: /api/v1/health' > ~/.hal/default/service-settings/gate.yml
cat <<EOF > ~/.hal/default/profiles/gate-local.yml
server:
  servlet:
    context-path: /api/v1
  tomcat:
    protocolHeader: X-Forwarded-Proto
    remoteIpHeader: X-Forwarded-For
    internalProxies: .*
    httpsServerPort: X-Forwarded-Port

security:
  basicform:
    enabled: true
  user:
    name: admin
    password: ${GATE_PASS}
EOF

# Install Spinnaker
hal -q deploy apply
until kubectl -n spinnaker wait --for=condition=Ready pod --all > /dev/null
do
	kubectl -n spinnaker get pods
done

# Install Ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/baremetal/service-nodeport.yaml
kubectl patch deployments -n ingress-nginx nginx-ingress-controller -p '{"spec":{"template":{"spec":{"containers":[{"name":"nginx-ingress-controller","ports":[{"containerPort":80,"hostPort":80},{"containerPort":443,"hostPort":443}]}],"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'
cat << EOF > spin-ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: spin
  name: spin-ingress
  namespace: spinnaker
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: spin-deck
          servicePort: 9000
        path: /
      - backend:
          serviceName: spin-gate
          servicePort: 8084
        path: /api/v1
EOF
kubectl -n spinnaker apply -f spin-ingress.yaml

