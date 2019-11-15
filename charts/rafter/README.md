# Rafter Umbrella Chart

This project contains the Helm chart for the ecosystem of the Rafter. It includes:

- [Controller Manager](../rafter-controller-manager)
- [Upload Service](../rafter-upload-service)
- [Front Matter Service](../rafter-front-matter-service)
- [AsyncAPI Service](../rafter-asyncapi-service)

## Prerequisites

- Kubernetes v1.14 or higher
- Helm v2.10 or higher

## Details

Read how to install, uninstall, and configure the chart.

### Install the chart

Use this command to install the chart:

``` bash
helm install incubator/rafter
```

To install the chart with the `rafter` release name, use:

``` bash
helm install --name rafter incubator/rafter
```

The command deploys the ecosystem of the Rafter on the Kubernetes cluster with the default configuration. For more information, please see the [**Configuration**](#configuration) section.

> **TIP:** To list all releases, use `helm list`.

### Uninstall the chart

To uninstall the `rafter` release, run:

``` bash
helm delete rafter
```

That command removes all the Kubernetes components associated with the chart and deletes the release.

### Configuration

For list of the parameters that you can configure during installation appropriate component, please read **Configuration** section of main doc of relevant component.

> **NOTE:** Remember about using values for appropriate component in an object named as a component. For example, for overriding values for [**Controller Manager**](../rafter-controller-manager) use `rafter-controller-manager` object.

Specify each parameter using the `--set key=value[,key=value]` argument for `helm install`. See this example:

``` bash
helm install --name rafter \
  --set rafter-controller-manager.serviceMonitor.create=true,rafter-controller-manager.serviceMonitor.name="rafter-controller-manager-service-monitor" \
    incubator/rafter
```

That command installs the release with the `rafter-controller-manager-service-monitor` name for the ServiceMonitor custom resource of Rafter Controller Manager.

Alternatively, use the default values in [values.yaml](./values.yaml) or provide a YAML file while installing the chart to specify the values for configurable parameters. See this example:

``` bash
helm install --name rafter -f values.yaml incubator/rafter
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

You can define values for all **envs.** parameters as objects by providing the parameters as the inline `value` or the `valueFrom` object. See the following example:

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
