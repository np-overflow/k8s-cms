#/bin/bash
set -e
# k8s-cms
# Worker Entrypoint Script
# 

# determine worker shard no.
CMS_WORKER_SHARD=${CMS_WORKER_SHARD:-"0"}
exec ./scripts/cmsWorker $CMS_WORKER_SHARD
