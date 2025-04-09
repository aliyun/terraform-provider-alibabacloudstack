---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_sitemonitor"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-sitemonitor"
description: |- 
  Provides a cloudmonitorservice Sitemonitor resource.
---

# alibabacloudstack_cloudmonitorservice_sitemonitor
-> **NOTE:** Alias name has: `alibabacloudstack_cms_site_monitor`

Provides a cloudmonitorservice Sitemonitor resource.

## Example Usage

Basic Usage

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

## Argument Reference

The following arguments are supported:

* `address` - (Required) The URL or IP address monitored by the site monitoring task. It must be a valid address that can be accessed over the internet.
* `task_name` - (Required) The name of the site monitoring task. The name must be 4 to 100 characters in length and can contain letters, digits, and underscores (`_`).
* `task_type` - (Required, ForceNew) The protocol type for the site monitoring task. Valid values include: `HTTP`, `PING`, `TCP`, `UDP`, `DNS`, `SMTP`, `POP3`, and `FTP`.
* `alert_ids` - (Optional) A list of alarm rule IDs associated with the site monitoring task.
* `interval` - (Optional) The monitoring interval for the site monitoring task. Unit: minutes. Valid values: `1`, `5`, and `15`. Default value: `1`.
* `isp_cities` - (Optional) A JSON array specifying the detection points (ISPs and cities) used for monitoring. If this parameter is not specified, three detection points will be chosen randomly. For example:
  ```json
  [
    {"city":"546","isp":"465"},
    {"city":"572","isp":"465"},
    {"city":"738","isp":"465"}
  ]
  ```
  You can call the [DescribeSiteMonitorISPCityList](https://www.alibabacloud.com/help/en/doc-detail/115045.htm) operation to query available detection point information.
  * `city` - (Required) The ID of the city where the detection point is located.
  * `isp` - (Required) The ID of the Internet Service Provider (ISP) for the detection point.
* `options_json` - (Optional) A JSON string containing extended options specific to the protocol of the site monitoring task. The options vary depending on the `task_type`. For example:
  ```json
  {
    "Dnstype": "A",       // DNS record type (for DNS tasks)
    "Failurerate": 0.5,   // Failure rate threshold (for PING tasks)
    "Pingnum": 10         // Number of PING packets sent (for PING tasks)
  }
  ```
* `task_state` - (Optional) The current state of the site monitoring task.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (ID) of the site monitoring task.
* `task_state` - The current state of the site monitoring task.
* `create_time` - The timestamp indicating when the site monitoring task was created.
* `update_time` - The timestamp indicating when the site monitoring task was last updated.