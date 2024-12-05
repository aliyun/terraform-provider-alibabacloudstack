---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_instance"
sidebar_current: "docs-alibabacloudstack-resource-gpdb-instance"
description: |-
  Provides a AnalyticDB for PostgreSQL instance resource.
---

# alibabacloudstack\_gpdb\_instance

Provides a AnalyticDB for PostgreSQL instance resource supports replica set instances only. the AnalyticDB for PostgreSQL provides stable, reliable, and automatic scalable database services. 



## Example Usage

### Create a Gpdb instance

```
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "Gpdb"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  name              = "vpc-123456"
}

resource "alibabacloudstack_gpdb_instance" "example" {
  description          = "tf-gpdb-test"
  engine               = "gpdb"
  engine_version       = "4.3"
  instance_class       = "gpdb.group.segsdx2"
  instance_group_count = "2"
  vswitch_id           = alibabacloudstack_vswitch.default.id
  security_ip_list     = ["10.168.1.12", "100.69.7.112"]
}
```

## Argument Reference

The following arguments are supported:

* `engine` (Required, ForceNew) Database engine: gpdb. System Default value: gpdb.
* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/86908.htm) `EngineVersion`.
* `instance_class` - (Required) Instance specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/86942.htm).
* `instance_group_count` - (Required) The number of groups. Valid values: [2,4,8,16,32]
* `instance_inner_connection` - (Optional, ForceNew)The endpoint of the cluster.
* `instance_inner_port` - (Optional, ForceNew)The endpoint's port of the cluster.
* `instance_vpc_id` - (Optional, ForceNew)The Vpc id.
* `description` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`.
* `availability_zone` - (Optional, ForceNew) The Zone to launch the DB instance. it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alibabacloudstack_zones`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Instance.
* `instance_id` - Alias of key `id`.
* `status` - The status of the instance.
* `instance_network_type` - Classic network or VPC.
* `region_id` - Region ID the instance belongs to.


## Import

AnalyticDB for PostgreSQL can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_gpdb_instance.example gp-bp1291daeda44194
```
