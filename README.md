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

The OpenCGA configuration parameters are in `example/full_setup/main.tf` in the `provider` block.
By default this is set to configure the local OpenCGA service running.

Run the following command to initialize the workspace and apply the sample configuration.

```shell
make run_example
```

This will create a project, two studies attached to the project and a variable set for
the second study. The variable set definitions can be found in `sample.json`.

## Updating Terraform docs

Run `go generate`. 

## Limitations

Currently the provider cannot accurately compare existing state with the configuration of 
resources in the terraform `main.tf`. This causes terraform to try and replace the resources.

The provider and OpenCGA API do not reliably support delete or update operations, so more
work is needed in order to enable terraform to accurately detect changes.

Therefore this provider should only be used to create new resources.

## Connecting to other OpenCGA instances

💀 **Warning - this is untested software so use with caution for now.**

Set the required username and base url for the OpenCGA instance you wish to configure.

```
provider "opencga" {
  username = "bertha"
  base_url = "https://opencgainternal.test.aws.gel.ac/opencga/webservices"
}
```

Export the user password in the shell used to run terraform cmds.

```shell
export OPENCGA_PASSWORD=xxxxxx
```
