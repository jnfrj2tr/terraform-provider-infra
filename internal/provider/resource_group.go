package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
	n}

func resource d *schema.ResourceData, meta interface{}) diag.Diagnostlient := meta.(*ClientConfig)

	name := d.Get("name").(string)
	group, err := createGroup(ctx, client, name)
	if err != nil {
		retErr(err)
	}

	d.SetId(group.ID)
	_ = d.Set("name", group.Name)

	return nil
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig)

	group, err := readGroup(ctx, client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if group == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", group.Name)

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig)

	if err := deleteGroup(ctx, client, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
