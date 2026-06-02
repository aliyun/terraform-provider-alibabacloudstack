---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_project"
description: |-
  Provides a Alibabacloudstack maxcompute project resource.
---

# alibabacloudstack_maxcompute_project

The project is the basic unit of operation in maxcompute. 

## Example Usage

Basic Usage

```terraform
variable "name" {
}

data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}


resource "alibabacloudstack_maxcompute_cu" "default" {
	cu_name =      "${var.name}"
	cu_num =       2
	cluster_name = "${data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster}"
}


resource "alibabacloudstack_maxcompute_project" "default" {
  account_pk = ""
  account = ""
  quota_id = "${alibabacloudstack_maxcompute_cu.default.id}"
  external_table = "true"
  vpc_ids = [
              "${alibabacloudstack_vpc_vpc.default.id}"
            ]
  name = "${var.name}"
  disk = "50"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Required, ForceNew) It has been deprecated from provider version 1.110.0 and `project_name` instead.
* `quota_id` - (Required) The quota ID of the maxcompute project. 
* `disk` - (Required) The disk size of the maxcompute project. 
* `account` - (Required, ForceNew) The account of the maxcompute project.
* `account_pk` - (Required, ForceNew) The account pk of the maxcompute project.
* `external_table` - (Optional) Whether to enable joint computing.
* `vpc_ids` - (Optional) The vpc ids of the maxcompute project.
* `core_arch` - (Optional) The core arch of the maxcompute project.
* `cpu_type` - (Optional) The cpu type of the maxcompute project.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the maxcompute project. It is the same as its name.
* `name` - The name of the maxcompute project. 

## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```