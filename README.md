# Terraform Provider ONTAP

WARNING : This is just a proof of concept, if you ended up here looking for a full Terraform provider for ONTAP, roll up your sleeves and start coding !

In the current state, it is just capable of CRUD qtrees and provides a SVM Data Source

## Build

Run the following command to build the provider

```shell
go build -o terraform-provider-ontap
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```