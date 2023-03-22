package apigatewayv2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceAPIs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAPIsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tftags.TagsSchema(),
		},
	}
}

func dataSourceAPIsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	tagsToMatch := tftags.New(d.Get("tags").(map[string]interface{})).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	apis, err := FindAPIs(conn, &apigatewayv2.GetApisInput{})

	if err != nil {
		return fmt.Errorf("error reading API Gateway v2 APIs: %w", err)
	}

	var ids []*string

	for _, api := range apis {
		if v, ok := d.GetOk("name"); ok && v.(string) != aws.StringValue(api.Name) {
			continue
		}

		if v, ok := d.GetOk("protocol_type"); ok && v.(string) != aws.StringValue(api.ProtocolType) {
			continue
		}

		if len(tagsToMatch) > 0 && !KeyValueTags(api.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).ContainsAll(tagsToMatch) {
			continue
		}

		ids = append(ids, api.ApiId)
	}

	d.SetId(meta.(*conns.AWSClient).Region)

	if err := d.Set("ids", flex.FlattenStringSet(ids)); err != nil {
		return fmt.Errorf("error setting ids: %w", err)
	}

	return nil
}
