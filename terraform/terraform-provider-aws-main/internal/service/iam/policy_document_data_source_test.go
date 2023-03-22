package iam_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccIAMPolicyDocumentDataSource_basic(t *testing.T) {
	// This really ought to be able to be a unit test rather than an
	// acceptance test, but just instantiating the AWS provider requires
	// some AWS API calls, and so this needs valid AWS credentials to work.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test", "json",
						testAccPolicyDocumentExpectedJSON(),
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_singleConditionValue(t *testing.T) {
	dataSourceName := "data.aws_iam_policy_document.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_singleConditionValue,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrEquivalentJSON(dataSourceName, "json", testAccPolicyDocumentConfig_SingleConditionValue_ExpectedJSON),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_source(t *testing.T) {
	// This really ought to be able to be a unit test rather than an
	// acceptance test, but just instantiating the AWS provider requires
	// some AWS API calls, and so this needs valid AWS credentials to work.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_deprecated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test_source", "json",
						testAccPolicyDocumentSourceExpectedJSON(),
					),
				),
			},
			{
				Config: testAccPolicyDocumentDataSourceConfig_blankDeprecated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test_source_blank", "json",
						testAccPolicyDocumentSourceBlankExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_sourceList(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_list,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test_source_list", "json",
						testAccPolicyDocumentSourceListExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_sourceConflicting(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_conflictingDeprecated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test_source_conflicting", "json",
						testAccPolicyDocumentSourceConflictingExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_sourceListConflicting(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccPolicyDocumentDataSourceConfig_listConflicting,
				ExpectError: regexp.MustCompile(`duplicate Sid (.*?)`),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_override(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_overrideDeprecated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test_override", "json",
						testAccPolicyDocumentOverrideExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_overrideList(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_overrideList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test_override_list", "json",
						testAccPolicyDocumentOverrideListExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_noStatementMerge(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_noStatementMergeDeprecated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.yak_politik", "json",
						testAccPolicyDocumentNoStatementMergeExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_noStatementOverride(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_noStatementOverrideDeprecated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.yak_politik", "json",
						testAccPolicyDocumentNoStatementOverrideExpectedJSON,
					),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_duplicateSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccPolicyDocumentDataSourceConfig_duplicateSid,
				ExpectError: regexp.MustCompile(`duplicate Sid`),
			},
			{
				Config: testAccPolicyDocumentDataSourceConfig_duplicateBlankSid,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test", "json",
						testAccPolicyDocumentDuplicateBlankSidExpectedJSON,
					),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/10777
func TestAccIAMPolicyDocumentDataSource_StatementPrincipalIdentifiers_stringAndSlice(t *testing.T) {
	dataSourceName := "data.aws_iam_policy_document.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_statementPrincipalIdentifiersStringAndSlice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "json", testAccPolicyDocumentExpectedJSONStatementPrincipalIdentifiersStringAndSlice()),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/10777
func TestAccIAMPolicyDocumentDataSource_StatementPrincipalIdentifiers_multiplePrincipals(t *testing.T) {
	dataSourceName := "data.aws_iam_policy_document.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartition(t, endpoints.AwsPartitionID) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_statementPrincipalIdentifiersMultiplePrincipals,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "json", testAccPolicyDocumentExpectedJSONStatementPrincipalIdentifiersMultiplePrincipals()),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_StatementPrincipalIdentifiers_multiplePrincipalsGov(t *testing.T) {
	dataSourceName := "data.aws_iam_policy_document.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartition(t, endpoints.AwsUsGovPartitionID) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDocumentDataSourceConfig_statementPrincipalIdentifiersMultiplePrincipals,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "json", testAccPolicyDocumentExpectedJSONStatementPrincipalIdentifiersMultiplePrincipalsGov()),
				),
			},
		},
	})
}

func TestAccIAMPolicyDocumentDataSource_version20081017(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccPolicyDocumentDataSourceConfig_version20081017ConversionCondition,
				ExpectError: regexp.MustCompile(`found \&\{ sequence in \(.+\), which is not supported in document version 2008-10-17`),
			},
			{
				Config:      testAccPolicyDocumentDataSourceConfig_version20081017ConversionNotPrincipals,
				ExpectError: regexp.MustCompile(`found \&\{ sequence in \(.+\), which is not supported in document version 2008-10-17`),
			},
			{
				Config:      testAccPolicyDocumentDataSourceConfig_version20081017ConversionNotRe,
				ExpectError: regexp.MustCompile(`found \&\{ sequence in \(.+\), which is not supported in document version 2008-10-17`),
			},
			{
				Config:      testAccPolicyDocumentDataSourceConfig_version20081017ConversionPrincipal,
				ExpectError: regexp.MustCompile(`found \&\{ sequence in \(.+\), which is not supported in document version 2008-10-17`),
			},
			{
				Config:      testAccPolicyDocumentDataSourceConfig_version20081017ConversionRe,
				ExpectError: regexp.MustCompile(`found \&\{ sequence in \(.+\), which is not supported in document version 2008-10-17`),
			},
			{
				Config: testAccPolicyDocumentDataSourceConfig_version20081017,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_iam_policy_document.test", "json", testAccPolicyDocumentVersion20081017ExpectedJSONDataSourceConfig),
				),
			},
		},
	})
}

