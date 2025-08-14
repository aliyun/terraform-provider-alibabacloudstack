---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  AlibabacloudStack Provider 用于通过 Terraform 管理阿里云私有云平台的各类资源。使用前需配置访问云平台的正确凭证。
---

# AlibabacloudStack Provider

AlibabacloudStack Provider 用于通过 Terraform 管理阿里云私有云平台下的多种资源。在使用前，需配置该 Provider 访问云平台的正确凭证。

## 示例代码

### 静态配置

```hcl
# 声明 AlibabacloudStack Provider 来源与版本
terraform {
  required_providers {
    alibabacloudstack = {
      source  = "aliyun/alibabacloudstack"
      # version = ">= 3.18.0"
    }
  }
}

# 配置 AlibabacloudStack Provider
provider "alibabacloudstack" {
  access_key               = "Your Access Key"
  secret_key               = "Your Secret Key"
  role_arn                 = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
  # security_token         = "Your STS Token"
  region                   = "Region Name"
  insecure                 = true
  # proxy                    = "http://IP:Port"
  resource_group_set_name  = "Your Resource Group Set Name"
  popgw_domain             = "xxx.xxx.com"
  protocol                 = "HTTPS"
}
```

### 环境变量配置

> Provider 支持通过环境变量配置大部分参数。  
> 基础环境变量如 `ALIBABACLOUDSTACK_ACCESS_KEY`、 `ALIBABACLOUDSTACK_SECRET_KEY`和`ALIBABACLOUDSTACK_ASSUME_ROLE_ARN` 用于为 AlibabacloudStack Provider 提供平台访问凭证。  
> 其他可配置环境请参考 **[参数说明](#参数说明)** 章节。

+ `main.tf`配置
```hcl
provider "alibabacloudstack" {
    resource_group_set_name ="Your Resource Group Set Name"
}
```

+ 执行终端配置

```shell
export ALIBABACLOUDSTACK_ACCESS_KEY="Your Access Key"
export ALIBABACLOUDSTACK_SECRET_KEY="Your Secret Key"
export ALIBABACLOUDSTACK_ASSUME_ROLE_ARN = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
export ALIBABACLOUDSTACK_REGION="Region Name"
export ALIBABACLOUDSTACK_INSECURE= true
export ALIBABACLOUDSTACK_PROXY= "http://IP:Port"
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
terraform plan
```

## 参数说明

### 环境参数

| 参数名            | 环境变量名                        | 参数类型 | 参数含义                                | 获取方式                                                                                     | 备注                                                                 |
|-------------------|-----------------------------------|----------|-----------------------------------------|----------------------------------------------------------------------------------------------|----------------------------------------------------------------------|
| popgw_domain      | ALIBABACLOUDSTACK_POPGW_DOMAIN    | string   | 阿里云专有云飞天企业版服务地址标准后缀        | ASO平台 >> 顶部个人头像 >> *个人信息* >> **专有云 API 调用使用** >> `Internet Domain`         | **必填**                                                             |
| region            | ALIBABACLOUDSTACK_REGION          | string   | 平台 Region 信息                        | ASO平台 >> 顶部 Region 信息                                                                  | **必填**                                                             |
| is_center_region  | ALIBABACLOUDSTACK_CENTER_REGION   | bool     | 当前 Region 是否为中心 Region           | ASO平台 >> 顶部个人头像 >> *个人信息* >> **专有云 API 调用使用** >> `当前 Region 是否为中心 Region` | 默认值为 `true`                                                      |
| protocol          | ALIBABACLOUDSTACK_PROTOCOL        | string   | 网络协议（可选 `HTTP` 或 `HTTPS`）      | 根据环境实际情况确定                                                                         | 默认值 `HTTP`                                                        |
| insecure          | ALIBABACLOUDSTACK_INSECURE        | bool     | 是否跳过 HTTPS 证书校验                 | 根据环境实际情况确定                                                                         | 默认值 `false`<br>仅当 protocol 为 `HTTPS` 时生效                     |
| proxy             | ALIBABACLOUDSTACK_PROXY           | string   | 代理服务器地址                          | 根据环境实际情况确定                                                                         |                                                                      |

---

### 凭证参数

> AlibabacloudStack Provider 支持多种访问凭证，请根据实际需求选择。

#### 1. 账号扮演

> **注意**：  
> - Provider 通过是否配置 `role_arn` 判断是不需要进行帐号扮演。
> - `role_arn` 可在 ASCM 平台的*个人信息*页，通过点击**查看当前角色策略**，获取用户在特定组织下的`RAM Role`。
> - 账号扮演有效时间为3600秒。
> - 若Region下资源集名称不唯一，需使用 `department` 和 `resource_group` 替代 `resource_group_set_name`。

| 参数名                  | 环境变量名                        | 参数类型 | 参数含义                | 备注                                                                 |
|-------------------------|-----------------------------------|----------|-------------------------|----------------------------------------------------------------------|
| access_key              | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | 账号 AK             | **必填**                                                             |
| secret_key              | ALIBABACLOUDSTACK_SECRET_KEY      | string   | 账号 SK             | **必填**                                                             |
| role_arn                | ALIBABACLOUDSTACK_ASSUME_ROLE_ARN  | string  | 待扮演的角色(Ram Role)  | **必填**                                                             |
| department              | ALIBABACLOUDSTACK_DEPARTMENT      | string   | 凭证登录时的组织        | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group          | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | 凭证登录时的资源集      | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group_set_name | ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | 凭证登录时的资源集名称  |                                                                      |

#### 2. 账号 STS Token

> **注意**：  
> - Provider 通过是否配置 `security_token` 判断 AK/SK 类型。  
> - 需通过API "Sts 2015-04-01 AssumeRole"接口使用角色扮演获取临时凭证。
> - 若Region下资源集名称不唯一，需使用 `department` 和 `resource_group` 替代 `resource_group_set_name`。

| 参数名                  | 环境变量名                        | 参数类型 | 参数含义                | 备注                                                                 |
|-------------------------|-----------------------------------|----------|-------------------------|----------------------------------------------------------------------|
| access_key              | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | 临时凭证 AK             | **必填**                                                             |
| secret_key              | ALIBABACLOUDSTACK_SECRET_KEY      | string   | 临时凭证 SK             | **必填**                                                             |
| security_token          | ALIBABACLOUDSTACK_SECURITY_TOKEN  | string   | 临时凭证 Token          | **必填**                                                             |
| department              | ALIBABACLOUDSTACK_DEPARTMENT      | string   | 凭证登录时的组织        | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group          | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | 凭证登录时的资源集      | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group_set_name | ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | 凭证登录时的资源集名称  |                                                                      |

#### 3. 账号 AK/SK

> **注意**：  
> - 账号 AK/SK 相关参数可在 ASCM 页面上查询。  
> - 若Region下资源集名称不唯一，需使用 `department` 和 `resource_group` 替代 `resource_group_set_name`。

| 参数名                  | 环境变量名                        | 参数类型 | 参数含义                | 备注                                                                 |
|-------------------------|-----------------------------------|----------|-------------------------|----------------------------------------------------------------------|
| access_key              | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | 账号 AK                 | **必填**                                                             |
| secret_key              | ALIBABACLOUDSTACK_SECRET_KEY      | string   | 账号 SK                 | **必填**                                                             |
| department              | ALIBABACLOUDSTACK_DEPARTMENT      | string   | 凭证登录时的组织        | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group          | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | 凭证登录时的资源集      | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group_set_name | ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | 凭证登录时的资源集名称  |                                                                      |