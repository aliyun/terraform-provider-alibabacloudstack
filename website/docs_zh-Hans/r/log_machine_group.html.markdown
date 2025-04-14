---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_machine_group"
sidebar_current: "docs-alibabacloudstack-resource-log-machine-group"
description: |-
  编排日志告警的机器组
---

# alibabacloudstack_log_machine_group

使用Provider配置的凭证在指定的资源集编排日志告警的机器组。
日志服务通过 Logtail 客户端以机器组的形式管理所有需要采集日志的 ECS 实例。 [详情参考](https://www.alibabacloud.com/help/doc-detail/28966.htm)

## 示例用法

### 基础用法

```
resource "alibabacloudstack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alibabacloudstack_log_machine_group" "example" {
  project       = alibabacloudstack_log_project.example.name
  name          = "tf-machine-group"
  identify_type = "ip"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.2"]
}
```


## 参数说明

支持以下参数：

* `project` - (必填，变更时重建) 机器组所属的日志项目名称。
* `name` - (必填，变更时重建) 机器组名称，在同一项目中必须唯一。
* `identify_type` - (可选) 机器标识类型。例如，可以设置为 `ip` 表示通过 IP 地址标识机器。
* `topic` - (可选) 机器组主题。用于对日志进行分类或标记。

## 属性说明

导出以下属性：

* `id` - 日志机器组的 ID。格式为 `<project>:<name>`。
* `project` - 项目名称。
* `name` - 机器组名称。
* `identify_type` - 机器标识类型。
* `identify_list` - 机器标识列表。例如，当 `identify_type` 设置为 `ip` 时，该列表包含具体的 IP 地址。
* `topic` - 机器组主题。

## 导入

日志机器组可以使用 id 导入，例如

```bash
$ terraform import alibabacloudstack_log_machine_group.example tf-log:tf-machine-group
```