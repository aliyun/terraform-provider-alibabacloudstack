---
subcategory: "云服务总线 CSB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_csb_project"
sidebar_current: "docs-Alibabacloudstack-resource-csb-project"
description: |-
  编排CSB项目
---

# alibabacloudstack_csb_project

使用Provider配置的凭证在指定的资源集下编排CSB项目。

有关 CSB 项目的更多信息以及如何使用它，请参阅 [创建项目](https://help.aliyun.com/apsara/enterprise/v_3_18_0_30393230/csb/apsarastack-developer-guide/obtains-information-about-a-single-service-group.html?spm=a2c4g.14484438.10001.97)

## 示例用法

### 基础用法

```hcl
resource "alibabacloudstack_csb_project" "project" {
  csb_id            = "your-csb-id"
  project_name      = "example-project"
  owner_name        = "project-owner"
  owner_email       = "owner@example.com"
  owner_phone_num   = "13800138000"
  description       = "Example CSB project"
}
```

## 参数说明

支持以下参数：

* `csb_id` - (必填，变更时强制重建) CSB 实例的 ID。修改此参数会强制重新创建资源。
* `project_name` - (必填) CSB 项目的名称。长度限制为 1 到 128 个字符。
* `owner_name` - (必填) 项目所有者的名称。
* `owner_email` - (可选) 项目所有者的电子邮件地址。
* `owner_phone_num` - (可选) 项目所有者的电话号码。
* `description` - (可选) 项目的描述信息。

## 属性说明

导出以下属性：

* `id` - 资源的唯一标识，格式为 `csb_id:project_name`。
* `csb_id` - CSB 实例的 ID。
* `project_name` - CSB 项目的项目名称。
* `project_id` - CSB 项目的内部 ID。
* `owner_name` - CSB 项目的项目所有者名称。
* `owner_email` - 项目所有者的电子邮件地址。
* `owner_phone_num` - 项目所有者的电话号码。
* `description` - 项目的描述信息。
* `owner_id` - CSB 项目的拥有者 ID。
* `api_num` - CSB 项目中已发布的 API 数量。

## Import

CSB Project 可以使用 `csb_id` 和 `project_name` 的组合进行导入，格式为 `csb_id:project_name`，例如：

```
$ terraform import alibabacloudstack_csb_project.example <csb_id>:<project_name>
```
