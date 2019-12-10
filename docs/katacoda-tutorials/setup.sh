echo "Upgrade to Helm 1.16.0"
curl -LO https://get.helm.sh/helm-"v2.16.0"-linux-amd64.tar.gz
tar -xzvf helm-"v2.16.0"-linux-amd64.tar.gz
mv ./linux-amd64/{helm,tiller} /usr/local/bin

echo "Waiting for Minikube to be up..."
LIMIT=50
COUNTER=0

while [ ${COUNTER} -lt ${LIMIT} ] && [ -z "$MINIKUBE_STATUS" ]; do
  (( COUNTER++ ))
  echo "Minikube is almost ready, $LIMIT checks left and so far it is attempt number $COUNTER"
  MINIKUBE_STATUS="$(kubectl get pod -n kube-system | grep kube-apiserver-minikube || :)"
  sleep 5
done

# In case apiserver is not available get Minikube logs
if [[ -z "$MINIKUBE_STATUS" ]]; then
  exit 1
fi

echo "Minikube is up"

echo "Tiller installation"
kubectl -n kube-system create serviceaccount tiller
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account=tiller

echo "Waiting for Tiller to be up..."
LIMIT=30
COUNTER=0


while [ ${COUNTER} -lt ${LIMIT} ] && [ -z "$TILLER_STATUS" ]; do
  (( COUNTER++ ))
  echo "Tiller is almost ready, $LIMIT checks left and so far it is attempt number $COUNTER"
  TILLER_STATUS="$(kubectl get deploy tiller-deploy -n kube-system -o jsonpath='{.status.availableReplicas}' || :)"
  sleep 3
done

# In case apiserver is not available get Minikube logs
if [[ -z "$TILLER_STATUS" ]]; then
  exit 1
fi

echo "Tiller is up"

clear 

echo "Kubernetes with Minikube and proper Helm setup are in place to start with the tutorial"
