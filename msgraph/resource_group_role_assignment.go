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

func resourceGroupRoleAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRoleAssignmentCreate,
		ReadContext:   resourceGroupRoleAssignmentRead,
		DeleteContext: resourceGroupRoleAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_app_role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupRoleAssignmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	groupId := d.Get("group_id").(string)
	enterpriseAppId, roleId := client.EnterpriseAppRoleDecodeId(d.Get("enterprise_app_role_id").(string))
	newGroupRoleAssignment := client.GroupRoleAssignment{
		GroupId:     groupId,
		ResourceId:  enterpriseAppId,
		PrincipalId: groupId,
		AppRoleId:   roleId,
	}
	err := json.NewEncoder(&buf).Encode(newGroupRoleAssignment)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GroupRoleAssignmentPath, groupId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.GroupRoleAssignment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.GroupId = retVal.PrincipalId
	d.SetId(retVal.GroupRoleAssignmentEncodeId())
	return diags
}

func resourceGroupRoleAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, assignmentId := client.GroupRoleAssignmentDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupRoleAssignmentPathGet, groupId, assignmentId)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.GroupRoleAssignment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("group_id", groupId)
	ear := client.EnterpriseAppRole{
		Id:              retVal.AppRoleId,
		EnterpriseAppId: retVal.ResourceId,
	}
	d.Set("enterprise_app_role_id", ear.EnterpriseAppRoleEncodeId())
	return diags
}

func resourceGroupRoleAssignmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, assignmentId := client.GroupRoleAssignmentDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GroupRoleAssignmentPathGet, groupId, assignmentId)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
