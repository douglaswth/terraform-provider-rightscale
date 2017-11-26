package rs

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rs/rs/rsc"
)

// Example:
//
// resource "rs_cm_security_group" "ssh" {
//     cloud_href = ${data.rs_cm_cloud.ec2_us_east_1.id}
//     network_href = ${resource.network.my_network.id}
//     description = "my security group"
// }

func resourceCMSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateFunc("rs_cm", "security_groups", securityGroupWriteFields),

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "ID of the security group cloud",
				Required:    true,
				ForceNew:    true,
			},
			"deployment_href": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "ID of the security group cloud",
				Required:    true,
				ForceNew:    true,
			},
			"network_href": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Read-only fields
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
		},
	}
}

func securityGroupWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"cloud_href": d.Get("cloud_href"),
		"name":       d.Get("name"),
	}
	for _, f := range []string{
		"deployment_href", "description", "network_href",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return fields
}