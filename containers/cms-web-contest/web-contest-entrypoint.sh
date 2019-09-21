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
    # poll until a contest has been discovered
    while [ $? -eq 0 ]
    do
        # autoselect default contest
        printf "\n" | ./scripts/cmsContestWebServer -c ALL 0

        echo "ContestWebServer waiting for a contest to be created..."
        sleep $POLL_INTERVAL
    done
else
    exec ./scripts/cmsContestWebServer --contest-id $CONTEST_ID 0
fi
