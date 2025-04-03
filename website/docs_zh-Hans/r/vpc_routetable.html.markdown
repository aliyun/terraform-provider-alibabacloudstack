---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_routetable"
sidebar_current: "docs-Alibabacloudstack-vpc-routetable"
description: |- 
  编排VPC的路由表
---

# alibabacloudstack_vpc_routetable
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_route_table`

使用Provider配置的凭证在指定的资源集编排VPC的路由表。

## 示例用法

### 基础用法

```hcl
variable "name" {
    default = "tf-testaccvpcroute_table19406"
}

resource "alibabacloudstack_vpc" "default" {
    cidr_block = "172.16.0.0/12"
    name       = var.name
}

resource "alibabacloudstack_vpc_routetable" "default" {
    vpc_id      = alibabacloudstack_vpc.default.id
    name        = var.name
    description = "A detailed description of the route table."
    tags = {
        Environment = "Test"
    }
}
```

## 参数参考

支持以下参数：

* `vpc_id` - (必填，变更时重建) 路由表所属的VPC的ID。此字段在创建后无法修改。
* `name` - (可选) 路由表的名称。如果不指定，Terraform将自动生成一个唯一名称。
* `description` - (可选) 路由表的详细描述。这有助于识别路由表的目的或用途。
* `tags` - (可选) 分配给路由表的标签映射。这些标签可用于分类和成本分配。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 路由表实例的ID。
* `route_table_name` - 由`name`参数指定或自动生成的路由表名称。