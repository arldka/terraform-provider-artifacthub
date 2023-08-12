[![Go Report Card](https://goreportcard.com/badge/github.com/arldka/terraform-provider-artifacthub)](https://goreportcard.com/report/github.com/arldka/terraform-provider-artifacthub)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.[label](https://app.fossa.com/projects/git%2Bgithub.com/arldka/terraform-provider-artifacthub?ref%3Dbadge_large)com%2Farldka%2Fterraform-provider-artifacthub.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Farldka%2Fterraform-provider-artifacthub?ref=badge_shield)
[![codecov](https://codecov.io/gh/arldka/terraform-provider-artifacthub/branch/main/graph/badge.svg?token=DYZPVO90SH)](https://codecov.io/gh/arldka/terraform-provider-artifacthub)
![Visitors](https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2Farldka%2Fterraform-provider-artifacthub&label=Visitors&countColor=%23d9e3f0&style=flat)
[![Latest release](https://img.shields.io/github/v/release/arldka/terraform-provider-artifacthub)](https://github.com/arldka/terraform-provider-artifacthub/releases/latest)



# Artifacthub Terraform Provider

Terraform provider for [Artifact Hub](https://artifacthub.io).

## License

This repository is open source, please refer to the [License](https://github.com/arldka/terraform-provider-artifacthub/blob/main/LICENSE) for more information.

## Getting Started & Documentation

If you're new to Terraform and want to get started creating infrastructure, please check out the Terraform official [Getting Started guides](https://learn.hashicorp.com/terraform#getting-started) on HashiCorp's learning platform. There are also [additional guides](https://learn.hashicorp.com/terraform#operations-and-development) to continue your learning.

### Use the provider

```
terraform {
  required_providers {
    artifacthub = {
      version = "0.1.0"
      source  = "arldka/artifacthub"
    }
  }
}

provider "artifacthub" {
  # Configuration options
  api_key = "xxxx"
  api_key_secret = "xxxx"
  # OR USE THE ARTIFACTHUB_API_KEY and ARTIFACTHUB_API_KEY_SECRET ENVIRONMENT VARIABLES
}
```

## Developing the provider

To learn more about how to contribute to the development of this provider please refer to the [community guidelines](https://github.com/arldka/terraform-provider-arldka/blob/main/CONTRIBUTING.md).

**Did you find a vulnerability?** Please refer to the [Security Policy](https://github.com/arldka/terraform-provider-artifacthub/security/policy) for more information.

### Makefile

The make command provides an easy way to access commands for local development of the provider.

- To build the provider locally, you can execute the make command at the root of the project.

  ```
  make build
  ```

  *Alternatively: 'go build -o terraform-provider-artifacthub'*
- To format all the go files in the project.

  ```
  make format
  ```

  *Alternatively: 'gofmt -l -s -w .'*
- To generate documentation using the tfplugindocs plugin.

  ```
  make docs
  ```

  *Alternatively: 'go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs'*


### Use locally built provider

1. To build and install the provider locally, you can execute the following command at the root of the project.

   ```
   make install
   ```
2. You can use the provider now by setting the source to "hashicorp.local/devoteamgcloud/looker"

  ```
  terraform {
    required_providers {
      artifacthub = {
        version = "0.1.0"
        source  = "arldka.cloud/dev/artifacthub"
     }
   } 
  }
  ```

# Provider Resource Coverage

## Data Sources

### Packages

- [ ] container
- [ ] coredns
- [X] helm
- [ ] helm-plugin
- [ ] falco
- [ ] gatekeeper
- [ ] keda-scaler
- [ ] keptn
- [ ] krew
- [ ] kubewarden
- [ ] kyverno
- [ ] opa
- [ ] olm
- [ ] tbaction
- [ ] tekton-pipeline
- [ ] tekton-task

## Resources

- [X] User Webhooks
- [X] Org Webhooks

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Farldka%2Fterraform-provider-artifacthub.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Farldka%2Fterraform-provider-artifacthub?ref=badge_large)