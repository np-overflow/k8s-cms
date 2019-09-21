#
# k8s-cms
# Project Makefile
# 

# vars
VERSION:=0.1.0a

DOCKER:=docker
TAG_PREFIX:=mrzzy

IMG_NAMES:=$(notdir $(wildcard containers/*))
DEP_BASE_NAMES:=$(filter-out cms-base,$(filter-out cms-db, $(IMG_NAMES)))
CMS_SRC_DIR:=deps/cms/

BASE_IMAGE:=$(TAG_PREFIX)/cms-base
IMAGES:=$(foreach img,$(IMG_NAMES),$(TAG_PREFIX)/$(img))
DEP_BASE_IMAGES:=$(foreach img,$(DEP_BASE_NAMES),$(TAG_PREFIX)/$(img))


# phony rules
.PHONY: all
.DEFAULT: all 

all: $(IMAGES)

# image rules
# image deps
$(DEP_BASE_IMAGES): $(BASE_IMAGE)

# docker build rule
$(TAG_PREFIX)/%: containers/%/Dockerfile $(CMS_SRC_DIR)
	$(DOCKER) build -f $< -t $@:$(VERSION) .
