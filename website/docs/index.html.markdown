---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  The AlibabacloudStack provider is used to interact with many resources supported by AlibabacloudStack. The provider needs to be configured with the proper credentials before it can be used.
---

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
  access_key               = var.access_key
  secret_key               = var.secret_key
  region                   = var.region
  role_arn                 = var.role_arn
  # security_token         = var.security_token
  insecure                 = true
  proxy                    = var.proxy
  resource_group_set_name  = var.resource_group_set_name
  popgw_domain             = var.domain
  protocol                 = "HTTPS"
}
```

### Environment Variable Configuration

> The Provider supports configuring most parameters through environment variables.  
> Basic environment variables such as `ALIBABACLOUDSTACK_ACCESS_KEY` and `ALIBABACLOUDSTACK_SECRET_KEY` provide platform access credentials for the AlibabacloudStack Provider.  
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
terraform plan
```

## Parameter Specifications

### Environment Parameters

| Parameter Name       | Environment Variable              | Type     | Description                              | How to Obtain                                                                                 | Remarks                                                              |
|----------------------|------------------------------------|----------|------------------------------------------|------------------------------------------------------------------------------------------------|----------------------------------------------------------------------|
| popgw_domain         | ALIBABACLOUDSTACK_POPGW_DOMAIN    | string   | Private cloud platform standard suffix   | ASO Platform >> Top Profile Icon >> *Personal Info* >> **Private Cloud API Usage** >> `Internet Domain` | **Required**                                                        |
| region               | ALIBABACLOUDSTACK_REGION          | string   | Platform Region information              | ASO Platform >> Top Region Information                                                        | **Required**                                                        |
| is_center_region     | ALIBABACLOUDSTACK_CENTER_REGION   | bool     | Whether current region is central region | ASO Platform >> Top Profile Icon >> *Personal Info* >> **Private Cloud API Usage** >> `Is Current Region a Central Region` | Default: `true`                                                     |
| protocol             | ALIBABACLOUDSTACK_PROTOCOL        | string   | Network protocol (`HTTP` or `HTTPS`)     | Determined by environment configuration                                                       | Default: `HTTP`                                                     |
| insecure             | ALIBABACLOUDSTACK_INSECURE        | bool     | Skip HTTPS certificate verification      | Determined by environment configuration                                                       | Default: `false`<br>Effective only when protocol is `HTTPS`          |
| proxy                | ALIBABACLOUDSTACK_PROXY           | string   | Proxy server address                     | Determined by environment configuration                                                       |                                                                      |

### Credential Parameters

> AlibabacloudStack Provider supports multiple credential types. Choose based on requirements.

#### 1. Account Role Assumption

> **Note**:  
> - Provider determines whether role assumption is needed based on `role_arn` configuration.  
> - The `role_arn` can be obtained by accessing the *User Information* page in ASCM, clicking **View Current Role Policy**, and retrieving the `RAM Role` for the user under a specific organization.
> - Role assumption validity period: 3600 seconds.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| role_arn               | ALIBABACLOUDSTACK_ASSUME_ROLE_ARN | string   | RAM Role to be assumed          | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name      |   

#### 2. Account STS Token

> **Note**:  
> - Provider identifies AK/SK type by `security_token` configuration.  
> - Requires invoking the "Sts 2015-04-01 AssumeRole" API to obtain temporary credentials through role assumption.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Temporary Access Key            | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Temporary Secret Key            | **Required**                                                     |
| security_token         | ALIBABACLOUDSTACK_SECURITY_TOKEN  | string   | Temporary Security Token        | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name      |                                                                  |

#### 3. Account AK/SK

> **Note**:  
> - AK/SK parameters can be queried in ASCM interface.  
> - Use `department` and `resource_group` instead of `resource_group_set_name` when resource group names are not unique.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|------------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name      |                                                                  |
