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
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

# configure user/groups
RUN useradd --user-group --home-dir /cms cmsuser

# install python dependencies
RUN ln -s /usr/bin/python3 /usr/bin/python
RUN mkdir /cms
COPY deps/cms/requirements.txt /cms/requirements.txt
RUN pip3 install -r /cms/requirements.txt

# setup cms source directory
COPY deps/cms /cms
ENV PYTHONPATH=/cms
RUN chown -R cmsuser /cms
RUN chmod -R u+x /cms/scripts/
RUN bash -c "cd /cms && python setup.py install"

# copy scripts
RUN mkdir /scripts
COPY containers/cms-base/cms-entrypoint.sh /scripts/
COPY deps/wait-for-it/wait-for-it.sh /scripts/
RUN chmod -R u+x /scripts

# copy cms config
COPY config /cms/config

# configure entrypoint script
ENV HOME /cms
WORKDIR /cms
ENTRYPOINT ["/scripts/cms-entrypoint.sh"]
CMD bash
