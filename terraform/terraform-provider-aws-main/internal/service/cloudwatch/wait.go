package cloudwatch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	MetricStreamDeleteTimeout = 2 * time.Minute
	MetricStreamReadyTimeout  = 1 * time.Minute

	StateRunning = "running"
	StateStopped = "stopped"
)

func WaitMetricStreamDeleted(ctx context.Context, conn *cloudwatch.CloudWatch, name string) (*cloudwatch.GetMetricStreamOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			StateRunning,
			StateStopped,
		},
		Target:  []string{},
		Refresh: StatusMetricStreamState(ctx, conn, name),
		Timeout: MetricStreamDeleteTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*cloudwatch.GetMetricStreamOutput); ok {
		return v, err
	}

	return nil, err
}

func WaitMetricStreamReady(ctx context.Context, conn *cloudwatch.CloudWatch, name string) (*cloudwatch.GetMetricStreamOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			StateStopped,
		},
		Target: []string{
			StateRunning,
		},
		Refresh: StatusMetricStreamState(ctx, conn, name),
		Timeout: MetricStreamReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*cloudwatch.GetMetricStreamOutput); ok {
		return v, err
	}

	return nil, err
}
