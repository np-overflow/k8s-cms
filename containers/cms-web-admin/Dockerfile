#
# k8s-cms
# Admin Web Server Container
# 

FROM mrzzy/cms-base


# copy source
COPY /containers/cms-web-admin /scripts
RUN chmod a+x /scripts

# copy config
COPY config /cms/config

# run admin web server
EXPOSE 8889
CMD bash /scripts/web-admin-entrypoint.sh
