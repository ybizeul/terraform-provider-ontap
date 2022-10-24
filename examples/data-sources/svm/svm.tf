terraform {
  required_providers {
    ontap = {
      version = "0.1"
      source  = "netapp/com/ontap"
    }
  }
}

provider "ontap" {
  hostname = "cluster.company.lan"
  username = "admin"
  password = "Netapp01"
}

data "ontap_svm" "svm1" {
  uuid="794e3659-21ed-11e8-97df-00a0985f3751"
}

output "test" {
  value = data.ontap_svm.svm1
}