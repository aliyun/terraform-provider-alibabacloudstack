---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vpc"
sidebar_current: "docs-Alibabacloudstack-vpc-vpc"
description: |- 
  编排VPC实例
---

# alibabacloudstack_vpc_vpc
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpc`

使用Provider配置的凭证在指定的资源集编排VPC实例。

## 示例用法

### 基础用法

```hcl
variable "name" {
    default = "tf-testaccvpcvpc48306"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block           = "172.16.0.0/12"
  vpc_name             = var.name
  description         = "RDK更新"
  enable_ipv6         = true
  resource_group_id   = "rg-abc123xyz"
  secondary_cidr_blocks = ["192.168.0.0/16"]
  user_cidrs          = ["10.0.0.0/8"]
}
```

## 参数说明

支持以下参数：

* `cidr_block` - (必填, 变更时重建) VPC的网段。您可指定以下CIDR块或其子集作为VPC的主IPv4 CIDR块：
  * RFC文档定义的标准私有CIDR块：192.168.0.0/16、172.16.0.0/12 和 10.0.0.0/8。子网掩码长度必须在8到28位之间。
  * 除以下范围外的自定义CIDR块：100.64.0.0/10、224.0.0.0/4、127.0.0.0/8、169.254.0.0/16及其子网。
* `vpc_name` - (可选) VPC的名称。名称长度为1到128个字符，不能以`http://`或`https://`开头。默认值为null。
* `description` - (可选) VPC的描述信息。描述长度为1到256个字符，不能以`http://`或`https://`开头。默认值为null。
* `dry_run` - (可选, 变更时重建) 是否进行预检请求。有效值：
  * `true`：执行预检请求。系统会检查必需参数、请求语法和限制。如果请求未能通过预检，将返回错误消息；如果请求通过预检，将返回`DryRunOperation`错误代码。
  * `false`（默认值）：发送正常请求并执行操作。如果请求通过预检，将返回2xx HTTP状态码并执行操作。
* `enable_ipv6` - (可选, 变更时重建) 是否启用IPv6 CIDR块。有效值：
  * `false`（默认值）：禁用IPv6 CIDR块。
  * `true`：启用IPv6 CIDR块。当此参数设置为`true`时，系统将自动为您创建一个免费版本的IPv6网关，并分配一个/56的IPv6网络段。
* `resource_group_id` - (可选) 要移入云资源实例的资源组ID。资源组是阿里云账号下进行资源分组管理的一种机制，资源组能够帮助您解决单个云账号内的资源分组和授权管理等复杂性问题。更多信息，请参见[什么是资源管理](https://help.aliyun.com/document_detail/94475.html)。
* `secondary_cidr_blocks` - (可选) VPC的附加CIDR块列表。**注意**：从provider版本1.185.0开始，该字段已被废弃，并将在未来的版本中移除。请改用新资源`alicloud_vpc_ipv4_cidr_block`。`secondary_cidr_blocks`属性和`alicloud_vpc_ipv4_cidr_block`资源不能同时使用。
* `user_cidrs` - (可选, 变更时重建) 用户定义的CIDR列表。
* `status` - (可选) VPC的状态。有效值：
  * `Pending`：VPC正在配置中。
  * `Available`：VPC可用。
* `tags` - (可选) 要分配给资源的标签映射。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - VPC的ID。
* `router_id` - 在VPC创建时默认创建的路由器ID。
* `route_table_id` - 在VPC创建时默认创建的路由器的路由表ID。
* `ipv6_cidr_block` - VPC的IPv6 CIDR块。仅当`enable_ipv6`设置为`true`时，此属性才可用。
* `resource_group_id` - VPC所属的资源组ID。
* `status` - VPC的状态。有效值：
  * `Pending`：VPC正在配置中。
  * `Available`：VPC可用。
* `vpc_name` - VPC的名称。