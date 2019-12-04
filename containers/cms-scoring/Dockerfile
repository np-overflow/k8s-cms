#
# k8s-cms
# Scoring Service Container
# 

FROM mrzzy/cms-base:0.3.0b

# copy config
COPY config /cms/config

# run scoring service
EXPOSE 28500
CMD ["/cms/scripts/cmsScoringService", "0"]
