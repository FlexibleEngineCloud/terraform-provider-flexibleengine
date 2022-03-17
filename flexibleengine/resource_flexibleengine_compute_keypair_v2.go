package flexibleengine

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/keypairs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeKeypairV2Create,
		Read:   resourceComputeKeypairV2Read,
		Delete: resourceComputeKeypairV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_key_path"},
			},
			"private_key_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeKeypairV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	var privateKeyPath string
	pk, isImport := d.GetOk("public_key")
	if !isImport {
		privateKeyPath, err = getKeyFilePath(d)
		if err != nil {
			return fmt.Errorf("private_key_path is invalid: %s", err)
		}
	}

	createOpts := keypairs.CreateOpts{
		Name:      d.Get("name").(string),
		PublicKey: pk.(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	kp, err := keypairs.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine keypair: %s", err)
	}

	d.SetId(kp.Name)

	if !isImport {
		if err = writeToPemFile(privateKeyPath, kp.PrivateKey); err != nil {
			return fmt.Errorf("Unable to generate private key: %s", err)
		}
		d.Set("private_key_path", privateKeyPath)
	}

	return resourceComputeKeypairV2Read(d, meta)
}

func getKeyFilePath(d *schema.ResourceData) (string, error) {
	if path, ok := d.GetOk("private_key_path"); ok {
		name := path.(string)
		f, err := os.Stat(name)
		if err != nil {
			return "", err
		}
		if f.IsDir() {
			return "", fmt.Errorf("%s must be a file", name)
		}

		return name, nil
	}

	keypairName := d.Get("name").(string)
	return fmt.Sprintf("%s.pem", keypairName), nil
}

func writeToPemFile(path, privateKey string) error {
	var err error
	// If the private key exists, give it write permission for editing (-rw-------) for root user.
	if _, err = ioutil.ReadFile(path); err == nil {
		os.Chmod(path, 0600)
	}
	if err = ioutil.WriteFile(path, []byte(privateKey), 0600); err != nil {
		return err
	}
	os.Chmod(path, 0400) // read-only permission (-r--------).
	return nil
}

func resourceComputeKeypairV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	kp, err := keypairs.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "keypair")
	}

	d.Set("name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeKeypairV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	privateKey := d.Get("private_key_path").(string)

	err = keypairs.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine keypair: %s", err)
	}
	d.SetId("")

	if privateKey != "" {
		log.Printf("try to remove the private key (%s) after the keypair is deleted", privateKey)
		err = os.Remove(privateKey)
		if err != nil {
			log.Printf("failed to remove private key %s: %s", privateKey, err)
		}
	}
	return nil
}
