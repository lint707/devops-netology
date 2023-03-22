package kendra_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfkendra "github.com/hashicorp/terraform-provider-aws/internal/service/kendra"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccPreCheck(t *testing.T) {
	acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)

	conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn

	input := &kendra.ListIndicesInput{}

	_, err := conn.ListIndices(context.Background(), input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func TestAccKendraIndex_basic(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	description := "basic"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.0.query_capacity_units", "0"),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.0.storage_capacity_units", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "document_metadata_configuration_updates.#", "13"),
					resource.TestCheckResourceAttr(resourceName, "edition", string(types.IndexEditionEnterpriseEdition)),
					resource.TestCheckResourceAttr(resourceName, "index_statistics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index_statistics.0.faq_statistics.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "index_statistics.0.faq_statistics.0.indexed_question_answers_count"),
					resource.TestCheckResourceAttr(resourceName, "index_statistics.0.text_document_statistics.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "index_statistics.0.text_document_statistics.0.indexed_text_bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "index_statistics.0.text_document_statistics.0.indexed_text_documents_count"),
					resource.TestCheckResourceAttr(resourceName, "name", rName3),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.access_cw", "arn"),
					resource.TestCheckResourceAttr(resourceName, "status", string(types.IndexStatusActive)),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "user_context_policy", "ATTRIBUTE_FILTER"),
					resource.TestCheckResourceAttr(resourceName, "user_group_resolution_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKendraIndex_serverSideEncryption(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_serverSideEncryption(rName, rName2, rName3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "server_side_encryption_configuration.0.kms_key_id", "data.aws_kms_key.this", "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKendraIndex_updateCapacityUnits(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	originalQueryCapacityUnits := 2
	updatedQueryCapacityUnits := 3
	originalStorageCapacityUnits := 1
	updatedStorageCapacityUnits := 2
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_capacityUnits(rName, rName2, rName3, originalQueryCapacityUnits, originalStorageCapacityUnits),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.0.query_capacity_units", strconv.Itoa(originalQueryCapacityUnits)),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.0.storage_capacity_units", strconv.Itoa(originalStorageCapacityUnits)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_capacityUnits(rName, rName2, rName3, updatedQueryCapacityUnits, updatedStorageCapacityUnits),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.0.query_capacity_units", strconv.Itoa(updatedQueryCapacityUnits)),
					resource.TestCheckResourceAttr(resourceName, "capacity_units.0.storage_capacity_units", strconv.Itoa(updatedStorageCapacityUnits)),
				),
			},
		},
	})
}
func TestAccKendraIndex_updateDescription(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	originalDescription := "original description"
	updatedDescription := "updated description"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, originalDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "description", originalDescription),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccKendraIndex_updateName(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName4 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	description := "description"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", rName3),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName4, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", rName4),
				),
			},
		},
	})
}

func TestAccKendraIndex_updateUserTokenJSON(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	originalGroupAttributeField := "groups"
	originalUserNameAttributeField := "username"
	updatedGroupAttributeField := "groupings"
	updatedUserNameAttributeField := "usernames"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_userTokenJSON(rName, rName2, rName3, originalGroupAttributeField, originalUserNameAttributeField),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.0.group_attribute_field", originalGroupAttributeField),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.0.user_name_attribute_field", originalUserNameAttributeField),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_userTokenJSON(rName, rName2, rName3, updatedGroupAttributeField, originalUserNameAttributeField),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.0.group_attribute_field", updatedGroupAttributeField),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.0.user_name_attribute_field", originalUserNameAttributeField),
				),
			},
			{
				Config: testAccIndexConfig_userTokenJSON(rName, rName2, rName3, updatedGroupAttributeField, updatedUserNameAttributeField),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.0.group_attribute_field", updatedGroupAttributeField),
					resource.TestCheckResourceAttr(resourceName, "user_token_configurations.0.json_token_type_configuration.0.user_name_attribute_field", updatedUserNameAttributeField),
				),
			},
		},
	})
}

