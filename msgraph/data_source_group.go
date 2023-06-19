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
			"is_public": {
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
	waitUntilExists := d.Get("wait_until_exists").(bool)
	waitTimeout := d.Get("wait_timeout").(int)
	waitPollingInterval := d.Get("wait_polling_interval").(int)
	var retVal *client.Group
	var err error
	if waitUntilExists {
		stateConf := &resource.StateChangeConf{
			Timeout:        time.Duration(waitTimeout) * time.Second,
			PollInterval:   time.Duration(waitPollingInterval) * time.Second,
			Pending:        []string{client.WaitNotExists},
			Target:         []string{client.WaitFound},
			NotFoundChecks: math.MaxInt,
			Refresh: func() (interface{}, string, error) {
				output, numGroups, err := checkGroupExists(ctx, c, requestPath, requestQuery, requestHeaders)
				if err != nil {
					return nil, client.WaitError, err
				} else if numGroups > 1 {
					filters = append(filters, searchDisplayName.(string))
					err = fmt.Errorf("Filter criteria does not result in a single group: %s", filters)
					return nil, client.WaitError, err
				} else if numGroups == 0 {
					tflog.Warn(ctx, "[WAIT]  Not exists.  Will try again...", map[string]interface{}{"searchDisplayName": searchDisplayName, "filters": filters})
					return nil, client.WaitNotExists, nil
				} else {
					return output, client.WaitFound, nil
				}
			},
		}
		output, err2 := stateConf.WaitForStateContext(context.Background())
		if output != nil {
			retVal = output.(*client.Group)
		}
		err = err2
	} else {
		output, numGroups, err2 := checkGroupExists(ctx, c, requestPath, requestQuery, requestHeaders)
		if err2 != nil {
			err = err2
		} else if numGroups != 1 {
			filters = append(filters, searchDisplayName.(string))
			err = fmt.Errorf("Filter criteria does not result in a single group: %s", filters)
		} else {
			retVal = output
		}
	}
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("display_name", retVal.DisplayName)
	d.Set("mail_nickname", retVal.MailNickname)
	d.Set("description", retVal.Description)
	d.Set("mail", retVal.Mail)
	d.Set("security_enabled", retVal.SecurityEnabled)
	d.Set("mail_enabled", retVal.MailEnabled)
	_, hasUnified := find(retVal.GroupTypes, client.Unified)
	d.Set("is_unified", hasUnified)
	d.Set("is_public", retVal.GroupIsPublic())
	d.SetId(retVal.Id)
	return diags
}

func checkGroupExists(ctx context.Context, c *client.Client, requestPath string, requestQuery url.Values, requestHeaders http.Header) (*client.Group, int, error) {
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, requestQuery, requestHeaders, &bytes.Buffer{})
	if err != nil {
		return nil, -1, err
	}
	retVal := &client.GroupCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, -1, err
	}
	numGroups := len(retVal.Groups)
	if numGroups != 1 {
		return nil, numGroups, nil
	} else {
		return &(retVal.Groups[0]), numGroups, nil
	}
}
