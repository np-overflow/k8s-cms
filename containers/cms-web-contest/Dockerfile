#
# k8s-cms
# Contest Web Server Container
# 

FROM mrzzy/cms-base:0.2.0a

# copy container source
COPY containers/cms-web-contest /scripts
RUN chown -R cmsuser /scripts
RUN chmod -R u+x /scripts

# copy config
COPY config /cms/config

# run contest web server
EXPOSE 21000
EXPOSE 8888
CMD /scripts/web-contest-entrypoint.sh
