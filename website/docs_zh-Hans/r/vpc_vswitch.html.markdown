---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vswitch"
sidebar_current: "docs-Alibabacloudstack-vpc-vswitch"
description: |- 
  编排VPC虚拟交换机
---

# alibabacloudstack_vpc_vswitch
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vswitch`

使用Provider配置的凭证在指定的资源集编排VPC虚拟交换机。

## 示例用法

### 基础用法

```hcl
variable "name" {
    default = "tf-testaccvpcvswitch97984"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  description     = "modify_description"
  vswitch_name   = "tf-testaccvpcvswitch97984"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block     = "172.16.0.0/24"
  enable_ipv6    = true
}
```

启用IPv6用法

```hcl
resource "alibabacloudstack_vpc_vpc" "ipv6_example" {
  vpc_name       = "ipv6_vpc"
  cidr_block     = "192.168.0.0/16"
  enable_ipv6    = true
}

resource "alibabacloudstack_vpc_vswitch" "ipv6_vswitch" {
  vswitch_name   = "ipv6_vswitch"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.ipv6_example.id}"
  cidr_block     = "192.168.0.0/24"
  enable_ipv6    = true
}
```

## 参数参考

支持以下参数：

* `zone_id` - (必填) 交换机所属的可用区ID。您可以通过调用 [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/36064.html) 操作查询最新的可用区列表。
* `vpc_id` - (必填，变更时重建) 交换机所属的虚拟私有云(VPC)的ID。
* `cidr_block` - (必填，变更时重建) 交换机的CIDR块。交换机网段要求如下：
  * 交换机的网段的掩码长度范围为16～29位。
  * 交换机的网段必须从属于所在VPC的网段。
  * 交换机的网段不能与所在VPC中路由条目的目标网段相同，但可以是目标网段的子集。
  * 交换机的网段不能是100.64.0.0/10及其子网网段。
* `enable_ipv6` - (可选，变更时重建) 指定是否启用交换机IPv6 CIDR块。有效值：
  * `false`(默认)：禁用IPv6 CIDR块。
  * `true`：启用IPv6 CIDR块。如果 `enable_ipv6` 为 `true`，则通过 `vpc_id` 指向的VPC也必须启用IPv6。系统将自动为您创建免费版本的IPv6网关，并分配一个/56的IPv6网络段。
* `vswitch_name` - (可选) 交换机的名称。默认为null。
* `description` - (可选) 交换机的描述。描述必须是1到256个字符的长度，并且不能以 `http://` 或 `https://` 开头。
* `tags` - (可选，映射)交换机的标签。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 交换机的ID。
* `availability_zone` - 交换机的可用区。
* `zone_id` - 交换机所属的可用区ID。
* `cidr_block` - 交换机的CIDR块。
* `ipv6_cidr_block` - 交换机的IPv6 CIDR块。
* `vpc_id` - VPC ID。
* `vswitch_name` - 交换机的名称。
* `description` - 交换机的描述。