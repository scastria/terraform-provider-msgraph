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

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"search_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"surname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_principal_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	givenName, ok := d.GetOk("given_name")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "givenName", givenName.(string)))
	}
	surname, ok := d.GetOk("surname")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "surname", surname.(string)))
	}
	jobTitle, ok := d.GetOk("job_title")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "jobTitle", jobTitle.(string)))
	}
	mail, ok := d.GetOk("mail")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "mail", mail.(string)))
	}
	userPrincipalName, ok := d.GetOk("user_principal_name")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "userPrincipalName", userPrincipalName.(string)))
	}
	if len(filters) > 0 {
		requestQuery[client.Filter] = []string{strings.Join(filters, client.FilterAnd)}
	}
	requestHeaders := http.Header{
		client.ConsistencyLevel: []string{client.ConsistencyLevelEventual},
	}
	requestPath := fmt.Sprintf(client.UserPath)
	body, err := c.HttpRequest(http.MethodGet, requestPath, requestQuery, requestHeaders, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.UserCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if len(retVal.Users) != 1 {
		d.SetId("")
		return diag.Errorf("Filter criteria does not result in a single user")
	}
	d.Set("display_name", retVal.Users[0].DisplayName)
	d.Set("given_name", retVal.Users[0].GivenName)
	d.Set("surname", retVal.Users[0].Surname)
	d.Set("job_title", retVal.Users[0].JobTitle)
	d.Set("mail", retVal.Users[0].Mail)
	d.Set("user_principal_name", retVal.Users[0].UserPrincipalName)
	d.SetId(retVal.Users[0].Id)
	return diags
}
