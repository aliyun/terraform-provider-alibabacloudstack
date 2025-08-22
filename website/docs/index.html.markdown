---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  The AlibabacloudStack provider is used to interact with many resources supported by Alibaba Cloud ApsaraStack. The provider must be configured with valid authentication credentials before it can be used.
---

# AlibabacloudStack Provider

The AlibabacloudStack provider is used to interact with resources supported by Alibaba Cloud ApsaraStack. The provider must be configured with valid authentication credentials before it can be used.

## Example Code

### Static Configuration

```hcl
# Declare AlibabacloudStack Provider source and version
terraform {
  required_providers {
    alibabacloudstack = {
      source  = "aliyun/alibabacloudstack"
      # version = ">= 3.18.0"
    }
  }
}

# Configure AlibabacloudStack Provider
provider "alibabacloudstack" {
  access_key               = "Your Access Key"
  secret_key               = "Your Secret Key"
  role_arn                 = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
  # security_token         = "Your STS Token"
  region                   = "Region Name"
  insecure                 = true
  # proxy                  = "http://IP:Port"
  resource_group_set_name  = "Your Resource Group Set Name"
  popgw_domain             = "xxx.xxx.com"
  protocol                 = "HTTPS"
}
```

### Environment Variable Configuration

> The Provider supports configuring parameters through environment variables.  
> Environment variables such as `ALIBABACLOUDSTACK_ACCESS_KEY`, `ALIBABACLOUDSTACK_SECRET_KEY` and `ALIBABACLOUDSTACK_ASSUME_ROLE_ARN` provide platform access credentials for the AlibabacloudStack Provider.  
> For other configurable environment variables, please refer to the **[Parameter Specifications](#parameter-specifications)** section.


+ `main.tf` Configuration
```hcl
provider "alibabacloudstack" {
    resource_group_set_name ="${var.resource_group_set_name}"
}
```

+ Terminal Environment Configuration

```shell
export ALIBABACLOUDSTACK_ACCESS_KEY="Your Access Key"
export ALIBABACLOUDSTACK_SECRET_KEY="Your Asecret Key"
export ALIBABACLOUDSTACK_ASSUME_ROLE_ARN = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
export ALIBABACLOUDSTACK_REGION="Region Name"
export ALIBABACLOUDSTACK_INSECURE= true
export ALIBABACLOUDSTACK_PROXY= "http://IP:Port"
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
terraform plan
```

## Parameter Specifications

### Environment Parameters

| Parameter Name       | Environment Variable              | Type     | Description                              | How to Obtain                                                                                 | Remarks                                                              |
|----------------------|------------------------------------|----------|------------------------------------------|------------------------------------------------------------------------------------------------|----------------------------------------------------------------------|
| popgw_domain         | ALIBABACLOUDSTACK_POPGW_DOMAIN    | string   | AlibabaCloud ApsaraStack platform Service Endpoint suffix   | Apsara Uni-manager Operations Console >> Top Profile Icon >> *User Information* >> **Apsara Stack API calls** >> **Internet Domain** | **Required**                                                        |
| region               | ALIBABACLOUDSTACK_REGION          | string   | Platform Region information              | Apsara Uni-manager Operations Console >> Top Region Information                                                        | **Required**                                                        |
| is_center_region     | ALIBABACLOUDSTACK_CENTER_REGION   | bool     | Specifies whether current region is central region | Apsara Uni-manager Operations Console >> Top Profile Icon >> *User Information* >> **Apsara Stack API calls** >> **Central Region or Not** | Default: `true`                                                     |
| protocol             | ALIBABACLOUDSTACK_PROTOCOL        | string   | The Network protocol used to access the environment      | Determined by environment configuration                                                       | Default: `HTTP`    (Valid values: `HTTP` or `HTTPS`)                                                 |
| insecure             | ALIBABACLOUDSTACK_INSECURE        | bool     | Specifies whether to ignore insecure HTTPS certificates      | Determined by environment configuration                                                       | Default: `false`<br>Effective only when protocol is `HTTPS`          |
| proxy                | ALIBABACLOUDSTACK_PROXY           | string   | The proxy endpoint that is used to access the environment    | Determined by environment configuration                                                       |                                                                      |

### Credential Parameters

> AlibabacloudStack Provider supports multiple credential types. Choose based on requirements.

#### 1. Role Assumption

> **Note**:  
> - Provider determines whether role assumption is needed based on `role_arn` configuration.  
> - The `role_arn` can be obtained by accessing the *User Information* page in Apsara Uni-manager Management Console, clicking *View Current Role Policy*, and retrieving the **RAM Role** for the user under a specific organization.
> - Role assumption validity period: 3600 seconds.
> - When resource set names within a Region lack uniqueness, implement the `department` and `resource_group parameters` as replacements for `resource_group_set_name`.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| role_arn               | ALIBABACLOUDSTACK_ASSUME_ROLE_ARN | string   | RAM Role to be assumed          | **Required**, e.g. `acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx` |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |

#### 2. Account STS Token

> **Note**:  
> - Provider identifies AK/SK type by `security_token` configuration.  
> - Requires invoking the "Sts 2015-04-01 AssumeRole" API to obtain temporary credentials through role assumption.
> - When resource set names within a Region lack uniqueness, implement the `department` and `resource_group` parameters as replacements for `resource_group_set_name`.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Temporary Access Key            | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Temporary Secret Key            | **Required**                                                     |
| security_token         | ALIBABACLOUDSTACK_SECURITY_TOKEN  | string   | Temporary Security Token        | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |

#### 3. Account AK/SK

> **Note**:  
> - AK/SK parameters can be queried in Apsara Uni-manager Management Console.  
> - Use `department` and `resource_group` instead of `resource_group_set_name` when resource group names are not unique.
> - When resource set names within a Region lack uniqueness, implement the `department` and `resource_group` parameters as replacements for `resource_group_set_name`.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|-----------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |
