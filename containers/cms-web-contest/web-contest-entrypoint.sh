#/bin/bash
set -e
#
# k8s-cms
# Contest Web Server Entrypoint
#

# figure out the contest id that this contest server targets.
CONTEST_ID=${CMS_CONTEST_ID:-"DEFAULT"}
if [ "$CONTEST_ID" = "DEFAULT" ]
then
    # autoselect default contest
    printf "\n" | exec ./scripts/cmsContestWebServer 0
else
    exec ./scripts/cmsContestWebServer --contest-id $CONTEST_ID 0
fi
