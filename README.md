# Terraform Provider For AlibabaCloud ApsaraStack

<img src="https://www.datocms-assets.com/2885/1506527326-color.svg" width="400px">

## Documentation Language

---

[简体中文](./README_zh-Hans.md) | English

## Official Sites

---


- [GitHub](https://github.com/aliyun/terraform-provider-alibabacloudstack)
- [Alibaba Cloud Help Center](https://help.aliyun.com/apsara/index.html)
- [Terraform Public Registry](https://registry.terraform.io/providers/aliyun/alibabacloudstack)

## Environment Configuration

---

### Install Dependencies

Download TF Core from [Terraform](https://www.terraform.io/downloads.html) or [OpenTofu](https://opentofu.org/docs/intro/install/) official websites, then install or extract it locally. It is recommended to add the installation path to the system `PATH` variable.

> **Note**: Terraform requires version 0.13.x or higher. OpenTofu has no minimum version requirement.

### Install Terraform Provider

> Select the appropriate AlibabacloudStack version based on your proprietary cloud version.

> If the proprietary cloud version is 3.16.2, the AlibabacloudStack version should be 3.16.x. Use `< 3.18.0` to declare the version number when obtaining the latest version.

> If the proprietary cloud version is 3.18.x (including 3.18.0, 3.18.1, 3.18.2, 3.18.6), the AlibabacloudStack version should be 3.18.x. Use `>= 3.18.0` to declare the version number when obtaining the latest version.

| AlibabaCloud ApsaraStack Version | AlibabacloudStack Version |
| ---  | ---  |
| v3.16.2 | < 3.18.0 |
| v3.18.x | >= 3.18.0 |

**Option 1: Automatic Installation**

> **Note**: Automatic installation requires your execution environment to have access to GitHub. This method requires no additional configuration, and TF Core will automatically install the Provider later.

**Option 2: Mirror Site Installation**

> **Description**: Mirror site installation effectively resolves installation failures caused by network isolation or instability.

1. Create a `.terraformrc` or `terraform.rc` configuration file. The file location depends on the host operating system:
   
   > **Note**:
   > 
   > - On Windows: The file must be named `terraform.rc` and placed in the `%APPDATA%` directory of the relevant user. Use `$env:APPDATA` in PowerShell to locate this directory.
   > 
   > - On other systems: The file must be named `.terraformrc` and placed directly in the user's home directory.
   > 
   > - Alternatively, use the `TF_CLI_CONFIG_FILE` environment variable to specify the Terraform CLI configuration file location. Any such file should follow the naming pattern `*.tfrc`.

2. Configure Mirror Site Information

> **Note**: The following example uses [Alibaba Cloud Open Source Mirror Site](https://developer.aliyun.com/mirror/terraform)

```hcl
provider_installation {
  network_mirror {
    url = "https://mirrors.aliyun.com/terraform/"
    // Restrict only AlibabacloudStack downloads from mirror
    include = [
      "registry.terraform.io/aliyun/alibabacloudstack",
      "registry.terraform.io/hashicorp/alibabacloudstack",
    ]
  }
  direct {
    // Other providers maintain original download paths
    exclude = [
      "registry.terraform.io/aliyun/alibabacloudstack",
      "registry.terraform.io/hashicorp/alibabacloudstack",
    ]
  }
}
```

**Option 3: Manual Installation**

Download the appropriate version of AlibabacloudStack from [GitHub](https://github.com/aliyun/terraform-provider-alibabacloudstack/releases/), create the directory structure according to the specified format under the selected installation path, and extract the Provider.

> **Warning**: Incorrect directory structure will cause Provider loading failure.

Standard Provider directory format:

```
XX(Plugin root path, e.g., ./terraform.d/providers/)
└── <hostname>(Use registry.terraform.io for Terraform, registry.opentofu.org for Opentofu)
    └── <Namespace>(e.g., hashicorp or aliyun)
        └── <Provider Name>(alibabacloudstack)
            └── <Provider Version>(e.g., 3.16.2)
                └── <OS Arch>(e.g., windows_amd64)
                    └── <Provider Bin>(terraform-provider-alibabacloudstack)
```

> **Note**: 
>
> - For the `namespace` layer, `hashicorp` is recommended. If using `aliyun`, ensure to explicitly declare it in subsequent Provider configurations.
>
> - Common system architectures: `windows_amd64`, `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64`
>
> - The plugin root path can use the official default locations (auto-loaded during execution) or a custom path (requires manual specification during execution).
>
> Terraform scans and loads Providers from the following OS-specific paths:
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
>    *Linux或其他类Unix系统*:
>
>        $HOME/.terraform.d/plugins
>
>        ~/.local/share/terraform/plugins
>
>        /usr/local/share/terraform/plugins
>
>        /usr/share/terraform/plugins


## Quick Start

---

### Initialize Project

1. **Create working directory**  
   Create a new working directory and create a `provider.tf` file with the following content:

```hcl
 terraform {
  required_providers {
    alibabacloudstack = {
      source = "aliyun/alibabacloudstack"
      #version = "< 3.18.0"
    }
  }
}
```

> **Note**: You can specify the version of the Provider according to your requirements. If not declared, the latest version will be used by default.

2. Initialize the directory

```bash
terraform init
```

> **Note**: If you customized the plugin root path in the [Install Terraform Provider](#install-terraform-provider) section, you need to specify the path during initialization:
>	```bash
>	terraform init -plugin-dir=<YOUR PLUGIN ROOT PATH>
>	```

### Configure Cluster Connection Information

**Option 1: Configuration File**

Create a `provider.tf` file in your working directory and configure according to your environment:

> AlibabacloudStack supports both AK/SK authentication and STS authentication. It is recommended to use STS authentication.
> + [STS authentication] When configuring `access_key`, `secret_key`, and `role_arn`, AlibabacloudStack will perform role assumption and use the generated STS Token for authentication;
> + [STS authentication] When configuring `access_key`, `secret_key`, and `security_token`, AlibabacloudStack will use the specified STS Token for authentication;
> + [AK/SK authentication] When configuring `access_key` and `secret_key`, AK/SK authentication will be used;

```hcl

provider "alibabacloudstack" {
  popgw_domain = "xxx.xxx.com" # AlibabaCloud ApsaraStack platform Service Endpoint
  domain = "xxx.xxx.com" # AlibabaCloud ApsaraStack platform ASAPI Service Endpoint  
  suffix
  access_key   = "xxxx" # Account Access Key
  secret_key   = "xxxx" # Account Secret Key
  # role_arn   = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx" # RAM Role to be assumed
  # security_token = "xxxxxxxx"
  region       = "xxxx" # Platform Region information
  proxy        = "HTTP://x.x.x.x:xxx" # Proxy server address
  protocol                = "HTTPS" # Network protocol (`HTTP` or `HTTPS`)
  insecure                = "true" # Skip HTTPS certificate verification
  resource_group_set_name = "ResourceSet(xxxx)" # Resource group set name
}

```

> For detailed parameter descriptions, refer to [AlibabacloudStack Provider Parameters Documentation](website/docs/index.html.markdown)

**Option 2: Environment Variables** 

Configure environment variables in the command execution terminal:

> AlibabacloudStack supports both AK/SK authentication and STS authentication. It is recommended to use STS authentication.
> + [STS authentication] When configuring `ALIBABACLOUDSTACK_ACCESS_KEY`, `ALIBABACLOUDSTACK_SECRET_KEY`, and `ALIBABACLOUDSTACK_ASSUME_ROLE_ARN`, AlibabacloudStack will perform role assumption and use the generated STS Token for authentication;
> + [STS authentication] When configuring `ALIBABACLOUDSTACK_ACCESS_KEY`, `ALIBABACLOUDSTACK_SECRET_KEY`, and `ALIBABACLOUDSTACK_SECURITY_TOKEN`, AlibabacloudStack will use the specified STS Token for authentication;
> + [AK/SK authentication] When configuring `ALIBABACLOUDSTACK_ACCESS_KEY` and `ALIBABACLOUDSTACK_SECRET_KEY`, AK/SK authentication will be used;

+ *Windows PowerShell*

``` powershell
$env:ALIBABACLOUDSTACK_POPGW_DOMAIN = "xxx.xxx.com"
$env:ALIBABACLOUDSTACK_REGION = "xxxx"
$env:ALIBABACLOUDSTACK_RESOURCE_GROUP_SET = "ResourceSet(xxxx)"
$env:ALIBABACLOUDSTACK_PROTOCOL = "HTTPS"
$env:ALIBABACLOUDSTACK_INSECURE = "true"
$env:ALIBABACLOUDSTACK_ACCESS_KEY = "xxxx"
$env:ALIBABACLOUDSTACK_SECRET_KEY = "xxxx"
```

+ *Unix-like Systems*

```bash
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
export ALIBABACLOUDSTACK_REGION="xxxx"
export ALIBABACLOUDSTACK_RESOURCE_GROUP_SET="ResourceSet(xxxx)"
export ALIBABACLOUDSTACK_PROTOCOL="HTTPS"
export ALIBABACLOUDSTACK_INSECURE="true"
export ALIBABACLOUDSTACK_ACCESS_KEY="xxxx"
export ALIBABACLOUDSTACK_SECRET_KEY="xxxx"
```

> For detailed parameter descriptions, refer to [AlibabacloudStack Provider Parameters Documentation](website/docs/index.html.markdown)

### Orchestrate Resources

1. Create a `main.tf` file in your working directory.

> For more examples, please refer to the [official documentation](https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs)

```hcl
resource "alibabacloudstack_vpc_vpc" "default_vpc" {
  name       = "vpc-test"
  cidr_block = "172.16.0.0/12"
}
```

2. Execute Orchestration

```bash
terraform plan  # View execution plan
terraform apply # Execute orchestration tasks
terraform show  # View orchestration results
terraform destroy # Destroy resources
```

## Changelog

---

Please refer to [ReleaseNote](./CHANGELOG.md)

## Development Guide

---

### Set Up Development Environment

+ Install [Golang](https://golang.org/doc/install)  
  > **Note**: Recommended version 1.21 for current project development.

+ Install [dlv](https://github.com/go-delve/delve/tree/master/Documentation/installation)  
  > **Note**: dlv is a Golang debugger (optional)

### Download Source Code and Compile

```bash
cd <YOUR WORKSPACE>
git clone https://github.com/aliyun/terraform-provider-alibabacloudstack.git
cd terraform-provider-alibabacloudstack
git checkout <appropriate TAG and branch>  # e.g. v3.16.16
go mod tidy
go mod vendor
go build
```

### Cross-compilation

> **Note**: When you need to compile AlibabacloudStack for target execution environments, perform cross-compilation after completing local compilation.

```bash
GOOS=<OS> GOARCH=<ARCH> go build
```

  > **Note**: GOOS options has `windows`, `linux`, `darwin`, `freebsd`, `openbsd`, `solaris`
  
  > **Note**：GOARCH options has `amd64`, `'386'`, `arm`, `arm64`
  
### Source Code Testing

1. Configure environment variables by following the *Environment Variables* method described in the [Configure Cluster Connection Information](#configure-cluster-connection-information) section
2. Execute tests:
```bash
   TF_ACC=1 TF_LOG=INFO go test ./alibabacloudstack -v -run=TestAccAlibabacloudStack -timeout=0
```

### Log Tracing

> **Note**: Enable logging to trace API requests initiated by AlibabacloudStack.

```bash
export DEBUG="terraform"
export TF_LOG="TRACE"
terraform apply
```

### Compatibility Statement

:white_check_mark: All features of Terraform Core are supported by the Provider

:warning: Partial features of Terraform Core are **not** supported by the Provider

| Rpc Name  | terraform-v1.0.11  | terraform-v1.1.9  | terraform-v1.2.9  | terraform-v1.3.10  | terraform-v1.4.7  | terraform-v1.5.7  | terraform-v1.6.6  | terraform-v1.7.5  | terraform-v1.8.5  | terraform-v1.9.3  | opentofu-v1.6.3  | opentofu-v1.7.3  | opentofu-v1.8.0 |
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