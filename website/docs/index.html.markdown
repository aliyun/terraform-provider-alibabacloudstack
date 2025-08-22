---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  AlibabaCloudStack Provider Configuration Guide.
---

# AlibabacloudStack Provider

The AlibabacloudStack Provider is a Terraform plugin for managing resources in Alibaba Cloud ApsaraStack. 

Before using this provider, you must properly configure parameters to set up environment information and access credentials for Alibaba Cloud ApsaraStack.

## Parameter Configuration

The AlibabacloudStack Provider supports two configuration methods, which can be used together:

- **Static Configuration**  
- **Environment Variable Configuration**  

> **Note:**  
> If a parameter is configured using both methods, the **static configuration takes precedence over environment variable configuration**.

For configurable parameters, please refer to the **[Parameter Specifications](#parameter-specifications)** section.

## Parameter Configuration

### Static Configuration

+ Create `provider.tf` file

```hcl
# Declare AlibabacloudStack Provider source and version
terraform {
  required_providers {
    alibabacloudstack = {
      source  = "aliyun/alibabacloudstack"
      # Unspecified versions default to the latest release.
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
  proxy                  = "http://IP:Port"
  resource_group_set_name  = "Your Resource Group Set Name"
  popgw_domain             = "xxx.xxx.com"
  protocol                 = "HTTPS"
}
```

### Environment Variable Configuration

+ Terminal Environment Configuration

```shell
export ALIBABACLOUDSTACK_ACCESS_KEY="Your Access Key"
export ALIBABACLOUDSTACK_SECRET_KEY="Your Asecret Key"
export ALIBABACLOUDSTACK_ASSUME_ROLE_ARN = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
export ALIBABACLOUDSTACK_REGION="Region Name"
export ALIBABACLOUDSTACK_INSECURE= true
export ALIBABACLOUDSTACK_PROXY= "http://IP:Port"
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
```

## Parameter Specifications

### Environment Parameters

| Parameter Name       | Environment Variable              | Type     | Description                              | How to Obtain                                                                                 | Remarks                                                              |
|----------------------|------------------------------------|----------|------------------------------------------|------------------------------------------------------------------------------------------------|----------------------------------------------------------------------|
| popgw_domain         | ALIBABACLOUDSTACK_POPGW_DOMAIN    | string   | Alibaba Cloud ApsaraStack Service Endpoint Suffix   | Apsara Uni-manager Operations Console >> Top Profile Icon >> *User Information* >> **Apsara Stack API calls** >> **Internet Domain** | **Required**                                                        |
| region               | ALIBABACLOUDSTACK_REGION          | string   | Platform Region information              | Apsara Uni-manager Operations Console >> Top Region Information                                                        | **Required**                                                        |
| is_center_region     | ALIBABACLOUDSTACK_CENTER_REGION   | bool     | Specifies whether current region is central region | Apsara Uni-manager Operations Console >> Top Profile Icon >> *User Information* >> **Apsara Stack API calls** >> **Central Region or Not** | Default: `true`                                                     |
| protocol             | ALIBABACLOUDSTACK_PROTOCOL        | string   | The Network protocol used to access the environment      | Determined by environment configuration                                                       | Default: `HTTP`    (Valid values: `HTTP` or `HTTPS`)                                                 |
| insecure             | ALIBABACLOUDSTACK_INSECURE        | bool     | Specifies whether to ignore insecure HTTPS certificates      | Determined by environment configuration                                                       | Default: `false`<br>Effective only when protocol is `HTTPS`          |
| proxy                | ALIBABACLOUDSTACK_PROXY           | string   | The proxy endpoint that is used to access the environment    | Determined by environment configuration                                                       |                                                                      |

### Credential Parameters

AlibabacloudStack supports two authentication methods: STS Token and AK/SK. **STS Token authentication is generally recommended**.

- STS Token Authentication
  STS Tokens are temporary credentials obtained through role assumption, valid only within their expiration period. Two acquisition methods are available:
  - **Automatic STS Token Retrieval**  
    When configuring `access_key`, `secret_key`, and `role_arn`, AlibabacloudStack automatically performs role assumption and authenticates using the generated STS Token.
  - **Manual STS Token Acquisition**  
    Manually call the "Sts 2015-04-01 AssumeRole" API to assume roles. Use the obtained temporary credentials (`access_key`, `secret_key`, and `security_token`) for STS Token authentication.

- AK/SK Authentication
  AK/SK (Access Key/Secret Key) are permanent credentials that pose **significant security risks** if compromised:
  - **Account AK/SK Authentication**  
    Activated when configuring `access_key` and `secret_key` directly.


Detailed Configuration:

#### Automatic STS Token Retrieval

> **Note:**  
> - Configure the `role_arn` parameter to enable automatic role assumption and generate STS Token for authentication.  
> - The `role_arn` can be obtained by accessing the Apsara Uni-manager Management Console >> Top Profile Icon >> *User Information* >> *View Current Role Policy*, and Switch organizations in *Organization* and retrieve the **RAM Role** value.  
> - Account AK/SK parameters can be queried on the Apsara Uni-manager Management Console >> Top Profile Icon >> *User Information* >> **AccessKey Pair** .
> - When resource set names within a Region lack uniqueness, implement the `department` and `resource_group parameters` as replacements for `resource_group_set_name`.


| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| role_arn               | ALIBABACLOUDSTACK_ASSUME_ROLE_ARN | string   | RAM Role to be assumed          | **Required**, e.g. `acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx` |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |

#### Manual STS Token Acquisition


> **Note**:  
> - Provider identifies AK/SK type by `security_token` configuration.  
> - Manually call the "Sts 2015-04-01 AssumeRole" API with `access_key`, `secret_key`, and `role_arn` to complete role assumption. Configure the obtained STS Token.  
> - When resource set names within a Region lack uniqueness, implement the `department` and `resource_group` parameters as replacements for `resource_group_set_name`.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Temporary Access Key            | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Temporary Secret Key            | **Required**                                                     |
| security_token         | ALIBABACLOUDSTACK_SECURITY_TOKEN  | string   | Temporary Security Token        | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |

#### Account AK/SK Authentication

> **Note**:  
> - Account AK/SK parameters can be queried on the Apsara Uni-manager Management Console >> Top Profile Icon >> *User Information* >> **AccessKey Pair** .
> - When resource set names within a Region lack uniqueness, implement the `department` and `resource_group` parameters as replacements for `resource_group_set_name`.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|-----------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |
