terraform {
  required_providers {
    artifacthub = {
      version = "0.1.0"
      source  = "arldka.cloud/dev/artifacthub"
    }
  }
}

resource "artifacthub_user_webhook" "test_single_package" {
  name        = "testSinglePackage"
  description = "test"
  url         = "https://test.com"
  packages {
    package_id = "75ee6e00-b4d5-429e-9d82-33ab730081ff"
  }
  event_kinds = [0]
  active      = false
}