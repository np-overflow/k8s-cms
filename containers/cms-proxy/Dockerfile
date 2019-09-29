#
# k8s-cms
# Proxy Service Container
# 

FROM mrzzy/cms-base

# copy container source
COPY containers/cms-proxy /scripts
RUN chmod -R a+x /scripts

# copy config
COPY config /cms/config

# run proxy service
EXPOSE 28600
CMD /scripts/proxy-entrypoint.sh
