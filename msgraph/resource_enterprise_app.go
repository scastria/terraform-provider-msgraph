package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-msgraph/msgraph/client"
	"net/http"
)

func resourceEnterpriseApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseAppCreate,
		ReadContext:   resourceEnterpriseAppRead,
		UpdateContext: resourceEnterpriseAppUpdate,
		DeleteContext: resourceEnterpriseAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_role": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_member_types": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"User", "Application"}, false),
							},
						},
						"display_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceEnterpriseAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newEntApp := client.EnterpriseApp{
		AppId: d.Get("app_id").(string),
		Tags:  []string{client.IntegratedApp},
	}
	fillEnterpriseApp(&newEntApp, d)
	err := json.NewEncoder(&buf).Encode(newEntApp)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnterpriseAppPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.EnterpriseApp{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Id)
	d.Set("display_name", retVal.DisplayName)
	return diags
}

func expandAppRoles(appRoles *schema.Set) []client.AppRole {
	retVal := make([]client.AppRole, len(appRoles.List()))
	for i, item := range appRoles.List() {
		itemMap := item.(map[string]interface{})
		retVal[i] = client.AppRole{
			Id:                 itemMap["id"].(string),
			DisplayName:        itemMap["display_name"].(string),
			Description:        itemMap["description"].(string),
			IsEnabled:          itemMap["is_enabled"].(bool),
			AllowedMemberTypes: convertSetToArray(itemMap["allowed_member_types"].(*schema.Set)),
		}
		if retVal[i].Id == "" {
			retVal[i].Id, _ = uuid.GenerateUUID()
			itemMap["id"] = retVal[i].Id
		}
	}
	return retVal
}

func collapseAppRoles(appRoles []client.AppRole) []map[string]interface{} {
	var retVal []map[string]interface{}
	for _, item := range appRoles {
		itemMap := map[string]interface{}{}
		itemMap["id"] = item.Id
		itemMap["display_name"] = item.DisplayName
		itemMap["description"] = item.Description
		itemMap["is_enabled"] = item.IsEnabled
		itemMap["allowed_member_types"] = item.AllowedMemberTypes
		retVal = append(retVal, itemMap)
	}
	return retVal
}

func fillEnterpriseApp(c *client.EnterpriseApp, d *schema.ResourceData) {
	appRoles, ok := d.GetOk("app_role")
	if ok {
		c.AppRoles = expandAppRoles(appRoles.(*schema.Set))
	} else {
		c.AppRoles = []client.AppRole{}
	}
}

func resourceEnterpriseAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnterpriseAppPathGet, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.EnterpriseApp{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("app_id", retVal.AppId)
	d.Set("display_name", retVal.DisplayName)
	if retVal.AppRoles != nil {
		d.Set("app_role", collapseAppRoles(retVal.AppRoles))
	}
	return diags
}

func resourceEnterpriseAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	//Set all appRoles to disabled to allow deletion
	if d.HasChange("app_role") {
		oldAppRoles, _ := d.GetChange("app_role")
		oldAppRoleList := expandAppRoles(oldAppRoles.(*schema.Set))
		for i, _ := range oldAppRoleList {
			oldAppRoleList[i].IsEnabled = false
		}
		buf := bytes.Buffer{}
		upAppRoles := client.AppRoles{
			AppRoles: oldAppRoleList,
		}
		err := json.NewEncoder(&buf).Encode(upAppRoles)
		if err != nil {
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.EnterpriseAppPathGet, d.Id())
		requestHeaders := http.Header{
			headers.ContentType: []string{client.ApplicationJson},
		}
		_, err = c.HttpRequest(ctx, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	buf := bytes.Buffer{}
	upEntApp := client.EnterpriseApp{
		Id:    d.Id(),
		AppId: d.Get("app_id").(string),
	}
	fillEnterpriseApp(&upEntApp, d)
	err := json.NewEncoder(&buf).Encode(upEntApp)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnterpriseAppPathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceEnterpriseAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnterpriseAppPathGet, d.Id())
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
