---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ehpc_job_template"
sidebar_current: "docs-alibabacloudstack-resource-ehpc-job-template"
description: |-
  编排弹性高性能计算(EHPC)作业模板
---

# alibabacloudstack_ehpc_job_template

使用Provider配置的凭证在指定的资源集下编排弹性高性能计算(EHPC)作业模板。

关于 EHPC 作业模板及其使用方法，请参阅 [什么是作业模板](https://www.alibabacloud.com/help/product/57664.html)。



## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_ehpc_job_template" "default" {
  job_template_name = "example_value"
  command_line      = "./LammpsTest/lammps.pbs"
}
```

## 参数说明

支持以下参数：

* `array_request` - (可选) 队列作业，形式为：1-10:2。用于定义批量任务的索引范围和步长。
* `clock_time` - (可选) 作业最大运行时间。格式为 HH:MM:SS，例如 "24:00:00" 表示最大运行时间为 24 小时。
* `command_line` - (必填) 作业命令。指定作业执行的具体命令行内容。
* `gpu` - (可选) 单个计算节点使用的 GPU 数量。可能的值：1~20000。用于分配 GPU 资源给任务。
* `job_template_name` - (必填) 作业模板名称。指定作业模板的唯一标识名称。
* `mem` - (可选) 单个计算节点最大内存。单位为 MB，例如 8192 表示 8GB 内存。
* `node` - (可选) 提交任务时所需的计算数据节点数量。可能的值：1~5000。用于定义需要分配的计算节点数。
* `package_path` - (可选) 作业命令所在目录。指定包含作业脚本或文件的路径。
* `priority` - (可选) 作业优先级。取值范围为 0 到 100，默认值为 50。数值越大，优先级越高。
* `queue` - (可选) 作业队列。指定作业提交到的队列名称。
* `re_runable` - (可选) 作业是否支持重新运行。取值为 true 或 false，默认值为 false。
* `runas_user` - (可选) 执行该作业的用户名。指定运行作业的用户身份。
* `stderr_redirect_path` - (可选) 错误输出路径。指定错误日志文件的存储路径。
* `stdout_redirect_path` - (可选) 标准输出路径。指定标准输出日志文件的存储路径。
* `task` - (可选) 单个计算节点所需的任务数。可能的值：1~20000。用于定义每个节点上运行的任务数。
* `thread` - (可选) 单个任务所需的线程数。用于定义每个任务所需的线程数。
* `variables` - (可选) 作业的环境变量。以键值对的形式定义作业运行时的环境变量。

## 属性说明

导出以下属性：

* `id` - Terraform 中作业模板的资源 ID。它是作业模板的唯一标识符。
* `re_runable` - 作业是否支持重新运行。返回 true 或 false，表示作业是否可以重新提交运行。

## 导入

EHPC 作业模板可以使用 id 导入，例如：

```bash
$ terraform import alibabacloudstack_ehpc_job_template.example <id>
```