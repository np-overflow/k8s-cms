#!/bin/bash
#
# k8s-cms
# Script to Configure Database
#

# configure database according to CMS docs
psql --username=postgres -c "CREATE USER cmsuser WITH PASSWORD '$POSTGRES_PASSWORD';" 
createdb --username=postgres --owner=cmsuser cmsdb 
psql --username=postgres --dbname=cmsdb --command='ALTER SCHEMA public OWNER TO cmsuser' 
psql --username=postgres --dbname=cmsdb --command='GRANT SELECT ON pg_largeobject TO cmsuser' 

# run db setup scripts
exec bash -c '/cms/scripts/cmsInitDB'
