# File: /Makefile
# Project: docker
# File Created: 14-08-2022 14:26:41
# Author: Clay Risser <email@clayrisser.com>
# -----
# Last Modified: 23-06-2023 15:27:46
# Modified By: Clay Risser <email@clayrisser.com>
# -----
# Risser Labs LLC (c) Copyright 2021 - 2022

include ../mkpm.mk
ifneq (,$(MKPM_READY))
include $(MKPM)/gnu

export DOCKER_COMPOSE := docker-compose
export CONTEXT := ..
export NAME := integration-operator
export REGISTRY := registry.gitlab.com/bitspur/internal
export VERSION := 0.1.3

include $(MKPM)/docker

endif
