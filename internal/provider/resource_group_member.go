package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages membership of a user in an Infra group.",
		CreateContext: resourceGroupMemberCreate,
		ReadContext:   resourceGroupMemberRead,
		DeleteContext: resourceGroupMemberDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"user_id": {
				Description: "The ID of the user to add to the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*http.Client)
	host := d.Get("host").(string)
	token := d.Get("token").(string)
	groupID := d.Get("group_id").(string)
	userID := d.Get("user_id").(string)

	member, err := addGroupMember(client, host, token, groupID, userID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", member.GroupID, member.UserID))
	return resourceGroupMemberRead(ctx, d, meta)
}

func resourceGroupMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*http.Client)
	host := d.Get("host").(string)
	token := d.Get("token").(string)
	groupID := d.Get("group_id").(string)
	userID := d.Get("user_id").(string)

	_, err := readGroupMember(client, host, token, groupID, userID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*http.Client)
	host := d.Get("host").(string)
	token := d.Get("token").(string)
	groupID := d.Get("group_id").(string)
	userID := d.Get("user_id").(string)

	if err := removeGroupMember(client, host, token, groupID, userID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func init() {
	registerResource("infra_group_member", resourceGroupMember())
}
