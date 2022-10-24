terraform {
  required_providers {
    ontap = {
      version = "0.1"
      source  = "netapp/com/ontap"
    }
  }
}

provider "ontap" {
  hostname = "cluster1.lab.tynsoe.org"
  username = "admin"
  password = "Netapp01"
}

data "ontap_svm" "svm1" {
  uuid = "bbb37711-8c11-11e7-933b-000c29cf11c2" // vs1
  //uuid = "1b7f0449-7329-11e7-a844-000c29cf11c2" // docker
}

output "test" {
  value = data.ontap_svm.svm1
}