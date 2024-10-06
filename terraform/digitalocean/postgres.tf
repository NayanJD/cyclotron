resource "digitalocean_tag" "volume_name" {
  name = "VOLUME_NAME:cyclotron-postgres-volume"
}

resource "digitalocean_droplet" "postgres" {
  image  = "ubuntu-24-04-x64"
  name   = "cyclotron-postgres"
  region = "blr1"
  size   = "m-2vcpu-16gb"
  ssh_keys = [var.ssh_key_id]
  tags     = [digitalocean_tag.volume_name.id]
  user_data = file("scripts/postgres.sh")
}

resource "digitalocean_volume" "postgres-volume" {
  region                   = "blr1"
  name                     = "cyclotron-postgres-volume" 
  size                     = 10 
  description              = "Block storage for cyclotron postgres"
  # initial_filesystem_label = var.block_storage_filesystem_label
  # initial_filesystem_type  = var.block_storage_filesystem_type
}

output "cyclotron-postgres-droplet-id" {
  value = digitalocean_droplet.postgres.ipv4_address
}

resource "digitalocean_volume_attachment" "postgres-volume-attachment" {
  droplet_id = digitalocean_droplet.postgres.id
  volume_id  = digitalocean_volume.postgres-volume.id
}
