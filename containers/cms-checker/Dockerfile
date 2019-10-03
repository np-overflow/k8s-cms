#
# k8s-cms
# Health Checker Container
# 

FROM mrzzy/cms-base:0.2.0a

# copy config
COPY config /cms/config

# run scoring service
EXPOSE 28500
CMD ["/cms/scripts/cmsChecker", "0"]
