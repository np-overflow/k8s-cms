#!/bin/bash
set -e
#
# k8s-cms
# Docker Entrypoint
#

## setup config
# determine path of configuration
CMS_CONFIG=${CMS_CONFIG:-"/cms/config/cms.conf"}
CMS_RANKING_CONFIG=${CMS_RANKING_CONFIG:-"/cms/config/cms.ranking.conf"}
# populates configuration with environment values
envsubst < $CMS_CONFIG > "/etc/$(basename $CMS_CONFIG)"
envsubst < $CMS_RANKING_CONFIG > "/etc/$(basename $CMS_RANKING_CONFIG)"
# update env vars pointing to populated config path
export CMS_CONFIG="/etc/$(basename $CMS_CONFIG)"
export CMS_RANKING_CONFIG="/etc/$(basename $CMS_RANKING_CONFIG)"

CMS_DB=${CMS_DB:-"0.0.0.0"}
if [ "$CMS_DB" != "0.0.0.0" ] # check if not running as database
then
    # database dependency check: wait for database to start
    CMS_DB_WAIT=${CMS_DB_WAIT:30} # how long to wait for the database
    /scripts/wait-for-it.sh -t $CMS_DB_WAIT -h $CMS_DB -p 5432

    # lose root privilege to tighten security
    exec su --preserve-environment -c "$@" cmsuser
else
    # database requires root permissions
    exec "$@"
fi
