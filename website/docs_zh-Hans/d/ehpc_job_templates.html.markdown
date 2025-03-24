---
subcategory: "弹性高性能计算(ehpc)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ehpc_job_templates"
sidebar_current: "docs-alibabacloudstack-datasource-ehpc-job-templates"
description: |-
  查询弹性高性能计算集群作业模板
---

# alibabacloudstack_ehpc_job_templates

根据指定过滤条件列出当前凭证权限可以访问的弹性高性能计算集群作业模板列表。


## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_ehpc_job_template" "default" {
  job_template_name = "example_value"
  command_line      = "./LammpsTest/lammps.pbs"
}
data "alibabacloudstack_ehpc_job_templates" "ids" {
  ids = [alibabacloudstack_ehpc_job_template.default.id]
}
output "ehpc_job_template_id_1" {
  value = data.alibabacloudstack_ehpc_job_templates.ids.id
}


```

## 参数参考

支持以下参数：

* `ids` - (可选，变更时重建) 作业模板ID列表。

## 参数参考

除了上述列出的参数外，还导出以下属性：

* `templates` - Ehpc作业模板列表。每个元素包含以下属性：
  * `array_request` - 队列作业，格式为：1-10:2。
  * `clock_time` - 作业最大运行时间。
  * `command_line` - 作业命令。
  * `gpu` - 单个计算节点使用的GPU数量。可能值：1~20000。
  * `id` - 作业模板的ID。
  * `job_template_id` - 资源的第一个ID。
  * `job_template_name` - 作业模板名称。
  * `mem` - 单个计算节点最大内存。
  * `node` - 提交任务时需要的数据节点数量。可能值：1~5000。
  * `package_path` - 作业命令所在的目录。
  * `priority` - 作业优先级。可能值：0~9。
  * `queue` - 作业队列。
  * `re_runable` - 作业是否支持重新运行。
  * `runas_user` - 执行作业的用户名。
  * `stderr_redirect_path` - 错误输出路径。
  * `stdout_redirect_path` - 标准输出路径。
  * `task` - 单个计算节点所需的任务数量。可能值：1~20000。
  * `thread` - 单个任务所需的线程数。可能值：1~20000。
  * `variables` - 作业的环境变量。