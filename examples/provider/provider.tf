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