func TestAccKendraIndex_updateTags(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	description := "description"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_tags(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2a"),
				),
			},
			{
				Config: testAccIndexConfig_tagsUpdated(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2b"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "Value3"),
				),
			},
		},
	})
}

func TestAccKendraIndex_updateRoleARN(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	description := "description"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.access_cw", "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_secretsManagerRole(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.access_sm", "arn"),
				),
			},
		},
	})
}

func TestAccKendraIndex_addDocumentMetadataConfigurationUpdates(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	resourceName := "aws_kendra_index.test"
	authorsFacetable := false
	longValDisplayable := true
	stringListValSearchable := true
	dateValSortable := false
	stringValImportance := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_documentMetadataConfigurationUpdatesBase(rName, rName2, rName3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "document_metadata_configuration_updates.#", "13"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_authors",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_category",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_created_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_data_source_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_document_title",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "2",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_excerpt_page_number",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "2",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_faq_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_file_type",
						"type":                   string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_language_code",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_last_updated_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_source_uri",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_version",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_view_count",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_documentMetadataConfigurationUpdatesAddNewMetadata(rName, rName2, rName3, authorsFacetable, longValDisplayable, stringListValSearchable, dateValSortable, stringValImportance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "document_metadata_configuration_updates.#", "17"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_authors",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     strconv.FormatBool(authorsFacetable),
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_category",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_created_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_data_source_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_document_title",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "2",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_excerpt_page_number",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "2",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_faq_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_file_type",
						"type":                   string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_language_code",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_last_updated_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_source_uri",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_version",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_view_count",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "example-string-value",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              strconv.Itoa(stringValImportance),
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "true",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-long-value",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   strconv.FormatBool(longValDisplayable),
						"search.0.facetable":     "true",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-string-list-value",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "true",
						"search.0.facetable":     "true",
						"search.0.searchable":    strconv.FormatBool(stringListValSearchable),
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-date-value",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "true",
						"search.0.facetable":     "true",
						"search.0.searchable":    "false",
						"search.0.sortable":      strconv.FormatBool(dateValSortable),
					}),
				),
			},
		},
	})
}

