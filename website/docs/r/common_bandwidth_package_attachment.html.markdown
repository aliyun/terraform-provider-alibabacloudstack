---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_common_bandwidth_package_attachment"
sidebar_current: "docs-apsarastack-resource-common-bandwidth-package-attachment"
description: |-
  Provides an Apsarastack Common  Attachment resource.
---

# apsarastack\_common\_bandwidth\_package\_attachment

Provides an apsarastack Common Bandwidth Package Attachment resource for associating Common Bandwidth Package to EIP Instance.

-> **NOTE:** Terraform will auto build common bandwidth package attachment while it uses `apsarastack_common_bandwidth_package_attachment` to build a common bandwidth package attachment resource.


## Example Usage

Basic Usage

```
resource "apsarastack_common_bandwidth_package" "foo" {
  bandwidth   = "2"
  name        = "test_common_bandwidth_package"
  description = "test_common_bandwidth_package"
}

resource "apsarastack_eip" "foo" {
  bandwidth            = "2"
}

resource "apsarastack_common_bandwidth_package_attachment" "foo" {
  bandwidth_package_id = "${apsarastack_common_bandwidth_package.foo.id}"
  instance_id          = "${apsarastack_eip.foo.id}"
}

```
## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Required, ForceNew) The bandwidth_package_id of the common bandwidth package attachment, the field can't be changed.
* `instance_id` - (Required, ForceNew) The instance_id of the common bandwidth package attachment, the field can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package attachment id and formates as `<bandwidth_package_id>:<instance_id>`.

