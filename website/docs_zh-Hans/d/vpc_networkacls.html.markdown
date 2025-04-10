---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_networkacls"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-networkacls"
description: |- 
  查询专有网络（VPC）网络访问许可列表
---

# alibabacloudstack_vpc_networkacls
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_network_acls`

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）网络访问许可列表。

## 示例用法

以下示例展示了如何使用 `alibabacloudstack_vpc_networkacls` 数据源来查询符合条件的网络 ACL 列表，并输出第一个网络 ACL 的 ID。

```hcl
variable "name" {	
	default = "tf-testAccNetworkAcl-18237"
}

resource "alibabacloudstack_vpc" "default" {
	name = var.name
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
	description = var.name
	network_acl_name = var.name
	vpc_id = alibabacloudstack_vpc.default.id
}

data "alibabacloudstack_vpc_networkacls" "example" {
  ids        = [alibabacloudstack_network_acl.default.id]
  name_regex = "tf-testAccNetworkAcl-18237"
  network_acl_name = "tf-testAccNetworkAcl-18237"
  status     = "Available"
  vpc_id     = alibabacloudstack_vpc.default.id
  resource_type = "NETWORKACL"
}

output "first_network_acl_id" {
  value = data.alibabacloudstack_vpc_networkacls.example.acls.0.id
}
```

## 参数说明

以下参数是支持的：

* `ids` - (选填, 变更时重建) 网络 ACL ID 列表。如果指定，结果将按这些 ID 进行过滤。
* `name_regex` - (选填, 变更时重建) 用于通过网络 ACL 名称筛选结果的正则表达式字符串。
* `network_acl_name` - (选填, 变更时重建) 网络 ACL 的名称。名称长度为 1～128 个字符，不能以 `http://` 或 `https://` 开头。
* `resource_id` - (选填, 变更时重建) 关联资源的 ID。如果指定了 `resource_type`，这是必填的。
* `resource_type` - (选填, 变更时重建) 关联资源的类型。有效值：`NETWORKACL`。`resource_type` 和 `resource_id` 需要同时指定才能生效。
* `status` - (选填, 变更时重建) 网络 ACL 的状态。有效值：`Available`、`Modifying`。
* `vpc_id` - (选填, 变更时重建) 关联 VPC 的 ID。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 网络 ACL 名称列表。
* `acls` - 网络 ACL 列表。每个元素包含以下属性：
  * `description` - 网络 ACL 的描述信息。描述长度为 1～256 个字符，不能以 `http://` 或 `https://` 开头。
  * `egress_acl_entries` - 出方向规则信息。每个条目包含：
    * `description` - 出方向规则的描述。
    * `destination_cidr_ip` - 目标 CIDR 块。
    * `network_acl_entry_name` - 出方向规则条目的名称。
    * `policy` - 授权策略（例如，`Allow`、`Deny`）。
    * `port` - 目标端口范围。
    * `protocol` - 传输层协议（例如，`tcp`、`udp`）。
  * `id` - 网络 ACL 的 ID。
  * `ingress_acl_entries` - 入方向规则信息。每个条目包含：
    * `source_cidr_ip` - 源 CIDR 块。
    * `description` - 入方向规则的描述。
    * `network_acl_entry_name` - 入方向规则条目的名称。
    * `policy` - 授权策略（例如，`Allow`、`Deny`）。
    * `port` - 源端口范围。
    * `protocol` - 传输层协议（例如，`tcp`、`udp`）。
  * `network_acl_id` - 网络 ACL 的 ID。
  * `network_acl_name` - 网络 ACL 的名称。名称长度为 1～128 个字符，不能以 `http://` 或 `https://` 开头。
  * `resources` - 关联的资源。每个条目包含：
    * `resource_id` - 关联资源的 ID。
    * `resource_type` - 关联资源的类型。
    * `status` - 关联资源的状态。
  * `status` - 网络 ACL 的状态。
  * `vpc_id` - 关联 VPC 的 ID。