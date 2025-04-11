---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_securitygroup"
sidebar_current: "docs-Alibabacloudstack-ecs-securitygroup"
description: |- 
  编排云服务器（Ecs）安全组
---

# alibabacloudstack_ecs_securitygroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_security_group`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）安全组。

## 示例用法

### 基础用法

```hcl
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}

resource "alibabacloudstack_security_group" "basic_group" {
  name        = "${var.name}_basic"
  description = "${var.name}_describe_basic"
}
```

### 基础用法(针对VPC)

```hcl
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_security_group" "vpc_group" {
  vpc_id         = alibabacloudstack_vpc.default.id
  name           = "${var.name}_vpc_group"
  description    = "${var.name}_describe_vpc_group"
  type          = "normal"
  inner_access_policy = "Accept"
}
```

### 高级用法(带标签)

```hcl
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_security_group" "tagged_group" {
  vpc_id         = alibabacloudstack_vpc.default.id
  name           = "${var.name}_tagged_group"
  description    = "${var.name}_describe_tagged_group"
  type          = "enterprise"
  inner_access_policy = "Drop"

  tags = {
    Environment = "Production"
    Owner      = "DevOps"
  }
}
```

## 参数说明

支持以下参数：

* `name` - (可选) 安全组的名称。长度必须为2到128个字符，可以包含字母、数字、下划线 (`_`)、点 (`.`) 和连字符 (`-`)。它必须以字母或中文字符开头，不能以 `http://` 或 `https://` 开头。如果不指定，Terraform 将自动生成一个唯一的名称。
  
* `description` - (可选) 安全组的描述。长度必须为2到256个字符，不能以 `http://` 或 `https://` 开头。默认情况下，此参数为空。

* `vpc_id` - (可选，强制更新) 要创建安全组的VPC的ID。如果要创建VPC类型的安全组，则此参数是必填的。在支持经典网络的区域中，您可以不指定 `vpc_id` 来创建经典网络类型的安全组。

* `type` - (可选，强制更新) 安全组的类型。有效值：
  * `normal`: 标准安全组(默认)。
  * `enterprise`: 企业级安全组。

* `inner_access_policy` - (可选) 安全组的内部访问策略。有效值：
  * `Accept`: 安全组中的所有实例可以相互通信。
  * `Drop`: 安全组中的所有实例相互隔离。
  此参数的值不区分大小写。默认值为 `Accept`。

* `tags` - (可选) 要分配给资源的标签映射。每个标签由键值对组成。标签键必须在资源内唯一。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 安全组的ID。

* `inner_access_policy` - 安全组的内部访问策略。有效值：
  * `Accept`: 安全组中的所有实例可以相互通信。
  * `Drop`: 安全组中的所有实例相互隔离。
  此参数的值不区分大小写。此属性反映了安全组的实际配置。