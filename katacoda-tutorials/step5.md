This scenario shows how to use ConfigMaps as an alternative asset source in Rafter. You create a ConfigMap with three files, `file1.md`, `file2.js`, and `file3.yaml`. Then will create a Bucket CR and an Asset CR that points to the previously created ConfigMap. By adding filtering to the Asset CR definition, Rafter will only upload the `.md` file from the ConfigMap content and upload it into the bucket.

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
        Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Mattis rhoncus urna neque viverra justo nec ultrices dui. Venenatis cras sed felis eget velit. Aliquam nulla facilisi cras fermentum odio. Nec ultrices dui sapien eget mi. Auctor elit sed vulputate mi sit amet mauris commodo quis. In pellentesque massa placerat duis ultricies lacus sed. Gravida arcu ac tortor dignissim convallis aenean et. Quisque sagittis purus sit amet. Nibh sit amet commodo nulla facilisi nullam vehicula ipsum. Rhoncus aenean vel elit scelerisque mauris pellentesque pulvinar pellentesque habitant. Auctor elit sed vulputate mi sit. Sed adipiscing diam donec adipiscing tristique risus. Nunc non blandit massa enim. Felis donec et odio pellentesque diam.

        Auctor elit sed vulputate mi sit amet mauris. Amet consectetur adipiscing elit duis tristique. Tellus rutrum tellus pellentesque eu. Nam libero justo laoreet sit amet cursus sit. Sagittis aliquam malesuada bibendum arcu vitae elementum. Amet tellus cras adipiscing enim eu turpis. Auctor urna nunc id cursus metus aliquam eleifend mi. Nec sagittis aliquam malesuada bibendum arcu vitae elementum curabitur. Consectetur lorem donec massa sapien faucibus et molestie ac. Sed risus pretium quam vulputate dignissim suspendisse in. Felis eget nunc lobortis mattis aliquam faucibus.

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

2. Create a bucket by applying the Bucket CR. Run:

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

3. Apply the Asset CR pointing to assets from the ConfigMap that have the `.md` extension. Run:

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

4. Make sure that the status of the Asset CR is `Ready` which means that fetching and filtering was completed. Run:

   `kubectl get assets sample-asset -o jsonpath='{.status.phase}'`{{execute}}

To make sure that the file is in storage and you can extract it, follow these steps:

1. Export the name of the remote bucket in storage as an environment variable. The name of the remote bucket is available in the Bucket CR status and differs from the name of the Bucket CR:

   `export BUCKET_NAME=$(kubectl get bucket sample-bucket -o jsonpath='{.status.remoteName}')`{{execute}}

2. Fetch the file content in the terminal window:

  `curl https://[[HOST_SUBDOMAIN]]-31311-[[KATACODA_HOST]].environments.katacoda.com/$BUCKET_NAME/sample-asset/file1.md`{{execute}}
