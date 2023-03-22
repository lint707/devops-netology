package kinesisanalyticsv2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// FindApplicationDetailByName returns the application corresponding to the specified name.
// Returns NotFoundError if no application is found.
func FindApplicationDetailByName(conn *kinesisanalyticsv2.KinesisAnalyticsV2, name string) (*kinesisanalyticsv2.ApplicationDetail, error) {
	input := &kinesisanalyticsv2.DescribeApplicationInput{
		ApplicationName: aws.String(name),
	}

	return FindApplicationDetail(conn, input)
}

// FindApplicationDetail returns the application details corresponding to the specified input.
// Returns NotFoundError if no application is found.
func FindApplicationDetail(conn *kinesisanalyticsv2.KinesisAnalyticsV2, input *kinesisanalyticsv2.DescribeApplicationInput) (*kinesisanalyticsv2.ApplicationDetail, error) {
	output, err := conn.DescribeApplication(input)

	if tfawserr.ErrCodeEquals(err, kinesisanalyticsv2.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.ApplicationDetail == nil {
		return nil, &resource.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.ApplicationDetail, nil
}

// FindSnapshotDetailsByApplicationAndSnapshotNames returns the application snapshot details corresponding to the specified application and snapshot names.
// Returns NotFoundError if no application snapshot is found.
func FindSnapshotDetailsByApplicationAndSnapshotNames(conn *kinesisanalyticsv2.KinesisAnalyticsV2, applicationName, snapshotName string) (*kinesisanalyticsv2.SnapshotDetails, error) {
	input := &kinesisanalyticsv2.DescribeApplicationSnapshotInput{
		ApplicationName: aws.String(applicationName),
		SnapshotName:    aws.String(snapshotName),
	}

	return FindSnapshotDetails(conn, input)
}

// FindSnapshotDetails returns the application snapshot details corresponding to the specified input.
// Returns NotFoundError if no application snapshot is found.
func FindSnapshotDetails(conn *kinesisanalyticsv2.KinesisAnalyticsV2, input *kinesisanalyticsv2.DescribeApplicationSnapshotInput) (*kinesisanalyticsv2.SnapshotDetails, error) {
	output, err := conn.DescribeApplicationSnapshot(input)

	if tfawserr.ErrCodeEquals(err, kinesisanalyticsv2.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if tfawserr.ErrMessageContains(err, kinesisanalyticsv2.ErrCodeInvalidArgumentException, "does not exist") {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.SnapshotDetails == nil {
		return nil, &resource.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.SnapshotDetails, nil
}
