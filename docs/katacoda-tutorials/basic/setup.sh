echo "Upgrade to Helm 1.16.0"
curl -LO https://get.helm.sh/helm-"v2.16.0"-linux-amd64.tar.gz
tar -xzvf helm-"v2.16.0"-linux-amd64.tar.gz
mv ./linux-amd64/{helm,tiller} /usr/local/bin

echo "Waiting for Minikube to start..."
LIMIT=50
COUNTER=0

while [ ${COUNTER} -lt ${LIMIT} ] && [ -z "$MINIKUBE_STATUS" ]; do
  (( COUNTER++ ))
  echo "Minikube is almost ready, there are $LIMIT checks left, and it is the $COUNTER attempt so far"
  MINIKUBE_STATUS="$(kubectl get pod -n kube-system | grep kube-apiserver-minikube || :)"
  sleep 5
done

# If the apiserver is not available, get Minikube log
if [[ -z "$MINIKUBE_STATUS" ]]; then
  exit 1
fi

echo "Minikube is up and running"

echo "Tiller installation"
kubectl -n kube-system create serviceaccount tiller
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account=tiller

echo "Waiting for Tiller to start..."
LIMIT=30
COUNTER=0


while [ ${COUNTER} -lt ${LIMIT} ] && [ -z "$TILLER_STATUS" ]; do
  (( COUNTER++ ))
  echo "Tiller is almost ready, there are $LIMIT checks left, and it is the $COUNTER attempt so far"
  TILLER_STATUS="$(kubectl get deploy tiller-deploy -n kube-system -o jsonpath='{.status.availableReplicas}' || :)"
  sleep 3
done

# If the apiserver is not available, get Minikube logs
if [[ -z "$TILLER_STATUS" ]]; then
  exit 1
fi

echo "Tiller is up and running"

clear 

echo "Kubernetes with Minikube and proper Helm setup are in place to start the tutorial"
