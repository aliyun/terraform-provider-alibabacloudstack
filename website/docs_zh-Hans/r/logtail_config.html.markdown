---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_logtail_config"
sidebar_current: "docs-alibabacloudstack-resource-logtail-config"
description: |-
  编排日志接入配置
---

# alibabacloudstack_logtail_config

使用Provider配置的凭证在指定的资源集编排日志接入配置。
Logtail接入服务是日志服务提供的日志采集Agent。
您可以使用Logtail在日志服务控制台中实时采集服务器例如阿里云弹性
计算服务(ECS)实例上的日志。[详情参考](https://www.alibabacloud.com/help/doc-detail/29058.htm
)

## 示例用法

### 基础用法

```
resource "alibabacloudstack_log_project" "example" {
  name        = "test-tf"
  description = "create by terraform"
}

resource "alibabacloudstack_log_store" "example" {
  project               = alibabacloudstack_log_project.example.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alibabacloudstack_logtail_config" "example" {
  project      = alibabacloudstack_log_project.example.name
  logstore     = alibabacloudstack_log_store.example.name
  input_type   = "file"
  log_sample   = "test"
  name         = "tf-log-config"
  output_type  = "LogService"
  input_detail = file("config.json")
}
```


## 参数参考

支持以下参数：

* `project` - (必填，变更时重建) 日志所属的项目名称。
* `logstore` - (必填，变更时重建) 查询索引所属的日志存储名称。
* `input_type` - (必填) 输入类型。目前仅支持文件和插件两种类型。
* `log_sample` - (可选) Logtail配置的日志样本。日志大小不能超过1,000字节。
* `name` - (必填，变更时重建) Logtail配置名称，在同一项目中必须唯一。
* `output_type` - (必填) 输出类型。目前仅支持LogService。
* `input_detail` - (必填) Logtail配置所需的JSON文件。([详情参考](https://www.alibabacloud.com/help/doc-detail/29058.htm))
  

## 属性参考

导出以下属性：

* `id` - 日志存储索引的ID。格式为 `<project>:<logstore>:<config_name>`。

## 导入

Logtial配置可以通过id导入，例如：

```bash
$ terraform import alibabacloudstack_logtail_config.example tf-log:tf-log-store:tf-log-config
```