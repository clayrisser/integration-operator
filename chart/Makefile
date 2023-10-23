include $(MKPM)/mkpm
include $(MKPM)/gnu

export DOCKER_DEFAULT_PLATFORM=linux/amd64
export PATCHES :=

include $(MKPM)/patch

.PHONY: prepack
prepack: patch-apply
