#
# k8s-cms
# Scoring Service Container
# 

FROM mrzzy/cms-base

# copy config
COPY config /cms/config

# run scoring service
EXPOSE 28500
CMD ./scripts/cmsScoringService