var testAccPolicyDocumentDataSourceConfig_basic = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "test" {
  policy_id = "policy_id"

  statement {
    sid = "1"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]
    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::*",
    ]
  }

  statement {
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::foo",
    ]

    condition {
      test     = "StringLike"
      variable = "s3:prefix"
      values = [
        "home/",
        "",
        "home/&{aws:username}/",
      ]
    }

    not_principals {
      type        = "AWS"
      identifiers = ["arn:blahblah:example"]
    }
  }

  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::foo/home/&{aws:username}",
      "arn:${data.aws_partition.current.partition}:s3:::foo/home/&{aws:username}/*",
    ]

    principals {
      type        = "AWS"
      identifiers = ["arn:blahblah:example"]
    }
  }

  statement {
    effect        = "Deny"
    not_actions   = ["s3:*"]
    not_resources = ["arn:${data.aws_partition.current.partition}:s3:::*"]
  }

  # Normalization of wildcard principals

  statement {
    effect  = "Allow"
    actions = ["kinesis:*"]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }

  statement {
    effect  = "Allow"
    actions = ["firehose:*"]

    principals {
      type        = "*"
      identifiers = ["*"]
    }
  }
}
`

func testAccPolicyDocumentExpectedJSON() string {
	return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Id": "policy_id",
  "Statement": [
    {
      "Sid": "1",
      "Effect": "Allow",
      "Action": [
        "s3:ListAllMyBuckets",
        "s3:GetBucketLocation"
      ],
      "Resource": "arn:%[1]s:s3:::*"
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:%[1]s:s3:::foo",
      "NotPrincipal": {
        "AWS": "arn:blahblah:example"
      },
      "Condition": {
        "StringLike": {
          "s3:prefix": [
            "home/",
            "",
            "home/${aws:username}/"
          ]
        }
      }
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": [
        "arn:%[1]s:s3:::foo/home/${aws:username}/*",
        "arn:%[1]s:s3:::foo/home/${aws:username}"
      ],
      "Principal": {
        "AWS": "arn:blahblah:example"
      }
    },
    {
      "Sid": "",
      "Effect": "Deny",
      "NotAction": "s3:*",
      "NotResource": "arn:%[1]s:s3:::*"
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "kinesis:*",
      "Principal": {
        "AWS": "*"
      }
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "firehose:*",
      "Principal": "*"
    }
  ]
}`, acctest.Partition())
}

const testAccPolicyDocumentDataSourceConfig_singleConditionValue = `
data "aws_iam_policy_document" "test" {
  statement {
    effect = "Deny"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions = ["elasticfilesystem:*"]

    resources = ["*"]

    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values   = ["false"]
    }
  }
}
`

const testAccPolicyDocumentConfig_SingleConditionValue_ExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Deny",
      "Action": "elasticfilesystem:*",
      "Resource": "*",
      "Principal": {
        "AWS": "*"
      },
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      }
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_deprecated = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "test" {
  policy_id = "policy_id"

  statement {
    sid = "1"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]
    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::*",
    ]
  }

  statement {
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::foo",
    ]

    condition {
      test     = "StringLike"
      variable = "s3:prefix"
      values = [
        "home/",
        "home/&{aws:username}/",
      ]
    }

    not_principals {
      type        = "AWS"
      identifiers = ["arn:blahblah:example"]
    }
  }

  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::foo/home/&{aws:username}",
      "arn:${data.aws_partition.current.partition}:s3:::foo/home/&{aws:username}/*",
    ]

    principals {
      type = "AWS"
      identifiers = [
        "arn:blahblah:example",
        "arn:blahblahblah:example",
      ]
    }
  }

  statement {
    effect        = "Deny"
    not_actions   = ["s3:*"]
    not_resources = ["arn:${data.aws_partition.current.partition}:s3:::*"]
  }

  # Normalization of wildcard principals

  statement {
    effect  = "Allow"
    actions = ["kinesis:*"]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }

  statement {
    effect  = "Allow"
    actions = ["firehose:*"]

    principals {
      type        = "*"
      identifiers = ["*"]
    }
  }
}

