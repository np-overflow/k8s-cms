#
# k8s-cms
# Admin Web Server Container
# 

FROM /cms-base:0.3.0b

# copy source
COPY /containers/cms-web-admin/web-admin-entrypoint.sh /scripts
RUN chown -R cmsuser /scripts
RUN chmod -R u+x /scripts

# copy config
COPY config /cms/config

# run admin web server
EXPOSE 8889
EXPOSE 21100
ENTRYPOINT [ "/scripts/cms-entrypoint.sh", "/scripts/web-admin-entrypoint.sh" ]
CMD [ "cmsAdminWebServer" ]