func TestAccKendraIndex_inplaceUpdateDocumentMetadataConfigurationUpdates(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	resourceName := "aws_kendra_index.test"
	originalAuthorsFacetable := true
	originalLongValDisplayable := true
	originalStringListValSearchable := true
	originalDateValSortable := false
	originalStringValImportance := 1

	updatedAuthorsFacetable := false
	updatedLongValDisplayable := false
	updatedStringListValSearchable := false
	updatedDateValSortable := true
	updatedStringValImportance := 2

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_documentMetadataConfigurationUpdatesAddNewMetadata(rName, rName2, rName3, originalAuthorsFacetable, originalLongValDisplayable, originalStringListValSearchable, originalDateValSortable, originalStringValImportance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "document_metadata_configuration_updates.#", "17"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_authors",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     strconv.FormatBool(originalAuthorsFacetable),
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_category",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_created_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_data_source_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_document_title",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "2",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_excerpt_page_number",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "2",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_faq_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_file_type",
						"type":                   string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_language_code",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_last_updated_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_source_uri",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_version",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_view_count",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "example-string-value",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              strconv.Itoa(originalStringValImportance),
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "true",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-long-value",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   strconv.FormatBool(originalLongValDisplayable),
						"search.0.facetable":     "true",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-string-list-value",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "true",
						"search.0.facetable":     "true",
						"search.0.searchable":    strconv.FormatBool(originalStringListValSearchable),
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-date-value",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "true",
						"search.0.facetable":     "true",
						"search.0.searchable":    "false",
						"search.0.sortable":      strconv.FormatBool(originalDateValSortable),
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIndexConfig_documentMetadataConfigurationUpdatesAddNewMetadata(rName, rName2, rName3, updatedAuthorsFacetable, updatedLongValDisplayable, updatedStringListValSearchable, updatedDateValSortable, updatedStringValImportance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "document_metadata_configuration_updates.#", "17"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_authors",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     strconv.FormatBool(updatedAuthorsFacetable),
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_category",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_created_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_data_source_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_document_title",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "2",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_excerpt_page_number",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "2",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_faq_id",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_file_type",
						"type":                   string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_language_code",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_last_updated_at",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_source_uri",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "_version",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              "1",
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "false",
						"search.0.facetable":                  "false",
						"search.0.searchable":                 "false",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "_view_count",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "false",
						"search.0.facetable":     "false",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                                "example-string-value",
						"type":                                string(types.DocumentAttributeValueTypeStringValue),
						"relevance.#":                         "1",
						"relevance.0.importance":              strconv.Itoa(updatedStringValImportance),
						"relevance.0.values_importance_map.%": "0",
						"search.#":                            "1",
						"search.0.displayable":                "true",
						"search.0.facetable":                  "true",
						"search.0.searchable":                 "true",
						"search.0.sortable":                   "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-long-value",
						"type":                   string(types.DocumentAttributeValueTypeLongValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   strconv.FormatBool(updatedLongValDisplayable),
						"search.0.facetable":     "true",
						"search.0.searchable":    "false",
						"search.0.sortable":      "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-string-list-value",
						"type":                   string(types.DocumentAttributeValueTypeStringListValue),
						"relevance.#":            "1",
						"relevance.0.importance": "1",
						"search.#":               "1",
						"search.0.displayable":   "true",
						"search.0.facetable":     "true",
						"search.0.searchable":    strconv.FormatBool(updatedStringListValSearchable),
						"search.0.sortable":      "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "document_metadata_configuration_updates.*", map[string]string{
						"name":                   "example-date-value",
						"type":                   string(types.DocumentAttributeValueTypeDateValue),
						"relevance.#":            "1",
						"relevance.0.freshness":  "false",
						"relevance.0.importance": "1",
						"relevance.0.duration":   "25920000s",
						"relevance.0.rank_order": "ASCENDING",
						"search.#":               "1",
						"search.0.displayable":   "true",
						"search.0.facetable":     "true",
						"search.0.searchable":    "false",
						"search.0.sortable":      strconv.FormatBool(updatedDateValSortable),
					}),
				),
			},
		},
	})
}

