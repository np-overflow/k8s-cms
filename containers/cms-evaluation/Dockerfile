#
# k8s-cms
# Evaluation Service Container
# 

FROM mrzzy/cms-base:0.3.0b

# copy config
COPY config /cms/config

# run ranking web server
EXPOSE 25000
CMD ["/cms/scripts/cmsEvaluationService", "0"]
