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

resource "ontap_qtree" "qtree1" {
 volume_uuid = "023b1262-6a22-4134-8972-33b370e5474c"
 svm_uuid = "ee8c5c70-71a7-11e9-abb6-000c29cf11c2"
 name = "myQtree3"
 #unix_permissions = 755
 security_style = "ntfs"
}

resource "ontap_qtree" "qtree2" {
 volume_uuid = "023b1262-6a22-4134-8972-33b370e5474c"
 svm_uuid = "ee8c5c70-71a7-11e9-abb6-000c29cf11c2"
 name = "myQtree4"
 unix_permissions = 755
 security_style = "unix"
}

output "qtree" {
 value = "(${resource.ontap_qtree.qtree1.id}) ${resource.ontap_qtree.qtree1.name} ${resource.ontap_qtree.qtree1.path}\n(${resource.ontap_qtree.qtree2.id}) ${resource.ontap_qtree.qtree2.name} ${resource.ontap_qtree.qtree2.path}"
}

# data "ontap_qtree" "qtree1" {
#   uuid = "023b1262-6a22-4134-8972-33b370e5474c/2"
# }

# output "test" {
#   value = data.ontap_qtree.qtree1
# }