func TestAccKendraIndex_disappears(t *testing.T) {
	var index kendra.DescribeIndexOutput

	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName3 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	description := "disappears"
	resourceName := "aws_kendra_index.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.KendraEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig_basic(rName, rName2, rName3, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExists(resourceName, &index),
					acctest.CheckResourceDisappears(acctest.Provider, tfkendra.ResourceIndex(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckIndexDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_kendra_index" {
			continue
		}

		input := &kendra.DescribeIndexInput{
			Id: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeIndex(context.Background(), input)

		if err == nil {
			if aws.ToString(resp.Id) == rs.Primary.ID {
				return fmt.Errorf("Index '%s' was not deleted properly", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIndexExists(name string, index *kendra.DescribeIndexOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn
		input := &kendra.DescribeIndexInput{
			Id: aws.String(rs.Primary.ID),
		}
		resp, err := conn.DescribeIndex(context.Background(), input)

		if err != nil {
			return err
		}

		*index = *resp

		return nil
	}
}

func testAccIndexConfigBase(rName, rName2 string) string {
	// Kendra IAM policies: https://docs.aws.amazon.com/kendra/latest/dg/iam-roles.html
	return fmt.Sprintf(`
data "aws_region" "current" {}
data "aws_caller_identity" "current" {}
data "aws_partition" "current" {}
data "aws_kms_key" "this" {
  key_id = "alias/aws/kendra"
}
data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["kendra.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "access_cw" {
  name               = %[1]q
  assume_role_policy = data.aws_iam_policy_document.test.json

  inline_policy {
    name = "access_cw"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = ["cloudwatch:PutMetricData"]
          Effect   = "Allow"
          Resource = "*"
          Condition = {
            StringEquals = {
              "cloudwatch:namespace" = "Kendra"
            }
          }
        },
        {
          Action   = ["logs:DescribeLogGroups"]
          Effect   = "Allow"
          Resource = "*"
        },
        {
          Action   = ["logs:CreateLogGroup"]
          Effect   = "Allow"
          Resource = "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/kendra/*"
        },
        {
          Action = [
            "logs:DescribeLogStreams",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ]
          Effect   = "Allow"
          Resource = "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/kendra/*:log-stream:*"
        },
      ]
    })
  }
}

resource "aws_iam_role" "access_sm" {
  name               = %[2]q
  assume_role_policy = data.aws_iam_policy_document.test.json

  inline_policy {
    name = "access_sm"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = ["cloudwatch:PutMetricData"]
          Effect   = "Allow"
          Resource = "*"
          Condition = {
            StringEquals = {
              "cloudwatch:namespace" = "Kendra"
            }
          }
        },
        {
          Action   = ["logs:DescribeLogGroups"]
          Effect   = "Allow"
          Resource = "*"
        },
        {
          Action   = ["logs:CreateLogGroup"]
          Effect   = "Allow"
          Resource = "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/kendra/*"
        },
        {
          Action = [
            "logs:DescribeLogStreams",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ]
          Effect   = "Allow"
          Resource = "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/kendra/*:log-stream:*"
        },
        {
          Action   = ["secretsmanager:GetSecretValue"]
          Effect   = "Allow"
          Resource = "arn:${data.aws_partition.current.partition}:secretsmanager:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:secret:example"
        },
        {
          Action   = ["kms:Decrypt"]
          Effect   = "Allow"
          Resource = "arn:${data.aws_partition.current.partition}:kms:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:key/example"
          Condition = {
            StringLike = {
              "kms:ViaService" = ["secretsmanager.*.amazonaws.com"]
            }
          }
        }
      ]
    })
  }
}
`, rName, rName2)
}

func testAccIndexConfig_basic(rName, rName2, rName3, description string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name        = %[1]q
  description = %[2]q
  role_arn    = aws_iam_role.access_cw.arn

  tags = {
    "Key1" = "Value1"
  }
}
`, rName3, description))
}

func testAccIndexConfig_capacityUnits(rName, rName2, rName3 string, queryCapacityUnits, storageCapacityUnits int) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.access_cw.arn

  capacity_units {
    query_capacity_units   = %[2]d
    storage_capacity_units = %[3]d
  }

  tags = {
    "Key1" = "Value1"
  }
}
`, rName3, queryCapacityUnits, storageCapacityUnits))
}

func testAccIndexConfig_secretsManagerRole(rName, rName2, rName3, description string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name        = %[1]q
  description = %[2]q
  role_arn    = aws_iam_role.access_sm.arn

  tags = {
    "Key1" = "Value1"
  }
}
`, rName3, description))
}

func testAccIndexConfig_serverSideEncryption(rName, rName2, rName3 string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.access_cw.arn

  server_side_encryption_configuration {
    kms_key_id = data.aws_kms_key.this.arn
  }
}
`, rName3))
}

func testAccIndexConfig_userTokenJSON(rName, rName2, rName3, groupAttributeField, userNameAttributeField string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.access_cw.arn

  user_token_configurations {
    json_token_type_configuration {
      group_attribute_field     = %[2]q
      user_name_attribute_field = %[3]q
    }
  }
}
`, rName3, groupAttributeField, userNameAttributeField))
}

func testAccIndexConfig_tags(rName, rName2, rName3, description string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name        = %[1]q
  description = %[2]q
  role_arn    = aws_iam_role.access_cw.arn

  tags = {
    "Key1" = "Value1"
    "Key2" = "Value2a",
  }
}
`, rName3, description))
}

