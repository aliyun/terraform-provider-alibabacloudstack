---
subcategory: "Simple Log Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_logtail_attachment"
description: |-
  Provides a Alibabacloudstack logtail attachment resource.
---

# alibabacloudstack_logtail_attachment

The Logtail access service is a log collection agent provided by Log Service.
You can use Logtail to collect logs from servers such as Alibaba Cloud Elastic
Compute Service (ECS) instances in real time in the Log Service console. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29058.htm)

This resource amis to attach one logtail configure to a machine group.

-> **NOTE:** One logtail configure can be attached to multiple machine groups and one machine group can attach several logtail configures.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_log_project" "test" {
  name        = "test-tf2"
  description = "create by terraform"
}

resource "alibabacloudstack_log_store" "test" {
  project               = alibabacloudstack_log_project.test.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alibabacloudstack_log_machine_group" "test" {
  project       = alibabacloudstack_log_project.test.name
  name          = "tf-log-machine-group"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}

resource "alibabacloudstack_logtail_config" "test" {
  project      = alibabacloudstack_log_project.test.name
  logstore     = alibabacloudstack_log_store.test.name
  input_type   = "file"
  log_sample   = "test"
  name         = "tf-log-config"
  output_type  = "LogService"
  input_detail = <<DEFINITION
  	{
		"logPath": "/logPath",
		"filePattern": "access.log",
		"logType": "json_log",
		"topicFormat": "default",
		"discardUnmatch": false,
		"enableRawLog": true,
		"fileEncoding": "gbk",
		"maxDepth": 10
	}
	
DEFINITION

}

resource "alibabacloudstack_logtail_attachment" "test" {
  project             = alibabacloudstack_log_project.test.name
  logtail_config_name = alibabacloudstack_logtail_config.test.name
  machine_group_name  = alibabacloudstack_log_machine_group.test.name
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The name of the Log Service project to which the Logtail configuration and machine group belong.

* `logtail_config_name` - (Required, ForceNew) The name of the Logtail configuration. It must be unique within the same project.

* `machine_group_name` - (Required, ForceNew) The name of the machine group. It must be unique within the same project.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Logtail attachment. It is formatted as `<project>:<logtail_config_name>:<machine_group_name>`.

## Import

Logtail Attachment can be imported using the project name, logtail config name, and machine group name separated by colons, e.g.

```bash
$ terraform import alibabacloudstack_logtail_attachment.example tf-log:tf-log-config:tf-log-machine-group
```