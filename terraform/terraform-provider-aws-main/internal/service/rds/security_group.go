package rds

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Managed by Terraform",
			},

			"ingress": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"security_group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"security_group_owner_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceSecurityGroupIngressHash,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,

		DeprecationMessage: `With the retirement of EC2-Classic the aws_db_security_group resource has been deprecated and will be removed in a future version.`,
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	var err error
	var errs []error

	opts := rds.CreateDBSecurityGroupInput{
		DBSecurityGroupName:        aws.String(d.Get("name").(string)),
		DBSecurityGroupDescription: aws.String(d.Get("description").(string)),
		Tags:                       Tags(tags.IgnoreAWS()),
	}

	log.Printf("[DEBUG] DB Security Group create configuration: %#v", opts)
	_, err = conn.CreateDBSecurityGroup(&opts)
	if err != nil {
		return fmt.Errorf("Error creating DB Security Group: %s", err)
	}

	d.SetId(d.Get("name").(string))

	log.Printf("[INFO] DB Security Group ID: %s", d.Id())

	sg, err := resourceSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	ingresses := d.Get("ingress").(*schema.Set)
	for _, ing := range ingresses.List() {
		err := resourceSecurityGroupAuthorizeRule(ing, *sg.DBSecurityGroupName, conn)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}

	log.Println(
		"[INFO] Waiting for Ingress Authorizations to be authorized")

	stateConf := &resource.StateChangeConf{
		Pending: []string{"authorizing"},
		Target:  []string{"authorized"},
		Refresh: resourceSecurityGroupStateRefreshFunc(d, meta),
		Timeout: 10 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	sg, err := resourceSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	d.Set("name", sg.DBSecurityGroupName)
	d.Set("description", sg.DBSecurityGroupDescription)

	// Create an empty schema.Set to hold all ingress rules
	rules := &schema.Set{
		F: resourceSecurityGroupIngressHash,
	}

	for _, v := range sg.IPRanges {
		rule := map[string]interface{}{"cidr": *v.CIDRIP}
		rules.Add(rule)
	}

	for _, g := range sg.EC2SecurityGroups {
		rule := map[string]interface{}{}
		if g.EC2SecurityGroupId != nil {
			rule["security_group_id"] = aws.StringValue(g.EC2SecurityGroupId)
		}
		if g.EC2SecurityGroupName != nil {
			rule["security_group_name"] = aws.StringValue(g.EC2SecurityGroupName)
		}
		if g.EC2SecurityGroupOwnerId != nil {
			rule["security_group_owner_id"] = aws.StringValue(g.EC2SecurityGroupOwnerId)
		}
		rules.Add(rule)
	}

	d.Set("ingress", rules)

	conn := meta.(*conns.AWSClient).RDSConn

	arn := aws.StringValue(sg.DBSecurityGroupArn)
	d.Set("arn", arn)

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return fmt.Errorf("error listing tags for RDS DB Security Group (%s): %s", d.Get("arn").(string), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating RDS DB Security Group (%s) tags: %s", d.Get("arn").(string), err)
		}
	}

	if d.HasChange("ingress") {
		sg, err := resourceSecurityGroupRetrieve(d, meta)
		if err != nil {
			return err
		}

		oi, ni := d.GetChange("ingress")
		if oi == nil {
			oi = new(schema.Set)
		}
		if ni == nil {
			ni = new(schema.Set)
		}

		ois := oi.(*schema.Set)
		nis := ni.(*schema.Set)
		removeIngress := ois.Difference(nis).List()
		newIngress := nis.Difference(ois).List()

		// DELETE old Ingress rules
		for _, ing := range removeIngress {
			err := resourceSecurityGroupRevokeRule(ing, *sg.DBSecurityGroupName, conn)
			if err != nil {
				return err
			}
		}

		// ADD new/updated Ingress rules
		for _, ing := range newIngress {
			err := resourceSecurityGroupAuthorizeRule(ing, *sg.DBSecurityGroupName, conn)
			if err != nil {
				return err
			}
		}
	}

	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	log.Printf("[DEBUG] DB Security Group destroy: %v", d.Id())

	opts := rds.DeleteDBSecurityGroupInput{DBSecurityGroupName: aws.String(d.Id())}

	log.Printf("[DEBUG] DB Security Group destroy configuration: %v", opts)
	_, err := conn.DeleteDBSecurityGroup(&opts)

	if err != nil {
		if tfawserr.ErrCodeEquals(err, "InvalidDBSecurityGroup.NotFound") {
			return nil
		}
		return err
	}

	return nil
}

