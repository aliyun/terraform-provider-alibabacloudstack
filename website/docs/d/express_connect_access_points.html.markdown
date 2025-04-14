---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_express_connect_access_points"
sidebar_current: "docs-alibabacloudstack-datasource-express-connect-access-points"
description: |-
  Provides a list of Express Connect Access Points to the user.
---

# alibabacloudstack_express_connect_access_points
-> **NOTE:** Alias name has: `alibabacloudstack_expressconnect_accesspoints`

This data source provides the Express Connect Access Points of the current Alibaba Cloud user.


## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_express_connect_access_points" "ids" {
  ids = ["ap-cn-hangzhou-yh-C"]
}
output "express_connect_access_point_id_1" {
  value = data.alibabacloudstack_express_connect_access_points.ids.points.0.id
}

data "alibabacloudstack_express_connect_access_points" "nameRegex" {
  name_regex = "^杭州-"
}
output "express_connect_access_point_id_2" {
  value = data.alibabacloudstack_express_connect_access_points.nameRegex.points.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew)  A list of Access Point IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Access Point name.
* `status` - (Optional, ForceNew) The Physical Connection to Which the Access Point State. Valid values: `disabled`, `full`, `hot`, `recommended`.
* `names` - (Optional, ForceNew) A list of Access Point names.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Access Point names.
* `points` - A list of Express Connect Access Points. Each element contains the following attributes:
  * `access_point_id` - The Access Point ID.
  * `access_point_name` - Access Point Name.
  * `attached_region_no` - The Access Point Is Located an ID.
  * `description` - The Access Point Description.
  * `host_operator` - The Access Point Belongs to the Operator.
  * `id` - The ID of the Access Point.
  * `location` - The Location of the Access Point.
  * `status` - The Physical Connection to Which the Access Point State.
  * `type` - The Physical Connection to Which the Network Type.
  * `attached_region_no` - The region number where the access point is attached.
  * `description` - The description of the access point.
  * `host_operator` - The operator to which the access point belongs.
  * `location` - The location of the access point.