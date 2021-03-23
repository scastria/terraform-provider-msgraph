package msgraph

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-msgraph/msgraph/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSGRAPH_TENANT_ID", nil),
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("MSGRAPH_ACCESS_TOKEN", nil),
				ConflictsWith: []string{"client_id", "client_secret"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("MSGRAPH_CLIENT_ID", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"client_secret"},
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("MSGRAPH_CLIENT_SECRET", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"client_id"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"msgraph_group":       resourceGroup(),
			"msgraph_group_owner": resourceGroupOwner(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"msgraph_user": dataSourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tenantId := d.Get("tenant_id").(string)
	accessToken := d.Get("access_token").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	//Check for valid authentication
	if (clientId == "") && (clientSecret == "") && (accessToken == "") {
		return nil, diag.Errorf("You must specify either client_id/client_secret for Client Credentials Authentication or access_token")
	}

	var diags diag.Diagnostics
	c, err := client.NewClient(tenantId, accessToken, clientId, clientSecret)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
