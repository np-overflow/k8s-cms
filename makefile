#
# k8s-cms
# Project Makefile
# 

# vars
DOCKER:=docker
TAG_PREFIX:=mrzzy

IMG_NAMES:=cms-base cms-db cms-web-admin
DEP_BASE_NAMES:=cms-web-admin
CMS_SRC_DIR:=deps/cms/

BASE_IMAGE:=$(TAG_PREFIX)/cms-base
IMAGES:=$(foreach img,$(IMG_NAMES),$(TAG_PREFIX)/$(img))
DEP_BASE_IMAGES:=$(foreach img,$(DEP_BASE_NAMES),$(TAG_PREFIX)/$(img))

# phony rules
.PHONY: all clean push pull
.DEFAULT: all

all: $(IMAGES)

# image rules
# image deps
$(DEP_BASE_IMAGES): $(BASE_IMAGE)

# docker build rule
$(TAG_PREFIX)/%: containers/%/Dockerfile $(CMS_SRC_DIR)
	$(DOCKER) build -f $< -t $@ .
