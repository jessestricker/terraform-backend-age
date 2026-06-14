# terraform-backend-age

A Terraform backend with age encryption.

## Usage

Configure your Terraform workspace to use terraform-backend-age:

```terraform
terraform {
  backend "http" {
    address = "http://localhost:4321"
  }
}
```

Create the key which will be used for encryption:

```shell
age-keygen -o state-key.txt
```

Save this at a secure place.

Then, whenever you are working on your Terraform code, start the terraform-backend-age:

```shell
terraform-backend-age &
```

> [!TIP]
> You can either specify the path to the key file with the `-key-file` CLI option,
> or by saving the contents of the key file in the environment variable `TF_BACKEND_AGE_KEY`.
>
> For more information, run `terraform-backend-age --help`.

Finally, when you are done, shut it down by sending it an interrupt:

```shell
killall --signal INT --wait terraform-backend-age
```
