In this scenario, you will create a Bucket and then push an Asset to it. The Asset is a package containing all static files needed to build a website, such as HTML, JS and CSS files. Follow these steps:

1. Export a URL to ready-to-use sources of a simple website as an environment variable:

   `export GH_WEBPAGE_URL=https://github.com/kyma-project/examples/archive/master.zip`{{execute}}

2. Create a Bucket by applying a Bucket custom resource (CR). The Bucket will contain the website files. Run:

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

3. Create an Asset by applying an Asset CR. The Rafter Controller Manager fetches the Asset from the location provided in **spec.source.url**. In this example, you can see that the fetched item is a package with a specific directory filtered.

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
       filter: /asset-store/webpage/
     bucketRef:
       name: pages
   EOF
   ```{{execute}}

4. Make sure that the Asset status is `Ready` which means that fetching, unpacking, and filtering is completed. Run:

   `kubectl get assets webpage -o jsonpath='{.status.phase}'`{{execute}}

5. Export the name of the Bucket available in the Bucket status as an environment variable. The name of the Bucket in storage is not exactly the same as the name of a specific Bucket CR:

   `export BUCKET_NAME=$(kubectl get bucket pages -o jsonpath='{.status.remoteName}')`{{execute}}

6. Echo the link and open it in a browser to access the website:

   `echo https://[[HOST_SUBDOMAIN]]-31311-[[KATACODA_HOST]].environments.katacoda.com/$BUCKET_NAME/webpage/examples-master/asset-store/webpage/index.html`{{execute}}
