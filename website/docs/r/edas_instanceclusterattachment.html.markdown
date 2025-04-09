---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_instanceclusterattachment"
sidebar_current: "docs-Alibabacloudstack-edas-instanceclusterattachment"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Instanceclusterattachment resource.
---

# alibabacloudstack_edas_instanceclusterattachment
-> **NOTE:** Alias name has: `alibabacloudstack_edas_instance_cluster_attachment`

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Instanceclusterattachment resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-edasicattachment19002"
}

variable "password" {
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
} 

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id              = alibabacloudstack_vpc.default.id
  cidr_block          = "10.1.1.0/24"
  availability_zone   = data.alibabacloudstack_zones.default.zones[0].id
  name               = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name       = var.name
  description= "New security group"
  vpc_id     = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  vswitch_id         = alibabacloudstack_vswitch.default.id
  image_id           = "centos_7_7_x64_20G_alibase_20200426.vhd"
  availability_zone  = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_efficiency"
  system_disk_size  = 60
  instance_type      = "ecs.n4v2.xlarge"

  security_groups    = [alibabacloudstack_security_group.default.id]
  instance_name      = var.name
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = var.name
  cluster_type = 2
  network_mode = 2
  vpc_id       = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
  cluster_id = alibabacloudstack_edas_cluster.default.id
  instance_ids = [
                   alibabacloudstack_instance.default.id
                 ]
  pass_word = var.password
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the cluster that you want to attach instances to.
* `instance_ids` - (Required, ForceNew) A list of ECS instance IDs that will be attached to the specified cluster.
* `pass_word` - (Required, ForceNew) The login password for the ECS instances in the cluster. This is required during the attachment process to ensure secure access.
* `status_map` -  (Optional) A map indicating the status of each instance in the cluster. The keys are instance IDs, and the values represent the status: `1` (Running), `0` (Converting), `-1` (Failed), and `-2` (Offline).
* `ecu_map` -  (Optional) A map linking each instance to its corresponding ECU (Elastic Compute Unit). The keys are instance IDs, and the values are ECU IDs.
* `cluster_member_ids` -  (Optional) A map of cluster member IDs associated with each instance. The keys are instance IDs, and the values are the cluster member IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource, which combines `<cluster_id>:<instance_id1,instance_id2,...>`.
* `status_map` -  A map indicating the status of each instance in the cluster. The keys are instance IDs, and the values represent the status: `1` (Running), `0` (Converting), `-1` (Failed), and `-2` (Offline).
* `ecu_map` -  A map linking each instance to its corresponding ECU (Elastic Compute Unit). The keys are instance IDs, and the values are ECU IDs.
* `cluster_member_ids` -  A map of cluster member IDs associated with each instance. The keys are instance IDs, and the values are the cluster member IDs.