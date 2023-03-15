<!--
SPDX-FileCopyrightText: The RamenDR authors
SPDX-License-Identifier: Apache-2.0
-->

# Metrics

Metrics are collected using Prometheus, and registered with its global
metrics registry in each controller. There are two ways where you can
look at the metrics in ramen. One way is use prometheus stack(recommended)
and the other way is to use curl or postman. More details on
each of these in the below sections.

More information on metrics is [here](https://book.kubebuilder.io/reference/metrics.html)

## 1. Using Prometheus Stack

We recommend to use [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus).
Installing this will give you a containerized stack that
includes Prometheus, AlertManager and Grafana.
Quickstart instructions are [here](https://github.com/prometheus-operator/kube-prometheus#quickstart).

For more detailed information and querying,
consider using the Prometheus Operator, but in summary:

### Setup

Follow the next steps before configuring ramen.

#### Installing kube-prometheus

```bash
git clone https://github.com/prometheus-operator/kube-prometheus.git

cd kube-prometheus

# install CRDs
kubectl create -f manifests/setup

# install resources
kubectl create -f manifests/

# verify everything is running
kubectl get all -n monitoring
```

#### Grant permission for prometheus to scrape metrics

Go to `ramen/config/hub/default/k8s/kustomizations.yaml`
and uncomment `../../../prometheus`
and `metrics_role_binding.yaml` under `Kustomization` section.
Next is to install and configure ramen.

#### Dashboard Access

To access GUI dashboards for each component, use `port-forward`:

* Prometheus:
 `$ kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090`
* Grafana: `$ kubectl --namespace monitoring port-forward svc/grafana 3000`
* Alert Manager:
 `$ kubectl --namespace monitoring port-forward svc/alertmanager-main 9093`

Then navigate to the appropriate port on localhost, e.g.
[http://localhost:9090](http://localhost:9090) for Prometheus.

Ramen metrics can be found in namespace "ramen-system"; e.g. PromQL query
`{namespace="ramen-system"}`.

#### Delete stack

```bash
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup
```

## 2. Basic testing (no Prometheus required)

If running from minikube or a container, expose the port using `port-forward`
on the hub. The endpoint exposed is `localhost:8443/metrics`.

>*This way does not use kube-rbac-proxy, use this method only for any quick debugging*

```bash
kubectl port-forward -n ramen-system \
deployment/ramen-hub-operator 8443:9289
```

Verify that the metrics endpoint is exposed with curl:
`curl http://localhost:8443/metrics`.

If curl can connect, search for your metrics in the output.

### Metrics List

All metrics are prefixed with `ramen_`. This makes them easier to find.

To get the list of all the Ramen metrics available and their descriptions,
run the Ramen code, then run this command:
`curl http://localhost:8443/metrics -s | grep "# HELP ramen_"`.
