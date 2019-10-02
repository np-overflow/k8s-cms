#!/bin/bash
#
# k8s-cms
# Web Admin Entrypoint
#

# determine if we are execing a command or web admin server proper
if [ "$@" != "cmsAdminWebServer" ]
then
    exec $@
fi

# add admin user if not prosent present
python3 ./cmscontrib/AddAdmin.py -p "$CMS_ADMIN_PASSWORD" "$CMS_ADMIN_USER"

# run admin web server
exec ./scripts/cmsAdminWebServer 0
