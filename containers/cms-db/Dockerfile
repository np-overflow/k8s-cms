#
# k8s-cms
# Database Container
# 

FROM postgres:11.5

# install packages
RUN apt-get update && apt-get install -y \
    python3-pip  \
    libpq-dev \
    libcups2-dev \
    gettext-base \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

# configure user/groups
RUN useradd --user-group --home-dir /project cmsuser

# install python dependencies
RUN ln -s /usr/bin/python3 /usr/bin/python
RUN mkdir /cms
COPY deps/cms/requirements.txt /cms/requirements.txt
RUN pip3 install -r /cms/requirements.txt

# copy scripts
COPY deps/cms /cms
COPY config /cms/config
ENV PYTHONPATH=/cms
RUN chmod -R a+x /cms/scripts
COPY containers/cms-db /scripts
COPY deps/wait-for-it/wait-for-it.sh /scripts/wait-for-it.sh
COPY containers/cms-base/cms-entrypoint.sh /scripts/cms-entrypoint.sh
RUN chmod -R a+x /scripts
# insert modified docker-entrypoint.sh
RUN cp /scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
 
# configure postgres db
RUN mkdir -p /docker-entrypoint-initdb.d
RUN cp /scripts/setup_db.sh /docker-entrypoint-initdb.d/

# configure entrypoint script
ENV CMS_DB=0.0.0.0
ENTRYPOINT ["/scripts/cms-entrypoint.sh"]
CMD [ "docker-entrypoint.sh", "postgres" ]
