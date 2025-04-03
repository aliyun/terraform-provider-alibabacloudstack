---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_dbinstance"
sidebar_current: "docs-Alibabacloudstack-polardb-dbinstance"
description: |-
  Provides a polardb Dbinstance resource.
---

# alibabacloudstack_polardb_dbinstance

Provides a polardb Dbinstance resource.

## Example Usage
```
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}



resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}


resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}



variable "name" {
	default = "tf-testaccdbinstanceconfig"
}

resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}


resource "alibabacloudstack_polardb_dbinstance" "default" {
  instance_storage = "5"
  instance_name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type = "local_ssd"
  engine = "MySQL"
  engine_version = "5.7"
  instance_type = "rds.mysql.t1.small"
}
```

## Argument Reference

The following arguments are supported:
  * `engine` - (Required, ForceNew) - Database type. Return value:* MySQL* PostgreSQL* SQLServer* MariaDB
  * `engine_version` - (Required, ForceNew) - Database version.
  * `zone_id_slave1` - (Optional) - The ID of the availability zone for the first slave instance.
  * `zone_id_slave2` - (Optional) - The ID of the availability zone for the second slave instance.
  * `tde_status` - (Optional) - The Transparent Data Encryption (TDE) status of the instance.
  * `enable_ssl` - (Optional, ForceNew) - Whether SSL is enabled for the instance.
  * `storage_type` - (Optional, ForceNew) - Field `storage_type` is deprecated and will be removed in a future release. Please use new field `db_instance_storage_type` instead.
  * `db_instance_storage_type` - (Optional, ForceNew) - The storage type of the instance. Valid values:* **local_ssd**, **ephemeral_ssd**: local SSD disk.* **cloud_ssd**:SSD disk.* **cloud_essd**:ESSD cloud disk.
  * `encryption_key` - (Optional) - The key ID of the disk encryption in the same region. This parameter indicates that cloud disk encryption is enabled and cannot be turned off after it is enabled, and **RoleARN** is required * *.You can view the key ID in the key management service console or create a new key. For more information, see [Create Key](~~ 181610 ~~).
  * `encryption` - (Optional, ForceNew) -  Whether encryption is enabled for the instance.
  * `instance_type` - (Optional) - Field `instance_type` is deprecated and will be removed in a future release. Please use new field `db_instance_class` instead.
  * `db_instance_class` - (Optional) - Instance type. For more information, see [instance type table](~~ 26312 ~~).
  * `instance_storage` - (Optional) - Field `instance_storage` is deprecated and will be removed in a future release. Please use new field `db_instance_storage` instead.
  * `db_instance_storage` - (Optional) - The storage capacity of the database instance.
  * `instance_charge_type` - (Optional) - Field `instance_charge_type` is deprecated and will be removed in a future release. Please use new field `payment_type` instead.
  * `payment_type` - (Optional) - The instance payment method. Valid values:**PayAsYouGo**: PayAsYouGo.**Subscription**: Subscription.**Serverless**:Serverless paid type, which only supports MySQL instances.
  * `period` - (Optional) - Specify the subscription instance as the year or month type. Valid values:Year: Year of Package.Month: monthly.
  * `monitoring_period` - (Optional) - The monitoring period for the instance.
  * `auto_renew` - (Optional) - Whether the instance is automatically renewed. It is only passed in when creating a subscription instance. Valid values:-**true**-**false**> * If you purchase on a monthly basis, the automatic renewal period is 1 month. * If you purchase on an annual basis, the automatic renewal period is 1 year.
  * `auto_renew_period` - (Optional) - The auto-renewal period for the instance.
  * `zone_id` - (Optional, ForceNew) - The ID of the availability zone for the instance.
  * `vswitch_id` - (Optional, ForceNew) - The ID of the VSwitch for the instance.
  * `instance_name` - (Optional) - Field `instance_name` is deprecated and will be removed in a future release. Please use new field `db_instance_description` instead.
  * `db_instance_description` - (Optional) - The description of the database instance.
  * `security_ip_mode` - (Optional) - The security IP mode for the instance.
  * `maintain_time` - (Optional) - Instance maintenance time period, which is the UTC time, +8 hours is the maintenance time period displayed on the console.
  * `role_arn` - (Optional) - The Global Resource Descriptor (ARN) that authorizes the POLARDB cloud service account to access KMS. You can view the ARN information by using the [CheckCloudResourceAuthorized](~~ 446261 ~~) interface.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `storage_type` - The storage type of the instance. Valid values: Valid values:* **local_ssd**, **ephemeral_ssd**: local SSD disk.* **cloud_ssd**:SSD disk.* **cloud_essd**:ESSD cloud disk.
  * `db_instance_storage_type` - The storage type of the instance. Valid values:* **local_ssd**, **ephemeral_ssd**: local SSD disk.* **cloud_ssd**:SSD disk.* **cloud_essd**:ESSD cloud disk.
  * `instance_type` - The type of the instance.
  * `db_instance_class` - Instance type. For more information, see [instance type table](~~ 26312 ~~).
  * `instance_storage` - The storage capacity of the instance.
  * `db_instance_storage` -  The storage capacity of the database instance.
  * `instance_charge_type` - The charge type of the instance.
  * `payment_type` - The instance payment method. Valid values:**PayAsYouGo**: PayAsYouGo.**Subscription**: Subscription.**Serverless**:Serverless paid type, which only supports MySQL instances.
  * `monitoring_period` -  The monitoring period for the instance.
  * `zone_id` - The ID of the availability zone for the instance.
  * `instance_name` - The name of the instance.
  * `db_instance_description` - The description of the database instance.
  * `connection_string` - The connection string for the instance.
  * `port` - Connection port.
  * `maintain_time` - Instance maintenance time period, which is the UTC time, +8 hours is the maintenance time period displayed on the console.
  * `role_arn` - The Global Resource Descriptor (ARN) that authorizes the POLARDB cloud service account to access KMS. You can view the ARN information by using the [CheckCloudResourceAuthorized](~~ 446261 ~~) interface.
