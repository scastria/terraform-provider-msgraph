package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-msgraph/msgraph/client"
	"net/http"
)

func resourceAppRegistrationOwner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppRegistrationOwnerCreate,
		ReadContext:   resourceAppRegistrationOwnerRead,
		DeleteContext: resourceAppRegistrationOwnerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"app_registration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAppRegistrationOwnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newAppRegistrationOwner := client.AppRegistrationOwner{
		AppRegistrationId: d.Get("app_registration_id").(string),
		OwnerId:           d.Get("owner_id").(string),
	}
	newAppRegistrationOwner.OdataId = c.RequestPath(fmt.Sprintf(client.UserPathGet, newAppRegistrationOwner.OwnerId))
	err := json.NewEncoder(&buf).Encode(newAppRegistrationOwner)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.AppRegistrationOwnerPathCreate, newAppRegistrationOwner.AppRegistrationId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newAppRegistrationOwner.AppRegistrationOwnerEncodeId())
	return diags
}

func resourceAppRegistrationOwnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	appRegistrationId, ownerId := client.AppRegistrationOwnerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.AppRegistrationOwnerPath, appRegistrationId)
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.UserCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Search owners looking for ownerId
	var foundOwner bool
	foundOwner = false
	for _, u := range retVal.Users {
		if u.Id == ownerId {
			foundOwner = true
			break
		}
	}
	if !foundOwner {
		d.SetId("")
		return diags
	}
	d.Set("app_registration_id", appRegistrationId)
	d.Set("owner_id", ownerId)
	return diags
}

func resourceAppRegistrationOwnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	appRegistrationId, ownerId := client.AppRegistrationOwnerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.AppRegistrationOwnerPathDelete, appRegistrationId, ownerId)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
