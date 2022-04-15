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
	"net/url"
	"strings"
	"time"
)

func dataSourceEnterpriseApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseAppRead,
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
			"login_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logout_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEnterpriseAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	requestPath := fmt.Sprintf(client.EnterpriseAppPath)
	waitUntilExists := d.Get("wait_until_exists").(bool)
	waitTimeout := d.Get("wait_timeout").(int)
	waitPollingInterval := d.Get("wait_polling_interval").(int)
	var retVal *client.EnterpriseApp
	var err error
	if waitUntilExists {
		stateConf := &resource.StateChangeConf{
			Timeout:        time.Duration(waitTimeout) * time.Second,
			PollInterval:   time.Duration(waitPollingInterval) * time.Second,
			Pending:        []string{client.WaitNotExists},
			Target:         []string{client.WaitFound},
			NotFoundChecks: math.MaxInt,
			Refresh: func() (interface{}, string, error) {
				output, numEnterpriseApps, err := checkEnterpriseAppExists(ctx, c, requestPath, requestQuery, requestHeaders)
				if err != nil {
					return nil, client.WaitError, err
				} else if numEnterpriseApps > 1 {
					filters = append(filters, searchDisplayName.(string))
					err = fmt.Errorf("Filter criteria does not result in a single enterprise app: %s", filters)
					return nil, client.WaitError, err
				} else if numEnterpriseApps == 0 {
					tflog.Warn(ctx, "[WAIT]  Not exists.  Will try again...", "searchDisplayName", searchDisplayName, "filters", filters)
					return nil, client.WaitNotExists, nil
				} else {
					return output, client.WaitFound, nil
				}
			},
		}
		output, err2 := stateConf.WaitForStateContext(context.Background())
		if output != nil {
			retVal = output.(*client.EnterpriseApp)
		}
		err = err2
	} else {
		output, numEnterpriseApps, err2 := checkEnterpriseAppExists(ctx, c, requestPath, requestQuery, requestHeaders)
		if err2 != nil {
			err = err2
		} else if numEnterpriseApps != 1 {
			filters = append(filters, searchDisplayName.(string))
			err = fmt.Errorf("Filter criteria does not result in a single enterprise app: %s", filters)
		} else {
			retVal = output
		}
	}
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("display_name", retVal.DisplayName)
	d.Set("app_id", retVal.AppId)
	d.Set("login_url", retVal.LoginUrl)
	d.Set("logout_url", retVal.LogoutUrl)
	d.SetId(retVal.Id)
	return diags
}

func checkEnterpriseAppExists(ctx context.Context, c *client.Client, requestPath string, requestQuery url.Values, requestHeaders http.Header) (*client.EnterpriseApp, int, error) {
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, requestQuery, requestHeaders, &bytes.Buffer{})
	if err != nil {
		return nil, -1, err
	}
	retVal := &client.EnterpriseAppCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, -1, err
	}
	numEnterpriseApps := len(retVal.EnterpriseApps)
	if numEnterpriseApps != 1 {
		return nil, numEnterpriseApps, nil
	} else {
		return &(retVal.EnterpriseApps[0]), numEnterpriseApps, nil
	}
}
