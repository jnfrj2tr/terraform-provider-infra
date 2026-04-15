package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDestination() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an Infra destination.",
		CreateContext: resourceDestinationCreate,
		ReadContext:   resourceDestinationRead,
		DeleteContext: resourceDestinationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the destination.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"kind": {
				Description: "Kind of the destination (e.g. kubernetes).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"unique_id": {
				Description: "Unique identifier for the destination.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(map[string]interface{})["client"].(*http.Client)
	host := meta.(map[string]interface{})["host"].(string)
	accessKey := meta.(map[string]interface{})["access_key"].(string)

	dest, err := createDestination(client, host, accessKey, createDestinationRequest{
		Name:     d.Get("name").(string),
		Kind:     d.Get("kind").(string),
		UniqueID: d.Get("unique_id").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dest.ID)
	return resourceDestinationRead(ctx, d, meta)
}

func resourceDestinationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(map[string]interface{})["client"].(*http.Client)
	host := meta.(map[string]interface{})["host"].(string)
	accessKey := meta.(map[string]interface{})["access_key"].(string)

	dest, err := readDestination(client, host, accessKey, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", dest.Name)
	_ = d.Set("kind", dest.Kind)
	_ = d.Set("unique_id", dest.UniqueID)
	return nil
}

func resourceDestinationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(map[string]interface{})["client"].(*http.Client)
	host := meta.(map[string]interface{})["host"].(string)
	accessKey := meta.(map[string]interface{})["access_key"].(string)

	if err := deleteDestination(client, host, accessKey, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
