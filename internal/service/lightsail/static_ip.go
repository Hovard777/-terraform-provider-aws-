package lightsail

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceStaticIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceStaticIPCreate,
		Read:   resourceStaticIPRead,
		Delete: resourceStaticIPDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStaticIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LightsailConn()

	name := d.Get("name").(string)
	log.Printf("[INFO] Allocating Lightsail Static IP: %q", name)
	_, err := conn.AllocateStaticIp(&lightsail.AllocateStaticIpInput{
		StaticIpName: aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("creating Lightsail Static IP: %w", err)
	}

	d.SetId(name)

	return resourceStaticIPRead(d, meta)
}

func resourceStaticIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LightsailConn()

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Lightsail Static IP: %q", name)
	out, err := conn.GetStaticIp(&lightsail.GetStaticIpInput{
		StaticIpName: aws.String(name),
	})
	if err != nil {
		if tfawserr.ErrCodeEquals(err, lightsail.ErrCodeNotFoundException) {
			log.Printf("[WARN] Lightsail Static IP (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Lightsail Static IP (%s):%w", d.Id(), err)
	}

	d.Set("arn", out.StaticIp.Arn)
	d.Set("ip_address", out.StaticIp.IpAddress)
	d.Set("support_code", out.StaticIp.SupportCode)

	return nil
}

func resourceStaticIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LightsailConn()

	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting Lightsail Static IP: %q", name)
	_, err := conn.ReleaseStaticIp(&lightsail.ReleaseStaticIpInput{
		StaticIpName: aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("deleting Lightsail Static IP (%s):%w", d.Id(), err)
	}
	return nil
}
