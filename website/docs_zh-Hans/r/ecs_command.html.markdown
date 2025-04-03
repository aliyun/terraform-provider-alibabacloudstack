---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_command"
sidebar_current: "docs-Alibabacloudstack-ecs-command"
description: |- 
  编排云服务器（Ecs）命令
---

# alibabacloudstack_ecs_command

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）命令。

## 示例用法

### 基础用法：

```terraform
variable "name" {
    default = "tf-testaccecscommand49325"
}

resource "alibabacloudstack_ecs_command" "default" {
  command_name      = var.name
  command_content   = "bHMK" # Base64-encoded content: "ls\n"
  type              = "RunShellScript"
  description       = "Command for listing files in the root directory"
  enable_parameter  = false
  timeout           = 120
  working_dir       = "/root"
}
```

高级用法(包含参数和复杂命令)：

```terraform
resource "alibabacloudstack_ecs_command" "complex" {
  command_name      = "StopKubeletAndRemoveK8S"
  command_content   = "c3lzdGVtY3RscyBzdG9wIGt1YmVsZXQuc2VydmljZTsgc3lzdGVtY3QgcXVldWUga3ViZWxldC5zZXJ2aWNlOyBzeXN0ZW1jdGxzIGRlbW9uLXJlbG9hZDsgeXVtIC15IHJlbW92ZSBrdWJlYWRtIGt1YmVsZXQga3VibGV0IGt1YnVsZXRjIg=="
  # Base64-encoded content: "systemctl stop kubelet.service; systemctl disable kubelet.service; systemctl daemon-reload; yum -y remove kubeadm kubelet kubectl;"
  type              = "RunShellScript"
  description       = "Stops and removes K8S components from the system."
  enable_parameter  = true
  timeout           = 300
  working_dir       = "/root"
}
```

## 参数参考

支持以下参数：

* `command_name` - (必填，变更时重建) 命令的名称。它必须是唯一的，并且可以包含多达 128 个字符。
* `command_content` - (必填，变更时重建) 命令的内容，经过 Base64 编码。这是将在 ECS 实例上执行的实际脚本或命令。
* `type` - (必填，变更时重建) 命令的类型。有效值：
  * `RunShellScript`: 用于 Linux 系统，执行 shell 脚本。
  * `RunBatScript`: 用于 Windows 系统，执行批处理脚本。
  * `RunPowerShellScript`: 用于 Windows 系统，执行 PowerShell 脚本。
* `description` - (可选，变更时重建) 命令的简要描述。它有助于识别命令的目的。
* `enable_parameter` - (可选，变更时重建) 指定命令是否支持自定义参数。默认值为 `false`。如果设置为 `true`，则在调用命令时可以传递参数。
* `timeout` - (可选，变更时重建) 命令在 ECS 实例上执行的超时期间。单位：秒。默认值为 `60`。
* `working_dir` - (可选，变更时重建) 命令将在 ECS 实例上执行的工作目录。默认值为 Linux 系统的 `/root` 和 Windows 系统的 `C:\Windows\system32`。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - ECS Command 资源的唯一标识符。
* `command_id` - 创建的 ECS Command 的 ID。