data "aws_iam_policy_document" "test_source" {
  source_json = data.aws_iam_policy_document.test.json

  statement {
    sid       = "SourceJSONTest1"
    actions   = ["*"]
    resources = ["*"]
  }
}
`

func testAccPolicyDocumentSourceExpectedJSON() string {
	return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Id": "policy_id",
  "Statement": [
    {
      "Sid": "1",
      "Effect": "Allow",
      "Action": [
        "s3:ListAllMyBuckets",
        "s3:GetBucketLocation"
      ],
      "Resource": "arn:%[1]s:s3:::*"
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:%[1]s:s3:::foo",
      "NotPrincipal": {
        "AWS": "arn:blahblah:example"
      },
      "Condition": {
        "StringLike": {
          "s3:prefix": [
            "home/",
            "home/${aws:username}/"
          ]
        }
      }
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": [
        "arn:%[1]s:s3:::foo/home/${aws:username}/*",
        "arn:%[1]s:s3:::foo/home/${aws:username}"
      ],
      "Principal": {
        "AWS": [
          "arn:blahblahblah:example",
          "arn:blahblah:example"
        ]
      }
    },
    {
      "Sid": "",
      "Effect": "Deny",
      "NotAction": "s3:*",
      "NotResource": "arn:%[1]s:s3:::*"
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "kinesis:*",
      "Principal": {
        "AWS": "*"
      }
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "firehose:*",
      "Principal": "*"
    },
    {
      "Sid": "SourceJSONTest1",
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*"
    }
  ]
}`, acctest.Partition())
}

