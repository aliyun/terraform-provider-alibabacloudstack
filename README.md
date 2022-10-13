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
$ mkdir -p $GOPATH/src/github.com/aliyun; cd $GOPATH/src/github.com/aliyun
$ git clone git@github.com:aliyun/terraform-provider-alibabacloudstack.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/aliyun/terraform-provider-alibabacloudstack
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
cd $GOPATH/src/github.com/aliyun/terraform-provider-alibabacloudstack
export ALIBABACLOUDSTACK_ACCESS_KEY=xxx
export ALIBABACLOUDSTACK_SECRET_KEY=xxx
export ALIBABACLOUDSTACK_REGION=xxx
export ALIBABACLOUDSTACK_DOMAIN=xxx
export ALIBABACLOUDSTACK_RESOURCE_GROUP_SET=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./alibabacloudstack -v -run=TestAccAlibabacloudStack -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```


## Refer

AlibabacloudStack Cloud Provider [Official Docs](https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs)
