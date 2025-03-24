---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6gateway"
sidebar_current: "docs-Alibabacloudstack-vpc-ipv6gateway"
description: |- 
  编排VPC的IPv6网关。
---

# alibabacloudstack_vpc_ipv6gateway
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpc_ipv6_gateway`

使用Provider配置的凭证在指定的资源集编排VPC的IPv6网关。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testaccvpcipv6gateway88979"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
}

resource "alibabacloudstack_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = var.name
  vpc_id            = alibabacloudstack_vpc.default.id
  description       = var.name
}
```

## 参数参考

支持以下参数：

* `description` - (可选) IPv6 网关的描述。描述必须是 **2 到 256** 个字符长度，不能以 `http://` 或 `https://` 开头。
* `ipv6_gateway_name` - (可选) IPv6 网关的名称。名称必须是 **2 到 128** 个字符长度，可以包含字母、数字、下划线 (`_`) 和连字符 (`-`)。名称必须以字母开头，但不能以 `http://` 或 `https://` 开头。
* `spec` - (可选) IPv6 网关的版本。尽管该参数在早期版本中用于区分规格，但现在已不再使用，所有 IPv6 网关均采用统一规格。有效值包括：`Small`(默认值)，表示免费版；`Medium`，表示企业版；`Large`，表示增强企业版。需要注意的是，虽然参数仍然存在，但其实际功能已被废弃。
* `vpc_id` - (必填，变更时重建) 要为其创建 IPv6 网关的虚拟私有云 (VPC) 的 ID。

### 超时时间

`timeouts` 块允许您为某些操作指定 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为 1 分钟)用于创建 IPv6 网关时。
* `update` - (默认为 1 分钟)用于更新 IPv6 网关时。
* `delete` - (默认为 5 分钟)用于删除 IPv6 网关时。

## 属性参考

除了上述参数外，还导出以下属性：

* `id` - Terraform 中的 IPv6 网关资源 ID。
* `status` - 资源的状态。有效值：`Available`(可用)、`Pending`(等待中)和 `Deleting`(删除中)。
* `spec` - IPv6 网关不区分规格，该参数不再使用，但为了兼容性仍保留。

### 导入

VPC IPv6 网关可以通过 id 导入，例如。

```bash
$ terraform import alibabacloudstack_vpc_ipv6_gateway.example <id>
```