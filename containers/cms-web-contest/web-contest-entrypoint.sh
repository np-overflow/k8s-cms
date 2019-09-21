#/bin/bash
set -e
#
# k8s-cms
# Contest Web Server Entrypoint
#

# figure out the contest id that this contest server targets.
CONTEST_ID=${CMS_CONTEST_ID:-"DEFAULT"}
POLL_INTERVAL=${CMS_POLL_INTERVAL:-"15"}
if [ "$CONTEST_ID" = "DEFAULT" ]
then
    # configure contest web server to host all contests
    ./scripts/cmsContestWebServer -c ALL 0
else
    exec ./scripts/cmsContestWebServer --contest-id $CONTEST_ID 0
fi
