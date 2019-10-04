#
# k8s-cms
# Resource Usage Container
# 

FROM mrzzy/cms-base:0.2.1b

# copy config
COPY config /cms/config

# run resource usage
EXPOSE 28000
CMD ["/cms/scripts/cmsResourceService", "0"]
