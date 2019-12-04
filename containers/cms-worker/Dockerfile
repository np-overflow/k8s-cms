#
# k8s-cms
# Worker Container
# 

FROM mrzzy/cms-base:0.3.0b

## install packages for language support
# no user interaction.
ARG DEBIAN_FRONTEND=noninteractive 
RUN apt-get update && apt-get install -y \
    build-essential \
    openjdk-8-jdk-headless \   
    fp-compiler \
    cgroup-lite zip \
    php7.2-cli php7.2-fpm \
    texlive-latex-base a2ps \
    phppgadmin a2ps haskell-platform rustc \
    mono-mcs \
    libcap-dev \
    libyaml-dev \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*
    
# build &install isolate sandbox
RUN python /cms/prerequisites.py --as-root -y build_isolate
RUN python /cms/prerequisites.py --as-root -y install_isolate
RUN rm -rf /cms/isolate/*.o /cms/isolate/isolate

# copy container source
COPY containers/cms-worker/worker-entrypoint.sh /scripts/
RUN chown -R cmsuser /scripts/
RUN chmod -R a+x /scripts/

# copy config
COPY config /cms/config

# run worker
EXPOSE 26000
ENTRYPOINT [ "/scripts/cms-entrypoint.sh", "/scripts/worker-entrypoint.sh" ]
CMD [ "cmsWorker" ]
