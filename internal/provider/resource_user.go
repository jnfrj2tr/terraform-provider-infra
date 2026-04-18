package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an Infra user.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email address of the user. Must be a valid email format.",
				ValidateFunc: validateEmail,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(map[string]interface{})["client"].(*http.Client)
	host := meta.(map[string]interface{})["host"].(string)
	token := meta.(map[string]interface{})["token"].(string)

	email := d.Get("email").(string)

	user, err := createUser(client, host, token, email)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.ID)
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(map[string]interface{})["client"].(*http.Client)
	host := meta.(map[string]interface{})["host"].(string)
	token := meta.(map[string]interface{})["token"].(string)

	user, err := readUser(client, host, token, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if user == nil {
		d.SetId("")
		return nil
	}

	// NOTE: the API returns the email in the Name field, not an Email field
	if err := d.Set("email", user.Name); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(map[string]interface{})["client"].(*http.Client)
	host := meta.(map[string]interface{})["host"].(string)
	token := meta.(map[string]interface{})["token"].(string)

	if err := deleteUser(client, host, token, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
