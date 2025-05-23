---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_secrets"
sidebar_current: "docs-alibabacloudstack-datasource-kms-secrets"
description: |-
    Provides a list of available KMS Secrets.
---

# alibabacloudstack_kms_secrets

This data source provides a list of KMS Secrets in an Alibabacloudstack Cloud account according to the specified filters.
 

## Example Usage

```
# Declare the data source
data "alibabacloudstack_kms_secrets" "kms_secrets_ds" {
  fetch_tags = true
  name_regex = "name_regex"
  tags = {
    "k-aa" = "v-aa",
    "k-bb" = "v-bb"
  }
}

output "first_secret_id" {
  value = "${data.alibabacloudstack_kms_secrets.kms_secrets_ds.secrets.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `fetch_tags` - (Optional) Whether to include the predetermined resource tag in the return value. Default to `false`.
* `ids` - (Optional) A list of KMS Secret ids. The value is same as KMS secret_name.
* `name_regex` - (Optional) A regex string to filter the results by the KMS secret_name.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `names` - (Optional) A list of KMS Secret names.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of Kms Secret ids. The value is same as KMS secret_name. 
* `names` -  A list of KMS Secret names.
* `secrets` - A list of KMS Secrets. Each element contains the following attributes:
  * `id` - ID of the Kms Secret. The value is same as KMS secret_name.
  * `secret_name` - Name of the KMS Secret.
  * `planned_delete_time` - Schedule deletion time.
  * `tags` - A mapping of tags to assign to the resource.