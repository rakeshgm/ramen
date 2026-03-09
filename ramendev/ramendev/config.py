# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import drenv
from drenv import kubectl
from drenv import minio

from . import command


def register(commands):
    parser = commands.add_parser(
        "config",
        help="Configure ramen hub operator",
    )
    parser.set_defaults(func=run)
    command.add_common_arguments(parser)
    command.add_ramen_arguments(parser)


def _config_context(env, args):
    """
    Normalize metro-dr vs regional-dr into a single context.
    All topology branching lives here; the rest of the code is topology-agnostic.

    """
    dr_clusters = env["clusters"]
    storage_cluster = env.get("storage_cluster")
    ns = args.ramen_namespace

    storage_clusters = [storage_cluster] if storage_cluster else dr_clusters
    s3_profile_names = [{"profile_name": f"minio-on-{cluster}", "cluster": cluster} 
                        for cluster in storage_clusters]

    ocm_policy = []
    for dr_cluster in dr_clusters:
        cluster = storage_clusters[0] if storage_cluster else dr_cluster
        ocm_policy.append((dr_cluster, f"{ns}.ramen-s3-secret-{cluster}"))

    return {
        "dr_clusters": dr_clusters,
        "storage_clusters": storage_clusters,
        "s3_profile_names": s3_profile_names,
        "ocm_policy": ocm_policy,
    }


def run(args):
    env = command.env_info(args)
    ctx = _config_context(env, args)

    s3_secrets = generate_ramen_s3_secrets(ctx["storage_clusters"], args)

    if env["hub"]:
        hub_cm = generate_config_map("hub", env, args, ctx)

        create_ramen_s3_secrets(env["hub"], s3_secrets)

        create_ramen_config_map(env["hub"], hub_cm)
        create_hub_dr_resources(env["hub"], env["topology"], ctx)

        wait_for_secret_propagation(env["hub"], ctx["ocm_policy"], args)
        wait_for_dr_clusters(env["hub"], ctx["dr_clusters"], args)
        wait_for_dr_policy(env["hub"], args)

        # ramen-ops namespace is created by the drpolicy controller. It should
        # exist when the dr policy is validated.
        create_ramen_ops_binding(env["hub"])
    else:
        dr_cluster_cm = generate_config_map("dr-cluster", env, args, ctx)

        for cluster in ctx["dr_clusters"]:
            create_ramen_s3_secrets(cluster, s3_secrets)
            create_ramen_config_map(cluster, dr_cluster_cm)


def generate_ramen_s3_secrets(clusters, args):
    template = drenv.template(command.resource("ramen-s3-secret.yaml"))
    return [
        template.substitute(namespace=args.ramen_namespace, cluster=cluster)
        for cluster in clusters
    ]


def create_ramen_s3_secrets(cluster, secrets):
    command.info("Creating ramen s3 secrets in cluster '%s'", cluster)
    for secret in secrets:
        kubectl.apply("--filename=-", input=secret, context=cluster, log=command.debug)


def generate_config_map(controller, env, args, ctx):
    volsync = env["features"].get("volsync", True)
    substitute_kw = {
        "name": f"ramen-{controller}-operator-config",
        "auto_deploy": "true",
        "volsync_disabled": "false" if volsync else "true",
        "namespace": args.ramen_namespace,
    }
    s3_profile_names = ctx["s3_profile_names"]
    topology = env["topology"]
    template = drenv.template(command.resource(f"{topology}/configmap.yaml"))

    if len(s3_profile_names) == 1:
        substitute_kw["storage_cluster"] = s3_profile_names[0]["cluster"]
        substitute_kw["minio_url"] = minio.service_url(s3_profile_names[0]["cluster"])
    else:
        substitute_kw["cluster1"] = s3_profile_names[0]["cluster"]
        substitute_kw["cluster2"] = s3_profile_names[1]["cluster"]
        substitute_kw["minio_url_cluster1"] = minio.service_url(s3_profile_names[0]["cluster"])
        substitute_kw["minio_url_cluster2"] = minio.service_url(s3_profile_names[1]["cluster"])

    return template.substitute(**substitute_kw)


def create_ramen_config_map(cluster, yaml):
    command.info("Updating ramen config map in cluster '%s'", cluster)
    kubectl.apply("--filename=-", input=yaml, context=cluster, log=command.debug)


def create_ramen_ops_binding(cluster):
    command.info("Creating ramen-ops managedclustersetbinding in cluster '%s'", cluster)
    resource = command.resource("ramen-ops-binding.yaml")
    kubectl.apply(f"--filename={resource}", context=cluster, log=command.debug)


def create_hub_dr_resources(hub, topology, ctx):
    dr_clusters = ctx["dr_clusters"]
    s3_profile_names = ctx["s3_profile_names"]
    substitute_kw = {
        "cluster1": dr_clusters[0],
        "cluster2": dr_clusters[1],
    }
    if len(s3_profile_names) == 1:
        substitute_kw["storage_cluster"] = s3_profile_names[0]["cluster"]

    for name in ["dr-clusters", "dr-policy"]:
        command.info("Creating %s for %s", name, topology)
        template = drenv.template(command.resource(f"{topology}/{name}.yaml"))
        yaml = template.substitute(**substitute_kw)
        kubectl.apply("--filename=-", input=yaml, context=hub, log=command.debug)


def wait_for_secret_propagation(hub, ocm_policy_per_cluster, args):
    command.info("Waiting until s3 secrets are propagated to managed clusters")
    for cluster, policy in ocm_policy_per_cluster:
        command.debug("Waiting until policy '%s' reports status", policy)
        drenv.wait_for(
            f"policy/{policy}",
            output="jsonpath={.status}",
            namespace=cluster,
            timeout=60,
            profile=hub,
            log=command.debug,
        )
        command.debug("Waiting until policy %s is compliant", policy)
        kubectl.wait(
            f"policy/{policy}",
            "--for=jsonpath={.status.compliant}=Compliant",
            "--timeout=60s",
            f"--namespace={cluster}",
            context=hub,
            log=command.debug,
        )


def wait_for_dr_clusters(hub, clusters, args):
    command.info("Waiting until DRClusters report phase")
    for name in clusters:
        drenv.wait_for(
            f"drcluster/{name}",
            output="jsonpath={.status.phase}",
            timeout=180,
            profile=hub,
            log=command.debug,
        )

    command.info("Waiting until DRClusters phase is available")
    kubectl.wait(
        "drcluster",
        "--all",
        "--for=jsonpath={.status.phase}=Available",
        context=hub,
        log=command.debug,
    )


def wait_for_dr_policy(hub, args):
    command.info("Waiting until DRPolicy is validated")
    kubectl.wait(
        "drpolicy/dr-policy",
        "--for=condition=Validated",
        context=hub,
        log=command.debug,
    )
