#!/bin/bash
#
# k8s-cms
# Web Admin Entrypoint
#

# add admin user if present
python3 ./cmscontrib/AddAdmin.py -p "$CMS_ADMIN_PASSWORD" "$CMS_ADMIN_USER"

# run admin web server
exec ./scripts/cmsAdminWebServer 
