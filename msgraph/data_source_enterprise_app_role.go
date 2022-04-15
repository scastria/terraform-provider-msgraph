package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-msgraph/msgraph/client"
	"math"
	"net/http"
	"strings"
	"time"
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
			"wait_until_exists": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"wait_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"wait_polling_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntAtLeast(0),
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
	waitUntilExists := d.Get("wait_until_exists").(bool)
	waitTimeout := d.Get("wait_timeout").(int)
	waitPollingInterval := d.Get("wait_polling_interval").(int)
	var retVal *client.EnterpriseAppRole
	var err error
	if waitUntilExists {
		stateConf := &resource.StateChangeConf{
			Timeout:        time.Duration(waitTimeout) * time.Second,
			PollInterval:   time.Duration(waitPollingInterval) * time.Second,
			Pending:        []string{client.WaitNotExists},
			Target:         []string{client.WaitFound},
			NotFoundChecks: math.MaxInt,
			Refresh: func() (interface{}, string, error) {
				output, numEnterpriseAppRoles, err := checkEnterpriseAppRoleExists(ctx, d, c, requestPath)
				if err != nil {
					return nil, client.WaitError, err
				} else if numEnterpriseAppRoles > 1 {
					err = fmt.Errorf("Filter criteria does not result in a single enterprise app role: %s", getFilterString(d))
					return nil, client.WaitError, err
				} else if numEnterpriseAppRoles == 0 {
					tflog.Warn(ctx, "[WAIT]  Not exists.  Will try again...", "filters", getFilterString(d))
					return nil, client.WaitNotExists, nil
				} else {
					return output, client.WaitFound, nil
				}
			},
		}
		output, err2 := stateConf.WaitForStateContext(context.Background())
		if output != nil {
			retVal = output.(*client.EnterpriseAppRole)
		}
		err = err2
	} else {
		output, numEnterpriseAppRoles, err2 := checkEnterpriseAppRoleExists(ctx, d, c, requestPath)
		if err2 != nil {
			err = err2
		} else if numEnterpriseAppRoles != 1 {
			err = fmt.Errorf("Filter criteria does not result in a single enterprise app role: %s", getFilterString(d))
		} else {
			retVal = output
		}
	}
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.EnterpriseAppId = enterpriseAppId
	d.Set("display_name", retVal.DisplayName)
	d.Set("description", retVal.Description)
	d.SetId(retVal.EnterpriseAppRoleEncodeId())
	return diags
}

func checkEnterpriseAppRoleExists(ctx context.Context, d *schema.ResourceData, c *client.Client, requestPath string) (*client.EnterpriseAppRole, int, error) {
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return nil, -1, err
	}
	retVal := &client.EnterpriseAppRoleCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, -1, err
	}
	//Check for a quick exit
	if len(retVal.EnterpriseAppRoles) == 0 {
		return nil, 0, nil
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
	numRoles := len(filteredRoles2)
	if numRoles != 1 {
		return nil, numRoles, nil
	} else {
		return &(filteredRoles2[0]), numRoles, nil
	}
}

func getFilterString(d *schema.ResourceData) []string {
	searchDisplayName := d.Get("search_display_name")
	displayName := d.Get("display_name")
	description := d.Get("description")
	retVal := []string{
		"search_display_name:",
		searchDisplayName.(string),
		"display_name:",
		displayName.(string),
		"description:",
		description.(string),
	}
	return retVal
}
