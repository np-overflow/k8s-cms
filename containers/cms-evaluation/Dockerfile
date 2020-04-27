#
# k8s-cms
# Evaluation Service Container
# 

FROM npoverflow/cms-base:latest

# copy config
COPY config /cms/config

# run ranking web server
EXPOSE 25000
CMD ["/cms/scripts/cmsEvaluationService", "0"]
