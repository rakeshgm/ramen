# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import drenv
from drenv import kubectl

from . import command


def register(commands):
    parser = commands.add_parser(
        "unconfig",
        help="Unconfigure ramen hub operator",
    )
    parser.set_defaults(func=run)
    command.add_common_arguments(parser)
    command.add_ramen_arguments(parser)


def run(args):
    env = command.env_info(args)

    # Note: We keep the ramen config map since we do not own it.

    if env["hub"]:
        from .config import _config_context

        ctx = _config_context(env, args)
        delete_hub_dr_resources(env["hub"], env["topology"], ctx)
        s3_secrets = generate_ramen_s3_secrets(ctx["storage_clusters"], args)
        delete_s3_secrets(env["hub"], s3_secrets)


def delete_hub_dr_resources(hub, topology, ctx):
    # Deleting in reverse order.
    dr_clusters = ctx["dr_clusters"]
    s3_profile_names = ctx["s3_profile_names"]
    substitute_kw = {
        "cluster1": dr_clusters[0],
        "cluster2": dr_clusters[1],
    }
    if len(s3_profile_names) == 1:
        substitute_kw["storage_cluster"] = s3_profile_names[0]["cluster"]

    for name in ["dr-policy", "dr-clusters"]:
        command.info("Deleting %s for %s", name, topology)
        template = drenv.template(command.resource(f"{topology}/{name}.yaml"))
        yaml = template.substitute(**substitute_kw)
        kubectl.delete(
            "--filename=-",
            "--ignore-not-found",
            input=yaml,
            context=hub,
            log=command.debug,
        )


def generate_ramen_s3_secrets(clusters, args):
    template = drenv.template(command.resource("ramen-s3-secret.yaml"))
    return [
        template.substitute(namespace=args.ramen_namespace, cluster=cluster)
        for cluster in clusters
    ]


def delete_s3_secrets(cluster, secrets):
    command.info("Deleting s3 secrets in cluster '%s'", cluster)
    for secret in secrets:
        kubectl.delete(
            "--filename=-",
            "--ignore-not-found",
            input=secret,
            context=cluster,
            log=command.debug,
        )
