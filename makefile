#
# k8s-cms
# Project Makefile
# 

# vars
VERSION:=latest

DOCKER:=docker
TAG_PREFIX:=mrzzy

CMS_SRC_DIR:=deps/cms/
EXPORT_DIR:=build/containers

# names of the docker images
IMG_NAMES:=$(notdir $(wildcard containers/*))
IMAGES:=$(foreach img,$(IMG_NAMES),$(TAG_PREFIX)/$(img))
BASE_IMAGE:=$(TAG_PREFIX)/cms-base
# names of the images that depends on base image
DEP_BASE_NAMES:=$(filter-out cms-base,$(filter-out cms-db, $(IMG_NAMES)))
DEP_BASE_IMAGES:=$(foreach img,$(DEP_BASE_NAMES),$(TAG_PREFIX)/$(img))

PUSH_TARGETS:=$(foreach img,$(IMAGES),push/$(img))
EXPORT_TARGETS:=$(foreach img,$(IMAGES),export/$(img))
LOAD_TARGETS:=$(foreach img,$(IMAGES),load/$(img))

# phony rules
.PHONY: all push clean clean-version export load
.DEFAULT: all 

all: $(DEP_BASE_IMAGES)

# image rules
# image deps
$(DEP_BASE_IMAGES): $(BASE_IMAGE)

# docker build rule
$(TAG_PREFIX)/%: containers/%/Dockerfile $(CMS_SRC_DIR)
	# latest tag
	$(DOCKER) build -f $< -t $@ .
	# versioned tag
	$(DOCKER) tag $@ $@:$(VERSION) 

# docker push rule
push: $(PUSH_TARGETS)

push/%:
	docker push $(subst push/,,$@)

# export docker images as tar archive
export: $(EXPORT_TARGETS)

export/%:
	mkdir -p $(EXPORT_DIR)
	docker save $(subst export/,,$@):$(VERSION) -o $(EXPORT_DIR)/$(notdir $@).tar

# load docker images from a tar archive
load: $(LOAD_TARGETS)

load/%:
	docker load -i $(EXPORT_DIR)/$(notdir $@).tar

# cleans docker images
# clean all docker images
clean: clean-version
	$(foreach img,$(IMAGES),docker rmi -f $(img);)

# clean versioned docker images
clean-version:
	$(foreach img,$(IMAGES),docker rmi -f $(img):$(VERSION);)
