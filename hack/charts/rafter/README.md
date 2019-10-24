# Rafter

This helm chart installs Rafter version v1.0.0 https://github.com/kyma-project/rafter/tree/v1.0.0

## Overview

This project contains the chart for the Rafter Controller Manager and his extensions.

## Prerequisites

- Kubernetes 1.12+
- Helm 2.11+ or Helm 3.0-beta3+

## Installing the Chart

To install the chart with the release name `rafter-release`:

``` bash
$ helm install --name rafter-release incubator/rafter
```

The command deploys Rafter on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `rafter-release` deployment:

``` bash
$ helm delete rafter-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Drupal chart and their default values.

| Parameter                         | Description                                | Default                                                   |
| --------------------------------- | ------------------------------------------ | --------------------------------------------------------- |
| `global.imageRegistry`            | Global Docker image registry               | `nil`                                                     |
| `global.imagePullSecrets`         | Global Docker registry secret names as an array | `[]` (does not add image pull secrets to deployed pods) |
| `global.storageClass`                     | Global storage class for dynamic provisioning                                               | `nil`                                                        |
| `image.registry`                  | Drupal image registry                      | `docker.io`                                               |
| `image.repository`                | Drupal Image name                          | `bitnami/drupal`                                          |
| `image.tag`                       | Drupal Image tag                           | `{TAG_NAME}`                                              |
| `image.pullPolicy`                | Drupal image pull policy                   | `IfNotPresent`                                            |
| `image.pullSecrets`               | Specify docker-registry secret names as an array | `[]` (does not add image pull secrets to deployed pods)  |
| `nameOverride`                    | String to partially override drupal.fullname template with a string (will prepend the release name) | `nil` |
| `fullnameOverride`                | String to fully override drupal.fullname template with a string                                     | `nil` |
| `drupalProfile`                   | Drupal installation profile                | `standard`                                                |
| `drupalUsername`                  | User of the application                    | `user`                                                    |
| `drupalPassword`                  | Application password                       | _random 10 character long alphanumeric string_            |
| `drupalEmail`                     | Admin email                                | `user@example.com`                                        |
| `allowEmptyPassword`              | Allow DB blank passwords                   | `yes`                                                     |
| `extraVars`                       | Extra environment variables                | `nil`                                                     |
| `ingress.enabled`                 | Enable ingress controller resource         | `false`                                                   |
| `ingress.hosts[0].name`           | Hostname to your Drupal installation       | `drupal.local`                                            |
| `ingress.hosts[0].path`           | Path within the url structure              | `/`                                                       |
| `ingress.hosts[0].tls`            | Utilize TLS backend in ingress             | `false`                                                   |
| `ingress.hosts[0].certManager`    | Add annotations for cert-manager           | `false`                                                   |
| `ingress.hosts[0].tlsSecret`      | TLS Secret (certificates)                  | `drupal.local-tls-secret`                                 |
| `ingress.hosts[0].annotations`    | Annotations for this host's ingress record | `[]`                                                      |
| `ingress.secrets[0].name`         | TLS Secret Name                            | `nil`                                                     |
| `ingress.secrets[0].certificate`  | TLS Secret Certificate                     | `nil`                                                     |
| `ingress.secrets[0].key`          | TLS Secret Key                             | `nil`                                                     |
| `externalDatabase.host`           | Host of the external database              | `nil`                                                     |
| `externalDatabase.user`           | Existing username in the external db       | `bn_drupal`                                               |
| `externalDatabase.password`       | Password for the above username            | `nil`                                                     |
| `externalDatabase.database`       | Name of the existing database              | `bitnami_drupal`                                          |
| `mariadb.enabled`                 | Whether to use the MariaDB chart           | `true`                                                    |
| `mariadb.rootUser.password`       | MariaDB admin password                     | `nil`                                                     |
| `mariadb.db.name`                 | Database name to create                    | `bitnami_drupal`                                          |
| `mariadb.db.user`                 | Database user to create                    | `bn_drupal`                                               |
| `mariadb.db.password`             | Password for the database                  | _random 10 character long alphanumeric string_            |
| `service.type`                    | Kubernetes Service type                    | `LoadBalancer`                                            |
| `service.port`                    | Service HTTP port                          | `80`                                                      |
| `service.httpsPort`               | Service HTTPS port                         | `443`                                                     |
| `service.externalTrafficPolicy`   | Enable client source IP preservation       | `Cluster`                                                 |
| `service.nodePorts.http`          | Kubernetes http node port                  | `""`                                                      |
| `service.nodePorts.https`         | Kubernetes https node port                 | `""`                                                      |
| `persistence.enabled`             | Enable persistence using PVC               | `true`                                                    |
| `persistence.drupal.storageClass` | PVC Storage Class for Drupal volume        | `nil` (uses alpha storage class annotation)               |
| `persistence.drupal.accessMode`   | PVC Access Mode for Drupal volume          | `ReadWriteOnce`                                           |
| `persistence.drupal.existingClaim`| An Existing PVC name                       | `nil`                                                     |
| `persistence.drupal.hostPath`     | Host mount path for Drupal volume          | `nil` (will not mount to a host path)                     |
| `persistence.drupal.size`         | PVC Storage Request for Drupal volume      | `8Gi`                                                     |
| `resources`                       | CPU/Memory resource requests/limits        | Memory: `512Mi`, CPU: `300m`                              |
| `volumeMounts.drupal.mountPath`   | Drupal data volume mount path              | `/bitnami/drupal`                                         |
| `podAnnotations`                  | Pod annotations                            | `{}`                                                      |
| `affinity`                        | Map of node/pod affinities                 | `{}`                                                      |
| `metrics.enabled`                 | Start a side-car prometheus exporter       | `false`                                                   |
| `metrics.image.registry`          | Apache exporter image registry             | `docker.io`                                               |
| `metrics.image.repository`        | Apache exporter image name                 | `bitnami/apache-exporter`                                 |
| `metrics.image.tag`               | Apache exporter image tag                  | `{TAG_NAME}`                                              |
| `metrics.image.pullPolicy`        | Image pull policy                          | `IfNotPresent`                                            |
| `metrics.image.pullSecrets`       | Specify docker-registry secret names as an array | `[]` (does not add image pull secrets to deployed pods)      |
| `metrics.podAnnotations`          | Additional annotations for Metrics exporter pod  | `{prometheus.io/scrape: "true", prometheus.io/port: "9117"}` |
| `metrics.resources`               | Exporter resource requests/limit           | {}                                                         |

The above parameters map to the env variables defined in [bitnami/drupal](http://github.com/bitnami/bitnami-docker-drupal). For more information please refer to the [bitnami/drupal](http://github.com/bitnami/bitnami-docker-drupal) image documentation.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

``` bash
$ helm install --name rafter-release \
  --set drupalUsername=admin,drupalPassword=password,mariadb.mariadbRootPassword=secretpassword \
    incubator/rafter
```

The above command sets the Drupal administrator account username and password to `admin` and `password` respectively. Additionally, it sets the MariaDB `root` user password to `secretpassword`.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

``` bash
$ helm install --name rafter-release -f values.yaml incubator/rafter
```

> **Tip**: You can use the default [values.yaml](./values.yaml)
