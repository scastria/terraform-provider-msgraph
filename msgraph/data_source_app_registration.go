package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-msgraph/msgraph/client"
	"net/http"
	"net/url"
	"strings"
)

func dataSourceAppRegistration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppRegistrationRead,
		Schema: map[string]*schema.Schema{
			"search_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAppRegistrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestQuery := url.Values{}
	//Handle searching
	searchDisplayName, ok := d.GetOk("search_display_name")
	if ok {
		requestQuery[client.Search] = []string{fmt.Sprintf(client.SearchValue, "displayName", searchDisplayName.(string))}
	}
	//Handle filtering
	filters := []string{}
	displayName, ok := d.GetOk("display_name")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "displayName", displayName.(string)))
	}
	appId, ok := d.GetOk("app_id")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "appId", appId.(string)))
	}
	if len(filters) > 0 {
		requestQuery[client.Filter] = []string{strings.Join(filters, client.FilterAnd)}
	}
	requestHeaders := http.Header{
		client.ConsistencyLevel: []string{client.ConsistencyLevelEventual},
	}
	requestPath := fmt.Sprintf(client.AppRegistrationPath)
	body, err := c.HttpRequest(http.MethodGet, requestPath, requestQuery, requestHeaders, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.AppRegistrationCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if len(retVal.AppRegistrations) != 1 {
		d.SetId("")
		return diag.Errorf("Filter criteria does not result in a single app registration")
	}
	d.Set("display_name", retVal.AppRegistrations[0].DisplayName)
	d.Set("app_id", retVal.AppRegistrations[0].AppId)
	d.SetId(retVal.AppRegistrations[0].Id)
	return diags
}
