---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_natgateway"
sidebar_current: "docs-Alibabacloudstack-natgateway-natgateway"
description: |- 
  编排专有网络的NAT网关实例
---

# alibabacloudstack_natgateway_natgateway
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nat_gateway`

使用Provider配置的凭证在指定的资源集编排专有网络的NAT网关实例。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccNatGatewayConfig13663"
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
  vpc_id             = "${alibabacloudstack_vswitch.default.vpc_id}"
  specification     = "Small"
  nat_gateway_name  = "${var.name}"
  description       = "This is a test NAT Gateway"
  bandwidth_packages = [
    {
      ip_count         = 2
      bandwidth        = 10
      zone            = "${data.alibabacloudstack_zones.default.zones.0.id}"
      public_ip_addresses = ["10.0.0.1", "10.0.0.2"]
    }
  ]
  tags = {
    CreatedBy = "Terraform"
    Env       = "Test"
  }
}
```

## 参数参考

支持以下参数：
  * `vpc_id` - (必填, 变更时重建) - 所属的VPC ID。
  * `specification` - (选填) - NAT网关的规格。有效值为 `Small`、`Middle` 和 `Large`。默认为 `Small`。
  * `spec` - (选填) - NAT网关的规格，与`specification`功能相同。
  * `nat_gateway_name` - (选填) - NAT网关的名称。该值可以是一个包含 2 到 128 个字符的字符串，必须仅包含字母数字字符或连字符，例如“-”、“.”、“_”，并且不能以连字符开头或结尾，也不能以 http:// 或 https:// 开头。默认为 null。
  * `name` - (选填) - NAT网关的名称，与`nat_gateway_name`功能相同。
  * `description` - (选填) - NAT网关的描述。此描述可以有一个包含 2 到 256 个字符的字符串，并且它不能以 http:// 或 https:// 开头。默认为 null。
  * `bandwidth_packages` - (选填) - NAT网关的带宽包列表。仅支持在 2017 年 11 月 4 日 00:00 之前创建的 NAT 网关。
    * `ip_count` - (必填) - 当前带宽包中的 IP 地址数量。其值范围是从 1 到 50。
    * `bandwidth` - (必填) - 当前带宽包的带宽值。其值范围是从 5 到 5000 Mbps。
    * `zone` - (选填) - 当前带宽包的可用区。如果未指定此值，Terraform 将设置一个随机 AZ。
    * `public_ip_addresses` - (选填) - 带宽包的公共 IP 地址。公共 IP 的数量等于 `ip_count`，多个 IP 以逗号分隔，例如“10.0.0.1,10.0.0.2”。
  * `tags` - (选填, Map) - 要分配给资源的标签映射。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - NAT网关的ID。
  * `specification` - NAT网关的规格。
  * `spec` - NAT网关的规格，与`specification`功能相同。
  * `name` - NAT网关的名称，与`nat_gateway_name`功能相同。
  * `nat_gateway_name` - 实例名称。
  * `description` - NAT网关的描述。
  * `vpc_id` - NAT网关的VPC ID。
  * `bandwidth_package_ids` - 带宽包ID的列表，以逗号分隔。
  * `snat_table_ids` - NAT网关自动创建的SNAT表的ID。
  * `forward_table_ids` - NAT网关自动创建的目的地网络地址转换(DNAT)表的ID。
  * `bandwidth_packages` - 与NAT网关关联的带宽包详细信息。