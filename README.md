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

resource "cloudatcost_machine" "servera" {
  //currently not able of removing label, only updating it.
  "label"="serverd",
  "cpu"="1",
  "ram"="512",
  "storage"="10",
  "os"="26",
  "datacenter"="Developer-DC-3",
  "runmode"="NorMal",
}

resource "cloudatcost_machine" "serverb" {
  "cpu"="1",
  "ram"="512",
  "storage"="20",
  "os"="26",
  "datacenter"="Developer-DC-2",
  "runmode"="nOrmal"
  //Without depends terraform makes those at the same time
  //This making it impossible for the API to know which is which server
  depends_on = ["cloudatcost_machine.servera"]
}
```