func resourceSecurityGroupRetrieve(d *schema.ResourceData, meta interface{}) (*rds.DBSecurityGroup, error) {
	conn := meta.(*conns.AWSClient).RDSConn

	opts := rds.DescribeDBSecurityGroupsInput{
		DBSecurityGroupName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] DB Security Group describe configuration: %#v", opts)

	resp, err := conn.DescribeDBSecurityGroups(&opts)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving DB Security Groups: %s", err)
	}

	if len(resp.DBSecurityGroups) != 1 ||
		aws.StringValue(resp.DBSecurityGroups[0].DBSecurityGroupName) != d.Id() {
		return nil, fmt.Errorf("Unable to find DB Security Group: %#v", resp.DBSecurityGroups)
	}

	return resp.DBSecurityGroups[0], nil
}

// Authorizes the ingress rule on the db security group
func resourceSecurityGroupAuthorizeRule(ingress interface{}, dbSecurityGroupName string, conn *rds.RDS) error {
	ing := ingress.(map[string]interface{})

	opts := rds.AuthorizeDBSecurityGroupIngressInput{
		DBSecurityGroupName: aws.String(dbSecurityGroupName),
	}

	if attr, ok := ing["cidr"]; ok && attr != "" {
		opts.CIDRIP = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_name"]; ok && attr != "" {
		opts.EC2SecurityGroupName = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_id"]; ok && attr != "" {
		opts.EC2SecurityGroupId = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_owner_id"]; ok && attr != "" {
		opts.EC2SecurityGroupOwnerId = aws.String(attr.(string))
	}

	log.Printf("[DEBUG] Authorize ingress rule configuration: %#v", opts)

	_, err := conn.AuthorizeDBSecurityGroupIngress(&opts)

	if err != nil {
		return fmt.Errorf("Error authorizing security group ingress: %s", err)
	}

	return nil
}

// Revokes the ingress rule on the db security group
func resourceSecurityGroupRevokeRule(ingress interface{}, dbSecurityGroupName string, conn *rds.RDS) error {
	ing := ingress.(map[string]interface{})

	opts := rds.RevokeDBSecurityGroupIngressInput{
		DBSecurityGroupName: aws.String(dbSecurityGroupName),
	}

	if attr, ok := ing["cidr"]; ok && attr != "" {
		opts.CIDRIP = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_name"]; ok && attr != "" {
		opts.EC2SecurityGroupName = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_id"]; ok && attr != "" {
		opts.EC2SecurityGroupId = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_owner_id"]; ok && attr != "" {
		opts.EC2SecurityGroupOwnerId = aws.String(attr.(string))
	}

	log.Printf("[DEBUG] Revoking ingress rule configuration: %#v", opts)

	_, err := conn.RevokeDBSecurityGroupIngress(&opts)

	if err != nil {
		return fmt.Errorf("Error revoking security group ingress: %s", err)
	}

	return nil
}

func resourceSecurityGroupIngressHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["cidr"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_owner_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return create.StringHashcode(buf.String())
}

func resourceSecurityGroupStateRefreshFunc(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := resourceSecurityGroupRetrieve(d, meta)

		if err != nil {
			log.Printf("Error on retrieving DB Security Group when waiting: %s", err)
			return nil, "", err
		}

		statuses := make([]string, 0, len(v.EC2SecurityGroups)+len(v.IPRanges))
		for _, ec2g := range v.EC2SecurityGroups {
			statuses = append(statuses, *ec2g.Status)
		}
		for _, ips := range v.IPRanges {
			statuses = append(statuses, *ips.Status)
		}

		for _, stat := range statuses {
			// Not done
			if stat != "authorized" {
				return nil, "authorizing", nil
			}
		}

		return v, "authorized", nil
	}
}
