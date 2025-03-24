---
subcategory: "ACK"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_account"
sidebar_current: "docs-Alibabacloudstack-datasource-account"
description: |- 
  Provides Id of AlibabacloudStack account.
---

# alibabacloudstack_account

This data source provides the ID of the Alibaba Cloud Stack account.

## Example Usage

```hcl
data "alibabacloudstack_account" "current" {}

output "account_id" {
  value = data.alibabacloudstack_account.current.id
}
```

## Argument Reference
This data source has no configurable arguments. It retrieves information based on the authenticated account.

## Attributes Reference
The following attributes are exported:

`id` - The unique identifier of the Alibaba Cloud Stack account.

## Common Notes
Ensure that your provider configuration is correct and includes the necessary credentials.
This data source can be used to fetch various details about your Alibaba Cloud Stack account.

## Notes
This data source is primarily used to fetch the account ID for the currently authenticated user.
Ensure that the provider is properly configured with valid credentials to access the Alibaba Cloud Stack API.

## Import
The alibabacloudstack_account data source does not support import as it directly retrieves the account information based on the authenticated session.