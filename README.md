Terraform Provider For AlibabacloudStack Cloud
==================



- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px"> 


<img src="https://www.datocms-assets.com/2885/1506527326-color.svg" width="400px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)
-   [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports):
    ```
    go get golang.org/x/tools/cmd/goimports
    ```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/aliyun/terraform-provider-alibabacloudstack`

```sh
$ mkdir -p $GOPATH/src/github.com/apsara-stack; cd $GOPATH/src/github.com/apsara-stack
$ git clone git@github.com:aliyun/terraform-provider-alibabacloudstack.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/apsara-stack/terraform-provider-alibabacloudstack
$ go build -o terraform-provider-alibabacloudstack
```

Using the provider
----------------------
### Create the main.tf on working directory & add following portion to configure provider

````
 terraform {
  required_providers {
    alibabacloudstack = {
      source = "aliyun/alibabacloudstack"
      version = "1.0.1"
    }
  }
}

# Configure the AlibabacloudStack Provider
 provider "alibabacloudstack" {
  access_key = "ckhCs1K*********"
  secret_key = "2lY9uNh***********************"
  region =  "cn-xxxxxx-env00-d01"
  proxy = "http://100.1.1.1:5001"
  insecure = true
  resource_group_set_name= "ResourceSet(wzw)"
  domain = "server.asapi.cn-xxxxx-envXX-d01.intra.envXX.shuguang.com/asapi/v3"
  protocol = "HTTP"
}
````                                               
- Add following data in main.tf to create the resource vpc from terraform
```
resource "alibabacloudstack_vpc" "default_vpc" {
  name       = "vpc-test"
  cidr_block = "172.16.0.0/12"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ go build -o terraform-provider-alibabacloudstack
...
$ $GOPATH/bin/terraform-provider-alibabacloudstack
...
```

Running `make dev` or `make devlinux` or `devwin` will only build the specified developing provider which matches the local system.
And then, it will unarchive the provider binary and then replace the local provider plugin.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Acceptance Testing
Before making a release, the resources and data sources are tested automatically with acceptance tests (the tests are located in the alibabacloudstack/*_test.go files).
You can run them by entering the following instructions in a terminal:
```
cd $GOPATH/src/github.com/apsara-stack/terraform-provider-alibabacloudstack
export ALIBABACLOUDSTACK_ACCESS_KEY=xxx
export ALIBABACLOUDSTACK_SECRET_KEY=xxx
export ALIBABACLOUDSTACK_REGION=xxx
export ALIBABACLOUDSTACK_POPGW_DOMAIN=xxx
export ALIBABACLOUDSTACK_RESOURCE_GROUP_SET=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./alibabacloudstack -v -run=TestAccAlibabacloudStack -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```


## Refer

AlibabacloudStack Cloud Provider [Official Docs](https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs)


## 当前Provider兼容性
<!-- INSERT TABLE HERE -->


:white_check_mark::当前功能被Provider支持
:x::当前功能在该Provider存在风险
:no_entry_sign::当前功能在该Provider下不可用
| Rpc Name  | terraform-v1.0.11  | terraform-v1.1.9  | terraform-v1.2.9  | terraform-v1.3.10  | terraform-v1.4.7  | terraform-v1.5.7  | terraform-v1.6.6  | terraform-v1.7.5  | terraform-v1.8.5  | terraform-v1.9.3  | opentofu-v1.6.3  | opentofu-v1.7.3  | opentofu-v1.8.0 |
| ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | ---  | --- |
| GetSchema  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :x:  | :white_check_mark:  | :x:  | :x: |
| PrepareProviderConfig  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ValidateResourceTypeConfig  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ValidateDataSourceConfig  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| UpgradeResourceState  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| Configure  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ReadResource  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| PlanResourceChange  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ApplyResourceChange  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ImportResourceState  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ReadDataSource  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| Stop  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| GetMetadata  |    |    |    |    |    |    | :white_check_mark:  | :white_check_mark:  | :x:  | :x:  |    | :x:  | :x: |
| MoveResourceState  |    |    |    |    |    |    |    |    | :no_entry_sign:  | :no_entry_sign:  |    | :no_entry_sign:  | :no_entry_sign: |
| GetFunctions  |    |    |    |    |    |    |    |    | :no_entry_sign:  | :no_entry_sign:  |    | :no_entry_sign:  | :no_entry_sign: |
| CallFunction  |    |    |    |    |    |    |    |    | :no_entry_sign:  | :no_entry_sign:  |    | :no_entry_sign:  | :no_entry_sign: |
| GetMetadata  |    |    |    |    |    |    | :white_check_mark:  | :white_check_mark:  | :x:  | :x:  |    | :x:  | :x: |
| GetSchema  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :x:  | :white_check_mark:  | :x:  | :x: |
| Configure  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ReadResource  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| PlanResourceChange  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ImportResourceState  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
| ReadDataSource  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :x:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark: |
