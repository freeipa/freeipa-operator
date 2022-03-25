# When running in github actions the env var CI is set to true
# Here we set customizations when running Makefile on it
ifeq (true,$(CI))
$(info For Git Hub Actions the test for the controller are disabled)
# Disable the controller tests
TEST_DISABLE_LIST += controller
endif

# When running in ci-operator the env var OPENSHIFT_CI is set
# to true. Here we set customizations when running Makefile
# on it.
ifeq (true,$(OPENSHIFT_CI))
$(warning This should be enabled in the future for the controller)
# Disable the controller tests
TEST_DISABLE+=controller
export TEST_DISABLE

# At 'vendor/sigs.k8s.io/controller-runtime/pkg/internal/testing/addr/manager.go:51'
# we can read the block:
#
# func init() {
#     baseDir, err := os.UserCacheDir()
#     if err != nil {
#         baseDir = os.TempDir()
#     }
#     cacheDir = filepath.Join(baseDir, "kubebuilder-envtest")
#     if err := os.MkdirAll(cacheDir, 0750); err != nil {
#         panic(err)
#     }
# }
#
# If it is not specified when the test library is initialized
# it tries to create the '/.cache' directory, evoking the
# 'make test' rule fails because no permissions to create it
XDG_CACHE_HOME:=/tmp/.cache
export XDG_CACHE_HOME

# Indicate to the cluster to use an existing cluster so the
# functional tests does not try to instanciate a
# Kubernetes API
USE_EXISTING_CLUSTER=1
export USE_EXISTING_CLUSTER

HOME=/tmp
export HOME

# FIXME At the moment of writing this lines, KUBECONFIG variable
#       was empty in ci-operator job, so there are no access to
#       the cluster for e2e tests and the controller tests.
endif

# Print out message to warn about the fact that some tests
# are not being executed
ifneq (,$(TEST_DISABLE))
$(info The following tests are disabled: $(TEST_DISABLE))
endif