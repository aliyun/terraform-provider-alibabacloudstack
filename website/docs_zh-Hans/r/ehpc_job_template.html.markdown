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

## 参数参考

支持以下参数：

* `array_request` - (可选) 队列作业，形式为：1-10:2。
* `clock_time` - (可选) 作业最大运行时间。
* `command_line` - (必填) 作业命令。
* `gpu` - (可选) 单个计算节点使用的 GPU 数量。可能的值：1~20000。
* `job_template_name` - (必填) 作业模板名称。
* `mem` - (可选) 单个计算节点最大内存。
* `node` - (可选) 提交任务时所需的计算数据节点数量。可能的值：1~5000。
* `package_path` - (可选) 作业命令所在目录。
* `priority` - (可选) 作业优先级。
* `queue` - (可选) 作业队列。
* `re_runable` - (可选) 作业是否支持重新运行。
* `runas_user` - (可选) 执行该作业的用户名。
* `stderr_redirect_path` - (可选) 错误输出路径。
* `stdout_redirect_path` - (可选) 标准输出路径。
* `task` - (可选) 单个计算节点所需的任务数。可能的值：1~20000。
* `thread` - (可选) 单个任务所需的线程数。
* `variables` - (可选) 作业的环境变量。

## 属性参考

导出以下属性：

* `id` - Terraform 中作业模板的资源 ID。
* `re_runable` - 作业是否支持重新运行。

## 导入

EHPC 作业模板可以使用 id 导入，例如：

```bash
$ terraform import alibabacloudstack_ehpc_job_template.example <id>
```