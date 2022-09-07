---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_common_bandwidth_packages"
sidebar_current: "docs-alibabacloudstack-datasource-common-bandwidth-packages"
description: |-
    Provides a list of Common Bandwidth Packages owned by an Alibabacloudstack Cloud account.
---

# alibabacloudstack\_common\_bandwidth\_packages

This data source provides a list of Common Bandwidth Packages owned by an Alibabacloudstack Cloud account.


## Example Usage

```

resource "alibabacloudstack_common_bandwidth_package" "foo" {
  bandwidth   = "2"
  name        = "tf-testAccCommonBandwidthPackage"
  description = "tf-testAcc-CommonBandwidthPackage"
}

data "alibabacloudstack_common_bandwidth_packages" "foo" {
  name_regex = "^tf-testAcc.*"
  ids        = ["${alibabacloudstack_common_bandwidth_package.foo.id}"]
}

output "common_bandwidth_packages" {
  value = "${data.alibabacloudstack_common_bandwidth_packages.foo.packages}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Common Bandwidth Packages IDs.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Common Bandwidth Packages IDs.
* `names` - A list of Common Bandwidth Packages names.
* `packages` - A list of Common Bandwidth Packages. Each element contains the following attributes:
  * `id` - ID of the Common Bandwidth Package.
  * `bandwidth` - The peak bandwidth of the Internet Shared Bandwidth instance.
  * `status` - Status of the Common Bandwidth Package.
  * `name` - Name of the Common Bandwidth Package.
  * `description` - The description of the Common Bandwidth Package instance.
  * `business_status` - The business status of the Common Bandwidth Package instance.
  * `isp` - ISP of the Common Bandwidth Package.
  * `creation_time` - Time of creation.
  * `public_ip_addresses` - Public ip addresses that in the Common Bandwidth Package.
 
## Public ip addresses Block
  
  The public ip addresses mapping supports the following:
  
  * `ip_address`   - The address of the EIP.
  * `allocation_id` - The ID of the EIP instance.
