# File: /Makefile
# Project: docker
# File Created: 14-08-2022 14:26:41
# Author: Clay Risser <email@clayrisser.com>
# -----
# Last Modified: 02-07-2023 11:49:19
# Modified By: Clay Risser <email@clayrisser.com>
# -----
# BitSpur (c) Copyright 2021 - 2022

include ../mkpm.mk
ifneq (,$(MKPM_READY))
include $(MKPM)/gnu

export DOCKER_COMPOSE := docker-compose
export CONTEXT := ..
export NAME := integration-operator
export REGISTRY := registry.gitlab.com/bitspur/rock8s
export VERSION := 0.2.0

include $(MKPM)/docker

endif
