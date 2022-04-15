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

func resourceEnterpriseAppOwner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseAppOwnerCreate,
		ReadContext:   resourceEnterpriseAppOwnerRead,
		DeleteContext: resourceEnterpriseAppOwnerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enterprise_app_id": {
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

func resourceEnterpriseAppOwnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newEnterpriseAppOwner := client.EnterpriseAppOwner{
		EnterpriseAppId: d.Get("enterprise_app_id").(string),
		OwnerId:         d.Get("owner_id").(string),
	}
	newEnterpriseAppOwner.OdataId = c.RequestPath(fmt.Sprintf(client.UserPathGet, newEnterpriseAppOwner.OwnerId))
	err := json.NewEncoder(&buf).Encode(newEnterpriseAppOwner)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnterpriseAppOwnerPathCreate, newEnterpriseAppOwner.EnterpriseAppId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newEnterpriseAppOwner.EnterpriseAppOwnerEncodeId())
	return diags
}

func resourceEnterpriseAppOwnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	enterpriseAppId, ownerId := client.EnterpriseAppOwnerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnterpriseAppOwnerPath, enterpriseAppId)
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
	d.Set("enterprise_app_id", enterpriseAppId)
	d.Set("owner_id", ownerId)
	return diags
}

func resourceEnterpriseAppOwnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	enterpriseAppId, ownerId := client.EnterpriseAppOwnerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnterpriseAppOwnerPathDelete, enterpriseAppId, ownerId)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
