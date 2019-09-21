#!/bin/bash
set -e
#
# k8s-cms
# Docker Entrypoint
#

# populates configuration with environment values
for CONFIG_FILE in /cms/config/*
do
    envsubst < $CONFIG_FILE > /tmp/$(basename $CONFIG_FILE)
    mv /tmp/$(basename $CONFIG_FILE) $CONFIG_FILE
done

exec "$@"
