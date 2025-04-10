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

* `ids` - (可选，变更时重建) 作业模板ID列表。此参数用于指定需要查询的作业模板的具体ID列表。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `templates` - Ehpc作业模板列表。每个元素包含以下属性：
  * `array_request` - 队列作业，格式为：1-10:2。表示作业数组请求，定义多个任务的范围和步长。
  * `clock_time` - 作业最大运行时间。设置作业允许运行的最大时间，超过该时间将被终止。
  * `command_line` - 作业命令。指定作业执行时运行的命令或脚本。
  * `gpu` - 单个计算节点使用的GPU数量。可能值范围为1~20000，表示每个计算节点可以分配的GPU数量。
  * `id` - 作业模板的ID。唯一标识一个作业模板。
  * `job_template_id` - 资源的第一个ID。与`id`类似，用于标识作业模板。
  * `job_template_name` - 作业模板名称。用户为作业模板指定的名称，便于管理和识别。
  * `mem` - 单个计算节点最大内存。表示每个计算节点可用的最大内存容量。
  * `node` - 提交任务时需要的数据节点数量。可能值范围为1~5000，表示作业所需的计算节点数。
  * `package_path` - 作业命令所在的目录。指定作业命令或脚本所在的路径。
  * `priority` - 作业优先级。可能值范围为0~9，数字越大优先级越高。
  * `queue` - 作业队列。指定作业提交到的队列名称。
  * `re_runable` - 作业是否支持重新运行。布尔值，表示作业失败后是否可以重新运行。
  * `runas_user` - 执行作业的用户名。指定作业运行时使用的用户账户。
  * `stderr_redirect_path` - 错误输出路径。指定作业错误日志的输出路径。
  * `stdout_redirect_path` - 标准输出路径。指定作业标准输出的日志路径。
  * `task` - 单个计算节点所需的任务数量。可能值范围为1~20000，表示每个计算节点需要运行的任务数。
  * `thread` - 单个任务所需的线程数。可能值范围为1~20000，表示每个任务需要的线程数。
  * `variables` - 作业的环境变量。以键值对的形式定义作业运行时所需的环境变量。