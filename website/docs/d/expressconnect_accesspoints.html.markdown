---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_express_connect_access_points"
sidebar_current: "docs-Alibabacloudstack-datasource-express-connect-access-points"
description: |-
  Provides a list of expressconnect accesspoints owned by an alibabacloudstack account.
---

# alibabacloudstack_express_connect_access_points
-> **NOTE:** Alias name has: `alibabacloudstack_expressconnect_accesspoints`

This data source provides a list of express connect access points in an alibabacloudstack account according to the specified filters.

## Example Usage
```
data "alibabacloudstack_express_connect_access_points" "default" {
  ids = ["accessPointId1", "accessPointId2"]
  name_regex = "accessPointName"
  status = "normal"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional, ForceNew) - A list of Access Point IDs to filter by.
  * `name_regex` - (Optional, ForceNew) - A regex pattern to filter Access Points by name.
  * `status` - (Optional, ForceNew) - The status of the resource
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `points` - A list of Access Points. Each element in the list is a map with the following keys:
    * `id` - The ID of the Access Point.
    * `access_point_id` - The Access Point ID
    * `access_point_name` - Access Point Name
    * `attached_region_no` - The Access Point Is Located an ID
    * `description` - The Access Point Description
    * `host_operator` - The Access Point Belongs to the Operator
    * `location` - The Location of the Access Point
    * `status` - The status of the resource
    * `type` - The Physical Connection to Which the Network Type
