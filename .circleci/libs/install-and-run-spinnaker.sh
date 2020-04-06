#!/bin/bash -xe
curl -O https://raw.githubusercontent.com/spinnaker/halyard/master/install/debian/InstallHalyard.sh
USERNAME=`whoami`
sudo bash InstallHalyard.sh --user $USERNAME -y
hal -v
kind create cluster
kubectl config use-context kind-kind
kubectl cluster-info
hal config provider kubernetes enable
CONTEXT=$(kubectl config current-context)
hal config provider kubernetes account add my-k8s-v2-account --provider-version v2 --context $CONTEXT
hal config features edit --artifacts true
hal config deploy edit --type distributed --account-name my-k8s-v2-account
kubectl create namespace spinnaker
kubectl -n spinnaker create -f https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/kubernetes/minio-standalone-pvc.yaml
kubectl -n spinnaker create -f https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/kubernetes/minio-standalone-deployment.yaml
curl -O https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/kubernetes/minio-standalone-service.yaml
sed -i 's/LoadBalancer/ClusterIP/g' minio-standalone-service.yaml
kubectl -n spinnaker create -f minio-standalone-service.yaml
mkdir -p ~/.hal/default/profiles
echo "spinnaker.s3.versioning: false" >> ~/.hal/default/profiles/front50-local.yml
export MINIO_ACCESS_KEY=minio
export MINIO_SECRET_KEY=minio123
echo $MINIO_SECRET_KEY | hal config storage s3 edit --path-style-access=true --endpoint "http://minio-service:9000" --access-key-id $MINIO_ACCESS_KEY --secret-access-key
hal config storage edit --type s3
hal version list
hal config version edit --version 1.18.8
hal deploy apply
until kubectl -n spinnaker wait --for=condition=Ready pod --all > /dev/null
do
	kubectl -n spinnaker get pods
done
