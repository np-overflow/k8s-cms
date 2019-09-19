#/bin/bash
set -e
#
# k8s-cms
# Proxy Service Entrypoint
#

# figure out the contest id that this proxy service targets.
CONTEST_ID=${CMS_CONTEST_ID:-"DEFAULT"}
if [ "$CONTEST_ID" = "DEFAULT" ]
then
    # autoselect default contest
    printf "\n" | exec ./scripts/cmsProxyService 0
else
    exec ./scripts/cmsProxyService --contest-id $CONTEST_ID 0
fi
