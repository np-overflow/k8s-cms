#
# k8s-cms
# Printing Service Container
# 

FROM mrzzy/cms-base:0.2.1b

# copy config
COPY config /cms/config

# run ranking
EXPOSE 25123
CMD [ "/cms/scripts/cmsPrintingService", "0" ] 
