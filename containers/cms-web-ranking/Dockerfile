#
# k8s-cms
# Ranking Web Server Container
# 

FROM mrzzy/cms-base:0.2.1b

# copy config
COPY config /cms/config

# run ranking web server
EXPOSE 8890
CMD ./scripts/cmsRankingWebServer
