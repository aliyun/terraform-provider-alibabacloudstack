---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  AlibabacloudStack Provider配置说明。
---

# AlibabacloudStack Provider

AlibabacloudStack Provider作为Terraform的插件，用于管理阿里云专有云飞天企业版的各类资源。

在使用该Provider前，需正确配置参数来设定阿里云专有云飞天企业版环境信息与访问凭证。

## 参数配置

AlibabacloudStack Privder支持静态配置和环境变量配置两种方式，两种配置方式可以同时混用。

> **注意**：  
> 如果某参数同时使用两种配置方式，仅生效静态配置。

可配置的参数请参考 **[参数说明](#参数说明)** 章节。

### 静态配置

+ 创建`terraform.tf`配置文件

```hcl
# 声明 AlibabacloudStack Provider 来源与版本
terraform {
  required_providers {
    alibabacloudstack = {
      source  = "aliyun/alibabacloudstack"
      # 如果不声明version，则默认使用最新版本。
      # version = "< 3.19.0"
    }
  }
}

+ 创建`provider.tf`配置文件

# 配置 AlibabacloudStack Provider
provider "alibabacloudstack" {
  access_key               = "Your Access Key"
  secret_key               = "Your Secret Key"
  role_arn                 = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
  # security_token         = "Your STS Token"
  region                   = "Region Name"
  insecure                 = true
  proxy                    = "http://IP:Port"
  resource_group_set_name  = "Your Resource Group Set Name"
  popgw_domain             = "xxx.xxx.com"
  protocol                 = "HTTPS"
}
```

### 环境变量配置

+ 终端配置环境变量

```shell
export ALIBABACLOUDSTACK_ACCESS_KEY="Your Access Key"
export ALIBABACLOUDSTACK_SECRET_KEY="Your Secret Key"
export ALIBABACLOUDSTACK_ASSUME_ROLE_ARN="acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
export ALIBABACLOUDSTACK_REGION="Region Name"
export ALIBABACLOUDSTACK_INSECURE=true
export ALIBABACLOUDSTACK_PROXY="http://IP:Port"
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
```

## 参数说明

### 环境参数

| 参数名            | 环境变量名                        | 参数类型 | 参数含义                                | 获取方式                                                                                     | 备注                                                                 |
|-------------------|-----------------------------------|----------|-----------------------------------------|----------------------------------------------------------------------------------------------|----------------------------------------------------------------------|
| popgw_domain      | ALIBABACLOUDSTACK_POPGW_DOMAIN    | string   | 阿里云专有云飞天企业版服务地址标准后缀        | 专有云运维平台 >> 顶部个人头像 >> *个人信息* >> **专有云 API 调用使用** >> **Internet Domain**         | **必填**                                                             |
| region            | ALIBABACLOUDSTACK_REGION          | string   | 平台 Region 信息                        | 专有云运维平台 >> 顶部 Region 信息                                                                  | **必填**                                                             |
| is_center_region  | ALIBABACLOUDSTACK_CENTER_REGION   | bool     | 当前 Region 是否为中心 Region           | 专有云运维平台 >> 顶部个人头像 >> *个人信息* >> **专有云 API 调用使用** >> **当前 Region 是否为中心 Region** | 默认值为 `true`                                                      |
| protocol          | ALIBABACLOUDSTACK_PROTOCOL        | string   | 访问环境时使用的网络协议                  | 根据环境实际情况确定                                                                         | 默认值 `HTTP`, 可选 `HTTP` 或 `HTTPS`                                                        |
| insecure          | ALIBABACLOUDSTACK_INSECURE        | bool     | 访问环境时是否跳过 HTTPS 证书校验          | 根据环境实际情况确定                                                                         | 默认值 `false`<br>仅当 protocol 为 `HTTPS` 时生效                     |
| proxy             | ALIBABACLOUDSTACK_PROXY           | string   | 访问环境时的代理服务器地址                 | 根据环境实际情况确定                                                                         |                                                                      |

---

### 凭证参数

AlibabacloudStack支持STS Token和AK/SK两种认证方式，**一般建议使用STS Token认证**。

- **STS Token**：STS Token为临时凭证，通过角色扮演获取，仅在有效期内有效，可通过以下两种方式获取。
  
  - **STS Token（自动获取）**
    
    当配置`access_key`、`secret_key`和`role_arn`时，由AlibabacloudStack进行角色扮演，并使用扮演产生的STS Token进行认证。
    
  - **STS Token（手动接口获取）**
    
    手动调用API "Sts 2015-04-01 AssumeRole"接口完成角色扮演后获得STS Token。使用获取的STS Token配置AlibabacloudStack进行STS Token认证。
    
- **AK/SK**：AK/SK是永久凭证，一旦泄漏会存在较大的安全风险。
  
  - **账号 AK/SK**
    
    当配置`access_key`和`secret_key`时使用AK/SK认证。

详细配置如下：

#### STS Token（自动获取）

> **注意**：  
> - 通过配置`role_arn`参数开启自动角色扮演，并生成STS Token进行认证。
> - `role_arn`可在统一云管平台>>顶部个人头像 >>*个人信息*>>*查看当前角色策略*，在*选择组织*中切换组织后获取**RAM Role**的值。
> - 账号 AK/SK 可在统一云管平台>>顶部个人头像 >>*个人信息*>>**AccessKey**处获取
> - 若Region下资源集名称不唯一，需使用 `department` 和 `resource_group` 替代 `resource_group_set_name`。

| 参数名                  | 环境变量名                        | 参数类型 | 参数含义                | 备注                                                                 |
|-------------------------|-----------------------------------|----------|-------------------------|----------------------------------------------------------------------|
| access_key              | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | 账号 AK             | **必填**                                                             |
| secret_key              | ALIBABACLOUDSTACK_SECRET_KEY      | string   | 账号 SK             | **必填**                                                             |
| role_arn                | ALIBABACLOUDSTACK_ASSUME_ROLE_ARN  | string  | 待扮演的角色(Ram Role)  | **必填**，示例值：`acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx`         |
| department              | ALIBABACLOUDSTACK_DEPARTMENT      | string   | 凭证登录时的组织ID        | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group          | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | 凭证登录时的资源集ID      | `resource_group_set_name` 不可用或未配置时必填                       |
| resource_group_set_name | ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | 凭证登录时的资源集名称  |  示例值：`ResourceSet(xxxx)`                                        |

#### STS Token（手动接口获取）

> **注意**：  
> - 通过手动调用API "Sts 2015-04-01 AssumeRole"接口并传入`access_key`、`secret_key`和`role_arn`完成角色扮演，获得临时凭证的AK、SK和Token进行配置，从而启用STS Token认证。
> - 若Region下资源集名称不唯一，需使用 `department` 和 `resource_group` 替代 `resource_group_set_name`。

| 参数名                  | 环境变量名                        | 参数类型 | 参数含义                | 备注                                                                 |
|-------------------------|-----------------------------------|----------|-------------------------|----------------------------------------------------------------------|
| access_key              | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | 临时凭证 AK             | **必填**                                                             |
| secret_key              | ALIBABACLOUDSTACK_SECRET_KEY      | string   | 临时凭证 SK             | **必填**                                                             |
| security_token          | ALIBABACLOUDSTACK_SECURITY_TOKEN  | string   | 临时凭证 Token          | **必填**                                                             |
| department              | ALIBABACLOUDSTACK_DEPARTMENT      | string   | 凭证登录时的组织ID        | `resource_group_set_name`不可用(冲突)或未配置时必填                     |
| resource_group          | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | 凭证登录时的资源集ID      | `resource_group_set_name`不可用(冲突)或未配置时必填                      |
| resource_group_set_name | ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | 凭证登录时的资源集名称  |  示例值：`ResourceSet(xxxx)`                                            |

#### 账号 AK/SK

> **注意**：  
> - 账号 AK/SK 可在统一云管平台>>顶部个人头像>>*个人信息*>>**AccessKey**处获取
> - 若Region下资源集名称不唯一，需使用 `department` 和 `resource_group` 替代 `resource_group_set_name`。

| 参数名                  | 环境变量名                        | 参数类型 | 参数含义                | 备注                                                                 |
|-------------------------|-----------------------------------|----------|-------------------------|----------------------------------------------------------------------|
| access_key              | ALIBABACLOUDSTACK_ACCESS_KEY      | string   | 账号 AK                 | **必填**                                                             |
| secret_key              | ALIBABACLOUDSTACK_SECRET_KEY      | string   | 账号 SK                 | **必填**                                                             |
| department              | ALIBABACLOUDSTACK_DEPARTMENT      | string   | 凭证登录时的组织ID        | `resource_group_set_name` 不可用(冲突)或未配置时必填                     |
| resource_group          | ALIBABACLOUDSTACK_RESOURCE_GROUP  | string   | 凭证登录时的资源集ID      | `resource_group_set_name` 不可用(冲突)或未配置时必填                      |
| resource_group_set_name | ALIBABACLOUDSTACK_RESOURCE_GROUP_SET | string | 凭证登录时的资源集名称  |   示例值：`ResourceSet(xxxx)`                                            |
