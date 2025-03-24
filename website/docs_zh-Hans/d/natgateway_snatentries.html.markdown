---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_snatentries"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-snatentries"
description: |- 
  查询专有网络的NAT网关SNAT表规则
---

# alibabacloudstack_natgateway_snatentries
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_snat_entries`

根据指定过滤条件列出当前凭证权限可以访问的专有网络的NAT网关SNAT表规则列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccForSnatEntriesDatasource"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = "${alibabacloudstack_vpc.default.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
  name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = "${alibabacloudstack_eip.default.id}"
  instance_id   = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_snat_entry" "default" {
  snat_table_id     = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  snat_ip           = "${alibabacloudstack_eip.default.ip_address}"
}

data "alibabacloudstack_natgateway_snatentries" "default" {
  snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
  source_cidr   = "172.16.0.0/21"
  ids           = ["${alibabacloudstack_snat_entry.default.id}"]
  output_file   = "snat_entries_output.txt"
}

output "snat_entries" {
  value = "${data.alibabacloudstack_natgateway_snatentries.default.entries}"
}
```

## 参数参考

以下参数是支持的：

* `snat_table_id` - (必填，变更时重建) SNAT条目所属的SNAT表ID。
* `source_cidr` - (选填)SNAT条目的源网段。
* `ids` - (选填)SNAT条目ID的列表。这可以用于过滤结果。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - SNAT条目ID列表。
* `entries` - SNAT条目列表。每个元素包含以下属性：
  * `id` - SNAT条目的ID。
  * `snat_ip` - SNAT条目的公网IP。
  * `source_cidr` - SNAT条目的源网段。
  * `status` - SNAT条目的状态。可能的值包括 `available`、`pending` 和 `inactive`。