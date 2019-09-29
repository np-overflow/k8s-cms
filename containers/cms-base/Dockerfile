#
# k8s-cms
# Database Container
# 

FROM ubuntu:18.04

# install packages
RUN apt-get update && apt-get install -y \
    python3-dev \
    python3-pip  \
    libpq-dev \
    libcups2-dev \
    gettext-base \
    iputils-ping \
    netcat \
    && rm -rf /var/lib/apt/lists/*

# configure user/groups
RUN useradd --user-group --home-dir /home/cms cmsuser

# install python dependencies
RUN ln -s /usr/bin/python3 /usr/bin/python
RUN mkdir /cms
COPY deps/cms/requirements.txt /cms/requirements.txt
RUN pip3 install -r /cms/requirements.txt

# setup cms source directory
COPY deps/cms /cms
ENV PYTHONPATH=/cms
RUN chmod -R a+x /cms/scripts
RUN bash -c "cd /cms && python setup.py install"

# copy scripts
COPY containers/cms-base /scripts
COPY deps/wait-for-it/wait-for-it.sh /scripts/wait-for-it.sh
RUN chmod -R a+x /scripts

# configure entrypoint script
ENTRYPOINT ["/scripts/cms-entrypoint.sh"]
WORKDIR /cms
CMD bash
