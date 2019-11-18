# AsyncAPI Service

This project contains the Helm chart for the AsyncAPI Service.

## Prerequisites

- Kubernetes v1.14 or higher
- Helm v2.15 or higher

## Details

Read how to install, uninstall, and configure the chart.

### Install the chart

Use this command to install the chart:

``` bash
helm install incubator/rafter-asyncapi-service
```

To install the chart with the `rafter-asyncapi-release` release name, use:

``` bash
helm install --name rafter-asyncapi-release incubator/rafter-asyncapi-service
```

The command deploys the AsyncAPI Service on the Kubernetes cluster with the default configuration. The [**Configuration**](#configuration) section lists the parameters that you can configure during installation.

> **TIP:** To list all releases, use `helm list`.

### Uninstall the chart

To uninstall the `rafter-asyncapi-release` release, run:

``` bash
helm delete rafter-asyncapi-release
```

That command removes all the Kubernetes components associated with the chart and deletes the release.

### Configuration

The following table lists the configurable parameters of the AsyncAPI Service chart and their default values.

| Parameter | Description | Default |
| --- | ---| ---|
| **image.repository** | AsyncAPI Service image repository | `eu.gcr.io/kyma-project/rafter-asyncapi-service` |
| **image.tag** | AsyncAPI Service image tag | `{TAG_NAME}` |
| **image.pullPolicy** | Pull policy for the AsyncAPI Service image | `IfNotPresent` |
| **nameOverride** | String that partially overrides the **rafterAsyncAPIService.name** template | `nil` |
| **fullnameOverride** | String that fully overrides the **rafterAsyncAPIService.fullname** template | `nil` |
| **deployment.labels** | Custom labels for the Deployment | `{}` |
| **deployment.annotations** | Custom annotations for the Deployment | `{}` |
| **deployment.replicas** | Number of AsyncAPI Service nodes | `1` |
| **deployment.extraProperties** | Additional properties injected in the Deployment | `{}` |
| **pod.labels** | Custom labels for the Pod | `{}` |
| **pod.annotations** | Custom annotations for the Pod | `{}` |
| **pod.extraProperties** | Additional properties injected in the Pod | `{}` |
| **pod.extraContainerProperties** | Additional properties injected in the container | `{}` |
| **service.name** | Service name. If not set, it is generated using the **rafterAsyncAPIService.fullname** template. | `nil` |
| **service.type** | Service type | `ClusterIP` |
| **service.port.name** |  Name of the Service port | `http` |
| **service.port.internal** | Internal port of the Service in the Pod | `3000` |
| **service.port.external** | Port on which the Service is exposed in Kubernetes | `80` |
| **service.port.protocol** | Protocol of the Service port | `TCP` |
| **service.labels** | Custom labels for the Service | `{}` |
| **service.annotations** | Custom annotations for the Service | `{}` |
| **serviceMonitor.create** | Parameter that defines whether to create a new ServiceMonitor custom resource for the Prometheus Operator | `false` |
| **serviceMonitor.name** | ServiceMonitor resource that the Prometheus Operator uses. If not set and the **serviceMonitor.create** parameter is set to `true`, the name is generated using the **rafterAsyncAPIService.fullname** template. If not set and **serviceMonitor.create** is set to `false`, the name is set to `default`. | `nil` |
| **serviceMonitor.scrapeInterval** | Scrape interval for the ServiceMonitor custom resource | `30s` |
| **serviceMonitor.labels** | Custom labels for the ServiceMonitor custom resource | `{}` |
| **serviceMonitor.annotations** | Custom annotations for the ServiceMonitor custom resource | `{}` |
| **envs.host** | AsyncAPI Service host | `0.0.0.0` |
| **envs.verbose** | Parameter that defines if logs from the AsyncAPI Service should be visible | `true` |

Specify each parameter using the `--set key=value[,key=value]` argument for `helm install`. See this example:

``` bash
helm install --name rafter-asyncapi-release \
  --set serviceMonitor.create=true,serviceMonitor.name="rafter-service-monitor" \
    incubator/rafter-asyncapi-service
```

That command installs the release with the `rafter-service-monitor` name for the ServiceMonitor custom resource.

Alternatively, use the default values in [values.yaml](./values.yaml) or provide a YAML file while installing the chart to specify the values for configurable parameters. See this example:

``` bash
helm install --name rafter-asyncapi-release -f values.yaml incubator/rafter-asyncapi-service
```

### values.yaml as a template

The `values.yaml` for the AsyncAPI Service chart serves as a template. Use such Helm variables as `.Release.*`, or `.Values.*`. See this example:

``` yaml
pod:
  annotations:
    sidecar.istio.io/inject: "{{ .Values.injectIstio }}"
    recreate: "{{ .Release.Time.Seconds }}"
``` 

### Change values for envs. parameters

You can define values for all **envs.** parameters as objects by providing the parameters as the inline `value` or the `valueFrom` object. See the following example:

``` yaml
envs:
  host:
    value: "0.0.0.0"
  verbose:
    valueFrom:
      configMapKeyRef:
        name: rafter-asyncapi-service-config
        key: RAFTER_ASYNCAPI_SERVICE_VERBOSE
```
