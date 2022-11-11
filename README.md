# Terraform Provider OpenCGA

Setup the project:

```shell
go mod init terraform-provider-opencga
go mod tidy
make build
```

## Run example

First, build and install the provider.

```shell
make install
```

## Updating things

For terraform documentation run `go generate`.  
For versioning, update `VERSION` in `Makefile`.

## Limitations

The provider and OpenCGA API do not reliably support delete or update operations, so more
work is needed in order to enable terraform to accurately detect changes.

Therefore this provider should only be used to create new resources.

## Connecting to other OpenCGA instances

💀 **Warning - this is untested software so use with caution for now.**

Set the required username and base url for the OpenCGA instance you wish to configure.

```
provider "opencga" {
  username = "user1"
  base_url = "https://opencga.co.uk/opencga/webservices"
}
```

Export the user password in the shell used to run terraform cmds.

```shell
export OPENCGA_PASSWORD=xxxxxx
```
