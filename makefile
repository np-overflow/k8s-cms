#
# k8s-cms
# Project Makefile
# 

# vars
VERSION:=0.2.0a

DOCKER:=docker
TAG_PREFIX:=mrzzy

CMS_SRC_DIR:=deps/cms/

# names of the docker images
IMG_NAMES:=$(notdir $(wildcard containers/*))
IMAGES:=$(foreach img,$(IMG_NAMES),$(TAG_PREFIX)/$(img))
BASE_IMAGE:=$(TAG_PREFIX)/cms-base
# names of the images that depends on base image
DEP_BASE_NAMES:=$(filter-out cms-base,$(filter-out cms-db, $(IMG_NAMES)))
DEP_BASE_IMAGES:=$(foreach img,$(DEP_BASE_NAMES),$(TAG_PREFIX)/$(img))

PUSH_TARGETS:=$(foreach img,$(IMAGES),push/$(img))

# phony rules
.PHONY: all push
.DEFAULT: all 

all: $(DEP_BASE_IMAGES)

# image rules
# image deps
$(DEP_BASE_IMAGES): $(BASE_IMAGE)

# docker build rule
$(TAG_PREFIX)/%: containers/%/Dockerfile $(CMS_SRC_DIR)
	# versioned tag
	$(DOCKER) build -f $< -t $@:$(VERSION) .
	# latest tag
	$(DOCKER) build -f $< -t $@ .

# docker push rule
push: $(PUSH_TARGETS)

push/%: %
	docker push $<
