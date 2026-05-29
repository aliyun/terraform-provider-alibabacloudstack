---
layout: "alicloud-doc"
page_title: "Resource: alibabacloudstack_bcmp_security_group"
subcategory: "裸金属算力平台 BMCP"
---

# alibabacloudstack_bcmp_security_group

提供 BMCP 安全组资源。

-> **注意：** 该资源也可以通过以下别名引用：
- `alibabacloudstack_bmcp_security_group`

## 示例

```terraform
variable "name" {
  default = "tf-example-bmcp-sg"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_bcmp_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = alibabacloudstack_vpc.default.id
}
```

## 参数说明

以下参数用于配置安全组：

- `vpc_id` - （必填，ForceNew）安全组所属的 VPC ID。修改此参数将强制创建新资源。
- `name` - （必填）安全组名称。长度必须在 2 到 128 个字符之间。
- `description` - （可选）安全组描述。长度必须在 0 到 256 个字符之间。

## 属性说明

以下属性由系统自动返回：

- `id` - 安全组 ID。
- `sg_id` - 安全组 ID。
- `name` - 安全组名称。
- `description` - 安全组描述。
- `vpc_id` - VPC ID。
- `resource_group` - 资源组 ID。
- `resource_group_name` - 资源组名称。
- `department` - 部门 ID。
- `department_name` - 部门名称。
- `region_id` - 地域 ID。
- `ascm_create_user` - 创建该安全组的 ASCM 用户。
- `create_time` - 安全组创建时间。
- `update_time` - 安全组最后一次更新时间。

## 超时

`timeouts` 块允许您为特定操作指定超时时间：

- `create` - （默认 10 分钟）用于创建 BMCP 安全组。
- `delete` - （默认 10 分钟）用于删除 BMCP 安全组。

## 导入

BMCP 安全组可以使用安全组 ID（sgId）进行导入，例如：

```
$ terraform import alibabacloudstack_bcmp_security_group.example sg-12345678
```
