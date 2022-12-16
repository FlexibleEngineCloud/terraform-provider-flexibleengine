# FULL Terraform flexibleengine Example

This script will create the following resources (if enabled):
* Volumes
* Floating IPs
* Neutron Ports
* Instances
* Keypair
* Network
* Subnet
* Router
* Router Interface
* Loadbalancer
* Templates
* Security Group (Allow ICMP, 80/tcp, 22/tcp)

## Resource Creation

This example will, by default not create Volumes. This is to show how to enable resources via parameters. To enable Volume creation, set the **disk__size_gb** variable to a value > 10.

## Available Variables

### Required

* **username** (your flexibleengine username)
* **password** (your flexibleengine password)
* **domain_name** (your flexibleengine domain name)
* You must have a **ssh_pub_key** file defined, or terraform will complain, see default path below.

### Optional
* **project** (this will prefix all your resources, _default=terraform_)
* **ssh_pub_key** (the path to the ssh public key you want to deploy, _default=~/.ssh/id_rsa.pub_)
* **instance_count** (affects the number of Floating IPs, Instances, Volumes and Ports, _default=1_)
* **flavor_name** (flavor of the created instances, _default=s1.medium_)
* **image_name** (image used for creating instances, _default=Standard_CentOS_7_latest_)
* **disk_size_gb** (size of the volumes in gigabytes, _default=None_)
