---
subcategory: "Prometheus 监控服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_contact"
sidebar_current: "docs-alibabacloudstack-resource-prometheus-v2-contact"
description: |-
  管理阿里云Prometheus v2联系人资源
---

# alibabacloudstack_prometheus_v2_contact

管理阿里云Prometheus v2监控服务的联系人资源，用于配置告警通知接收人信息。

## 示例用法

### 基础联系人配置

```hcl

variable "name" {
  default = "tfacc-10446"
}


resource "alibabacloudstack_prometheus_v2_contact" "default" {
  mail     = "test@example.com"
  username = var.name
  mobile   = "13812345678"
}
```

## 参数说明

支持以下参数，按类型排序（同类参数按字母序）：

* `mail` - (必填) 联系人的电子邮箱地址，用于接收告警通知。邮箱格式需符合标准规范。
* `mobile` - (必填) 联系人的手机号码，用于短信告警通知。需填写11位中国大陆手机号。
* `username` - (必填) 联系人的唯一标识名称，长度1-64个字符，支持中文、字母、数字及下划线。

## 属性说明

导出以下属性（按规则排序）：

* `id` - 联系人的唯一系统标识ID，由Prometheus服务自动生成。

## Import

联系人可以使用联系人 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_prometheus_v2_contact.example 12345
```