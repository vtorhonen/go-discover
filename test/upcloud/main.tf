provider "upcloud" {}

resource "upcloud_server" "test" {
  count    = 2
  zone     = "de-fra1"
  hostname = "go-discover-test-${count.index}.example.tld"

  cpu = "1"
  mem = "1024"

  network_interface {
    type = "utility"
  }

  storage_devices {
    size    = 10
    action  = "clone"
    tier    = "maxiops"
    storage = "Ubuntu Server 16.04 LTS (Xenial Xerus)"
  }
}

resource "upcloud_tag" "test" {
  name    = "go-discover-test-tag"
  servers = upcloud_server.test.*.id
}
