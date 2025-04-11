---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_remind"
sidebar_current: "docs-Alibabacloudstack-data-works-remind"
description: |- 
  编排Data Works提醒
---

# alibabacloudstack_data_works_remind

使用Provider配置的凭证在指定的资源集下编排Data Works提醒。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testaccdata_worksremind93940"
}

resource "alibabacloudstack_data_works_remind" "default" {
  alert_methods = "SMS"
  alert_unit = "OWNER"
  remind_name = var.name
  remind_type = "FINISHED"
  remind_unit = "PROJECT"
  project_id = "10023"
  dnd_end = "23:59"
  node_ids = "node1,node2"
  baseline_ids = "baseline1,baseline2"
  biz_process_ids = "bizprocess1,bizprocess2"
  max_alert_times = 5
  alert_interval = 1200
  detail = "{\"hour\":23,\"minu\":59}"
  alert_targets = "uid1,uid2"
  robot_urls = "https://robot.url1,https://robot.url2"
  use_flag = true
}
```

## 参数参考

支持以下参数：

* `alert_unit` - (必填) 报警接收粒度，包括：
  * `OWNER`: 任务拥有者
  * `OTHER`: 指定人员
* `remind_name` - (必填) 自定义提醒规则名称，不得超过 128 个字符。
* `remind_type` - (必填) 触发条件，包括：
  * `FINISHED`: 任务完成
  * `UNFINISHED`: 任务未完成
  * `ERROR`: 任务错误
  * `CYCLE_UNFINISHED`: 周期任务未完成
  * `TIMEOUT`: 任务超时
* `remind_unit` - (必填) 对象类型，包括：
  * `NODE`: 任务节点
  * `BASELINE`: 基线
  * `PROJECT`: 工作区
  * `BIZPROCESS`: 业务流程
* `dnd_end` - (选填) 勿扰截止时间，格式为 HH:MM。HH 的取值范围为 0-23，mm 的取值范围为 0-59，默认值为 `00:00`。
* `node_ids` - (选填) 当对象类型 (`remind_unit`) 为 `NODE` 时，监控的任务节点 id。多个 id 以逗号（`,`）分隔，一个规则最多可以监控 50 个节点。
* `baseline_ids` - (选填) 当对象类型 (`remind_unit`) 为 `BASELINE` 时，监控的基线 id。多个基线 id 以逗号（`,`）分隔，一个规则最多可以监控 5 个基线。
* `project_id` - (选填) 当对象类型 (`remind_unit`) 为 `PROJECT` 时，监控的工作区 id。一个规则只能监控一个工作区。
* `biz_process_ids` - (选填) 当对象类型 (`remind_unit`) 为 `BIZPROCESS` 时，监控的业务流程 id。多个业务流程 id 以逗号（`,`）分隔，一个规则最多可以监控 5 个业务流程。
* `max_alert_times` - (选填) 最大告警次数。最小值为 1，最大值为 10，默认值为 3。
* `alert_interval` - (选填) 最小告警间隔，单位为秒。最小值为 1200，默认值为 1800。
* `detail` - (选填) 不同触发条件的描述如下：
  * 当 `remind_type` 为 `FINISHED` 时，为空。
  * 当 `remind_type` 为 `UNFINISHED` 时，参数格式为 `{"hour":23,"minu":59}`。hour 的取值范围为 0-47，minu 的取值范围为 0-59。
  * 当 `remind_type` 为 `ERROR` 时，传递为空。
  * 当 `remind_type` 为 `CYCLE_UNFINISHED` 时，参数传递格式为 `{"1": "05:50", "2": "06:50", ...}`。JSON key 为周期编号，其取值范围为 1-288。value 为该周期未完成对应的时间，格式为 HH:mm。hh 的取值范围为 0-47，mm 的取值范围为 0-59。
  * 当 `remind_type` 为 `TIMEOUT` 时，参数格式为 1800 秒。即从操作开始运行超过 30 分钟将触发告警。
* `alert_methods` - (选填) 告警方式包括：
  * `MAIL`: 邮件
  * `SMS`: 短信
  * `PHONE`: 电话（仅 DataWorks 专业版及以上支持）
  多种告警方式以英文逗号（`,`）分隔。
* `alert_targets` - (选填)
  * 当 `alert_unit` 为 `OWNER`（节点任务拥有者）时，为空。
  * 当 `alert_unit` 为 `OTHER` 时，传入指定用户的 Alibaba Cloud UID。多个 Alibaba Cloud UID 以英文逗号（`,`）分隔，最大数量为 10。
* `robot_urls` - (选填) 钉钉机器人 webhook 地址，多个 webhook 地址以英文逗号（`,`）分隔。
* `use_flag` - (选填) 开关规则，包括：
  * `true`: 开启
  * `false`: 关闭

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `remind_id` - 自定义提醒规则 ID。格式为 `<remind_id>:<$.ProjectId>`。
* `use_flag` - 规则开关状态，表示当前规则是否启用。