#/bin/bash
set -e
#
# k8s-cms
# Worker Entrypoint Script
# 

# determine worker shard no.
if [ -n "$CMS_WORKER_NAME" ]
then
    CMS_WORKER_SHARD=$(printf "$CMS_WORKER_NAME" | \
        sed -e 's/cms-worker-\([0-9]*\)\(.cms-worker\)*/\1/g')
else
    CMS_WORKER_SHARD=${CMS_WORKER_SHARD:-"0"}
fi

exec ./scripts/cmsWorker $CMS_WORKER_SHARD
