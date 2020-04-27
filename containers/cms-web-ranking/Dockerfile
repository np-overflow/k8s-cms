#
# k8s-cms
# Ranking Web Server Container
# 

FROM npoverflow/cms-base:latest

# copy config
COPY config /cms/config

# run ranking web server
EXPOSE 8890
CMD ./scripts/cmsRankingWebServer
