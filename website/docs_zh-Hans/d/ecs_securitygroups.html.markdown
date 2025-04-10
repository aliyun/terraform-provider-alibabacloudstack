---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_securitygroups"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-securitygroups"
description: |- 
  查询云服务器安全组
---

# alibabacloudstack_ecs_securitygroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_security_groups`

根据指定过滤条件列出当前凭证权限可以访问的云服务器安全组列表。

## 示例用法

```hcl
# 过滤安全组并将结果打印到文件中
data "alibabacloudstack_ecs_securitygroups" "sec_groups_ds" {
  name_regex = "^web-"
}

# 结合 VPC 使用
data "alibabacloudstack_ecs_securitygroups" "primary_sec_groups_ds" {
  vpc_id = var.vpc_id
}

output "first_group_id" {
  value = "${data.alibabacloudstack_ecs_securitygroups.primary_sec_groups_ds.groups.0.id}"
}

# 使用标签过滤安全组
data "alibabacloudstack_ecs_securitygroups" "taggedSecurityGroups" {
  tags = {
    Environment = "Production"
    Department  = "Finance"
  }
}

# 示例：结合其他资源使用
variable "name" {
  default = "tf-testAlibabacloudstackEcsSecurityGroups76924"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip            = "172.16.0.0/24"
}

data "alibabacloudstack_ecs_securitygroups" "default" {
  ids = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (选填, 变更时重建) - 一个正则表达式字符串，用于按名称过滤结果中的安全组。
* `vpc_id` - (选填, 变更时重建) - 安全组所属VPC ID。仅当您想检索或创建VPC类型的安全组时需要此参数。在支持经典网络的区域中，不指定此参数即可检索或创建经典网络类型的安全组。
* `ids` - (选填) - 一个安全组 ID 列表，用于过滤结果。
* `tags` - (选填) - 分配给安全组的标签映射。它必须是以下格式：
  ```hcl
  tags = {
    tagKey1 = "tagValue1",
    tagKey2 = "tagValue2"
  }
  ```

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 安全组 ID 列表。
* `names` - 安全组名称列表。
* `groups` - 安全组列表。每个元素包含以下属性：
  * `id` - 安全组的 ID。
  * `name` - 安全组的名称。
  * `description` - 安全组的描述信息。长度为2~256个英文或中文字符，不能以`http://`和`https://`开头。默认值为空。
  * `vpc_id` - 拥有该安全组的 VPC 的 ID。仅当您想检索或创建 VPC 类型的安全组时需要此参数。在支持经典网络的区域中，不指定此参数即可检索或创建经典网络类型的安全组。
  * `creation_time` - 安全组的创建时间。
  * `tags` - 分配给安全组的标签映射。