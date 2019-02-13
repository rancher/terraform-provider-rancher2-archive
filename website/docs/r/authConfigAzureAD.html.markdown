---
layout: "rancher2"
page_title: "Rancher2: rancher2_auth_config_azuread"
sidebar_current: "docs-rancher2-auth-config-azuread"
description: |-
  Provides a Rancher v2 Auth Config AzureAD resource. This can be used to configure and enable Auth Config AzureAD for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_auth\_config\_azuread

Provides a Rancher v2 Auth Config AzureAD resource. This can be used to configure and enable Auth Config AzureAD for rancher v2 rke clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Auth Config AzureAD
resource "rancher2_auth_config_azuread" "azuread" {
  application_id = "<AZUREAD_APP_ID>"
  application_secret = "<AZUREAD_APP_SECRET>"
  auth_endpoint = "<AZUREAD_AUTH_ENDPOINT>"
  code = "<AZUREAD_AUTH_CODE>"
  graph_endpoint = "<AZUREAD_GRAPH_ENDPOINT>"
  rancher_url = "<RANCHER_URL>"
  tenant_id = "<AZUREAD_TENANT_ID>"
  token_endpoint = "<AZUREAD_TOKEN_ENDPOINT>"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) AzureAD auth application ID (string).
* `application_secret` - (Required) AzureAD auth application secret (string).
* `auth_endpoint` - (Required) AzureAD auth endpoint (string).
* `code` - (Required) AzureAD auth code. Generated from `https://login.microsoftonline.com/<TENANT_ID>/oauth2/authorize?client_id=<APPLICATION_ID>&response_type=code` (string).
* `graph_endpoint` - (Required) AzureAD graph endpoint (string).
* `rancher_url` - (Required) Rancher URL (string).
* `tenant_id` - (Required) AzureAD tenant ID (string).
* `token_endpoint` - (Required) AzureAD token endpoint (string).
* `endpoint` - (Optional) AzureAD endpoint. Default `https://login.microsoftonline.com/`.
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted`
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `azuread_user://<USER_ID>`  `azuread_group://<GROUP_ID>`
* `enabled` - (Optional) Enable auth config provider. Default `true`.
* `tls` - (Optional) Enable TLS connection. Default `true`.
* `annotations` - (Optional/Computed) Annotations of the resource (map).
* `labels` - (Optional/Computed) Labels of the resource (map).
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `name` - (Computed) The name of the resource.
* `type` - (Computed) The type of the resource.

