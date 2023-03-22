package workspaces

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"bundle_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				RequiredWith:  []string{"user_name"},
				ConflictsWith: []string{"workspace_id"},
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"computer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_volume_encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_name": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				RequiredWith:  []string{"directory_id"},
				ConflictsWith: []string{"workspace_id"},
			},
			"user_volume_encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"volume_encryption_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"directory_id", "user_name"},
			},
			"workspace_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_volume_size_gib": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"running_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"running_mode_auto_stop_timeout_in_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_volume_size_gib": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WorkSpacesConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	var workspace *workspaces.Workspace

	if workspaceID, ok := d.GetOk("workspace_id"); ok {
		resp, err := conn.DescribeWorkspaces(&workspaces.DescribeWorkspacesInput{
			WorkspaceIds: aws.StringSlice([]string{workspaceID.(string)}),
		})
		if err != nil {
			return err
		}

		if len(resp.Workspaces) != 1 {
			return fmt.Errorf("expected 1 result for WorkSpace  %q, found %d", workspaceID, len(resp.Workspaces))
		}

		workspace = resp.Workspaces[0]

		if workspace == nil {
			return fmt.Errorf("no WorkSpace with ID %q found", workspaceID)
		}
	}

	if directoryID, ok := d.GetOk("directory_id"); ok {
		userName := d.Get("user_name").(string)
		resp, err := conn.DescribeWorkspaces(&workspaces.DescribeWorkspacesInput{
			DirectoryId: aws.String(directoryID.(string)),
			UserName:    aws.String(userName),
		})
		if err != nil {
			return err
		}

		if len(resp.Workspaces) != 1 {
			return fmt.Errorf("expected 1 result for %q WorkSpace in %q directory, found %d", userName, directoryID, len(resp.Workspaces))
		}

		workspace = resp.Workspaces[0]

		if workspace == nil {
			return fmt.Errorf("no %q WorkSpace in %q directory found", userName, directoryID)
		}
	}

	d.SetId(aws.StringValue(workspace.WorkspaceId))
	d.Set("bundle_id", workspace.BundleId)
	d.Set("directory_id", workspace.DirectoryId)
	d.Set("ip_address", workspace.IpAddress)
	d.Set("computer_name", workspace.ComputerName)
	d.Set("state", workspace.State)
	d.Set("root_volume_encryption_enabled", workspace.RootVolumeEncryptionEnabled)
	d.Set("user_name", workspace.UserName)
	d.Set("user_volume_encryption_enabled", workspace.UserVolumeEncryptionEnabled)
	d.Set("volume_encryption_key", workspace.VolumeEncryptionKey)
	if err := d.Set("workspace_properties", FlattenWorkspaceProperties(workspace.WorkspaceProperties)); err != nil {
		return fmt.Errorf("error setting workspace properties: %w", err)
	}

	tags, err := ListTags(conn, d.Id())
	if err != nil {
		return fmt.Errorf("error listing tags: %w", err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
