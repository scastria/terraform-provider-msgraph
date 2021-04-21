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

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"search_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mail_nickname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mail_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_unified": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	mailNickname, ok := d.GetOk("mail_nickname")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "mailNickname", mailNickname.(string)))
	}
	mail, ok := d.GetOk("mail")
	if ok {
		filters = append(filters, fmt.Sprintf(client.FilterValue, "mail", mail.(string)))
	}
	if len(filters) > 0 {
		requestQuery[client.Filter] = []string{strings.Join(filters, client.FilterAnd)}
	}
	requestHeaders := http.Header{
		client.ConsistencyLevel: []string{client.ConsistencyLevelEventual},
	}
	requestPath := fmt.Sprintf(client.GroupPath)
	body, err := c.HttpRequest(http.MethodGet, requestPath, requestQuery, requestHeaders, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.GroupCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if len(retVal.Groups) != 1 {
		d.SetId("")
		filters = append(filters, searchDisplayName.(string))
		return diag.Errorf("Filter criteria does not result in a single group: %s", filters)
	}
	d.Set("display_name", retVal.Groups[0].DisplayName)
	d.Set("mail_nickname", retVal.Groups[0].MailNickname)
	d.Set("description", retVal.Groups[0].Description)
	d.Set("mail", retVal.Groups[0].Mail)
	d.Set("security_enabled", retVal.Groups[0].SecurityEnabled)
	d.Set("mail_enabled", retVal.Groups[0].MailEnabled)
	_, hasUnified := find(retVal.Groups[0].GroupTypes, client.Unified)
	d.Set("is_unified", hasUnified)
	d.SetId(retVal.Groups[0].Id)
	return diags
}
