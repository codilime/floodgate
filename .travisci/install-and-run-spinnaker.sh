#!/bin/bash -e

EXEC_DIR=$(dirname "$0")
HAL_VERSION=${HAL_VERSION:-1.35.0}

# Install packages
sudo apt update
sudo apt install -y jq

# Install Halyard
curl -O https://raw.githubusercontent.com/spinnaker/halyard/master/install/debian/InstallHalyard.sh
USERNAME=`whoami`
sudo bash InstallHalyard.sh --version ${HAL_VERSION} --user $USERNAME -y 
hal -v

# Create Kind cluster
kind create cluster --config="${EXEC_DIR}/kind-cluster-config.yaml"
kubectl config use-context kind-kind
kubectl cluster-info

# Configure Spinnaker installation
## Generate random password
GATE_PASS=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 32 ; echo '')

## Configure default provider
hal -q config provider kubernetes enable
CONTEXT=$(kubectl config current-context)
hal -q config provider kubernetes account add my-k8s-v2-account --provider-version v2 --context $CONTEXT
## Configure account for inner kind communication
cp ~/.kube/config ~/.kube/kind
sed -i "s/server:\ .*/server:\ https:\/\/10.96.0.1:443/g" ~/.kube/kind
hal -q config provider kubernetes account add inner-kind --provider-version v2 --context $CONTEXT --kubeconfig-file ~/.kube/kind
hal -q config deploy edit --type distributed --account-name my-k8s-v2-account

## Install minio
kubectl create namespace spinnaker
sed -i 's/LoadBalancer/ClusterIP/g' "${EXEC_DIR}/minio-standalone.yaml"
kubectl -n spinnaker create -f "${EXEC_DIR}/minio-standalone.yaml"
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
mkdir -p ~/.hal/default/service-settings
echo 'healthEndpoint: /api/v1/health' > ~/.hal/default/service-settings/gate.yml
sed -i "s/GATE_PASS/${GATE_PASS}/g" "${EXEC_DIR}/gate-local.yml"
cp "${EXEC_DIR}/gate-local.yml" ~/.hal/default/profiles/gate-local.yml

# Install Spinnaker
hal -q deploy apply
until kubectl -n spinnaker wait --for=condition=Ready pod --all > /dev/null
do
	kubectl -n spinnaker get pods
done

# Install Ingress controller
kubectl apply -f "${EXEC_DIR}/ingress-mandatory.yaml"
kubectl apply -f "${EXEC_DIR}/ingress-service-nodeport.yaml"
kubectl patch deployments -n ingress-nginx nginx-ingress-controller -p '{"spec":{"template":{"spec":{"containers":[{"name":"nginx-ingress-controller","ports":[{"containerPort":80,"hostPort":80},{"containerPort":443,"hostPort":443}]}],"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'
kubectl -n spinnaker apply -f "${EXEC_DIR}/spinnaker-ingress.yaml"

# Generate Floodgate config file
sed -i "s/GATE_PASS/${GATE_PASS}/g" "${EXEC_DIR}/floodgate-config.yaml"
cp "${EXEC_DIR}/floodgate-config.yaml" ~/floodgate.yaml
