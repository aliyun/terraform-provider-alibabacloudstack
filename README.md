# Alibaba Cloud ApsaraStack Terraform Provider

<img src="https://www.datocms-assets.com/2885/1506527326-color.svg" width="400px">

## Document languages

---

[简体中文](./README_zh-Hans.md)|English

## Websites

---

- [Alibaba Cloud Help Center](https://help.aliyun.com/apsara/index.html)
- [Terraform Public Registry](https://registry.terraform.io/providers/aliyun/alibabacloudstack)
- [GitHub](https://github.com/aliyun/terraform-provider-alibabacloudstack)

## Environment configuration

---

### Install dependencies

Download TF Core from the [Terraform](https://www.terraform.io/downloads.html) or [OpenTofu](https://opentofu.org/docs/intro/install/) website and install TF Core on your on-premises computer. We recommend that you add the installation path to the PATH system variable. 

> **Note** : The version of Terraform must be 0.13.x or later. OpenTofu does not require a minimum version. 

### Install the Terraform provider

> Select an appropriate AlibabacloudStack version based on your Apsara Stack version.
>
> + If you use Alibaba Cloud ApsaraStack V3.16.2, the compatible version of AlibabaCloudStack should correspondingly be V3.16.x. You can use `< 3.18.0` to declare the version number if you want to obtain the latest version. **Note**: Multiple Alibaba Cloud ApsaraStack versions are released. Refer to the documentation version corresponding to the Alibaba Cloud ApsaraStack version. 
>
> + If you use Alibaba Cloud ApsaraStack V3.18.x (such as V3.18.0, V3.18.1, V3.18.2, and V3.18.6), the compatible version of AlibabaCloudStack should correspondingly be 3.18.x. You can use `< 3.19.0` to declare the version number if you want to obtain the latest version. 
>
> + To obtain the Terraform provider versions for other Alibaba Cloud ApsaraStack versions, contact Apsara Stack customer service.

> + Pre-compiled providers are not provided for non-LTS versions. Please compile them yourself if needed.

| Alibaba Cloud ApsaraStack version | AlibabacloudStack version |
| ---  | ---  |
| v3.16.2 | LTS | < 3.18.0 | Bug Support |
| v3.18.x | LTS | < 3.19.0 | Standard Maintenance |
| v3.19.0 | non-LTS | < 3.20.0 | Bug Support |

**Solution 1: Automatic installation**

> **Note**: Automatic installation is supported only when your environment can access GitHub. This solution does not require additional configurations. You can skip the following content of this section and proceed with the [Getting Started](#getting-started) section to get started with the Terraform provider. TF Core will automatically install AlibabacloudStack. 

**Solution 2: Installation from the mirror site**

> **Note**: The mirror site provides local download sources to prevent installation errors caused by network isolation or instability.

1. Create a Terraform configuration file (The file location varies based on the operating system on your host)
   
   > **Note**:
   > On a Windows operating system, the file must be named `terraform.rc` and stored in the `%APPDATA%` directory of the user. This directory varies based on the Windows version and system configurations. You can run the `$env:APPDATA` command in PowerShell to find the file directory on your operating system. 
   >
   > On all other operating systems, the file must be named `.terraformrc` and stored in the home directory of the user. 
   >
   > In addition, you can specify the path of the Terraform configuration file in the `TF_CLI_CONFIG_FILE` environment variable. The configuration file specified in this way must be named with the `.tfrc` suffix. 

2. Configure the information about the image site

   > **Note**: In the following example, the [Alibaba Cloud open source image site](https://developer.aliyun.com/mirror/terraform) is used.

``` hcl
provider_installation {
  network_mirror {
    url = "https://mirrors.aliyun.com/terraform/"
    // Allows AlibabacloudStack to download only from the specified image source
    include = [
      "registry.terraform.io/aliyun/alibabacloudstack",
      "registry.terraform.io/hashicorp/alibabacloudstack",
    ]
  }
  direct {
    // All providers except AlibabacloudStack must maintain the original download URLs
    exclude = [
      "registry.terraform.io/aliyun/alibabacloudstack",
      "registry.terraform.io/hashicorp/alibabacloudstack",
    ]
  }
}
```

**Solution 3: Manual installation**

Download the required AlibabacloudStack version from [GitHub](https://github.com/aliyun/terraform-provider-alibabacloudstack/releases/). Create a directory structure in the desired installation path based on the specified format, and decompress the provider.

> **Note**: The download of the provider fails if you use an invalid directory structure.

Standard directory structure of the provider:

```
XX(The provider repository path, such as./terraform.d/providers/)
└── <hostname>(Terraform uses registry.terraform.io and OpenTofu uses registry.opentofu.org.)
    └── <namespace>(aliyun)
        └── <provider name>(alibabacloudstack)
            └── <provider version>(such as  3.18.0)
                └── <system architecture>(such as windows_amd64)
                    └── <provider file>(terraform-provider-alibabacloudstack)
```
> **Note**: We recommend that you use the aliyun namespace. In the [Initialize a project](#initialize-a-project) step, we recommend that you declare the `source` of AlibabacloudStack as aliyun/alibabacloudstack when you write the `provider.tf` file. 

> **Note**: Common system architectures include `windows_amd64`, `linux_amd64`, `linux_arm64`, `darwin_amd64`, and `darwin_arm64`.


> **Note**: You can specify the default path as the provider repository so that the provider can be automatically loaded each time you run Terraform. You can also install the provider repository to a custom path. In this case, you need to manually specify the path each time you run Terraform. 
> 
> Terraform scans and loads the AlibabacloudStack provider from multiple directories in sequence based on the operating system that you use.
> 
>    *Windows*:
>
>        %APPDATA%/terraform.d/plugins
>
>        %APPDATA%/HashiCorp/Terraform/plugins
> 
>    *Mac OS X*:
>
>        $HOME/.terraform.d/plugins
>
>        ~/Library/Application Support/io.terraform/plugins
>
>        /Library/Application Support/io.terraform/plugins
>    
>    *Linux or other Unix-like operating systems*:
>
>        $HOME/.terraform.d/plugins
>
>        ~/.local/share/terraform/plugins
>
>        /usr/local/share/terraform/plugins
>
>        /usr/share/terraform/plugins

## Getting Started

---

### Initialize a project

1. Create a working directory and a `terraform.tf` file

```hcl
terraform {
  required_providers {
    alibabacloudstack = {
      source = "aliyun/alibabacloudstack"
      version = "< 3.19.0"
    }
  }
}
```

> **Note**: You can specify a provider version based on your business requirements. If you do not specify a provider version, the latest version is used by default. 

2. Initialize the directory

``` bash
terraform init
```

> **Note**: If you specify a third-party provider repository path in the [Install the Terraform provider](#install-the-terraform-provider) section, you must specify the path during directory initialization.
>	```bash
>	terraform init -plugin-dir=<provider repository path>
>	```

### Configure the Alibaba Cloud ApsaraStack Cluster Info

**Solution 1: Use a configuration file**

Create a `provider.tf` file in the working directory and configure the parameters based on the environment that you use. 

> AlibabacloudStack supports authentication based on AccessKey pairs and Security Token Service (STS) tokens. We recommend that you use STS for authentication.
> + [STS authentication method 1]: When `access_key`, `secret_key`, and `role_arn` are configured, AlibabacloudStack assumes the role and uses the STS token generated by the role for authentication. 
> + [STS authentication method 2]: When `access_key`, `secret_key`, and `security_token` are configured, AlibabacloudStack uses the specified STS token for authentication. 
> + [Authentication based on AccessKey pairs]: When `access_key` and `secret_key` are configured, the AccessKey pair is used for authentication. 

```hcl

provider "alibabacloudstack" {
  popgw_domain = "xxx.xxx.com" # The standard suffix in the service address of Alibaba Cloud ApsaraStack.
  access_key   = "xxxx" # The AccessKey ID of the account.
  secret_key   = "xxxx" # The AccessKey secret of the account.
  role_arn   = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx" # The Resource Access Management (RAM) role to be assumed.
  # security_token = "xxxxxxxx"
  region       = "xxxx" # The region of Apsara Stack.
  proxy        = "HTTP://x.x.x.x:xxx" # The address of the proxy server. 
  protocol                = "HTTPS" # The network protocol. Valid values: HTTP and HTTPS.
  insecure                = "true" # Specify whether to skip HTTPS certificate verification. Valid values: true and false.
  resource_group_set_name = "ResourceSet(xxxx)" # The name of the resource set.
}

```

> For more information about the parameters, see [AlibabacloudStack provider parameters](website/docs/index.html.markdown)

**Solution 2: Use environment variables**

Configure environment variables for the terminal on which you run the commands. 


> AlibabacloudStack supports authentication based on AccessKey pairs and STS tokens. We recommend that you use STS for authentication.
> + [STS authentication method 1]: When `ALIBABACLOUDSTACK_ACCESS_KEY`, `ALIBABACLOUDSTACK_SECRET_KEY`, and `ALIBABACLOUDSTACK_ASSUME_ROLE_ARN` are configured, AlibabacloudStack assumes the role and uses the STS token generated by the role for authentication. 
> + [STS authentication method 2]: When `ALIBABACLOUDSTACK_ACCESS_KEY`, `ALIBABACLOUDSTACK_SECRET_KEY`, and `ALIBABACLOUDSTACK_SECURITY_TOKEN` are configured, AlibabacloudStack uses the specified STS token for authentication. 
> + [Authentication based on AccessKey pairs]: When `ALIBABACLOUDSTACK_ACCESS_KEY` and `ALIBABACLOUDSTACK_SECRET_KEY` are configured, the AccessKey pair is used for authentication. 

+ *Windows PowerShell*

``` powershell
$env:ALIBABACLOUDSTACK_POPGW_DOMAIN = "xxx.xxx.com"
$env:ALIBABACLOUDSTACK_REGION = "xxxx"
$env:ALIBABACLOUDSTACK_RESOURCE_GROUP_SET = "ResourceSet(xxxx)"
$env:ALIBABACLOUDSTACK_PROTOCOL = "HTTPS"
$env:ALIBABACLOUDSTACK_INSECURE = "true"
$env:ALIBABACLOUDSTACK_ACCESS_KEY = "xxxx"
$env:ALIBABACLOUDSTACK_SECRET_KEY = "xxxx"
$env:ALIBABACLOUDSTACK_ASSUME_ROLE_ARN = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
```

+ *Unix-like operating systems*

```bash
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
export ALIBABACLOUDSTACK_REGION="xxxx"
export ALIBABACLOUDSTACK_RESOURCE_GROUP_SET="ResourceSet(xxxx)"
export ALIBABACLOUDSTACK_PROTOCOL="HTTPS"
export ALIBABACLOUDSTACK_INSECURE="true"
export ALIBABACLOUDSTACK_ACCESS_KEY="xxxx"
export ALIBABACLOUDSTACK_SECRET_KEY="xxxx"
export ALIBABACLOUDSTACK_ASSUME_ROLE_ARN="acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
```

> For more information about the parameters, see [AlibabacloudStack provider parameters](website/docs/index.html.markdown)

### Orchestrate resources

1. Create a `main.tf` file in the working directory.


```hcl
resource "alibabacloudstack_vpc_vpc" "default_vpc" {
  name       = "vpc-test"
  cidr_block = "172.16.0.0/12"
}
```
> This example demonstrates the basic configuration and operational workflows of Terraform with a simple VPC resource.
> For more examples, see alibabacloudstack on the [Terraform website](https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs)

2. Perform resource orchestration

``` bash
terraform plan # View the resource plan.

terraform apply # Execute the orchestration task.

terraform show # View the orchestration result.

terraform destroy # Destroy resources.
```

## Change history

---

For more information, see [Change history.](./CHANGELOG.md)

## Developer guide

---

### Install a compilation environment

+ Install [Golang](https://golang.org/doc/install) 

> **Note**: We recommend that you use Golang 1.21 during project development. 

+ Install [dlv](https://github.com/go-delve/delve/tree/master/Documentation/installation)

> **Note**: Delve is a debugger for Golang. Delve is an optional tool.

### Download and compile the source code

```bash
cd <YOUR WORKSPACE>
git clone https://github.com/aliyun/terraform-provider-alibabacloudstack.git
cd terraform-provider-alibabacloudstack
git checkout <An appropriate tag and branch>
go mod tidy
go mod vendor
go build
```

### Perform cross-compilation

> **Note**: If you need to compile AlibabacloudStack in your environment, you can perform cross-compile after on-premises compilation is completed. 

```bash
GOOS=<OS> GOARCH=<ARCH> go build
```

> **Note**: You can select  `windows`, `linux`, `darwin`, `freebsd`, `openbsd`, `solaris` for GOOS.

> **Note**: You can select `amd64`, `'386'`, `arm`, `arm64` for GOARCH.

### Test the source code

1. For more information, see the environment variable configuration methods in the [Configure the Alibaba Cloud ApsaraStack Cluster Info](#configure-the-alibaba-cloud-apsarastack-cluster-info) section.
2. Run the testcase.
   ```bash
   TF_ACC=1 TF_LOG=INFO go test ./alibabacloudstack -v -run="TestAccAlibabacloudStackxxxxx" -timeout=0
   ```

### Debug the source code

1. For more information, see the environment variable configuration methods in the [Configure the Alibaba Cloud ApsaraStack Cluster Info](#configure-the-alibaba-cloud-apsarastack-cluster-info) section.
2. Debug the testcase.
   ```bash
   TF_ACC=1 dlv test ./alibabacloudstack -- -test.v -test.run="TestAccAlibabacloudStackxxxxx"
   ```

### Enable logging

> **Note**: If you need to trace API requests initiated by AlibabacloudStack, you can enable the logging feature. 

```bash
export DEBUG="terraform"
export TF_LOG="TRACE"
terraform apply
```

### Compatibility

:white_check_mark::All capabilities of Terraform Core are supported by the provider.

:warning::Some capabilities of Terraform Core are not supported by the provider.

| RPC Name  | terraform-v1.0.11  | terraform-v1.1.9  | terraform-v1.2.9  | terraform-v1.3.10  | terraform-v1.4.7  | terraform-v1.5.7  | terraform-v1.6.6  | terraform-v1.7.5  | terraform-v1.8.5  | terraform-v1.9.3  | opentofu-v1.6.3  | opentofu-v1.7.3  | opentofu-v1.8.0 |
| ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | --- |
| GetSchema  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :warning:  | :warning:  | :white_check_mark:  | :warning:  | :warning: |
| PrepareProviderConfig  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ValidateResourceTypeConfig  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ValidateDataSourceConfig  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| UpgradeResourceState  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| Configure  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :warning:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ReadResource  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :warning:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| PlanResourceChange  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :warning:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ApplyResourceChange  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ImportResourceState  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :warning:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ReadDataSource  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :warning:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| Stop  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
