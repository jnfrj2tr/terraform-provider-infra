package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGrant() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an Infra grant, assigning a privilege to a user or group on a resource.",
		CreateContext: resourceGrantCreate,
		ReadContext:   resourceGrantRead,
		DeleteContext: resourceGrantDelete,
		Schema: map[string]*schema.Schema{
			"user": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"user", "group"},
				Description:  "The user identity to grant access to.",
			},
			"group": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"user", "group"},
				Description:  "The group to grant access to.",
			},
			"privilege": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The privilege level to grant (e.g. admin, view).",
			},
			"resource": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource to grant access to (e.g. kubernetes.cluster-name).",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGrantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("invalid provider configuration")
	}
	grant := map[string]string{
		"user":      d.Get("user").(string),
		"group":     d.Get("group").(string),
		"privilege": d.Get("privilege").(string),
		"resource":  d.Get("resource").(string),
	}
	id, err := createGrant(ctx, client, grant)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating grant: %w", err))
	}
	d.SetId(id)
	return resourceGrantRead(ctx, d, meta)
}

func resourceGrantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("invalid provider configuration")
	}
	grant, err := readGrant(ctx, client, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading grant: %w", err))
	}
	if grant == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("user", grant["user"])
	_ = d.Set("group", grant["group"])
	_ = d.Set("privilege", grant["privilege"])
	_ = d.Set("resource", grant["resource"])
	return nil
}

func resourceGrantDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("invalid provider configuration")
	}
	if err := deleteGrant(ctx, client, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("deleting grant: %w", err))
	}
	return nil
}
