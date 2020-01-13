#!/bin/bash
#
# k8s-cms
# Script to Configure Database
#

# configure database according to CMS docs
createdb --username=$POSTGRES_USER --owner=$POSTGRES_USER cmsdb 
psql --username=$POSTGRES_USER --dbname=cmsdb --command="ALTER SCHEMA public OWNER TO $POSTGRES_USER"
psql --username=$POSTGRES_USER --dbname=cmsdb --command="GRANT SELECT ON pg_largeobject TO $POSTGRES_USER" 

# run db setup scripts
exec bash -c '/cms/scripts/cmsInitDB'
