# 阿里云专有云飞天企业版Terraform Provider

<img src="https://www.datocms-assets.com/2885/1506527326-color.svg" width="400px">

## 文档语言

---

简体中文|[English](./README.md)

## 官方站点

---

- [GitHub](https://github.com/aliyun/terraform-provider-alibabacloudstack)
- [阿里云帮助中心](https://help.aliyun.com/apsara/index.html)
- [Terraform Public Registry](https://registry.terraform.io/providers/aliyun/alibabacloudstack)

## 环境配置

---

### 安装依赖

从[Terraform](https://www.terraform.io/downloads.html)官网或[OpenTofu](https://opentofu.org/docs/intro/install/)官网下载TF Core，并安装或解压到本地。推荐将安装路径配置到系统`PATH`变量中。

> **注意**：Terraform需要版本不低于0.13.x，OpenTofu无最低版本要求。

### 安装Terraform Provider

> 请根据专有云版本选择合适的AlibabacloudStack版本。

> 如果专有云版本为3.16.2, AlibabacloudStack版本为 3.16.x，获取最新版本可用使用`< 3.18.0`声明版本号。

> 如果专有云版本为3.18.x(3.18.0,3.18.1,3.18.2,3.18.6), AlibabacloudStack版本为 3.18.x，获取最新版本可用使用`>= 3.18.0`声明版本号。

| 专有云版本 | AlibabacloudStack版本 |
| ---  | ---  |
| v3.16.2 | < 3.18.0 |
| v3.18.x | >= 3.18.0 |

**方案一 自动安装**

> **注意**：只有你的执行环境可以正常访问Github可以使用自动安装，该方案不需要额外的配置，直接跳转到[快速开始](#快速开始)开始使用，后续TF Core会自动安装Provider。

**方案二 镜像站点安装**

> **说明**：镜像站点安装可以有效解决网络隔离或网络不稳定导致的安装失败的问题。

1. 创建.terraformrc 或terraform.rc配置文件，文件位置取决于主机的操作系统。
   
   > **说明**：
   > 
   > 在 Windows 环境上，文件必须命名为terraform.rc，并放置在相关用户的%APPDATA%目录中。这个目录的物理位置取决于Windows 版本和系统配置；在 PowerShell 中使用 $env:APPDATA 可以找到其在系统上的位置。
   > 
   > 在所有其他系统上，必须将该文件命名为.terraformrc，并直接放在相关用户的主目录中。
   > 
   > 另外，你也可以使用TF_CLI_CONFIG_FILE环境变量指定 Terraform CLI 配置文件的位置，任何此类文件都应遵循命名模式`*.tfrc`。

2. 配置镜像站点信息

   > **说明**：下面以[阿里云开源镜像站](https://developer.aliyun.com/mirror/terraform)为例

``` hcl
provider_installation {
  network_mirror {
    url = "https://mirrors.aliyun.com/terraform/"
    // 限制只有AlibabacloudStack从镜像源下载
    include = [
      "registry.terraform.io/aliyun/alibabacloudstack",
      "registry.terraform.io/hashicorp/alibabacloudstack",
    ]
  }
  direct {
    // 声明除了AlibabacloudStack, 其它Provider保持原有的下载链路
    exclude = [
      "registry.terraform.io/aliyun/alibabacloudstack",
      "registry.terraform.io/hashicorp/alibabacloudstack",
    ]
  }
}
```

**方案三 手动安装**

从[GitHub](https://github.com/aliyun/terraform-provider-alibabacloudstack/releases/)下载合适版本的AlibabacloudStack，在选定安装路径下依照指定的格式创建目录架构，并解压Provider。

> **注意**： 错误的目录结构会导致Provider加载失败

Provider标准目录格式：

```
XX(插件根路径 如：./terraform.d/providers/)
└── <hostname>(Terraform用registry.terraform.io，Opentofu用registry.opentofu.org)
    └── <命名空间>(如： hashicorp或aliyun)
        └── <Provider名称>(alibabacloudstack)
            └── <Provider版本>(如：3.16.2)
                └── <系统架构>(如： windows_amd64)
                    └── <Provider文件>(terraform-provider-alibabacloudstack)
```

> **注意**： 命名空间层，建议使用hashicorp，如果使用aliyun，在后续步骤编写Provider信息时要注意申明aliyun。

> **注意**： 常见的系统架构有：windows_amd64，linux_amd64，linux_arm64，darwin_amd64，darwin_arm64

>
> **注意**： 插件根路径可以选择官方默认路径下，在每次执行时自动加载，也可以安装到自定义路径下，但需要在每次执行时手动指定路径。
> 
> 根据不同的系统，Terraform会依次尝试从多个路径中进行扫描并加载Provider。
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

## 快速开始

---

### 初始化项目

1. 新建工作目录，并创建`provider.tf`文件

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

> **注意**：你可以根据你的需求指定使用的Provider的版本，如果不声明，则默认使用最新版本。

2. 初始化目录

``` bash
terraform init
```

> **注意**：如果在《安装Terraform Provider》章节中，自定义了插件根路径，初始化时需要指定相关路径
>	```bash
>	terraform init -plugin-dir=<YOUR PLUGIN ROOT PATH>
>	```

### 配置集群链接信息

**方案一 配置文件**

在工作目录下创建`provider.tf`文件，并根据环境进行进行配置。

> AlibabacloudStack支持AK/SK验证和STS两种验证方式，建议使用STS验证。
> + [STS验证方式一] 当配置`access_key`、`secret_key`和`role_arn`时，由AlibabacloudStack进行角色扮演，并使用扮演产生的STS Token进行鉴权;
> + [STS验证方式二] 当配置`access_key`、`secret_key`和`security_token`时，由AlibabacloudStack使用指定的STS Token进行鉴权;
> + [AK/SK验证] 当配置`access_key`和`secret_key`时使用AK/SK验证;

```hcl

provider "alibabacloudstack" {
  popgw_domain = "xxx.xxx.com" # 阿里云专有云飞天企业版服务地址标准后缀
  # domain = "xxx.xxx.com"       # 阿里云专有云飞天企业版ASAPI网关服务地址
  access_key   = "xxxx" # 账号 AK
  secret_key   = "xxxx" # 账号 SK
  role_arn   = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx" # 待扮演角色的Ram Role
  # security_token = "xxxxxxxx"
  region       = "xxxx" # 平台 Region 信息
  proxy        = "HTTP://x.x.x.x:xxx" # 平台 Region 信息
  protocol                = "HTTPS" # 网络协议（可选 `HTTP` 或 `HTTPS`）
  insecure                = "true" # 是否跳过 HTTPS 证书校验
  resource_group_set_name = "ResourceSet(xxxx)" # 资源集名称
}

```

> 详细的参数说明可以参考[AlibabacloudStack Provider参数说明](website/docs_zh-Hans/index.html.markdown)

**方案二 环境变量**

在为命令的执行终端配置环境变量。

> AlibabacloudStack支持AK/SK验证和STS两种验证方式，建议使用STS验证。
> + [STS验证方式一] 当配置`ALIBABACLOUDSTACK_ACCESS_KEY`、`ALIBABACLOUDSTACK_SECRET_KEY`和`ALIBABACLOUDSTACK_ASSUME_ROLE_ARN`时，由AlibabacloudStack进行角色扮演，并使用扮演产生的STS Token进行鉴权;
> + [STS验证方式二] 当配置`ALIBABACLOUDSTACK_ACCESS_KEY`、`ALIBABACLOUDSTACK_SECRET_KEY`和`ALIBABACLOUDSTACK_SECURITY_TOKEN`时，由AlibabacloudStack使用指定的STS Token进行鉴权;
> + [AK/SK验证] 当配置`ALIBABACLOUDSTACK_ACCESS_KEY`和`ALIBABACLOUDSTACK_SECRET_KEY`时使用AK/SK验证;

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

+ *类Unix系统*

```bash
export ALIBABACLOUDSTACK_POPGW_DOMAIN="xxx.xxx.com"
export ALIBABACLOUDSTACK_REGION="xxxx"
export ALIBABACLOUDSTACK_RESOURCE_GROUP_SET="ResourceSet(xxxx)"
export ALIBABACLOUDSTACK_PROTOCOL="HTTPS"
export ALIBABACLOUDSTACK_INSECURE="true"
export ALIBABACLOUDSTACK_ACCESS_KEY="xxxx"
export ALIBABACLOUDSTACK_SECRET_KEY="xxxx"
export ALIBABACLOUDSTACK_ASSUME_ROLE_ARN = "acs:ram::xxxxxxx:role/ascm-role-x-x-xxxx"
```

> 详细的参数说明可以参考[AlibabacloudStack Provider参数说明](website/docs_zh-Hans/index.html.markdown)

### 编排资源

1. 在工作目录下创建`main.tf`文件。

> 更多示例请查看[官方手册](https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs)

```hcl
resource "alibabacloudstack_vpc_vpc" "default_vpc" {
  name       = "vpc-test"
  cidr_block = "172.16.0.0/12"
}
```

2. 执行编排

``` bash
terraform plan # 查看资源计划

terraform apply # 执行编排任务

terraform show # 查看编排结果

terraform destroy # 销毁资源
```

## 变更日志

---

请查看[ReleaseNode](./CHANGELOG_zh_Hans.md)

## 开发教程

---

### 安装编译环境

+ 安装[Golang](https://golang.org/doc/install) 

> **说明**：当前项目开发过程中建议使用1.21版本。

+ 安装[dlv](https://github.com/go-delve/delve/tree/master/Documentation/installation)

> **说明**：dlv为Golang调试工具，非必须

### 下载源码并编译

```bash
cd <YOUR WORKSPACE>
git clone https://github.com/aliyun/terraform-provider-alibabacloudstack.git
cd terraform-provider-alibabacloudstack
git checkout <合适的TAG和分支>
go mod tidy
go mod vendor
go build
```

### 交叉编译

> **说明**： 当你需要执行环境编译AlibabacloudStack时，可以在完成本地编译的基础上进行交叉编译。

```bash
GOOS=<OS> GOARCH=<ARCH> go build
```

> **说明**：GOOS可以选择 `windows`, `linux`, `darwin`, `freebsd`, `openbsd`, `solaris`

> **说明**：GOARCH可以选择 `amd64`, `'386'`, `arm`, `arm64`

### 源码测试

1. 参考《配置集群链接信息》章节中的*环境变量*配置方案，配置环境变量
2. 执行测试
   ```bash
   TF_ACC=1 TF_LOG=INFO go test ./alibabacloudstack -v -run=TestAccAlibabacloudStack -timeout=0
   ```

### 日志追踪

> **说明**： 当你对AlibabacloudStack所发起的API进行追逐时可以打开日志进行日志的追踪。

```bash
export DEBUG="terraform"
export TF_LOG="TRACE"
terraform apply
```

### 兼容性申明

:white_check_mark:：Terraform Core的所有能力被Provider支持

:warning:：Terraform Core的部分能力不被Provider支持

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
