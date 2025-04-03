---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_connection"
sidebar_current: "docs-Alibabacloudstack-data-works-connection"
description: |- 
  编排Data Works连接
---

# alibabacloudstack_data_works_connection

使用Provider配置的凭证在指定的资源集下编排Data Works连接。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testaccdata_worksconnection44523"
}

variable "password" {
  default = "inputYourCodeHere@ascm"
}

resource "alibabacloudstack_data_works_connection" "default" {
  project_id     = "10060"
  connection_type = "rds"
  content = {
    username      = "cxt_new"
    database     = "cxt_test_new"
    tag          = "rds"
    password     = var.password
    instanceName = "rm-6cq93i8k9q0045i5t"
    rdsOwnerId   = "1371730998580255"
  }
  env_type       = "1"
  sub_type       = "mysql"
  name           = var.name
  description    = "Description for ${var.name}"
}
```

## 参数参考

支持以下参数：

  * `connection_id` - (变更时重建) - 连接的ID。这是自动生成的，创建后无法修改。
  * `project_id` - (必填) - 要创建连接的项目的ID(工作空间ID)。
  * `connection_type` - (必填) - 连接类型。目前支持`rds`。
  * `content` - (必填) - 数据源的详细信息。这是一个包含以下键的映射：
    * `username` - (必填) - 用于连接到数据源的用户名。
    * `database` - (必填) - 要连接的数据库名称。
    * `tag` - (必填) - 与连接关联的标签。对于`rds`，这通常设置为`rds`。
    * `password` - (必填) - 用于连接到数据源的密码。
    * `instanceName` - (必填) - RDS实例的名称。
    * `rdsOwnerId` - (必填) - RDS实例的所有者ID。
  * `env_type` - (必填) - 数据源所属的环境。有效值为：
    * `0` - 开发环境。
    * `1` - 生产环境。
  * `sub_type` - (选填) - 字符串的子类型，适用于某些父类型包含子类型的场景。对于`rds`，有效的子类型为：
    * `mysql`
    * `sqlserver`
    * `postgresql`
  * `name` - (必填, 变更时重建) - 连接的名称。此值在项目内必须唯一。
  * `description` - (选填) - 连接的描述。这提供了有关连接的附加信息。

## 属性参考

除了上述所有参数外，还导出了以下属性：

  * `connection_id` - 连接的唯一标识符。此值格式为 `<connection_id>:<$.ProjectId>`。