# Terraform GSuite Provider

This is a [Terraform](https://github.com/hashicorp/terraform) provider for managing GSuite. It is based on [Seth Vargo's Google Calendar provider](https://github.com/sethvargo/terraform-provider-googlecalendar).


## Installation

1. Download the latest compiled binary from [GitHub releases](https://github.com/paultyng/terraform-provider-gsuite/releases).

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins
    $ mv terraform-provider-googlecalendar $HOME/.terraform.d/plugins/terraform-provider-googlecalendar
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```

    This will find the plugin locally.

