ifeq (,$(shell command -v yaml2json 2>/dev/null))
$(error 'yaml2json is required; you could install it by 'pip install yaml2json')
endif
ifeq (,$(shell command -v json2yaml 2>/dev/null))
$(error 'json2yaml is required; you could install it by 'pip install json2yaml')
endif

ifndef post-manifests
define post-manifests
@yaml2json config/rbac/role.yaml | jq '.kind="Role"' > config/rbac/role.json \
&& json2yaml config/rbac/role.json > config/rbac/role.yaml \
&& rm -f config/rbac/role.json \
&& yaml2json config/rbac/role_binding.yaml | jq '.kind="RoleBinding" | .roleRef.kind="Role"' > config/rbac/role_binding.json \
&& json2yaml config/rbac/role_binding.json > config/rbac/role_binding.yaml \
&& rm -f config/rbac/role_binding.json
endef
endif
