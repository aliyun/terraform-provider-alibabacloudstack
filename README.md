Terraform Provider For ApsaraStack Cloud
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

Clone repository to: `$GOPATH/src/github.com/aliyun/terraform-provider-apsarastack`

```sh
$ mkdir -p $GOPATH/src/github.com/apsara-stack; cd $GOPATH/src/github.com/apsara-stack
$ git clone git@github.com:apsara-stack/terraform-provider-apsarastack
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/apsara-stack/terraform-provider-apsarastack
$ make build
```

Using the provider
----------------------
### Create the main.tf on working directory & add following portion to configure provider

````
 terraform {
  required_providers {
    apsarastack = {
      source = "apsara-stack/apsarastack"
      version = "1.0.1"
    }
  }
}

# Configure the ApsaraStack Provider
 provider "apsarastack" {
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
resource "apsarastack_vpc" "default_vpc" {
  name       = "vpc-test"
  cidr_block = "172.16.0.0/12"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-apsarastack
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
Before making a release, the resources and data sources are tested automatically with acceptance tests (the tests are located in the apsarastack/*_test.go files).
You can run them by entering the following instructions in a terminal:
```
cd $GOPATH/src/github.com/apsara-stack/terraform-provider-apsarastack
export APSARASTACK_ACCESS_KEY=xxx
export APSARASTACK_SECRET_KEY=xxx
export APSARASTACK_REGION=xxx
export APSARASTACK_DOMAIN=xxx
export APSARASTACK_RESOURCE_GROUP_SET=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./apsarastack -v -run=TestAccApsaraStack -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```


## Refer

ApsaraStack Cloud Provider [Official Docs](https://registry.terraform.io/providers/apsara-stack/apsarastack/latest/docs)
