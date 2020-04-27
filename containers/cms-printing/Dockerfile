#
# k8s-cms
# Printing Service Container
# 

FROM npoverflow/cms-base:latest

# copy config
COPY config /cms/config

# run ranking
EXPOSE 25123
CMD [ "/cms/scripts/cmsPrintingService", "0" ] 
