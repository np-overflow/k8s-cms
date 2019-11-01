#
# k8s-cms
# Printing Service Container
# 

FROM /cms-base:0.3.0b

# copy config
COPY config /cms/config

# run ranking
EXPOSE 25123
CMD [ "/cms/scripts/cmsPrintingService", "0" ] 
