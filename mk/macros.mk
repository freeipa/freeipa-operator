ifndef post-manifests
define post-manifests
@[ -e .venv ] || $(MAKE) .venv
@(source .venv/bin/activate; yaml2json config/rbac/role.yaml) | jq '.kind="Role"' > config/rbac/role.json \
&& (source .venv/bin/activate; json2yaml config/rbac/role.json > config/rbac/role.yaml) \
&& rm -f config/rbac/role.json \
&& (source .venv/bin/activate; yaml2json config/rbac/role_binding.yaml) | jq '.kind="RoleBinding" | .roleRef.kind="Role"' > config/rbac/role_binding.json \
&& (source .venv/bin/activate; json2yaml config/rbac/role_binding.json) > config/rbac/role_binding.yaml \
&& rm -f config/rbac/role_binding.json
endef
endif

# $1 expected value
# $2 resource-reference
# $3 path to the value
ifndef kubectl-wait-for-value
define kubectl-wait-for-value
@while test "$1" != "$(shell kubectl get "$2" -o jsonpath='{$3}' 2>/dev/null)"; do sleep 1; done
endef
endif
