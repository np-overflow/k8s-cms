#
# k8s-cms
# Master Container
# 

FROM mrzzy/cms-base

# install module depdendencies
COPY src/master/requirements.txt  /tmp/requirements.txt
RUN pip3 install -r requirements.txt

# copy source code
COPY src/master /master
RUN chown -R cmsuser /master/

# container entrypoint
WORKDIR /master
EXPOSE 5000
CMD ["python3", "server.py"]
