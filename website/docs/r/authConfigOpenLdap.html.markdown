---
layout: "rancher2"
page_title: "Rancher2: rancher2_auth_config_openldap"
sidebar_current: "docs-rancher2-auth-config-openldap"
description: |-
  Provides a Rancher v2 Auth Config OpenLdap resource. This can be used to configure and enable Auth Config OpenLdap for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_auth\_config\_openldap

Provides a Rancher v2 Auth Config OpenLdap resource. This can be used to configure and enable Auth Config OpenLdap for rancher v2 rke clusters and retrieve their information.

Beside local, just one auth config provider could be enabled at once.

## Example Usage

```hcl
# Create a new rancher2 Auth Config OpenLdap
resource "rancher2_auth_config_openldap" "openldap" {
  servers = ["<OPENLDAP_SERVER>"]
  service_account_distinguished_name = "<SERVICE_DN>"
  service_account_password = "<SERVICE_PASSWORD>"
  user_search_base = "<SEARCH_BASE>"
  port = <OPENLDAP_PORT>
}
```

## Argument Reference

The following arguments are supported:

* `servers` - (Required) Openldap servers list (list)
* `service_account_distinguished_name` - (Required/Sensitive) Service account DN for access openldap service (string)
* `service_account_password` - (Required/Sensitive) Service account password for access openldap service (string)
* `user_search_base` - (Required) User search base DN (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `openldap_user://<DN>`  `openldap_group://<DN>` (list)
* `certificate` - (Optional/Sensitive) CA certificate for TLS if selfsigned (string)
* `connection_timeout` - (Optional) Openldap connection timeout. Default `5000` (int)
* `enabled` - (Optional) Enable auth config provider. Default `true` (bool)
* `group_dn_attribute` - (Optional/Computed) Group DN attribute. Default `entryDN` (string)
* `group_member_mapping_attribute` - (Optional/Computed) Group member mapping attribute. Default `member` (string)
* `group_member_user_attribute` - (Optional/Computed) Group member user attribute. Default `entryDN` (string)
* `group_name_attribute` - (Optional/Computed) Group name attribute. Default `cn` (string)
* `group_object_class` - (Optional/Computed) Group object class. Default `groupOfNames` (string)
* `group_search_attribute` - (Optional/Computed) Group search attribute. Default `cn` (string)
* `group_search_base` - (Optional/Computed) Group search base (string)
* `nested_group_membership_enabled` - (Optional/Computed) Nested group membership enable. Default `false` (bool)
* `port` - (Optional) Openldap port. Default `389` (int)
* `user_disabled_bit_mask` - (Optional/Computed) User disabled bit mask (int)
* `user_enabled_attribute` - (Optional/Computed) User enable attribute (string)
* `user_login_attribute` - (Optional/Computed) User login attribute. Default `uid` (string)
* `user_member_attribute` - (Optional/Computed) User member attribute. Default `memberOf` (string)
* `user_name_attribute` - (Optional/Computed) User name attribute. Default `givenName` (string)
* `user_object_class` - (Optional/Computed) User object class. Default `inetorgperson` (string)
* `user_search_attribute` - (Optional/Computed) User search attribute. Default `uid|sn|givenName` (string)
* `tls` - (Optional/Computed) Enable TLS connection (bool)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)

