Before you proceed with the steps, wait for the environment to finish setting up. Once ready, start playing with Rafter:

1. Add a new chart's repository to Helm. Run:

   `helm repo add rafter-charts https://rafter-charts.storage.googleapis.com`{{execute}}

2. Install Rafter:

   `helm install --name rafter --set rafter-controller-manager.minio.service.type=NodePort rafter-charts/rafter`{{execute}}
