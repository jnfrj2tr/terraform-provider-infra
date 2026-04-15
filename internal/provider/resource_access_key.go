package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessKeyCreate,
		ReadContext:   resourceAccessKeyRead,
		DeleteContext: resourceAccessKeyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identity_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config)

	key, err := createAccessKey(client, d.Get("name").(string), d.Get("identity_id").(string), d.Get("expires_at").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(key.ID)
	_ = d.Set("secret", key.Secret)

	return resourceAccessKeyRead(ctx, d, meta)
}

func resourceAccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config)

	key, err := readAccessKey(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if key == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", key.Name)
	_ = d.Set("identity_id", key.IdentityID)
	_ = d.Set("expires_at", key.ExpiresAt)

	return nil
}

func resourceAccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config)

	if err := deleteAccessKey(client, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
