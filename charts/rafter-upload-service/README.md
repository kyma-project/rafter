# Upload Service

This project contains the Helm chart for the Upload Service.

## Prerequisites

- Kubernetes v1.14 or higher
- Helm v2.10 or higher
- The `rafter-charts` repository added to your Helm client with this command:

```bash
helm repo add rafter-charts https://kyma-project.github.io/rafter
```

## Details

Read how to install, uninstall, and configure the chart.

### Install the chart

Use this command to install the chart:

``` bash
helm install rafter-charts/rafter-upload-service
```

To install the chart with the `rafter-upload-service` release name, use:

``` bash
helm install --name rafter-upload-service rafter-charts/rafter-upload-service
```

The command deploys the Upload Service on the Kubernetes cluster with the default configuration. The [**Configuration**](#configuration) section lists the parameters that you can configure during installation.

> **TIP:** To list all releases, use `helm list`.

### Uninstall the chart

To uninstall the `rafter-upload-service` release, run:

``` bash
helm delete rafter-upload-service
```

That command removes all the Kubernetes components associated with the chart and deletes the release.

### Configuration

The following table lists the configurable parameters of the Upload Service chart and their default values.

| Parameter | Description | Default |
| --- | ---| ---|
| **image.repository** | Upload Service image repository | `eu.gcr.io/kyma-project/rafter-upload-service` |
| **image.tag** | Upload Service image tag | `{TAG_NAME}` |
| **image.pullPolicy** | Pull policy for the Upload Service image | `IfNotPresent` |
| **nameOverride** | String that partially overrides the **rafterUploadService.name** template | `nil` |
| **fullnameOverride** | String that fully overrides the **rafterUploadService.fullname** template | `nil` |
| **minio.enabled** | Parameter that defines whether to deploy MinIO | `true` |
| **minio.persistence.enabled** | Use persistent volume to store data in MinIO. | `true` |
| **minio.persistence.size** | Size of persistent volume claim for MinIO. | `10Gi` |
| **deployment.labels** | Custom labels for the Deployment | `{}` |
| **deployment.annotations** | Custom annotations for the Deployment | `{}` |
| **deployment.replicas** | Number of Upload Service nodes | `1` |
| **deployment.extraProperties** | Additional properties injected in the Deployment | `{}` |
| **pod.labels** | Custom labels for the Pod | `{}` |
| **pod.annotations** | Custom annotations for the Pod | `{}` |
| **pod.extraProperties** | Additional properties injected in the Pod | `{}` |
| **pod.extraContainerProperties** | Additional properties injected in the container | `{}` |
| **service.name** | Service name. If not set, it is generated using the **rafterUploadService.fullname** template. | `nil` |
| **service.type** | Service type | `ClusterIP` |
| **service.port.name** |  Name of the Service port | `http` |
| **service.port.internal** | Internal port of the Service in the Pod | `3000` |
| **service.port.external** | Port on which the Service is exposed in Kubernetes | `80` |
| **service.port.protocol** | Protocol of the Service port | `TCP` |
| **service.labels** | Custom labels for the Service | `{}` |
| **service.annotations** | Custom annotations for the Service | `{}` |
| **serviceAccount.create** | Parameter that defines whether to create a new ServiceAccount for the Upload Service | `true` |
| **serviceAccount.name** | ServiceAccount resource that the Upload Service uses. If not set and the **serviceAccount.create** parameter is set to `true`, the name is generated using the **rafterUploadService.fullname** template. If not set and **serviceAccount.create** is set to `false`, the name is set to `default`. | `nil` |
| **serviceAccount.labels** | Custom labels for the ServiceAccount | `{}` |
| **serviceAccount.annotations** | Custom annotations for the ServiceAccount | `{}` |
| **rbac.clusterScope.create** | Parameter that defines whether to create a new ClusterRole and ClusterRoleBinding for the Upload Service | `true` |
| **rbac.clusterScope.role.name** | ClusterRole resource that the Upload Service uses. If not set and the **rbac.clusterScope.create** parameter is set to `true`, the name is generated using the **rafterUploadService.fullname** template. If not set and **rbac.clusterScope.create** is set to `false`, the name is set to `default`. | `nil` |
| **rbac.clusterScope.role.labels** | Custom labels for the ClusterRole | `{}` |
| **rbac.clusterScope.role.annotations** | Custom annotations for the ClusterRole | `{}` |
| **rbac.clusterScope.role.extraRules** | Additional rules injected in the ClusterRole | `[]` |
| **rbac.clusterScope.roleBinding.name** | ClusterRoleBinding resource that the Upload Service uses. If not set and the **rbac.clusterScope.create** parameter is set to `true`, the name is generated using the **rafterUploadService.fullname** template. If not set and **rbac.clusterScope.create** is set to `false`, the name is set to `default`. | `nil` |
| **rbac.clusterScope.roleBinding.labels** | Custom labels for the ClusterRoleBinding | `{}` |
| **rbac.clusterScope.roleBinding.annotations** | Custom annotations for the ClusterRoleBinding | `{}` |
| **serviceMonitor.create** | Parameter that defines whether to create a new ServiceMonitor custom resource for the Prometheus Operator | `false` |
| **serviceMonitor.name** | ServiceMonitor resource that the Prometheus Operator uses. If not set and the **serviceMonitor.create** parameter is set to `true`, the name is generated using the **rafterUploadService.fullname** template. If not set and **serviceMonitor.create** is set to `false`, the name is set to `default`. | `nil` |
| **serviceMonitor.scrapeInterval** | Scrape interval for the ServiceMonitor custom resource | `30s` |
| **serviceMonitor.labels** | Custom labels for the ServiceMonitor custom resource | `{}` |
| **serviceMonitor.annotations** | Custom annotations for the ServiceMonitor custom resource | `{}` |
| **envs.verbose** | Parameter that defines if logs from the Upload Service should be visible | `false` |
| **envs.kubeconfigPath** | Path to the `kubeconfig` file needed to run the Upload Service outside of a cluster | `nil` |
| **envs.upload.timeout** | File processing time-out | `30m` |
| **envs.upload.workers** | Maximum number of concurrent metadata extraction workers | `10` |
| **envs.upload.endpoint** | Address of the content storage server | `{{ .Release.Name }}-minio.{{ .Release.Namespace }}.svc.cluster.local` |
| **envs.upload.externalEndpoint** | External address of the content storage server | `http://{{ .Release.Name }}-minio.{{ .Release.Namespace }}.svc.cluster.local:9000` |
| **envs.upload.port** | Port on which the content storage server listens | `9000` |
| **envs.upload.accessKey** | Access key required to sign in to the content storage server | Value from `{{ .Release.Name }}-minio` ConfigMap |
| **envs.upload.secretKey** | Secret key required to sign in to the content storage server | Value from `{{ .Release.Name }}-minio` ConfigMap |
| **envs.upload.secure** | HTTPS connection with the content storage server | `false` |
| **envs.bucket.privatePrefix** | Prefix of the private system bucket | `system-private` |
| **envs.bucket.publicPrefix** | Prefix of the public system bucket | `system-public` |
| **envs.bucket.region** | Region of the system buckets | `us-east-1` |
| **envs.configMap.enabled** | Toggle used to save and load the configuration using the ConfigMap | `true` |
| **envs.configMap.name** | ConfigMap name | `rafter-upload-service` |
| **envs.configMap.namespace** | Namespace in which the ConfigMap is created | `{{ .Release.Namespace }}` |
| **migrator.pre.minioDeploymentRefName** | Name of the MinIO Deployment used in the pre-migration job. If not set, it is generated using the **Release.Name** template with the `-minio` suffix. | `nil` |
| **migrator.pre.minioSecretRefName** | Name of the MinIO Secret used in the pre-migration job. If not set, it is generated using the **Release.Name** template with the `-minio` suffix. | `nil` |
| **migrator.post.minioSecretRefName** | Name of the MinIO Secret used in the post-migration job. If not set, it is generated using the **Release.Name** template with the `-minio` suffix. | `nil` |

Specify each parameter using the `--set key=value[,key=value]` argument for `helm install`. See this example:

``` bash
helm install --name rafter-upload-service \
  --set serviceMonitor.create=true,serviceMonitor.name="rafter-service-monitor" \
    rafter-charts/rafter-upload-service
```

That command installs the release with the `rafter-service-monitor` name for the ServiceMonitor custom resource.

Alternatively, use the default values in [`values.yaml`](./values.yaml) or provide a YAML file while installing the chart to specify the values for configurable parameters. See this example:

``` bash
helm install --name rafter-upload-service -f values.yaml rafter-charts/rafter-upload-service
```

### values.yaml as a template

The `values.yaml` for the Upload Service chart serves as a template. Use such Helm variables as `.Release.*`, or `.Values.*`. See this example:

``` yaml
pod:
  annotations:
    sidecar.istio.io/inject: "{{ .Values.injectIstio }}"
    recreate: "{{ .Release.Time.Seconds }}"
``` 

### Change values for envs. parameters

You can define values for all **envs.** parameters as objects by providing the parameters as the inline `value` or the **valueFrom** object. See the following example:

``` yaml
envs:
  upload:
    workers:
      value: "0.0.0.0"
  verbose:
    valueFrom:
      configMapKeyRef:
        name: rafter-upload-service-config
        key: RAFTER_UPLOAD_SERVICE_VERBOSE
```

### Switch MinIO to Gateway mode

By default, you install the Upload Service with MinIO in stand-alone mode. If you want to switch MinIO to Gateway mode and you don't want to lose your buckets uploaded by the Upload Service, you must first override parameters for MinIO under the **minio** object and change:

-  **minio.persistence.enabled** to `false` 
-  **minio.podAnnotations.persistence** to `off`

> **NOTE:** If the names of deployments or secrets used before and after switching to Gateway mode differ, you must update parameters under **migrator.pre** and **migrator.post** objects.
