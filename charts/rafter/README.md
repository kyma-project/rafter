# Rafter

This project contains the main Helm chart for Rafter. It includes:

- [Controller Manager](../rafter-controller-manager)
- [Upload Service](../rafter-upload-service)
- [Front Matter Service](../rafter-front-matter-service)
- [AsyncAPI Service](../rafter-asyncapi-service)

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
helm install rafter-charts/rafter
```

To install the chart with the `rafter` release name, use:

``` bash
helm install --name rafter rafter-charts/rafter
```

The command deploys Rafter on the Kubernetes cluster with the default configuration. For more information, see the [**Configuration**](#configuration) section.

> **TIP:** To list all releases, use `helm list`.

### Uninstall the chart

To uninstall the `rafter` release, run:

``` bash
helm delete rafter
```

That command removes all the Kubernetes components associated with the chart and deletes the release.

### Configuration

For a list of the parameters that you can configure during the installation of a given component, read the **Configuration** section in the component's `README.md` document.

> **NOTE:** Use values for the appropriate component in an object named as this component. For example, to override values for the [**Controller Manager**](../rafter-controller-manager), use the **rafter-controller-manager** object.

Specify each parameter using the `--set key=value[,key=value]` argument for `helm install`. See this example:

``` bash
helm install --name rafter \
  --set rafter-controller-manager.serviceMonitor.create=true,rafter-controller-manager.serviceMonitor.name="rafter-controller-manager-service-monitor" \
    rafter-charts/rafter
```

That command installs the release with the `rafter-controller-manager-service-monitor` name for the ServiceMonitor custom resource created based on the template from the Rafter Controller Manager.

Alternatively, use the default values in [`values.yaml`](./values.yaml) or provide a YAML file while installing the chart to specify the values for configurable parameters. See this example:

``` bash
helm install --name rafter -f values.yaml rafter-charts/rafter
```

### values.yaml as a template

The `values.yaml` for the Rafter chart serves as a template. Use such Helm variables as `.Release.*`, or `.Values.*`. See this example:

``` yaml
rafter-controller-manager:
  pod:
    annotations:
      sidecar.istio.io/inject: "{{ .Values.injectIstio }}"
      recreate: "{{ .Release.Time.Seconds }}"
``` 

### Change values for envs. parameters

You can define values for all **envs.** parameters as objects by providing the parameters as the inline `value` or the **valueFrom** object. See the following example:

``` yaml
rafter-controller-manager:
  envs:
    clusterAssetGroup:
      relistInterval: 
        value: 5m
    assetGroup:
      valueFrom:
        configMapKeyRef:
          name: rafter-config
          key: RAFTER_ASSET_GROUP_RELIST_INTERVALL
```

### Switch MinIO to Gateway mode

By default, you install the Upload Service with MinIO in stand-alone mode. If you want to switch MinIO to Gateway mode and you don't want to lose your buckets uploaded by the Upload Service, you must override parameters for MinIO under the **rafter-controller-manager.minio** object and change these parameters:

- **rafter-upload-service.minio.persistence.enabled** to `false`
- **rafter-controller-manager.minio.persistence.enabled** to `false`
- **rafter-upload-service.minio.podAnnotations.persistence** to `off`
- **rafter-controller-manager.minio.podAnnotations.persistence** to `off`

> **NOTE:** If the names of deployments or secrets used before and after switching to Gateway mode differ, you must update parameters under **rafter-upload-service.migrator.pre** and **rafter-upload-service.migrator.post** objects.
