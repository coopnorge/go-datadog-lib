# vim: set noexpandtab fo-=t:
# https://www.gnu.org/software/make/manual/make.html
.PHONY: default
default:

########################################################################
# boiler plate
########################################################################
SHELL=bash
current_makefile:=$(lastword $(MAKEFILE_LIST))
current_makefile_dirname:=$(dir $(current_makefile))
current_makefile_dirname_abspath:=$(dir $(abspath $(current_makefile)))
current_makefile_dirname_realpath:=$(dir $(realpath $(current_makefile)))

ifneq ($(filter all vars,$(VERBOSE)),)
dump_var=$(info var $(1)=$($(1)))
dump_vars=$(foreach var,$(1),$(call dump_var,$(var)))
else
dump_var=
dump_vars=
endif

ifneq ($(filter all targets,$(VERBOSE)),)
__ORIGINAL_SHELL:=$(SHELL)
SHELL=$(warning Building $@$(if $<, (from $<))$(if $?, ($? newer)))$(TIME) $(__ORIGINAL_SHELL)
endif

define __newline


endef


skip=
# skipable makes the targets passed to it skipable with skip=foo%
# $(1): targets that should be skipable
skipable=$(filter-out $(skip),$(1))

go_install=GOFLAGS= go install $(if $(filter all commands,$(VERBOSE)),-v) $(go_install_flags)

########################################################################
# variables
########################################################################
GOPATH:=$(shell type go >/dev/null 2>&1 && go env GOPATH)
export PATH:=$(if $(GOPATH),$(GOPATH)/bin:,)$(PATH)

########################################################################
# targets
########################################################################

.PHONY: toolchain
toolchain:
	$(go_install) github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0

.PHONY: toolchain-update
toolchain-update: go_get_flags+=-u
toolchain-update: toolchain

.PHONY: fetch-deps
fetch-deps:
	go mod vendor

.PHONY: validate-static
validate-static:
	golangci-lint run $(if $(filter all commands,$(VERBOSE)),-v) ./...

.PHONY: validate-fix
validate-fix:
	golangci-lint run $(if $(filter all commands,$(VERBOSE)),-v) --fix ./...

.PHONY: test
test:
	go test -cover -race \
		-coverprofile=coverage.out -covermode=atomic \
		$(if $(filter all commands,$(VERBOSE)),-v) \
		$(if $(gotest_files),$(gotest_files),./...) \
		$(gotest_args) \

.PHONY: view-coverage
view-coverage:
	go tool cover -html=coverage.out

.PHONY: validate-dynamic
validate-dynamic: test

.PHONY: validate
validate: validate-static validate-dynamic

