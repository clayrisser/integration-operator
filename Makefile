# File: /Makefile
# Project: integration-operator
# File Created: 14-08-2022 14:24:57
# Author: Clay Risser <email@clayrisser.com>
# -----
# Last Modified: 14-08-2022 14:50:39
# Modified By: Clay Risser <email@clayrisser.com>
# -----
# Risser Labs LLC (c) Copyright 2021 - 2022

include mkpm.mk
ifneq (,$(MKPM_READY))
include $(MKPM)/gnu

.PHONY: of-% build generate manifests install uninstall start
build: of-build
start: of-run
generate: of-generate
install: of-install
manifests: generate of-manifests
uninstall: of-uninstall
of-%:
	@$(MAKE) -s -f ./operator-framework.mk $(subst of-,,$@)

.PHONY: docker/%
docker/%:
	@$(MAKE) -s -C docker $(subst docker/,,$@)

endif
