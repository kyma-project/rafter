This scenario shows how to use ConfigMaps as an alternative asset source in Rafter. You will create a ConfigMap with three files, `file1.md`, `file2.js`, and `file3.yaml`. Then you will create a Bucket CR and an Asset CR that points to the previously created ConfigMap. By adding filtering to the Asset CR definition, Rafter will only select the `.md` file from the ConfigMap content and push it into the bucket.

Follow these steps:

1. Create the `sample-configmap` ConfigMap that contains sources of files with three different extensions:

    ```yaml
    cat <<EOF | kubectl apply -f -
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: sample-configmap
      namespace: default
    data:
      file1.md: |
        Rafter is a solution for storing and managing different types of files called assets. It uses MinIO as object storage. The whole concept of Rafter relies on Kubernetes custom resources (CRs) managed by the Rafter Controller Manager. These CRs include:

          - Asset CR which manages single assets or asset packages from URLs or ConfigMaps
          - Bucket CR which manages buckets
          - AssetGroup CR which manages a group of Asset CRs of a specific type to make it easier to use and extract webhook information

      file2.js: |
        const http = require('http');
        const server = http.createServer();

        server.on('request', (request, response) => {
            let body = [];
            request.on('data', (chunk) => {
                body.push(chunk);
            }).on('end', () => {
                body = Buffer.concat(body).toString();

                console.log('> Headers');
                console.log(request.headers);

                console.log('> Body');
                console.log(body);
                response.end();
            });
        }).listen(8083);

      file3.yaml: |
        apiVersion: apiextensions.k8s.io/v1beta1
        kind: CustomResourceDefinition
        metadata:
          annotations:
            controller-gen.kubebuilder.io/version: v0.2.4
          creationTimestamp: null
          name: assetgroups.rafter.kyma-project.io
        spec:
          additionalPrinterColumns:
          - JSONPath: .status.phase
            name: Phase
            type: string
          - JSONPath: .metadata.creationTimestamp
    EOF
    ```{{execute}}

2. Create a bucket by applying the `sample-bucket` Bucket CR. Run:

    ```yaml
    cat <<EOF | kubectl apply -f -
    apiVersion: rafter.kyma-project.io/v1beta1
    kind: Bucket
    metadata:
      name: sample-bucket
      namespace: default
    spec:
      region: "us-east-1"
      policy: readonly
    EOF
    ```{{execute}}

3. Apply the `sample-asset` Asset CR that selects from the ConfigMap all assets with the `.md` extension. The **url** parameter specifies the Namespace and ConfigMap names and has the `{namespace}/{configMap-name}` format. Run:

    ```yaml
    cat <<EOF | kubectl apply -f -
    apiVersion: rafter.kyma-project.io/v1beta1
    kind: Asset
    metadata:
      name: sample-asset
      namespace: default
    spec:
      source:
        url: default/sample-configmap
        mode: configmap
        filter: \.md$
      bucketRef:
        name: sample-bucket
    EOF
    ```{{execute}}

4. Make sure that the status of the Asset CR is `Ready` which means that fetching and filtering were completed. Run:

   `kubectl get assets sample-asset -o jsonpath='{.status.phase}'`{{execute}}

To make sure that the file is in storage and you can extract it, follow these steps:

1. Export the name of the remote bucket in storage as an environment variable. The name of the remote bucket is available in the Bucket CR status and differs from the name of the Bucket CR:

   `export BUCKET_NAME=$(kubectl get bucket sample-bucket -o jsonpath='{.status.remoteName}')`{{execute}}

2. Fetch the file content in the terminal window:

  `curl https://[[HOST_SUBDOMAIN]]-31311-[[KATACODA_HOST]].environments.katacoda.com/$BUCKET_NAME/sample-asset/file1.md`{{execute}}
