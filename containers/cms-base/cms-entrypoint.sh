#!/bin/bash
set -e
#
# k8s-cms
# Docker Entrypoint
#

# determine path of configuration
CMS_CONFIG=${CMS_CONFIG:-"/cms/config/cms.conf"}
CMS_RANKING_CONFIG=${CMS_CONFIG:-"/cms/config/cms.conf"}

# populates configuration with environment values
envsubst < $CMS_CONFIG > "/etc/$(basename $CMS_CONFIG)"
envsubst < $CMS_RANKING_CONFIG > "/etc/$(basename $CMS_RANKING_CONFIG)"

# update env vars pointing to populated config path
CMS_CONFIG="/etc/$(basename $CMS_CONFIG)"
CMS_RANKING_CONFIG="/etc/$(basename $CMS_RANKING_CONFIG)"

exec "$@"
