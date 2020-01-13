#
# k8s-cms
# Master Container
# 

FROM mrzzy/cms-base

# install module depdendencies
COPY src/master/requirements.txt  /tmp/requirements.txt
RUN pip3 install -r /tmp/requirements.txt

# copy source code
WORKDIR /master
COPY src/master /master
RUN chown -R cmsuser /master/

# container entrypoint
EXPOSE 5000
CMD ["python3", "server.py"]
