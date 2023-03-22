package kinesis_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/kinesis"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccKinesisStreamConsumerDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_kinesis_stream_consumer.test"
	resourceName := "aws_kinesis_stream_consumer.test"
	streamName := "aws_kinesis_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kinesis.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamConsumerDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "stream_arn", streamName, "arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
				),
			},
		},
	})
}

func TestAccKinesisStreamConsumerDataSource_name(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_kinesis_stream_consumer.test"
	resourceName := "aws_kinesis_stream_consumer.test"
	streamName := "aws_kinesis_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kinesis.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamConsumerDataSourceConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "stream_arn", streamName, "arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
				),
			},
		},
	})
}

func TestAccKinesisStreamConsumerDataSource_arn(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_kinesis_stream_consumer.test"
	resourceName := "aws_kinesis_stream_consumer.test"
	streamName := "aws_kinesis_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kinesis.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamConsumerDataSourceConfig_arn(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "stream_arn", streamName, "arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
				),
			},
		},
	})
}

func testAccStreamConsumerBaseDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_kinesis_stream" "test" {
  name        = %q
  shard_count = 2
}
`, rName)
}

func testAccStreamConsumerDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccStreamConsumerBaseDataSourceConfig(rName),
		fmt.Sprintf(`
data "aws_kinesis_stream_consumer" "test" {
  stream_arn = aws_kinesis_stream_consumer.test.stream_arn
}

resource "aws_kinesis_stream_consumer" "test" {
  name       = %q
  stream_arn = aws_kinesis_stream.test.arn
}
`, rName))
}

func testAccStreamConsumerDataSourceConfig_name(rName string) string {
	return acctest.ConfigCompose(
		testAccStreamConsumerBaseDataSourceConfig(rName),
		fmt.Sprintf(`
data "aws_kinesis_stream_consumer" "test" {
  name       = aws_kinesis_stream_consumer.test.name
  stream_arn = aws_kinesis_stream_consumer.test.stream_arn
}

resource "aws_kinesis_stream_consumer" "test" {
  name       = %q
  stream_arn = aws_kinesis_stream.test.arn
}
`, rName))
}

func testAccStreamConsumerDataSourceConfig_arn(rName string) string {
	return acctest.ConfigCompose(
		testAccStreamConsumerBaseDataSourceConfig(rName),
		fmt.Sprintf(`
data "aws_kinesis_stream_consumer" "test" {
  arn        = aws_kinesis_stream_consumer.test.arn
  stream_arn = aws_kinesis_stream_consumer.test.stream_arn
}

resource "aws_kinesis_stream_consumer" "test" {
  name       = %q
  stream_arn = aws_kinesis_stream.test.arn
}
`, rName))
}
