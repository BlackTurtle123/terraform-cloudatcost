# terraform-cloudatcost
cloudatcost provider terraform

## References
https://github.com/cloudatcost/api

https://github.com/masayukioguni/go-cloudatcost/cloudatcost

https://github.com/Blackturtle123/go-cloudatcost/cloudatcost
## Features
Creating vm.

Choosing what run mode.

Update run mode.

Update vm, vm will be recreated.
## Example
```
provider "cloudatcost" {
  api_key     = "key"
  login       = "email"
}

resource "cloudatcost_instance" "servera" {
  //currently not able of removing label, only updating it.
  "label"="serverd",
  "cpu"="1",
  "ram"="512",
  "storage"="10",
  "os"="CentOS 6.7 64bit",
  "runmode"="safe",
}

resource "cloudatcost_instance" "serverb" {
  "cpu"="1",
  "ram"="512",
  "storage"="20",
  "os"="CentOS 6.7 64bit",
  "runmode"="normal"
  //Without depends terraform makes those at the same time
  //This making it impossible for the API to know which is which server
  depends_on = ["cloudatcost_instance.servera"]
}
```
