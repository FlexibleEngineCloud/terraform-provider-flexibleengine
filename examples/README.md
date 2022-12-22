# FULL Terraform flexibleengine Example

This script will create the following resources (if enabled):

* VPC
* VPC Subnet
* Keypair
* Instances
* Security Group (Allow ICMP, 80/tcp, 22/tcp)
* Volumes
* EIP
* Loadbalancer
* Backup Policies and Backups

## Usage

This example will, by default not create Volumes. This is to show how to enable resources via parameters.
To enable Volume creation, set the **disk__size_gb** variable to a value > 10.

Add your variables to a file called **terraform.tfvars** or set them as environment variables.

```shell
export TF_VAR_access_key=your_access_key
export TF_VAR_secret_key=your_secret_key
export TF_VAR_domain_name=your_domain_name
```

Then run the following commands:

```shell
terraform init
terraform plan
terraform apply
```

Observe the output of the apply command, it will show you the external IP address of the loadbalancer.

##  Destroy

```shell
terraform destroy
```

You can encounter an error when destroying the bucket, this is because the bucket is not empty. You can delet
e the objects manually in the console and run the destroy command again.

## Available Variables

### Required

* **access_key** (your flexibleengine access_key)
* **secret_key** (your flexibleengine secret_key)
* **domain_name** (your flexibleengine domain name)
* You must have a **ssh_pub_key** file defined, or terraform will complain, see default path below.

### Optional

* **project** (this will prefix all your resources, *default=terraform*)
* **ssh_pub_key** (the path to the ssh public key you want to deploy, *default=~/.ssh/id_rsa.pub*)
* **instance_count** (affects the number of Floating IPs, Instances, Volumes and Ports, *default=1*)
* **flavor_name** (flavor of the created instances, *default=s6.medium.2*)
* **image_name** (image used for creating instances, *default=OBS Ubuntu 22.04*)
* **disk_size_gb** (size of the volumes in gigabytes, *default=None*)
