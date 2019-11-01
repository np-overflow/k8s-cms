#
# k8s-cms
# Log Service Container
# 

FROM /cms-base:0.3.0b

# copy config
COPY config /cms/config

# run log service
EXPOSE 29000
CMD [ "./scripts/cmsLogService", "0" ]
