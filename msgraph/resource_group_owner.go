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

func resourceGroupOwner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupOwnerCreate,
		ReadContext:   resourceGroupOwnerRead,
		DeleteContext: resourceGroupOwnerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
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

func resourceGroupOwnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newGroupOwner := client.GroupOwner{
		GroupId: d.Get("group_id").(string),
		OwnerId: d.Get("owner_id").(string),
	}
	newGroupOwner.OdataId = c.RequestPath(fmt.Sprintf(client.UserPathGet, newGroupOwner.OwnerId))
	err := json.NewEncoder(&buf).Encode(newGroupOwner)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GroupOwnerPathCreate, newGroupOwner.GroupId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newGroupOwner.GroupOwnerEncodeId())
	return diags
}

func resourceGroupOwnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, ownerId := client.GroupOwnerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupOwnerPath, groupId)
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
	d.Set("group_id", groupId)
	d.Set("owner_id", ownerId)
	return diags
}

func resourceGroupOwnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, ownerId := client.GroupOwnerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupOwnerPathDelete, groupId, ownerId)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
