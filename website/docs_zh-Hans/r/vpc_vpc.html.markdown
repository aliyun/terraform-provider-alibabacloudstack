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

## 参数参考

支持以下参数：

* `cidr_block` - (必填) VPC的网段。建议您使用192.168.0.0/16、172.16.0.0/12、10.0.0.0/8三个RFC标准私网网段及其子网作为专有网络的主IPv4网段，网段掩码有效范围为8~28位。您也可以使用除100.64.0.0/10、224.0.0.0/4、127.0.0.0/8或169.254.0.0/16及其子网外的自定义地址段作为专有网络的主IPv4网段。
* `vpc_name` - (可选) 要修改的VPC名称。名称长度为1～128个字符，不能以`http://`或`https://`开头。
* `description` - (可选) 要修改的VPC描述信息。描述长度为1～256个字符，不能以`http://`或`https://`开头。
* `dry_run` - (可选，变更时重建) 是否只预检此次请求。取值：
  * **true**：发送检查请求，不会查询资源状况。检查项包括AccessKey是否有效、RAM用户的授权情况和是否填写了必填参数。如果检查不通过，则返回对应错误。如果检查通过，会返回错误码`DryRunOperation`。
  * **false**(默认值)：发送正常请求，通过检查后返回2xx HTTP状态码并直接查询资源状况。
* `enable_ipv6` - (可选，变更时重建) 是否启用IPv6 CIDR块。取值：
  * **false**(默认值)：不开启。
  * **true**：开启。当此参数设置为`true`时，系统将自动为您创建一个免费版本的IPv6网关，并分配一个/56的IPv6网络段。
* `resource_group_id` - (可选) 需要移入云资源实例的资源组ID。资源组是在阿里云账号下进行资源分组管理的一种机制，资源组能够帮助您解决单个云账号内的资源分组和授权管理等复杂性问题。更多信息，请参见[什么是资源管理](https://help.aliyun.com/document_detail/94475.html)。
* `secondary_cidr_blocks` - (可选) 附加网段信息。注意：从provider版本1.185.0开始，该字段已被废弃，并将在未来的版本中移除。请改用新资源`alicloud_vpc_ipv4_cidr_block`。`secondary_cidr_blocks`属性和`alicloud_vpc_ipv4_cidr_block`资源不能同时使用。
* `user_cidrs` - (可选，变更时重建) 用户侧网络的网段，如需定义多个网段请使用半角逗号隔开，最多支持3个网段。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - VPC的ID。
* `router_id` - 在VPC创建时默认创建的路由器ID。
* `route_table_id` - 在VPC创建时默认创建的路由器的路由表ID。
* `ipv6_cidr_block` - 默认VPC的IPv6网段。仅当`enable_ipv6`设置为`true`时，此属性才可用。
* `resource_group_id` - VPC所属的资源组ID。
* `status` - VPC的状态。有效值：
  * `Pending`: 正在配置VPC。
  * `Available`: VPC可用。
* `vpc_name` - VPC的名称。