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
	"strings"
)

func dataSourceEnterpriseAppRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseAppRoleRead,
		Schema: map[string]*schema.Schema{
			"enterprise_app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceEnterpriseAppRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	enterpriseAppId := d.Get("enterprise_app_id").(string)
	c := m.(*client.Client)
	//app roles do not support searching and filtering so do it manually after reading all roles
	requestPath := fmt.Sprintf(client.EnterpriseAppRolePath, enterpriseAppId)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.EnterpriseAppRoleCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Check for a quick exit
	if len(retVal.EnterpriseAppRoles) == 0 {
		d.SetId("")
		return diag.Errorf("Filter criteria does not result in a single enterprise app role")
	}
	//Do manual searching
	searchedRoles := []client.EnterpriseAppRole{}
	searchDisplayName, ok := d.GetOk("search_display_name")
	if ok {
		searchDisplayNameLower := strings.ToLower(searchDisplayName.(string))
		for _, ear := range retVal.EnterpriseAppRoles {
			if strings.Contains(strings.ToLower(ear.DisplayName), searchDisplayNameLower) {
				searchedRoles = append(searchedRoles, ear)
			}
		}
	} else {
		searchedRoles = retVal.EnterpriseAppRoles
	}
	//Do manual filtering
	filteredRoles1 := []client.EnterpriseAppRole{}
	displayName, ok := d.GetOk("display_name")
	if ok {
		displayNameLower := strings.ToLower(displayName.(string))
		for _, ear := range searchedRoles {
			if strings.ToLower(ear.DisplayName) == displayNameLower {
				filteredRoles1 = append(filteredRoles1, ear)
			}
		}
	} else {
		filteredRoles1 = searchedRoles
	}
	filteredRoles2 := []client.EnterpriseAppRole{}
	description, ok := d.GetOk("description")
	if ok {
		descriptionLower := strings.ToLower(description.(string))
		for _, ear := range filteredRoles1 {
			if strings.ToLower(ear.Description) == descriptionLower {
				filteredRoles2 = append(filteredRoles2, ear)
			}
		}
	} else {
		filteredRoles2 = filteredRoles1
	}
	if len(filteredRoles2) != 1 {
		d.SetId("")
		filters := []string{
			searchDisplayName.(string),
			displayName.(string),
			description.(string),
		}
		return diag.Errorf("Filter criteria does not result in a single enterprise app role: %s", filters)
	}
	filteredRoles2[0].EnterpriseAppId = enterpriseAppId
	d.Set("display_name", filteredRoles2[0].DisplayName)
	d.Set("description", filteredRoles2[0].Description)
	d.SetId(filteredRoles2[0].EnterpriseAppRoleEncodeId())
	return diags
}
