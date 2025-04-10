---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_natgateways"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-natgateways"
description: |- 
  查询专有网络的NAT网关实例
---

# alibabacloudstack_natgateway_natgateways
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nat_gateways`

根据指定过滤条件列出当前凭证权限可以访问的专有网络的NAT网关实例列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccNatGatewaysDatasource17110"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = alibabacloudstack_vpc.default.id
  specification = "Small"
  name          = var.name
  description   = "${var.name}_description"
}

data "alibabacloudstack_natgateway_natgateways" "default" {
  vpc_id     = alibabacloudstack_vpc.default.id
  name_regex = alibabacloudstack_nat_gateway.default.name
  ids        = [alibabacloudstack_nat_gateway.default.id]
}

output "nat_gateways" {
  value = data.alibabacloudstack_natgateway_natgateways.default.gateways
}
```

## 参数参考

以下参数是支持的：

* `vpc_id` - (可选，变更时重建) 所属的 VPC ID。如果指定，数据源将仅返回与此 VPC 关联的 NAT 网关。
* `ids` - (可选) NAT 网关 ID 列表。如果指定，数据源将仅返回与此 ID 匹配的结果。
* `name_regex` - (可选，变更时重建) 用于通过名称筛选 NAT 网关的正则表达式字符串。这允许基于命名约定进行更灵活的过滤。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 所有匹配的 NAT 网关的名称列表。
* `gateways` - NAT 网关列表。每个元素包含以下属性：
  * `id` - NAT 网关的 ID。
  * `name` - NAT 网关的名称。
  * `description` - NAT 网关的描述。
  * `creation_time` - NAT 网关创建的时间。
  * `spec` - NAT 网关的规格（例如，Small, Medium）。
  * `status` - NAT 网关的状态，取值范围及含义如下：
    - `Initiating`：初始化中
    - `Available`：可用
    - `Pending`：配置中
    - `Deleting`：删除中
  * `snat_table_id` - 与 NAT 网关关联的 SNAT 表的 ID。
  * `forward_table_id` - 与 NAT 网关关联的转发表的 ID。
  * `vpc_id` - NAT 网关所属的 VPC ID。