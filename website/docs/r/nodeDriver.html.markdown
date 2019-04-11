---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_driver"
sidebar_current: "docs-rancher2-resource-node_driver"
description: |-
  Provides a Rancher v2 Node Driver resource. This can be used to create Node Driver for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_node\_driver

Provides a Rancher v2 Node Driver resource. This can be used to create Node Driver for rancher v2 rke clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Node Driver
resource "rancher2_node_driver" "foo" {
    active = true
    builtin = false
    checksum = "0x0"
    description = "Foo description"
    external_id = "foo_external"
    name = "foo"
    ui_url = "local://ui"
    url = "local://"
    whitelist_domains = ["*.foo.com"]
}
```

## Argument Reference

The following arguments are supported:

* `active` - (Required) Specify if the node driver state (bool)
* `builtin` - (Required) Specify wheter the node driver is an internal node driver or not (bool)
* `checksum` - (Optional) Verify that the downloaded driver matches the expected checksum (string)
* `description` - (Optional) Description of the node driver (string)
* `external_id` - (Optional) External ID (string)
* `name` - (Required) Name of the node driver (string)
* `ui_url` - (Optional) The URL to load for customized Add Nodes screen for this driver (string)
* `url` - (Required) The URL to download the machine driver binary for 64-bit Linux (string)
* `whitelist_domains` - (Optional) Domains to whitelist for the ui (list)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_node_driver` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating node drivers.
- `update` - (Default `10 minutes`) Used for node driver modifications.
- `delete` - (Default `10 minutes`) Used for deleting node drivers.

## Import

Node Driver can be imported using the rancher Node Driver ID

```
$ terraform import rancher2_node_driver.foo <node_driver_id>
```
