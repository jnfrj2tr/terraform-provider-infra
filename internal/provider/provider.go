package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("INFRA_HOST", "https://api.infrahq.com"),
					Description: "The Infra API host URL. Defaults to https://api.infrahq.com if not set.",
				},
				"access_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("INFRA_ACCESS_KEY", nil),
					Description: "The Infra API access key.",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"infra_grant": resourceGrant(),
			},
			DataSourcesMap: map[string]*schema.Resource{},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

type apiClient struct {
	Host      string
	AccessKey string
}

func configure(version string, p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		host := d.Get("host").(string)
		if host == "" {
			host = os.Getenv("INFRA_HOST")
		}
		accessKey := d.Get("access_key").(string)
		if accessKey == "" {
			accessKey = os.Getenv("INFRA_ACCESS_KEY")
		}
		client := &apiClient{
			Host:      host,
			AccessKey: accessKey,
		}
		return client, nil
	}
}