var testAccPolicyDocumentDataSourceConfig_list = `
data "aws_iam_policy_document" "policy_a" {
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["foo:ActionOne"]
  }

  statement {
    sid     = "validSidOne"
    effect  = "Allow"
    actions = ["bar:ActionOne"]
  }
}

data "aws_iam_policy_document" "policy_b" {
  statement {
    sid     = "validSidTwo"
    effect  = "Deny"
    actions = ["foo:ActionTwo"]
  }
}

data "aws_iam_policy_document" "policy_c" {
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["bar:ActionTwo"]
  }
}

data "aws_iam_policy_document" "test_source_list" {
  version = "2012-10-17"

  source_policy_documents = [
    data.aws_iam_policy_document.policy_a.json,
    data.aws_iam_policy_document.policy_b.json,
    data.aws_iam_policy_document.policy_c.json
  ]
}
`
var testAccPolicyDocumentSourceListExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "foo:ActionOne"
    },
    {
      "Sid": "validSidOne",
      "Effect": "Allow",
      "Action": "bar:ActionOne"
    },
    {
      "Sid": "validSidTwo",
      "Effect": "Deny",
      "Action": "foo:ActionTwo"
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "bar:ActionTwo"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_blankDeprecated = `
data "aws_iam_policy_document" "test_source_blank" {
  source_json = ""

  statement {
    sid       = "SourceJSONTest2"
    actions   = ["*"]
    resources = ["*"]
  }
}
`

var testAccPolicyDocumentSourceBlankExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "SourceJSONTest2",
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_conflictingDeprecated = `
data "aws_iam_policy_document" "test_source" {
  statement {
    sid       = "SourceJSONTestConflicting"
    actions   = ["iam:*"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "test_source_conflicting" {
  source_json = data.aws_iam_policy_document.test_source.json

  statement {
    sid       = "SourceJSONTestConflicting"
    actions   = ["*"]
    resources = ["*"]
  }
}
`

var testAccPolicyDocumentSourceConflictingExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "SourceJSONTestConflicting",
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_listConflicting = `
data "aws_iam_policy_document" "policy_a" {
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["foo:ActionOne"]
  }

  statement {
    sid     = "conflictSid"
    effect  = "Allow"
    actions = ["bar:ActionOne"]
  }
}

data "aws_iam_policy_document" "policy_b" {
  statement {
    sid     = "validSid"
    effect  = "Deny"
    actions = ["foo:ActionTwo"]
  }
}

data "aws_iam_policy_document" "policy_c" {
  statement {
    sid     = "conflictSid"
    effect  = "Allow"
    actions = ["bar:ActionTwo"]
  }
}

data "aws_iam_policy_document" "test_source_list_conflicting" {
  version = "2012-10-17"

  source_policy_documents = [
    data.aws_iam_policy_document.policy_a.json,
    data.aws_iam_policy_document.policy_b.json,
    data.aws_iam_policy_document.policy_c.json
  ]
}
`

var testAccPolicyDocumentDataSourceConfig_overrideDeprecated = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "override" {
  statement {
    sid = "SidToOverwrite"

    actions   = ["s3:*"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "test_override" {
  override_json = data.aws_iam_policy_document.override.json

  statement {
    actions   = ["ec2:*"]
    resources = ["*"]
  }

  statement {
    sid = "SidToOverwrite"

    actions = ["s3:*"]

    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::somebucket",
      "arn:${data.aws_partition.current.partition}:s3:::somebucket/*",
    ]
  }
}
`

var testAccPolicyDocumentOverrideExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "ec2:*",
      "Resource": "*"
    },
    {
      "Sid": "SidToOverwrite",
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": "*"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_overrideList = `
data "aws_iam_policy_document" "policy_a" {
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["foo:ActionOne"]
  }

  statement {
    sid     = "overrideSid"
    effect  = "Allow"
    actions = ["bar:ActionOne"]
  }
}

data "aws_iam_policy_document" "policy_b" {
  statement {
    sid     = "validSid"
    effect  = "Deny"
    actions = ["foo:ActionTwo"]
  }
}

data "aws_iam_policy_document" "policy_c" {
  statement {
    sid     = "overrideSid"
    effect  = "Deny"
    actions = ["bar:ActionOne"]
  }
}

data "aws_iam_policy_document" "test_override_list" {
  version = "2012-10-17"

  override_policy_documents = [
    data.aws_iam_policy_document.policy_a.json,
    data.aws_iam_policy_document.policy_b.json,
    data.aws_iam_policy_document.policy_c.json
  ]
}
`

var testAccPolicyDocumentOverrideListExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "foo:ActionOne"
    },
    {
      "Sid": "overrideSid",
      "Effect": "Deny",
      "Action": "bar:ActionOne"
    },
    {
      "Sid": "validSid",
      "Effect": "Deny",
      "Action": "foo:ActionTwo"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_noStatementMergeDeprecated = `
data "aws_iam_policy_document" "source" {
  statement {
    sid       = ""
    actions   = ["ec2:DescribeAccountAttributes"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "override" {
  statement {
    sid       = "OverridePlaceholder"
    actions   = ["s3:GetObject"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "yak_politik" {
  source_json   = data.aws_iam_policy_document.source.json
  override_json = data.aws_iam_policy_document.override.json
}
`

var testAccPolicyDocumentNoStatementMergeExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "ec2:DescribeAccountAttributes",
      "Resource": "*"
    },
    {
      "Sid": "OverridePlaceholder",
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "*"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_noStatementOverrideDeprecated = `
data "aws_iam_policy_document" "source" {
  statement {
    sid       = "OverridePlaceholder"
    actions   = ["ec2:DescribeAccountAttributes"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "override" {
  statement {
    sid       = "OverridePlaceholder"
    actions   = ["s3:GetObject"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "yak_politik" {
  source_json   = data.aws_iam_policy_document.source.json
  override_json = data.aws_iam_policy_document.override.json
}
`

var testAccPolicyDocumentNoStatementOverrideExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "OverridePlaceholder",
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "*"
    }
  ]
}`

var testAccPolicyDocumentDataSourceConfig_duplicateSid = `
data "aws_iam_policy_document" "test" {
  statement {
    sid       = "1"
    effect    = "Allow"
    actions   = ["ec2:DescribeAccountAttributes"]
    resources = ["*"]
  }

  statement {
    sid       = "1"
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["*"]
  }
}
`

var testAccPolicyDocumentDataSourceConfig_duplicateBlankSid = `
data "aws_iam_policy_document" "test" {
  statement {
    sid       = ""
    effect    = "Allow"
    actions   = ["ec2:DescribeAccountAttributes"]
    resources = ["*"]
  }

  statement {
    sid       = ""
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["*"]
  }
}
`

var testAccPolicyDocumentDuplicateBlankSidExpectedJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "ec2:DescribeAccountAttributes",
      "Resource": "*"
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "*"
    }
  ]
}`

const testAccPolicyDocumentDataSourceConfig_version20081017 = `
data "aws_iam_policy_document" "test" {
  version = "2008-10-17"

  statement {
    actions   = ["ec2:*"]
    resources = ["*"]
  }
}
`

const testAccPolicyDocumentVersion20081017ExpectedJSONDataSourceConfig = `{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "ec2:*",
      "Resource": "*"
    }
  ]
}`

const testAccPolicyDocumentDataSourceConfig_version20081017ConversionCondition = `
data "aws_iam_policy_document" "test" {
  version = "2008-10-17"

  statement {
    actions   = ["*"]
    resources = ["*"]

    condition {
      test = "StringLike"
      values = [
        "home/",
        "home/&{aws:username}/",
      ]
      variable = "s3:prefix"
    }
  }
}
`

const testAccPolicyDocumentDataSourceConfig_version20081017ConversionNotPrincipals = `
data "aws_iam_policy_document" "test" {
  version = "2008-10-17"

  statement {
    actions   = ["*"]
    resources = ["*"]

    not_principals {
      identifiers = ["&{aws:username}"]
      type        = "AWS"
    }
  }
}
`

const testAccPolicyDocumentDataSourceConfig_version20081017ConversionNotRe = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "test" {
  version = "2008-10-17"

  statement {
    actions       = ["*"]
    not_resources = ["arn:${data.aws_partition.current.partition}:s3:::foo/home/&{aws:username}", ]
  }
}
`

const testAccPolicyDocumentDataSourceConfig_version20081017ConversionPrincipal = `
data "aws_iam_policy_document" "test" {
  version = "2008-10-17"

  statement {
    actions   = ["*"]
    resources = ["*"]

    principals {
      identifiers = ["&{aws:username}"]
      type        = "AWS"
    }
  }
}
`

const testAccPolicyDocumentDataSourceConfig_version20081017ConversionRe = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "test" {
  version = "2008-10-17"

  statement {
    actions   = ["*"]
    resources = ["arn:${data.aws_partition.current.partition}:s3:::foo/home/&{aws:username}", ]
  }
}
`

var testAccPolicyDocumentDataSourceConfig_statementPrincipalIdentifiersStringAndSlice = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "test" {
  statement {
    actions   = ["*"]
    resources = ["*"]
    sid       = "StatementPrincipalIdentifiersStringAndSlice"

    principals {
      identifiers = ["arn:${data.aws_partition.current.partition}:iam::111111111111:root"]
      type        = "AWS"
    }

    principals {
      identifiers = ["arn:${data.aws_partition.current.partition}:iam::222222222222:root", "arn:${data.aws_partition.current.partition}:iam::333333333333:root"]
      type        = "AWS"
    }
  }
}
`

func testAccPolicyDocumentExpectedJSONStatementPrincipalIdentifiersStringAndSlice() string {
	return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "StatementPrincipalIdentifiersStringAndSlice",
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*",
      "Principal": {
        "AWS": [
          "arn:%[1]s:iam::111111111111:root",
          "arn:%[1]s:iam::333333333333:root",
          "arn:%[1]s:iam::222222222222:root"
        ]
      }
    }
  ]
}`, acctest.Partition())
}

var testAccPolicyDocumentDataSourceConfig_statementPrincipalIdentifiersMultiplePrincipals = `
data "aws_partition" "current" {}

data "aws_iam_policy_document" "test" {
  statement {
    actions   = ["*"]
    resources = ["*"]
    sid       = "StatementPrincipalIdentifiersStringAndSlice"

    principals {
      identifiers = [
        "arn:${data.aws_partition.current.partition}:iam::111111111111:root",
        "arn:${data.aws_partition.current.partition}:iam::222222222222:root",
      ]
      type = "AWS"
    }

    principals {
      identifiers = [
        "arn:${data.aws_partition.current.partition}:iam::333333333333:root",
      ]
      type = "AWS"
    }

    principals {
      identifiers = [
        "arn:${data.aws_partition.current.partition}:iam::444444444444:root",
      ]
      type = "AWS"
    }
  }
}
`

func testAccPolicyDocumentExpectedJSONStatementPrincipalIdentifiersMultiplePrincipals() string {
	return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "StatementPrincipalIdentifiersStringAndSlice",
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*",
      "Principal": {
        "AWS": [
          "arn:%[1]s:iam::333333333333:root",
          "arn:%[1]s:iam::444444444444:root",
          "arn:%[1]s:iam::222222222222:root",
          "arn:%[1]s:iam::111111111111:root"
        ]
      }
    }
  ]
}`, acctest.Partition())
}

func testAccPolicyDocumentExpectedJSONStatementPrincipalIdentifiersMultiplePrincipalsGov() string {
	return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "StatementPrincipalIdentifiersStringAndSlice",
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*",
      "Principal": {
        "AWS": [
          "arn:%[1]s:iam::333333333333:root",
          "arn:%[1]s:iam::222222222222:root",
          "arn:%[1]s:iam::111111111111:root",
          "arn:%[1]s:iam::444444444444:root"
        ]
      }
    }
  ]
}`, acctest.Partition())
}
