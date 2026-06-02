---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
description: |-
  AlibabaCloudStack Provider Configuration Guide.
---

# AlibabacloudStack Provider

The AlibabacloudStack Provider is a Terraform plugin for managing resources in Alibaba Cloud ApsaraStack. 

Before using this provider, you must properly configure parameters including environment parameters and credential parameters.

## Parameter Configuration

The AlibabacloudStack Provider supports two configuration methods, which can be used together:

- **Static Configuration**  
- **Environment Variable Configuration**  

> **Note:**  
> If one parameter is configured using both methods, the **static configuration takes precedence over environment variable configuration**.

For configurable parameters, please refer to the **[Parameter Specifications](#parameter-specifications)** section.

### Static Configuration

+ Create `terraform.tf` file

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

+ Create `provider.tf` file

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
| popgw_domain         | ALIBABACLOUDSTACK_POPGW_DOMAIN    | string   | Standard address suffix of AlibabaCloud ApsaraStack   | Apsara Uni-manager Operations Console >> Top Profile Icon >> *User Information* >> **Apsara Stack API calls** >> **Internet Domain** | **Required**                                                        |
| region               | ALIBABACLOUDSTACK_REGION          | string   | Platform Region information              | Apsara Uni-manager Operations Console >> Top Region Information                                                        | **Required**                                                        |
| is_center_region     | ALIBABACLOUDSTACK_CENTER_REGION   | bool     | Whether current region is central region | Apsara Uni-manager Operations Console >> Top Profile Icon >> *User Information* >> **Apsara Stack API calls** >> **Central Region or Not** | Default: `true`                                                     |
| protocol             | ALIBABACLOUDSTACK_PROTOCOL        | string   | The Network protocol used to access the environment      | Determined by the actual situation                                                       | Default: `HTTP`    (Valid values: `HTTP` or `HTTPS`)                                                 |
| insecure             | ALIBABACLOUDSTACK_INSECURE        | bool     | Whether to skip verification of HTTPS certificates      | Determined by the actual situation                                                       | Default: `false`<br>Effective only when protocol is `HTTPS`          |
| proxy                | ALIBABACLOUDSTACK_PROXY           | string   | The proxy endpoint that is used to access the environment    | Determined by the actual situation                                                       |                                                                      |
| oss_endpoints      | ALIBABACLOUDSTACK_OSS_ENDPOINTS    | map   | Cluster names and corresponding addresses required for accessing OSS storage services in the environment    | Alibaba Cloud Platform >> Products >> OSS Service >> Bucket Creation Page  | Data structure: {cluster name: endpoint of the cluster}|

---
### Credential Parameters

AlibabacloudStack supports two authentication methods: STS Token and AK/SK. **STS Token authentication is generally recommended**.

- **STS Token**: STS Token is a kind of  temporary credential. It can be obtained through role assumption, which is valid only within their expiration period. Two acquisition methods are available:

  - **Automatic STS Token Retrieval**  

    When configuring `access_key`, `secret_key`, and `role_arn`, AlibabacloudStack performs role assumption and  automatically generates the STS Token.

  - **Manual STS Token Acquisition**  

    You can also manually call the "Sts 2015-04-01 AssumeRole" API to assume roles, which can generate temporary credentials (`access_key`, `secret_key`, and `security_token`) for STS Token authentication.

- **AK/SK**: AK/SK (Access Key/Secret Key) are permanent credentials that pose **significant security risks** if compromised:

  - **Account AK/SK Authentication**  

    Activated when configuring `access_key` and `secret_key` directly.


Detailed Configuration:

#### Automatic STS Token Retrieval

> **Note:**  
> - Configure the `role_arn` parameter to enable automatic role assumption and generate STS Token for authentication.  
> - The `role_arn` can be obtained by accessing the Apsara Uni-manager Management Console >> Top Profile Icon >> *User Information* >> *View Current Role Policy* >>  *Organization* ,and switching organizations to get the **RAM Role** value.  
> - Access Key/Secret Key parameters can be queried on the Apsara Uni-manager Management Console >> Top Profile Icon >> *User Information* >> **AccessKey Pair** .
> - IF names of resource set within one region lack uniqueness, you need to configure the `department` and the `resource_group parameters` as replacements for the `resource_group_set_name` parameter.


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
> 
> - STS Token can be also generated by manually calling the "Sts 2015-04-01 AssumeRole" API with `access_key`, `secret_key`, and `role_arn` upstairs.  It consists of temporary credentials (access_key, secret_key, and security_token) for STS Token , which  can be configured as below.
> - IF names of resource set within one region lack uniqueness, you need to configure the `department` and the `resource_group parameters` as replacements for the `resource_group_set_name` parameter.

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
> - Access Key/Secret Key parameters can be queried on the Apsara Uni-manager Management Console >> Top Profile Icon >> *User Information* >> **AccessKey Pair** .
> - IF names of resource set within one region lack uniqueness, you need to configure the `department` and the `resource_group parameters` as replacements for the `resource_group_set_name` parameter.

| Parameter Name         | Environment Variable              | Type     | Description                     | Remarks                                                          |
|------------------------|-----------------------------------|----------|---------------------------------|------------------------------------------------------------------|
| access_key             | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | Account Access Key              | **Required**                                                     |
| secret_key             | ALIBABACLOUDSTACK_SECRET_KEY      | string   | Account Secret Key              | **Required**                                                     |
| department             | ALIBABACLOUDSTACK_DEPARTMENT      | string   | Authentication organization     | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group         | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | Authentication resource group   | Required if `resource_group_set_name` is unavailable/unconfigured |
| resource_group_set_name| ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | Resource group set name        | e.g. `ResourceSet(xxxx)`                                    |
