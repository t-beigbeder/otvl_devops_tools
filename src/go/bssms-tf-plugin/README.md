# Opentofu provider for bootstrap secrets management service

## Requirements

- Opentofu
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

TODO: complete the user documentation.

Installable resource: has to be provisioned before OpenStack instances,
will generate a dedicated key-pair for each instance,
with a secret key to be provided to the instance user_data.

Secrets resource: has to be provisioned after OpenStack instances,
configured with instance specific secrets as key-value pairs,
will activate a `bssms` _provisioner_ that will transfer public-key encrypted secrets
to each instance when it becomes ready.

## Developing the Provider

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.
