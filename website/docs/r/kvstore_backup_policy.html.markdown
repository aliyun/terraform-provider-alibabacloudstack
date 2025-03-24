---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_backup_policy"
sidebar_current: "docs-alibabacloudstack-resource-kvstore-backup-policy"
description: |-
  Provides a backup policy for ApsaraDB Redis / Memcache instance resource.
---

# alibabacloudstack_kvstore_backup_policy

Provides a backup policy for ApsaraDB Redis / Memcache instance resource. 

## Example Usage

Basic Usage

```
variable "creation" {
  default = "KVStore"
}
variable "multi_az" {
  default = "false"
}
variable "name" {
  default = "kvstorebackuppolicyvpc"
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alibabacloudstack_kvstore_instance" "default" {
  instance_class = "Memcache"
  instance_name  = "${var.name}"
  vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "memcache.master.small.default"
  
}
resource "alibabacloudstack_kvstore_backup_policy" "default" {
  instance_id   = "${alibabacloudstack_kvstore_instance.default.id}"
  backup_period = ["Tuesday", "Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of ApsaraDB for Redis or Memcache intance.
* `backup_time` - (Optional) Backup time, in the format of HH:mmZ- HH:mm Z 
* `preferred_backup_time` - (Optional) Preferred backup time, in the format of HH:mmZ- HH:mm Z 
* `backup_period` - (Optional) Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
* `preferred_backup_period` - (Optional) Preferred backup cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday 

## Attributes Reference

The following attributes are exported:

* `id` - The id of the backup policy.
* `instance_id` - The id of ApsaraDB for Redis or Memcache intance.
* `backup_time` - Backup time, in the format of HH:mmZ- HH:mm Z
* `preferred_backup_time` - Preferred backup time, in the format of HH:mmZ- HH:mm Z 
* `backup_period` - Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
* `preferred_backup_period` - Preferred backup cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday 

## Import

KVStore backup policy can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_kvstore_backup_policy.example r-abc12345678
```