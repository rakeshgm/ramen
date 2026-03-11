<!--
SPDX-FileCopyrightText: The RamenDR authors
SPDX-License-Identifier: Apache-2.0
-->

# Metrics

Metrics are collected using Prometheus, and registered with its global
metrics registry in each controller. Ramen uses controller-runtime's native
metrics authentication and authorization, which secures the metrics endpoint
using Kubernetes RBAC without requiring additional sidecar containers.

There are two ways to access metrics in Ramen: using the Prometheus stack
(recommended) or direct access via curl/port-forward. More details on
each of these in the sections below.

More information on metrics is [here](https://book.kubebuilder.io/reference/metrics.html)

## 1. Using Prometheus Stack

We recommend to use [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus).
Installing this will give you a containerized stack that
includes Prometheus, AlertManager and Grafana.

For more detailed information and querying,
consider using the Prometheus Operator, but in summary:

### Setup

Follow the next steps before configuring ramen.

#### Installing kube-prometheus

Quickstart instructions are [here](https://github.com/prometheus-operator/kube-prometheus#quickstart).

#### Grant permission for prometheus to scrape metrics

Go to `ramen/config/hub/default/k8s/kustomizations.yaml`
and uncomment `../../../prometheus`
and `metrics_role_binding.yaml` under `Kustomization` section.
Next is to install and configure ramen.

#### Dashboard Access

[Accessing Graphical User Interfaces](https://github.com/prometheus-operator/kube-prometheus/blob/main/docs/access-ui.md#access-uis)

## 2. Basic testing (no Prometheus required)

If running from minikube or a container, expose the port using `port-forward`
on the hub. The endpoint exposed is `localhost:8443/metrics`.

```bash
kubectl port-forward -n ramen-system \
deployment/ramen-hub-operator 8443:8443
```

The metrics endpoint is secured with authentication and authorization.
To access it, you need a valid ServiceAccount token:

```bash
# Get a token for a ServiceAccount with proper permissions
TOKEN=$(kubectl create token <service-account-name> -n <namespace>)

# Access metrics with authentication
curl -k -H "Authorization: Bearer $TOKEN" https://localhost:8443/metrics
```

For quick debugging without authentication (not recommended for production),
you can temporarily disable metrics authentication by setting
`RamenConfig.Metrics.BindAddress` to `"0"` or by modifying the controller
configuration.

### Metrics List

All metrics are prefixed with `ramen_`. This makes them easier to find.

To get the list of all the Ramen metrics available and their descriptions,
run the Ramen code, then run this command:
`curl http://localhost:8443/metrics -s | grep "# HELP ramen_"`.
