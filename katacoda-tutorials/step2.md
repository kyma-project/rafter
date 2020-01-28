In this scenario, you will learn how to use Rafter to store static webpages. You will create a bucket, push an asset to it, and open the website from sources stored in the bucket. For the purpose of this scenario, the asset is a package containing all static files needed to build a website, such as HTML, CSS, and JS files.

Follow these steps:

1. Export a URL to ready-to-use sources of a simple website as an environment variable:

   `export GH_WEBPAGE_URL=https://github.com/kyma-project/examples/archive/master.zip`{{execute}}

2. Create a bucket by applying a Bucket custom resource (CR). Run:

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: rafter.kyma-project.io/v1beta1
   kind: Bucket
   metadata:
     name: pages
     namespace: default
   spec:
     region: "us-east-1"
     policy: readonly
   EOF
   ```{{execute}}

3. Create an asset by applying an Asset CR. The Rafter Controller Manager fetches the asset from the location provided in **spec.source.url**. In this example, you can see that the fetched item is a package with a specific directory filtered.

   ```yaml
   cat <<EOF | kubectl apply -f -
   apiVersion: rafter.kyma-project.io/v1beta1
   kind: Asset
   metadata:
     name: webpage
     namespace: default
   spec:
     source:
       url: ${GH_WEBPAGE_URL}
       mode: package
       filter: /rafter/webpage/
     bucketRef:
       name: pages
   EOF
   ```{{execute}}

4. Make sure that the status of the Asset CR is `Ready` which means that fetching, unpacking, and filtering were completed. Run:

   `kubectl get assets webpage -o jsonpath='{.status.phase}'`{{execute}}

5. Export the name of the remote bucket in storage as an environment variable. This name is available in the Bucket CR status and differs from the name of the Bucket CR:

   `export BUCKET_NAME=$(kubectl get bucket pages -o jsonpath='{.status.remoteName}')`{{execute}}

6. Echo the link and open it in a browser to access the website:

   `echo https://[[HOST_SUBDOMAIN]]-31311-[[KATACODA_HOST]].environments.katacoda.com/$BUCKET_NAME/webpage/examples-master/rafter/webpage/index.html`{{execute}}