func testAccIndexConfig_tagsUpdated(rName, rName2, rName3, description string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name        = %[1]q
  description = %[2]q
  role_arn    = aws_iam_role.access_cw.arn

  tags = {
    "Key1" = "Value1",
    "Key2" = "Value2b"
    "Key3" = "Value3"
  }
}
`, rName3, description))
}

func testAccIndexConfig_documentMetadataConfigurationUpdatesBase(rName, rName2, rName3 string) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.access_cw.arn
  document_metadata_configuration_updates {
    name = "_authors"
    type = "STRING_LIST_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = false
    }
    relevance {
      importance = 1
    }
  }

  document_metadata_configuration_updates {
    name = "_category"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_created_at"
    type = "DATE_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      freshness  = false
      importance = 1
      duration   = "25920000s"
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "_data_source_id"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_document_title"
    type = "STRING_VALUE"
    search {
      displayable = true
      facetable   = false
      searchable  = true
      sortable    = true
    }
    relevance {
      importance            = 2
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_excerpt_page_number"
    type = "LONG_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = false
    }
    relevance {
      importance = 2
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "_faq_id"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_file_type"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_language_code"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_last_updated_at"
    type = "DATE_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      freshness  = false
      importance = 1
      duration   = "25920000s"
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "_source_uri"
    type = "STRING_VALUE"
    search {
      displayable = true
      facetable   = false
      searchable  = false
      sortable    = false
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_version"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_view_count"
    type = "LONG_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance = 1
      rank_order = "ASCENDING"
    }
  }
}
`, rName3))
}

func testAccIndexConfig_documentMetadataConfigurationUpdatesAddNewMetadata(rName, rName2, rName3 string, authorsFacetable, longValDisplayable, stringListValSearchable, dateValSortable bool, stringValImportance int) string {
	return acctest.ConfigCompose(
		testAccIndexConfigBase(rName, rName2),
		fmt.Sprintf(`
resource "aws_kendra_index" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.access_cw.arn
  document_metadata_configuration_updates {
    name = "_authors"
    type = "STRING_LIST_VALUE"
    search {
      displayable = false
      facetable   = %[2]t
      searchable  = false
      sortable    = false
    }
    relevance {
      importance = 1
    }
  }

  document_metadata_configuration_updates {
    name = "_category"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_created_at"
    type = "DATE_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      freshness  = false
      importance = 1
      duration   = "25920000s"
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "_data_source_id"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_document_title"
    type = "STRING_VALUE"
    search {
      displayable = true
      facetable   = false
      searchable  = true
      sortable    = true
    }
    relevance {
      importance            = 2
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_excerpt_page_number"
    type = "LONG_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = false
    }
    relevance {
      importance = 2
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "_faq_id"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_file_type"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_language_code"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_last_updated_at"
    type = "DATE_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      freshness  = false
      importance = 1
      duration   = "25920000s"
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "_source_uri"
    type = "STRING_VALUE"
    search {
      displayable = true
      facetable   = false
      searchable  = false
      sortable    = false
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_version"
    type = "STRING_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance            = 1
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "_view_count"
    type = "LONG_VALUE"
    search {
      displayable = false
      facetable   = false
      searchable  = false
      sortable    = true
    }
    relevance {
      importance = 1
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "example-string-value"
    type = "STRING_VALUE"
    search {
      displayable = true
      facetable   = true
      searchable  = true
      sortable    = true
    }
    relevance {
      importance            = %[6]d
      values_importance_map = {}
    }
  }

  document_metadata_configuration_updates {
    name = "example-long-value"
    type = "LONG_VALUE"
    search {
      displayable = %[3]t
      facetable   = true
      searchable  = false
      sortable    = true
    }
    relevance {
      importance = 1
      rank_order = "ASCENDING"
    }
  }

  document_metadata_configuration_updates {
    name = "example-string-list-value"
    type = "STRING_LIST_VALUE"
    search {
      displayable = true
      facetable   = true
      searchable  = %[4]t
      sortable    = false
    }
    relevance {
      importance = 1
    }
  }

  document_metadata_configuration_updates {
    name = "example-date-value"
    type = "DATE_VALUE"
    search {
      displayable = true
      facetable   = true
      searchable  = false
      sortable    = %[5]t
    }
    relevance {
      freshness  = false
      importance = 1
      duration   = "25920000s"
      rank_order = "ASCENDING"
    }
  }
}
`, rName3, authorsFacetable, longValDisplayable, stringListValSearchable, dateValSortable, stringValImportance))
}
