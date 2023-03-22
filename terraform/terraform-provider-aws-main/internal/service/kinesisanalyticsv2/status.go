package kinesisanalyticsv2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// statusApplication fetches the ApplicationDetail and its Status
func statusApplication(conn *kinesisanalyticsv2.KinesisAnalyticsV2, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		applicationDetail, err := FindApplicationDetailByName(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return applicationDetail, aws.StringValue(applicationDetail.ApplicationStatus), nil
	}
}

// statusSnapshotDetails fetches the SnapshotDetails and its Status
func statusSnapshotDetails(conn *kinesisanalyticsv2.KinesisAnalyticsV2, applicationName, snapshotName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapshotDetails, err := FindSnapshotDetailsByApplicationAndSnapshotNames(conn, applicationName, snapshotName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return snapshotDetails, aws.StringValue(snapshotDetails.SnapshotStatus), nil
	}
}
