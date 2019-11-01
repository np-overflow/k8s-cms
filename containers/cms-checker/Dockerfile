#
# k8s-cms
# Health Checker Container
# 

FROM /cms-base:0.3.0b

# copy config
COPY config /cms/config

# run scoring service
EXPOSE 28500
CMD ["/cms/scripts/cmsChecker", "0"]
