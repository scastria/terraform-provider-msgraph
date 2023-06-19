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
	waitUntilExists := d.Get("wait_until_exists").(bool)
	waitTimeout := d.Get("wait_timeout").(int)
	waitPollingInterval := d.Get("wait_polling_interval").(int)
	var retVal *client.User
	var err error
	if waitUntilExists {
		stateConf := &resource.StateChangeConf{
			Timeout:        time.Duration(waitTimeout) * time.Second,
			PollInterval:   time.Duration(waitPollingInterval) * time.Second,
			Pending:        []string{client.WaitNotExists},
			Target:         []string{client.WaitFound},
			NotFoundChecks: math.MaxInt,
			Refresh: func() (interface{}, string, error) {
				output, numUsers, err := checkUserExists(ctx, c, requestPath, requestQuery, requestHeaders)
				if err != nil {
					return nil, client.WaitError, err
				} else if numUsers > 1 {
					filters = append(filters, searchDisplayName.(string))
					err = fmt.Errorf("Filter criteria does not result in a single user: %s", filters)
					return nil, client.WaitError, err
				} else if numUsers == 0 {
					tflog.Warn(ctx, "[WAIT]  Not exists.  Will try again...", map[string]interface{}{"searchDisplayName": searchDisplayName, "filters": filters})
					return nil, client.WaitNotExists, nil
				} else {
					return output, client.WaitFound, nil
				}
			},
		}
		output, err2 := stateConf.WaitForStateContext(context.Background())
		if output != nil {
			retVal = output.(*client.User)
		}
		err = err2
	} else {
		output, numUsers, err2 := checkUserExists(ctx, c, requestPath, requestQuery, requestHeaders)
		if err2 != nil {
			err = err2
		} else if numUsers != 1 {
			filters = append(filters, searchDisplayName.(string))
			err = fmt.Errorf("Filter criteria does not result in a single user: %s", filters)
		} else {
			retVal = output
		}
	}
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("display_name", retVal.DisplayName)
	d.Set("given_name", retVal.GivenName)
	d.Set("surname", retVal.Surname)
	d.Set("job_title", retVal.JobTitle)
	d.Set("mail", retVal.Mail)
	d.Set("user_principal_name", retVal.UserPrincipalName)
	d.SetId(retVal.Id)
	return diags
}

func checkUserExists(ctx context.Context, c *client.Client, requestPath string, requestQuery url.Values, requestHeaders http.Header) (*client.User, int, error) {
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, requestQuery, requestHeaders, &bytes.Buffer{})
	if err != nil {
		return nil, -1, err
	}
	retVal := &client.UserCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, -1, err
	}
	numUsers := len(retVal.Users)
	if numUsers != 1 {
		return nil, numUsers, nil
	} else {
		return &(retVal.Users[0]), numUsers, nil
	}
}
