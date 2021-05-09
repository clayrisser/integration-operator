NAME := integration-operator
REGISTRY := codejamninja
VERSION := 0.0.1
IMAGE := $(REGISTRY)/$(NAME)

docker-build:
	@$(MAKE) -f operator-framework.mk docker-build IMG="$(IMAGE):$(VERSION)"

docker-push:
	@$(MAKE) -f operator-framework.mk docker-push IMG="$(IMAGE):$(VERSION)"

operator-framework-%:
	@$(MAKE) -f operator-framework.mk $(shell echo $@ | sed "s/operator-framework-//")

.PHONY: generate manifests install-crds uninstall-crds run build
generate: operator-framework-generate
install-crds: generate operator-framework-install
manifests: generate operator-framework-manifests
run: operator-framework-run
uninstall-crds: generate operator-framework-uninstall
build: operator-framework-build
