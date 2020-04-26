#
# k8s-cms
# Scoring Service Container
# 

FROM npoverflow/cms-base:latest

# copy config
COPY config /cms/config

# run scoring service
EXPOSE 28500
CMD ["/cms/scripts/cmsScoringService", "0"]
