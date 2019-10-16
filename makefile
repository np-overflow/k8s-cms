#
# k8s-cms
# Project Makefile
# 

# vars
VERSION:=0.3.0b

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
.PHONY: all push clean clean-version
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

push/%: %
	docker push $<


# cleans docker images
# clean all docker images
clean: clean-version
	$(foreach img,$(IMAGES),docker rmi -f $(img);)

# clean versioned docker images
clean-version:
	$(foreach img,$(IMAGES),docker rmi -f $(img):$(VERSION);)
