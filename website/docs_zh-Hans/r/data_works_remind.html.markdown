---
subcategory: "一站式大数据开发治理平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_remind"
sidebar_current: "docs-Alibabacloudstack-data-works-remind"
description: |-
  提供 DataWorks 自定义报警规则资源。
---

# alibabacloudstack_data_works_remind

提供 DataWorks 自定义报警规则（Remind）资源。

-> **注意：** `node_ids`、`baseline_ids`、`project_id` 和 `biz_process_ids` 四个参数互斥。必须根据 `remind_unit` 类型指定其中至少一个。

## 示例

基本用法

```terraform
resource "alibabacloudstack_data_works_remind" "default" {
  remind_name   = "tf_test_remind"
  remind_type   = "ERROR"
  remind_unit   = "NODE"
  node_ids      = ["12345", "12346"]
  alert_unit    = "OWNER"
  alert_methods = ["MAIL", "SMS"]
}
```

## 参数说明

以下参数适用于此资源：

### 必填参数

* `remind_name` - （必填）自定义规则的名称，不能超过 128 个字符。

* `remind_type` - （必填）触发报警的条件。有效值：`FINISHED`（完成）、`UNFINISHED`（未完成）、`ERROR`（出错）、`CYCLE_UNFINISHED`（周期未完成）、`TIMEOUT`（运行超时）。

* `remind_unit` - （必填）监控对象的类型。有效值：`NODE`（任务节点）、`BASELINE`（基线）、`PROJECT`（工作空间）、`BIZPROCESS`（业务流程）。

* `alert_unit` - （必填）报警接收对象的粒度。有效值：`OWNER`（任务责任人）和 `OTHER`（指定的人）。

* `alert_methods` - （必填）报警方式。有效值：`MAIL`（邮件）、`SMS`（短信）、`PHONE`（电话，仅 DataWorks 专业版及以上版本支持）、`WEBHOOKS`（企业微信或飞书机器人）、`DINGROBOTS`（钉钉群机器人）。

### 可选参数

* `node_ids` - （可选）监控的任务节点 ID。当 `remind_unit` 为 `NODE` 时生效。多个 ID 之间使用英文逗号分隔，一个规则最多监控 50 个节点。与 `baseline_ids`、`project_id` 和 `biz_process_ids` 冲突。

* `baseline_ids` - （可选）监控的基线 ID。当 `remind_unit` 为 `BASELINE` 时生效。多个 ID 之间使用英文逗号分隔，一个规则最多监控 5 条基线。与 `node_ids`、`project_id` 和 `biz_process_ids` 冲突。

* `project_id` - （可选）监控的工作空间 ID。当 `remind_unit` 为 `PROJECT` 时生效。一个规则只能监控一个工作空间。与 `node_ids`、`baseline_ids` 和 `biz_process_ids` 冲突。

* `biz_process_ids` - （可选）监控的业务流程 ID。当 `remind_unit` 为 `BIZPROCESS` 时生效。多个 ID 之间使用英文逗号分隔，一个规则最多监控 5 个业务流程。与 `node_ids`、`baseline_ids` 和 `project_id` 冲突。

* `max_alert_times` - （可选）最大报警次数。有效值：1 到 10，默认值为 3。

* `alert_interval` - （可选）最小报警间隔，单位为秒。最小值为 1200，默认值为 1800。

* `detail` - （可选）不同触发条件的参数配置：
  - 当 `remind_type` 为 `FINISHED` 或 `ERROR` 时，传空。
  - 当 `remind_type` 为 `UNFINISHED` 时，配置为 JSON 字符串：`{"hour":23,"minu":59}`。`hour` 的有效值范围：[0,47]，`minu` 的有效值范围：[0,59]。
  - 当 `remind_type` 为 `CYCLE_UNFINISHED` 时，配置为 JSON 字符串，key 为周期号，value 为超时时间：`{"1":"05:50","2":"06:50"}`。周期号有效值：[1,288]。
  - 当 `remind_type` 为 `TIMEOUT` 时，传入超时时间，单位为秒（例如 `1800`）。

* `alert_targets` - （可选）当 `alert_unit` 为 `OTHER` 时，指定用户的阿里云 UID。多个 UID 之间使用英文逗号分隔。当 `alert_unit` 为 `OWNER` 时，此属性由 API 返回计算得出。

* `webhooks` - （可选）企业微信或飞书机器人的 webhook URL。仅当 `alert_methods` 包含 `WEBHOOKS` 时生效。

* `robot_urls` - （可选）钉钉群机器人的 webhook URL。仅当 `alert_methods` 包含 `DINGROBOTS` 时生效。

* `dnd_end` - （可选）免打扰截止时间，格式为 `hh:mm`。`hh` 的有效值范围：[0,23]，`mm` 的有效值范围：[0,59]。默认值：`00:00`。

* `use_flag` - （可选）是否启用报警规则。有效值：`true`、`false`。此属性由 API 返回计算得出。

## 属性说明

除上述所有参数外，还导出以下属性：

* `id` - 自定义报警规则的 ID（remind_id）。

* `remind_id` - 创建的报警规则 ID。

## 导入

DataWorks 自定义报警规则可以使用 `remind_id` 导入，例如：

```bash
$ terraform import alibabacloudstack_data_works_remind.example 12345
```
