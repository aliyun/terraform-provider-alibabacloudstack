---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_quotas"
sidebar_current: "docs-apsarastack-datasource-ascm-quotas"
description: |-
    Provides a list of quota to the user.
---

# apsarastack\_ascm_quotas

This data source provides the quota of the current Apsara Stack Cloud user.

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
    name = "Dummy_Test_1"
}

resource "apsarastack_ascm_quota" "default" {
    quota_type = "organization"
    quota_type_id = apsarastack_ascm_organization.default.parent_id
    product_name = "RDS"
    total_cpu = 1500
    total_disk = 1500
    total_mem = 1500
    target_type = "MySql"
}

data "apsarastack_ascm_quotas" "default" {
    quota_type = "organization"
    quota_type_id = apsarastack_ascm_organization.default.parent_id
    product_name = "RDS"
    target_type = "MySql"
    output_file = "Rds_quota"
}
output "quota" {
    value = data.apsarastack_ascm_quotas.default.*
}
```

## Argument Reference

The following arguments are supported:

  * `product_name` - (Required) The name of the service. Valid values: ECS, OSS, VPC, RDS, SLB, ODPS, GPDB, DDS, R-KVSTORE, and EIP.
  * `quota_type` - (Required) The type of the quota. Valid values: organization and resourceGroup.
  * `quota_type_id` - (Required) The ID of the quota type. Specify an organization ID when the QuotaType parameter is set to organization. Specify a resource set ID when the QuotaType parameter is set to resourceGroup.
  * `cluster_name` - (Optional) The name of the cluster. This reserved parameter is optional and can be left empty.
  * `target_type` - (Optional) This reserved parameter is optional and can be left empty. It will be used only for some products. Products where target_type are required with their values - RDS ("MySql"), R-KVSTORE ("redis") and DDS ("mongodb").
  * `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `quotas` - A list of Quota. Each element contains the following attributes:
  * `id` - ID of the quota.
  * `quota_type` - Name of an organization, or a Resource Group.
  * `quota_type_id` - ID of an organization, or a Resource Group.
  * `total_vip_internal` - Total vip internal.
  * `region` - name of the region where product belong.
  * `total_vip_public` - Total vip public.
  * `total_vpc` - Total Vpc.
  * `total_cpu` - Total Cpu.
  * `total_cu` - Total Cu.
  * `total_disk` - Total Disk.
  * `total_mem` - Total Mem.
  * `used_mem` - Consumed Mem.
  * `total_gpu` - Total Gpu.
  * `total_amount` - Total Amount.
  * `total_disk_cloud_ssd` - Total disk cloud ssd.
  * `used_disk` - Consumed Disk.
  * `allocate_disk` - Allocated Disk.
  * `allocate_cpu` - Allocated Cpu.
  * `total_eip` - Total Eip.
  * `total_disk_cloud_efficiency` - Total disk cloud efficiency.