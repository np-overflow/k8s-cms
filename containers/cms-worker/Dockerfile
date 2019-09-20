#
# k8s-cms
# Worker Container
# 

FROM mrzzy/cms-base

## install packages for language support
# no user interaction.
ARG DEBIAN_FRONTEND=noninteractive 
RUN apt-get update && apt-get install -y \
    build-essential \
    openjdk-8-jdk-headless \   
    fp-compiler \
    cppreference-doc-en-html \
    cgroup-lite zip \
    php7.2-cli php7.2-fpm \
    texlive-latex-base a2ps \
    phppgadmin a2ps haskell-platform rustc \
    mono-mcs \
    libcap-dev \
    libyaml-dev \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*
    
# install isolate sandbox
RUN python /cms/prerequisites.py --as-root -y build_isolate
RUN python /cms/prerequisites.py --as-root -y install_isolate
RUN rm -rf /cms/isolate/*.o /cms/isolate/isolate

# copy container source
COPY containers/cms-worker /scripts
RUN chmod -R u+x /scripts

# copy config
COPY config /cms/config

# run worker
EXPOSE 26000
CMD /scripts/worker-entrypoint.sh
