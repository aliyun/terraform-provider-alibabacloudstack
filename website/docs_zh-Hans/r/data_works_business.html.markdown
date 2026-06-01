---
subcategory: "一站式大数据开发治理平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_business"
sidebar_current: "docs-Alibabacloudstack-resource-data-works-business"
description: |-
  编排DataWorks业务流程
---

# alibabacloudstack_data_works_business

使用Provider配置的凭证在指定的资源集下编排DataWorks业务流程。


## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_data_works_business" "example" {
  project_id  = "12345"
  name        = "tf-testaccdataworksbusiness"
  description = "这是一个测试业务流程"
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 业务流程的名称。

* `project_id` - (选填，变更时重建) 要创建业务流程的数据工坊项目的ID。修改此参数会强制重新创建资源。

* `description` - (选填) 业务流程的描述。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 业务流程资源的唯一标识符。其值格式为 `<project_id>:<business_id>`。

* `business_id` - 数据工坊项目内业务流程的唯一标识符。

## Import

DataWorks业务流程可以通过 project_id 和 business_id 进行导入，两者之间用冒号分隔，例如：

```
$ terraform import alibabacloudstack_data_works_business.example 12345:67890
```
