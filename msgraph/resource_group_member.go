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

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMemberCreate,
		ReadContext:   resourceGroupMemberRead,
		DeleteContext: resourceGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newGroupMember := client.GroupMember{
		GroupId:  d.Get("group_id").(string),
		MemberId: d.Get("member_id").(string),
	}
	newGroupMember.OdataId = c.RequestPath(fmt.Sprintf(client.UserPathGet, newGroupMember.MemberId))
	err := json.NewEncoder(&buf).Encode(newGroupMember)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GroupMemberPathCreate, newGroupMember.GroupId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newGroupMember.GroupMemberEncodeId())
	return diags
}

func resourceGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, memberId := client.GroupMemberDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupMemberPath, groupId)
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
	//Search members looking for memberId
	var foundOwner bool
	foundOwner = false
	for _, u := range retVal.Users {
		if u.Id == memberId {
			foundOwner = true
			break
		}
	}
	if !foundOwner {
		d.SetId("")
		return diags
	}
	d.Set("group_id", groupId)
	d.Set("member_id", memberId)
	return diags
}

func resourceGroupMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, memberId := client.GroupMemberDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupMemberPathDelete, groupId, memberId)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
