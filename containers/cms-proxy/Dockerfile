#
# k8s-cms
# Proxy Service Container
# 

FROM mrzzy/cms-base:0.3.0b

# copy container source
COPY containers/cms-proxy/proxy-entrypoint.sh /scripts/
RUN chown -R cmsuser /scripts
RUN chmod -R u+x /scripts

# copy config
COPY config /cms/config

# run proxy service
EXPOSE 28600
ENTRYPOINT ["/scripts/cms-entrypoint.sh", "/scripts/proxy-entrypoint.sh"]
CMD [ "cmsProxyService" ]
