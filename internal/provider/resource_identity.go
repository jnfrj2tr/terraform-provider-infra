package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCreate,
		ReadContext:   resourceIdentityRead,
		DeleteContext: resourceIdentityDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "user",
			},
			"identity_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIdentityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("invalid provider configuration")
	}

	name := d.Get("name").(string)
	kind := d.Get("kind	identity, err := createIdent nil {
		retErr(err)
	}

	dn	_ = d.Set("identity_id", identity.ID)
	return nil
}

func resourceIdentityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("invalid provider configuration")
	}

	identity, err := readIdentity(ctx, client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if identity == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", identity.Name)
	_ = d.Set("kind", identity.Kind)
	_ = d.Set("identity_id", identity.ID)
	return nil
}

func resourceIdentityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("invalid provider configuration")
	}

	if err := deleteIdentity(ctx, client, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
