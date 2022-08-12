# File: /Makefile
# Project: integration-operator
# File Created: 23-06-2021 09:14:26
# Author: Clay Risser <email@clayrisser.com>
# -----
# Last Modified: 12-08-2022 09:43:06
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


export MAKE_CACHE := $(shell pwd)/.make
export PARENT := true
include blackmagic.mk

NAME := integration-operator
REGISTRY := registry.gitlab.com/risserlabs/internal/integration-operator
VERSION := $(shell echo $(shell $(GIT) describe --abbrev=0 --tags 2>/dev/null || $(GIT) describe --tags) | sed 's| .*$$||g' | sed 's|^v||g')
IMAGE := $(REGISTRY)/$(NAME)

ACTIONS += install
$(ACTION)/install:
	@go get

ACTIONS += build~install
BUILD_DEPS := $(call deps,build,$(shell $(GIT) ls-files 2>$(NULL) | \
	grep -E "\.go$$"))
$(ACTION)/build:
	@$(MAKE) -s operator-framework-build

ACTIONS += start~install
$(ACTION)/start:
	@$(MAKE) -s operator-framework-run

docker-build:
	@$(MAKE) -f operator-framework.mk docker-build IMG="$(IMAGE):$(VERSION)"

docker-push:
	@$(MAKE) -f operator-framework.mk docker-push IMG="$(IMAGE):$(VERSION)"

operator-framework-%:
	@$(MAKE) -f operator-framework.mk $(shell echo $@ | sed "s/operator-framework-//")

.PHONY: lint-chart
lint-chart:
	@helm lint charts/integration-operator
	@helm install --debug --dry-run --generate-name charts/integration-operator

.PHONY: debug ~debug +debug
debug: ~debug
~debug: ~lint +debug
+debug:
	@helm install --debug --dry-run --generate-name $(CHART)

.PHONY: generate manifests install-crds uninstall-crds run build
generate: operator-framework-generate
install-crds: generate operator-framework-install
manifests: generate operator-framework-manifests
uninstall-crds: generate operator-framework-uninstall

.PHONY: clean
clean:
	-@$(call clean)
	-@$(GIT) clean -fXd $(NOFAIL)

.PHONY: purge
purge: clean

-include $(patsubst %,$(_ACTIONS)/%,$(ACTIONS))

+%:
	@$(MAKE) -e -s $(shell echo $@ | $(SED) 's/^\+//g')

%: ;

CACHE_ENVS += 
