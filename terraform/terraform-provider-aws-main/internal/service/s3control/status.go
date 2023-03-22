package s3control

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// statusPublicAccessBlockConfigurationBlockPublicACLs fetches the PublicAccessBlockConfiguration and its BlockPublicAcls
func statusPublicAccessBlockConfigurationBlockPublicACLs(conn *s3control.S3Control, accountID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		publicAccessBlockConfiguration, err := findPublicAccessBlockConfiguration(conn, accountID)

		if err != nil {
			return nil, "false", err
		}

		if publicAccessBlockConfiguration == nil {
			return nil, "false", nil
		}

		return publicAccessBlockConfiguration, strconv.FormatBool(aws.BoolValue(publicAccessBlockConfiguration.BlockPublicAcls)), nil
	}
}

// statusPublicAccessBlockConfigurationBlockPublicPolicy fetches the PublicAccessBlockConfiguration and its BlockPublicPolicy
func statusPublicAccessBlockConfigurationBlockPublicPolicy(conn *s3control.S3Control, accountID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		publicAccessBlockConfiguration, err := findPublicAccessBlockConfiguration(conn, accountID)

		if err != nil {
			return nil, "false", err
		}

		if publicAccessBlockConfiguration == nil {
			return nil, "false", nil
		}

		return publicAccessBlockConfiguration, strconv.FormatBool(aws.BoolValue(publicAccessBlockConfiguration.BlockPublicPolicy)), nil
	}
}

// statusPublicAccessBlockConfigurationIgnorePublicACLs fetches the PublicAccessBlockConfiguration and its IgnorePublicAcls
func statusPublicAccessBlockConfigurationIgnorePublicACLs(conn *s3control.S3Control, accountID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		publicAccessBlockConfiguration, err := findPublicAccessBlockConfiguration(conn, accountID)

		if err != nil {
			return nil, "false", err
		}

		if publicAccessBlockConfiguration == nil {
			return nil, "false", nil
		}

		return publicAccessBlockConfiguration, strconv.FormatBool(aws.BoolValue(publicAccessBlockConfiguration.IgnorePublicAcls)), nil
	}
}

// statusPublicAccessBlockConfigurationRestrictPublicBuckets fetches the PublicAccessBlockConfiguration and its RestrictPublicBuckets
func statusPublicAccessBlockConfigurationRestrictPublicBuckets(conn *s3control.S3Control, accountID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		publicAccessBlockConfiguration, err := findPublicAccessBlockConfiguration(conn, accountID)

		if err != nil {
			return nil, "false", err
		}

		if publicAccessBlockConfiguration == nil {
			return nil, "false", nil
		}

		return publicAccessBlockConfiguration, strconv.FormatBool(aws.BoolValue(publicAccessBlockConfiguration.RestrictPublicBuckets)), nil
	}
}

func statusMultiRegionAccessPointRequest(conn *s3control.S3Control, accountID string, requestTokenARN string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findMultiRegionAccessPointOperationByAccountIDAndTokenARN(conn, accountID, requestTokenARN)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.RequestStatus), nil
	}
}
