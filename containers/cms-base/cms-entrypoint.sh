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

# database dependency check: wait for database to start
CMS_DB=${CMS_DB:-"localhost"}
if [ "$CMS_DB" != "localhost" ]
then
    /scripts/wait-for-it.sh -t 10 -h $CMS_DB -p 5432
fi

# lose root privileges
su cmsuser
exec "$@"
