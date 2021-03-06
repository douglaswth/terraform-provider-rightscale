package rightscale

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-rightscale/rightscale/rsc"
)

// resource "rightscale_instance" "my_instance" {
//   cloud_href = "/api/clouds/6"
//   image_href = "/api/clouds/6/images/3TRNL47PJB97N"
//   instance_type_href = "/api/clouds/6/instance_types/8SCHNH0JBHE1R"
//   deployment_href = "/api/deployments/934588004"
//   name = "My Instance"
// }

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceInstanceCreate,
		Update: resourceInstanceUpdate,

		// Note: none of the fields have "ForceNew" set because all
		// fields can be modified as long as the instance is not
		// running.
		Schema: map[string]*schema.Schema{
			"associate_public_ip_address": {
				Description: "Specify whether or not you want a public IP assigned when this Instance is launched. Only applies to Network-enabled Instances. If this is not specified, it will default to true.",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"cloud_href": {
				Description: "The ID of the instance cloud",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloud_specific_attributes": instanceCloudAttributes,
			"datacenter_href": {
				Description: "The ID of the instance datacenter",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"deployment_href": {
				Description: "The ID of the instance deployment",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"image_href": {
				Description: "The ID of the instance image",
				Type:        schema.TypeString,
				Required:    true,
			},
			"inputs": {
				Description: "Inputs associated with an instance when incarnated from a server or server array",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeMap},
				Optional:    true,
				ForceNew:    true,
				// Uncomment validation when tf adds supported operation on lists or sets.
				//ValidateFunc: validation.StringMatch(regexp.MustCompile("\\w+=\\w+:\\w+"), "values must be in format of 'key=type:value'"),
			},
			"instance_type_href": {
				Description: "The ID of the instance type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ip_forwarding_enabled": {
				Description: "Allows this Instance to send and receive network traffic when the source and destination IP addresses do not match the IP address of this Instance.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"kernel_image_href": {
				Description: "The ID of the instance kernel image.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The name of the instance",
				Type:        schema.TypeString,
				Required:    true,
			},
			"placement_group_href": {
				Description: "The placement group to launch the instance in. Not supported by all clouds & instance types.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ramdisk_image_href": {
				Description: "The ID of the ramdisk image",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"security_group_hrefs": {
				Description: "The IDs of the security groups",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"server_template_href": {
				Description: "The ID of the server template.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"ssh_key_href": {
				Description: "The ID of the SSH key to use",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"subnet_hrefs": {
				Description: "The IDs of the instance subnets",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"user_data": {
				Description: "User data that RightScale automatically passes to your instance at boot time",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"locked": {
				Description: "whether instance is locked, a locked instance cannot be terminated or deleted",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Description: "The private ip address for the instance",
				Optional:    true,
				ForceNew:    true,
			},

			// Read-only fields
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"public_ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"pricing_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

var instanceCloudAttributes = &schema.Schema{
	Description: "Cloud specific attributes that have no generic rightscale abstraction",
	Type:        schema.TypeList,
	MaxItems:    1,
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"admin_username": {
				Description: "The user that will be granted administrative privileges. Supported by AzureRM cloud only. For more information, review the documentation.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"automatic_instance_store_mapping": {
				Description:  "A flag indicating whether instance store mapping should be enabled. Not supported in all Clouds.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"availability_set": {
				Description: "Availability set for raw instance. Supported by Azure v2 cloud only.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"create_boot_volume": {
				Description:  "If enabled, the instance will launch into volume storage. Otherwise, it will boot to local storage.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"create_default_port_forwarding_rules": {
				Description:  "Automatically create default port forwarding rules (enabled by default). Supported by Azure cloud only.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"delete_boot_volume": {
				Description:  "If enabled, the associated volume will be deleted when the instance is terminated.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"disk_gb": {
				Description: "The size of root disk. Supported by UCA cloud only.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ebs_optimized": {
				Description:  "Whether the instance is able to connect to IOPS-enabled volumes.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"iam_instance_profile": {
				Description: "The name or ARN of the IAM Instance Profile (IIP) to associate with the instance (Amazon only)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"keep_alive_id": {
				Description: "The id of keep alive. Supported by UCA cloud only.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"keep_alive_url": {
				Description: "he ulr of keep alive. Supported by UCA cloud only.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"local_ssd_count": {
				Description: "Additional local SSDs. Supported by GCE cloud only",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"local_ssd_interface": {
				Description: "The type of SSD(s) to be created. Supported by GCE cloud only",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"max_spot_price": {
				Description: "Specify the max spot price you will pay for. Required when 'pricing_type' is 'spot'. Only applies to clouds which support spot-pricing and when 'spot' is chosen as the 'pricing_type'. Should be a Float value >= 0.001, eg: 0.095, 0.123, 1.23, etc...",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"memory_mb": {
				Description: "The size of instance memory. Supported by UCA cloud only.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"metadata": {
				Description: "Extra data used for configuration, in query string format.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"num_cores": {
				Description: "The number of instance cores. Supported by UCA cloud only.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"placement_tenancy": {
				Description:  "The tenancy of the server you want to launch. A server with a tenancy of dedicated runs on single-tenant hardware and can only be launched into a VPC.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"default", "dedicated"}, false),
			},
			"preemptible": {
				Description:  "Launch a preemptible instance. A preemptible instance costs much less, but lasts only 24 hours. It can be terminated sooner due to system demands. Supported by GCE cloud only.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"pricing_type": {
				Description:  "Specify whether or not you want to utilize 'fixed' (on-demand) or 'spot' pricing. Defaults to 'fixed' and only applies to clouds which support spot instances. Can only be set on when creating a new Instance, Server, or ServerArray, or when updating a Server or ServerArray's next_instance.WARNING: By using spot pricing, you acknowledge that your instance/server/array may not be able to be launched (and arrays may be unable to grow) as newly launched instances might be stuck in bidding, and/or existing instances may be terminated at any time, due to the cloud's spot pricing changes and availability.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"fixed", "spot"}, false),
			},
			"root_volume_performance": {
				Description: "The number of IOPS (I/O Operations Per Second) this root volume should support. Only available on clouds supporting performance provisioning.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"root_volume_size": {
				Description: "The size for root disk. Not supported in all Clouds.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"root_volume_type_uid": {
				Description: "The type of root volume for instance. Only available on clouds supporting root volume type.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"service_account": {
				Description: "Email of service account for instance. Scope will default to cloud-platform. Supported by GCE cloud only.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	},
}

func resourceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	var mustLock bool
	{
		locked, ok := d.GetOk("locked")
		mustLock = ok && locked.(bool)
	}

	client := m.(rsc.Client)
	res, err := client.Create("rs_cm", "instances", instanceWriteFields(d))
	if err != nil {
		return err
	}
	for k, v := range res.Fields {
		d.Set(k, v)
	}

	if mustLock {
		if err := updateLock(d, client, "instances"); err != nil {
			// Attempt to delete previously created instance, ignore errors
			client.Delete(res.Locator)
			return err
		}
		d.Set("locked", true)
	}

	// set ID last so Terraform does not assume the instance has been
	// created until all operations have completed successfully.
	d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	client := m.(rsc.Client)
	loc, err := locator(d)
	if err != nil {
		return err
	}

	// update lock
	if err := updateLock(d, client, "instances"); err != nil {
		return handleRSCError(d, err)
	}
	d.SetPartial("locked")

	// then the other fields
	// Skip updating instance if only lock status is changing
	for updateFields := range instanceUpdateFields(d) {
		if d.HasChange(updateFields) {
			if err := client.Update(loc, instanceUpdateFields(d)); err != nil {
				return handleRSCError(d, err)
			}
			break
		}
	}

	d.Partial(false)
	return nil
}

func instanceUpdateFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"deployment_href",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	if a, ok := d.GetOk("cloud_specific_attributes"); ok {
		fields["cloud_specific_attributes"] = a.([]interface{})[0]
	}
	// Note that this function is used only for raw instances, which do not support inputs.
	// Invoke instanceWriteFieldsFromMap from other resources for instances derived from rs objects that support inputs.
	return rsc.Fields{"cloud_href": d.Get("cloud_href"), "instance": fields}
}

func instanceWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"associate_public_ip_address", "datacenter_href",
		"deployment_href", "image_href", "instance_type_href",
		"ip_forwarding_enabled", "kernel_image_href", "name",
		"placement_group_href", "ramdisk_image_href",
		"security_group_hrefs", "ssh_key_href", "subnet_hrefs",
		"user_data", "server_template_href", "private_ip_address",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	if a, ok := d.GetOk("cloud_specific_attributes"); ok {
		fields["cloud_specific_attributes"] = a.([]interface{})[0]
	}
	// Note that this function is used only for raw instances, which do not support inputs.
	// Invoke instanceWriteFieldsFromMap from other resources for instances derived from rs objects that support inputs.
	return rsc.Fields{"cloud_href": d.Get("cloud_href"), "instance": fields}
}

func instanceWriteFieldsFromMap(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"associate_public_ip_address", "datacenter_href",
		"deployment_href", "image_href", "instance_type_href",
		"ip_forwarding_enabled", "kernel_image_href", "name",
		"placement_group_href", "ramdisk_image_href",
		"security_group_hrefs", "ssh_key_href", "subnet_hrefs",
		"user_data", "server_template_href", "cloud_href", "private_ip_address",
	} {
		if v, ok := d[f]; ok {
			fields[f] = v
		}
	}
	if a, ok := d["cloud_specific_attributes"]; ok && len(a.([]interface{})) > 0 {
		fields["cloud_specific_attributes"] = a.([]interface{})[0]
	}
	if a, ok := d["inputs"]; ok && len(a.([]interface{})) > 0 {
		if r, err := cmInputs(d["inputs"].([]interface{})); err != nil {
			log.Printf("[ERROR]: %v", err)
		} else {
			fields["inputs"] = r["inputs"]
		}
	}
	return fields
}
