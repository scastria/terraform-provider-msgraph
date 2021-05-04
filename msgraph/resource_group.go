package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-msgraph/msgraph/client"
	"net/http"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mail_nickname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"primary_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_unified": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mail_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	var visibility string
	if d.Get("is_public").(bool) {
		visibility = client.Public
	} else {
		visibility = client.Private
	}
	newGroup := client.Group{
		DisplayName:     d.Get("display_name").(string),
		Description:     d.Get("description").(string),
		MailNickname:    d.Get("mail_nickname").(string),
		MailEnabled:     d.Get("is_unified").(bool),
		SecurityEnabled: d.Get("security_enabled").(bool),
		Visibility:      visibility,
	}
	if (!newGroup.MailEnabled) && (!newGroup.SecurityEnabled) {
		d.SetId("")
		return diag.Errorf("A non-unified group MUST be security enabled")
	}
	if newGroup.MailEnabled {
		newGroup.GroupTypes = []string{client.Unified}
	} else {
		newGroup.GroupTypes = []string{}
	}
	primaryOwnerId, ok := d.GetOk("primary_owner_id")
	if ok {
		newGroup.Owners = []string{c.RequestPath(fmt.Sprintf(client.UserPathGet, primaryOwnerId.(string)))}
	}
	err := json.NewEncoder(&buf).Encode(newGroup)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GroupPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Group{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Id)
	d.Set("mail", retVal.Mail)
	d.Set("mail_enabled", retVal.MailEnabled)
	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupPathGet, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Group{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("display_name", retVal.DisplayName)
	d.Set("mail_nickname", retVal.MailNickname)
	d.Set("description", retVal.Description)
	d.Set("security_enabled", retVal.SecurityEnabled)
	_, hasUnified := find(retVal.GroupTypes, client.Unified)
	d.Set("is_unified", hasUnified)
	d.Set("mail", retVal.Mail)
	d.Set("mail_enabled", retVal.MailEnabled)
	d.Set("is_public", retVal.GroupIsPublic())
	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	var visibility string
	if d.Get("is_public").(bool) {
		visibility = client.Public
	} else {
		visibility = client.Private
	}
	upGroup := client.Group{
		Id:              d.Id(),
		DisplayName:     d.Get("display_name").(string),
		Description:     d.Get("description").(string),
		MailNickname:    d.Get("mail_nickname").(string),
		MailEnabled:     d.Get("is_unified").(bool),
		SecurityEnabled: d.Get("security_enabled").(bool),
		Visibility:      visibility,
	}
	if (!upGroup.MailEnabled) && (!upGroup.SecurityEnabled) {
		return diag.Errorf("A non-unified group MUST be security enabled")
	}
	if upGroup.MailEnabled {
		upGroup.GroupTypes = []string{client.Unified}
	} else {
		upGroup.GroupTypes = []string{}
	}
	err := json.NewEncoder(&buf).Encode(upGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GroupPathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupPathGet, d.Id())
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
