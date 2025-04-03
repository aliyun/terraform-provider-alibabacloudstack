---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_dbinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-dbinstances"
description: |-
  Provides a list of polardb dbinstances owned by an alibabacloudstack account.
---

# alibabacloudstack_polardb_dbinstances
-> **NOTE:** Alias name has: `alibabacloudstack_polardb_instances`

This data source provides a list of polardb dbinstances in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
		default = "PolarDB"
}

resource "alibabacloudstack_polardb_instance" "default" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "${var.name}"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_polardb_instances" "default" {
  db_instance_id        = "${alibabacloudstack_polardb_instance.default.id}"
  db_instance_class = "${alibabacloudstack_polardb_instance.default.db_instance_class}"
  status     = "Running"
  region_id  = "zhe env region id"
}

```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of DB instance IDs to filter by.
  * `network_type` - (Optional) - The network type of the instance. Valid values:* **Classic**: Classic network.* **VPC**: VPC.
  * `db_instance_type` - (Optional) - The type of the DB instance.
  * `vswitch_id` - (Optional) - The ID of the VSwitch.
  * `vpc_id` - (Optional) - The Id of the VPC.
  * `db_instance_class` - (Optional) - Instance type. For more information, see [instance type table](~~ 26312 ~~).
  * `payment_type` - (Optional) - The instance payment method. Valid values:**PayAsYouGo**: PayAsYouGo.**Subscription**: Subscription.**Serverless**:Serverless paid type, which only supports MySQL instances.
  * `status` - (Optional) - The status of the resource
  * `db_instance_id` - (Optional) - The ID of the instance.
  * `engine_version` - (Optional) - Database version.
  * `resource_group_id` - (Optional) - The ID of the resource group.
  * `region_id` - (Required) - The region ID of the resource
  * `engine` - (Optional) - Database type. Valid value:* MySQL* PostgreSQL* SQLServer* MariaDB
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `db_instances` - A list of DB instances. Each element in the list is a map with the following keys:
    * `id` - The ID of the DB instance.
    * `auto_pay` - Whether to pay automatically. Value range:-**true**: automatic payment. You need to ensure that your account balance is sufficient.-**false**: only orders are generated without deduction.> The default value is true. If the balance of your payment method is insufficient, you can set the AutoPay parameter to false. In this case, unpaid orders will be generated. You can log on to the POLARDB management console to pay by yourself.>>
    * `auto_renew` - Whether the instance is automatically renewed. It is only passed in when creating a subscription instance. Valid values:-**true**-**false**> * If you purchase on a monthly basis, the automatic renewal period is 1 month. * If you purchase on an annual basis, the automatic renewal period is 1 year.
    * `auto_upgrade_minor_version` - How to upgrade the minor version of the instance. Valid values:* **Auto**: automatically upgrade the minor version.* **Manual**: It is not automatically upgraded. It is only mandatory when the current version is offline.
    * `business_info` - Business extension parameters.
    * `category` - Instance series. Valid values:* **Basic**: Basic Edition.* **HighAvailability**: High availability version.* **AlwaysOn**: Cluster version.* **Finance**: Three-node Enterprise Edition.* **serverless_basic**:Serverless Basic edition.
    * `classic_expired_days` - The number of days to retain the classic network connection address. Valid values: **1 to 120**. Unit: days.
    * `commodity_code` - Commodity code
    * `connection_mode` - The access mode of the instance. Valid values:* **Standard**: The Standard access mode.* **Safe**: Database proxy mode.The default value is the POLARDB system allocation.> SQL Server 2012, 2016, and 2017 only support the standard access mode.
    * `connection_string_prefix` - The read-only address prefix name, which cannot be repeated. It consists of lowercase letters and underscores. It must start with a letter and be no longer than 30 characters in length.> by default, the prefix is formed by the string of "instance name + rw.
    * `connection_string_type` - Connection address type, value:* **Normal**: Normal connection* **ReadWriteSplitting**: read/write splitting connectionAll connections are returned by default.
    * `current_connection_string` - The current connection address of the instance, which can be an internal or external network connection address or a classic network connection address in mixed access mode.
    * `db_instance_description` - The description of the DB instance.
    * `db_instance_storage` - The storage capacity of the DB instance.
    * `db_instance_type` - The type of the DB instance.
    * `db_instance_class` - Instance type. For more information, see [instance type table](~~ 26312 ~~).
    * `db_instance_id` - The ID of the instance.
    * `db_instance_net_type` - Network type, value:* **Intranet**: Intranet.
    * `db_instance_storage_type` - The storage type of the instance. Valid values:* **local_ssd**, **ephemeral_ssd**: local SSD disk.* **cloud_ssd**:SSD disk.* **cloud_essd**:ESSD cloud disk.
    * `distribution_type` - Read weight allocation mode, value:-**Standard**: The weight is automatically assigned by specification.-**Custom**: Custom Allocation weight
    * `effective_time` - Effective time, value:* **Immediate**: Effective immediately.* **Mainaintime**: takes effect within the O & M period. For more information, see [modifydbinstancemaintime](~~ 26249 ~~).Default value: **Immediate * *.
    * `encryption_key` - The key ID of the disk encryption in the same region. This parameter indicates that cloud disk encryption is enabled and cannot be turned off after it is enabled, and **RoleARN** is required * *.You can view the key ID in the key management service console or create a new key. For more information, see [Create Key](~~ 181610 ~~).
    * `engine` - Database type. Return value:* MySQL* PostgreSQL* SQLServer* MariaDB
    * `engine_version` - Database version.
    * `expire_time` - Expiration time. <I> yyyy-MM-dd</I> T <I> HH:mm:ss</I> Z(UTC time).> Pay-as-you-go instances have no expiration time.
    * `guard_db_instance_id` - The ID of the guard DB instance.
    * `lock_mode` - Instance lock mode. Value:* **Unlock**: normal.* **ManualLock**: manually triggers the lock.* **LockByExpiration**: The instance is automatically locked when it expires.* **Lockbyrevolution**: automatically locked before instance rollback.* **LockByDiskQuota**: The instance is automatically locked when the space is full.* **LockReadInstanceByDiskQuota**: The read-only instance is automatically locked when the space is full.
    * `lock_reason` - Lock the cause.
    * `maintain_time` - Instance maintenance time period, which is the UTC time, +8 hours is the maintenance time period displayed on the console.
    * `master_instance_id` - The ID of the primary instance. If this parameter is not returned, the instance is the primary instance.
    * `max_delay_time` - The delay threshold, ranging from 0 to 7200, in seconds, is 30 by default.> When the delay of a read-only instance exceeds this threshold, read traffic is not sent to the instance.
    * `network_type` - The network type of the instance. Valid values:* **Classic**: Classic network.* **VPC**: VPC.
    * `payment_type` - The instance payment method. Valid values:**PayAsYouGo**: PayAsYouGo.**Subscription**: Subscription.**Serverless**:Serverless paid type, which only supports MySQL instances.
    * `period` - Specify the subscription instance as the year or month type. Valid values:Year: Year of Package.Month: monthly.
    * `port` - Connection port.
    * `private_ip_address` - No configuration is required to indicate the intranet IP address of the target instance. By default, the system automatically allocates VPCId and vSwitchId.
    * `read_write_splitting_classic_expired_days` - The number of days to retain the classic network address for read/write splitting. Valid values: **1-120**. Unit: days. Default value: **7 * *.> This parameter is valid when the instance has a classic network type read/write splitting address and **RetainClassic** = **True.
    * `read_write_splitting_private_ip_address` - The IP address of the intranet read/write splitting address of the instance must be within the IP address range of the specified switch. By default, the system automatically allocates by using **VPCId** and **VSwitchId.> This value is valid when the current instance has a read/write splitting address of the classic network type.
    * `record_total` - Total number of records.
    * `region_id` - The region ID of the resource
    * `resource_group_id` - The ID of the resource group.
    * `resource_type` - Resource type definition. Unique value: **INSTANCE * *.
    * `retain_classic` - Whether to retain the classic network address. Valid values:* **True**: Reserved* **False**: not reservedDefault value: **False * *.
    * `role_arn` - The Global Resource Descriptor (ARN) that authorizes the POLARDB cloud service account to access KMS. You can view the ARN information by using the [CheckCloudResourceAuthorized](~~ 446261 ~~) interface.
    * `security_ip_list` - A list of security IP addresses.
    * `security_ip_mode` - The mode of security IP addresses.
    * `sql_collector_status` - Enable or disable SQL insight (SQL audit). Valid values: **Enable | Disabled * *
    * `status` - The status of the resource
    * `table_meta` - Specifies the restored library table. Format:'''[{"type":"db","name":"<database 1 name>","newname":"<new database 1 name>","tables":[{"type":"table","name":"<table 1 name in database 1>","newname":"<New table 1 name>" },{ "type":"table","name":"<table 2 name in database 1>","newname":"<New table 2 name>"}]},{ "type":"db","name":"<Database 2 name>","newname":"<new database 2 name>","tables":[{"type":"table","name":"<table 3 name in database 2>","newname":"<new table 3 name>"},{"type":"table","name":"<table 4 name in database 2>","newname":"<New Table 4 name>"}]}]'''
    * `tags` - The tag of the resource
    * `temp_db_instance_id` - The ID of the temporary DB instance.
    * `used_time` - Specify the duration of the purchase. Valid values:When the parameter Period = Year, the value of UsedTime is 1 to 5.When Period = Month, the value of UsedTime is 1 to 11.
    * `vswitch_id` - The ID of the VSwitch.
    * `vpc_cloud_instance_id` - The ID of the VPC instance.
    * `vpc_id` -  The ID of the VPC.
    * `weight` - Read weight allocation, that is, the ratio of read requests to the primary instance and read-only instance. Increment by 100 with a maximum value of 10000.* POLARDB instance format: '{"<read-only instance ID>":<weight>,"master":<weight>,"slave":<weight>}'* MyBASE instance format: '[{"instanceName":"<master instance ID>","weight":<weight>,"role":"master"},{"instanceName":"<master instance ID>","weight":<weight>,"role":"slave"},{"instanceName":"<read-only instance ID>","weight":<weight>,"role":"master"}]'>-This parameter must be specified when **DistributionType** is **Custom.>-When **DisrtibutionType** is set to **Standard**, this parameter is invalid.
    * `zone_id_slave_one` - This parameter can be configured only when the original instance is a high-availability version. This parameter indicates the available zone ID of the target instance.POLARDB PostgreSQL allows you to configure the new standby instance to another zone in the same region as the original instance.You can view the zone ID by using the [DescribeRegions](~~ 26243 ~~) operation.
    * `zone_id_slave_two` - Standby Zone 2
