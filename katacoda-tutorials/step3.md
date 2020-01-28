In this scenario, you will create a bucket, push an asset to it, and communicate with a webhook service responsible for extracting metadata from Markdown files. The service will first extract selected metadata from the file specified in the Asset CR and then add it to the Asset CR status.

Follow these steps:

1. Export a URL to a single Markdown file as an environment variable:

   `export MARKDOWN_FILE_URL=https://gist.githubusercontent.com/derberg/01666184bac1ddb4b388c31739924dca/raw/b1d0aff9dcc5f5ee309c33d330b9ba23de470da0/sample-markdown.md`{{execute}}

2. Create a bucket for the Markdown file by applying a Bucket CR. Run:

    ```yaml
     cat <<EOF | kubectl apply -f -
     apiVersion: rafter.kyma-project.io/v1beta1
     kind: Bucket
     metadata:
       name: content
       namespace: default
     spec:
       region: "us-east-1"
       policy: readonly
     EOF
     ```{{execute}}

3. Apply an Asset CR that points to a single Markdown file. To extract metadata from the Markdown file before it goes to storage, specify communication with **metadataWebhookService** in the Asset CR.  As a result, the status of the Asset CR will be enriched with the file metadata. Run:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: rafter.kyma-project.io/v1beta1
   kind: Asset
   metadata:
     name: markdown-file
     namespace: default
   spec:
     source:
       url: ${MARKDOWN_FILE_URL}
       mode: single
       metadataWebhookService:
         - name: rafter-rafter-front-matter-service
           namespace: default
           endpoint: "/v1/extract"
     bucketRef:
       name: content
   EOF
   ```{{execute}}

4. Make sure that the status of the Asset CR is `Ready` which means that fetching and communication with the Front Matter Service were completed. Run:

   `kubectl get assets markdown-file -o jsonpath='{.status.phase}'`{{execute}}

5. Check the status of the Asset CR. Now it has an additional **metadata** object with a set of values. In this scenario, there should be the **title** key with a value.

   `kubectl get assets markdown-file -o jsonpath='{.status.assetRef.files[0].metadata.title}'`{{execute}}

To make sure that the file is in storage and you can extract it, follow these steps:

6. Export the file name and the name of the remote bucket in storage as environment variables. The name of the remote bucket is available in the Bucket CR status and differs from the name of the Bucket CR:

   `export FILE_NAME=$(kubectl get asset markdown-file -o jsonpath='{.status.assetRef.files[0].name}')`{{execute}}

   `export BUCKET_NAME=$(kubectl get bucket content -o jsonpath='{.status.remoteName}')`{{execute}}

7. Fetch the file content in the terminal window:

   `curl https://[[HOST_SUBDOMAIN]]-31311-[[KATACODA_HOST]].environments.katacoda.com/$BUCKET_NAME/markdown-file/$FILE_NAME`{{execute}}
