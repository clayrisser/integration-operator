# File: /Makefile
# Project: integration-operator
# File Created: 23-06-2021 09:14:26
# Author: Clay Risser <email@clayrisser.com>
# -----
# Last Modified: 27-06-2021 02:30:51
# Modified By: Clay Risser <email@clayrisser.com>
# -----
# Silicon Hills LLC (c) Copyright 2021
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

NAME := integration-operator
REGISTRY := codejamninja
VERSION := 0.0.1
IMAGE := $(REGISTRY)/$(NAME)

.PHONY: install
install:
	@go get

docker-build:
	@$(MAKE) -f operator-framework.mk docker-build IMG="$(IMAGE):$(VERSION)"

docker-push:
	@$(MAKE) -f operator-framework.mk docker-push IMG="$(IMAGE):$(VERSION)"

operator-framework-%:
	@$(MAKE) -f operator-framework.mk $(shell echo $@ | sed "s/operator-framework-//")

.PHONY: generate manifests install-crds uninstall-crds run build
build: operator-framework-build
generate: operator-framework-generate
install-crds: generate operator-framework-install
manifests: generate operator-framework-manifests
start: operator-framework-run
uninstall-crds: generate operator-framework-uninstall
