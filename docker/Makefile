# File: /docker/Makefile
# Project: integration-operator
# File Created: 13-10-2023 15:44:20
# Author: Clay Risser
# -----
# BitSpur (c) Copyright 2021 - 2023
#
# Licensed under the GNU Affero General Public License (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.gnu.org/licenses/agpl-3.0.en.html
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# You can be released from the requirements of the license by purchasing
# a commercial license. Buying such a license is mandatory as soon as you
# develop commercial activities involving this software without disclosing
# the source code of your own applications.

include $(MKPM)/mkpm
include $(MKPM)/gnu

export DOCKER_COMPOSE := docker-compose
export CONTEXT := ..
export NAME := bitspur/rock8s/integration-operator
export REGISTRY := registry.gitlab.com
export VERSION := 1.2.0
export GIT_COMMIT ?= $(shell git describe --tags --always --dirty 2>/dev/null)

include $(MKPM)/docker

BUILDX ?= $(call ternary,$(WHICH) buildx,buildx,$(call ternary,$(WHICH) docker-buildx,docker-buildx,$(DOCKER) buildx))

.PHONY: bake
bake:
	@$(BUILDX) bake --provenance false
