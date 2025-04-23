---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_sitemonitor"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-sitemonitor"
description: |- 
  编排云监控服务（CMS）站点监控
---

# alibabacloudstack_cloudmonitorservice_sitemonitor
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_site_monitor`

使用Provider配置的凭证在指定的资源集下编排云监控服务（CMS）站点监控。

## 示例用法

### 基础用法

```hcl
variable "name" {
    default = "tf-testacccloud_monitor_servicesite_monitor82409"
}

resource "alibabacloudstack_cloudmonitorservice_sitemonitor" "default" {
  address     = "www.aliyun.com"
  task_name   = var.name
  task_type   = "PING"
  interval    = "1"
  isp_cities = [
    {
      city = "546"
      isp  = "465"
    },
    {
      city = "572"
      isp  = "465"
    },
    {
      city = "738"
      isp  = "465"
    }
  ]
  options_json = <<JSON
  {
    "Dnstype": "A",
    "Failurerate": 0.5,
    "Pingnum": 10
  }
  JSON
}
```

## 参数参考

支持以下参数：

* `address` - (必填) - 站点监控任务监控的URL或IP地址。它必须是一个可以在互联网上访问的有效地址。
* `task_name` - (必填) - 站点监控任务的名称。名称长度必须为4到100个字符，并可以包含字母、数字和下划线(`_`)。
* `task_type` - (必填, 变更时重建) - 站点监控任务的协议类型。有效值包括：`HTTP`、`PING`、`TCP`、`UDP`、`DNS`、`SMTP`、`POP3` 和 `FTP`。
* `alert_ids` - (选填) - 与站点监控任务关联的告警规则ID列表。
* `interval` - (选填) - 站点监控任务的监控间隔。单位：分钟。有效值：`1`、`5` 和 `15`。默认值：`1`。
* `isp_cities` - (选填) - 一个JSON数组，指定用于监控的探测点(ISP和城市)。如果不指定此参数，将随机选择三个探测点。例如：
  ```json
  [
    {"city":"546","isp":"465"},
    {"city":"572","isp":"465"},
    {"city":"738","isp":"465"}
  ]
  ```
  您可以通过调用 [DescribeSiteMonitorISPCityList](https://www.alibabacloud.com/help/en/doc-detail/115045.htm) 操作查询可用的探测点信息。
  * `city` - (必填) - 探测点所在城市的ID。
  * `isp` - (必填) - 探测点的互联网服务提供商(ISP)的ID。
* `options_json` - (选填) - 一个JSON字符串，包含特定于站点监控任务协议的扩展选项。选项因 `task_type` 而异。例如：
  ```json
  {
    "Dnstype": "A",       // DNS记录类型(适用于DNS任务)
    "Failurerate": 0.5,   // 失败率阈值(适用于PING任务)
    "Pingnum": 10         // 发送的PING数据包数量(适用于PING任务)
  }
  ```

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 站点监控任务的唯一标识符(ID)。
* `task_state` - 站点监控任务的当前状态。
* `create_time` - 表示站点监控任务创建时间的时间戳。
* `update_time` - 表示站点监控任务最后更新时间的时间戳。
