package iam

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	groupName := d.Get("group_name").(string)

	req := &iam.GetGroupInput{
		GroupName: aws.String(groupName),
	}

	var users []*iam.User
	var group *iam.Group

	log.Printf("[DEBUG] Reading IAM Group: %s", req)
	err := conn.GetGroupPages(req, func(page *iam.GetGroupOutput, lastPage bool) bool {
		if group == nil {
			group = page.Group
		}
		users = append(users, page.Users...)
		return !lastPage
	})
	if err != nil {
		return fmt.Errorf("Error getting group: %w", err)
	}
	if group == nil {
		return fmt.Errorf("no IAM group found")
	}

	d.SetId(aws.StringValue(group.GroupId))
	d.Set("arn", group.Arn)
	d.Set("path", group.Path)
	d.Set("group_id", group.GroupId)
	if err := d.Set("users", dataSourceGroupUsersRead(users)); err != nil {
		return fmt.Errorf("error setting users: %w", err)
	}

	return nil
}

func dataSourceGroupUsersRead(iamUsers []*iam.User) []map[string]interface{} {
	users := make([]map[string]interface{}, 0, len(iamUsers))
	for _, i := range iamUsers {
		u := make(map[string]interface{})
		u["arn"] = aws.StringValue(i.Arn)
		u["user_id"] = aws.StringValue(i.UserId)
		u["user_name"] = aws.StringValue(i.UserName)
		u["path"] = aws.StringValue(i.Path)
		users = append(users, u)
	}
	return users
}
