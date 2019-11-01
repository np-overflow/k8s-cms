#
# k8s-cms
# Resource Usage Container
# 

FROM /cms-base:0.3.0b

# copy config
COPY config /cms/config

# run resource usage
EXPOSE 28000
CMD ["/cms/scripts/cmsResourceService", "0"]
