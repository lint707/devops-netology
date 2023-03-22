package memorydb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/memorydb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfmemorydb "github.com/hashicorp/terraform-provider-aws/internal/service/memorydb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccMemoryDBACL_basic(t *testing.T) {
	rName := "tf-test-" + sdkacctest.RandString(8)
	user1 := "tf-test-" + sdkacctest.RandString(8)
	resourceName := "aws_memorydb_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, memorydb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccACLConfig_basic(rName, []string{user1}, []string{user1}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "memorydb", "acl/"+rName),
					resource.TestCheckResourceAttrSet(resourceName, "minimum_engine_version"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Test", "test"),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "user_names.*", "aws_memorydb_user.test.0", "user_name"),
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

func TestAccMemoryDBACL_disappears(t *testing.T) {
	rName := "tf-test-" + sdkacctest.RandString(8)
	resourceName := "aws_memorydb_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, memorydb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccACLConfig_basic(rName, nil, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfmemorydb.ResourceACL(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMemoryDBACL_nameGenerated(t *testing.T) {
	resourceName := "aws_memorydb_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, memorydb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccACLConfig_noName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
				),
			},
		},
	})
}

func TestAccMemoryDBACL_namePrefix(t *testing.T) {
	resourceName := "aws_memorydb_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, memorydb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccACLConfig_namePrefix("tftest-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", "tftest-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tftest-"),
				),
			},
		},
	})
}

func TestAccMemoryDBACL_update_tags(t *testing.T) {
	rName := "tf-test-" + sdkacctest.RandString(8)
	resourceName := "aws_memorydb_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, memorydb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccACLConfig_tags0(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccACLConfig_tags2(rName, "Key1", "value1", "Key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccACLConfig_tags1(rName, "Key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccACLConfig_tags0(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "0"),
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

func TestAccMemoryDBACL_update_userNames(t *testing.T) {
	rName := "tf-test-" + sdkacctest.RandString(8)
	user1 := "tf-test1-" + sdkacctest.RandString(8)
	user2 := "tf-test2-" + sdkacctest.RandString(8)
	user3 := "tf-test3-" + sdkacctest.RandString(8)
	resourceName := "aws_memorydb_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, memorydb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckACLDestroy,
		Steps: []resource.TestStep{
			{
				// Empty ACL.
				Config: testAccACLConfig_basic(rName, []string{}, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Adding users.
				Config: testAccACLConfig_basic(rName, []string{user1, user2}, []string{user1, user2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user1),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Removing and adding a user.
				Config: testAccACLConfig_basic(rName, []string{user1, user2, user3}, []string{user1, user3}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user1),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user3),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Removing a user.
				Config: testAccACLConfig_basic(rName, []string{user1, user2, user3}, []string{user1}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccACLConfig_basic(rName, []string{user1, user2}, []string{user1, user2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user1),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Deleting a user before disassociating it.
				Config: testAccACLConfig_basic(rName, []string{user1}, []string{user1}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_names.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_names.*", user1),
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

func testAccCheckACLDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).MemoryDBConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_memorydb_acl" {
			continue
		}

		_, err := tfmemorydb.FindACLByName(context.Background(), conn, rs.Primary.Attributes["name"])

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("MemoryDB ACL %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckACLExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MemoryDB ACL ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).MemoryDBConn

		_, err := tfmemorydb.FindACLByName(context.Background(), conn, rs.Primary.Attributes["name"])

		return err
	}
}

func testAccACLConfigUsers(names ...string) string {
	var userNames string
	for i, name := range names {
		if i > 0 {
			userNames += ", "
		}
		userNames += fmt.Sprintf("%q", name)
	}

	return fmt.Sprintf(`
locals {
  user_names = [%[1]s]
}

resource "aws_memorydb_user" "test" {
  count         = length(local.user_names)
  access_string = "on ~* &* +@all"
  user_name     = local.user_names[count.index]

  authentication_mode {
    type      = "password"
    passwords = ["aaaaaaaaaaaaaaaa"]
  }
}
`, userNames)
}

func testAccACLConfig_basic(rName string, userNames []string, usersInACL []string) string {
	var userNamesInACL string
	for i, userName := range usersInACL {
		if i > 0 {
			userNamesInACL += ", "
		}
		userNamesInACL += fmt.Sprintf("%q", userName)
	}

	return acctest.ConfigCompose(
		testAccACLConfigUsers(userNames...),
		fmt.Sprintf(`
resource "aws_memorydb_acl" "test" {
  depends_on = [aws_memorydb_user.test]

  name       = %[1]q
  user_names = [%[2]s]

  tags = {
    Test = "test"
  }
}
`, rName, userNamesInACL),
	)
}

func testAccACLConfig_noName() string {
	return `
resource "aws_memorydb_acl" "test" {}
`
}

func testAccACLConfig_namePrefix(namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_memorydb_acl" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}

func testAccACLConfig_tags0(rName string) string {
	return fmt.Sprintf(`
resource "aws_memorydb_acl" "test" {
  name = %[1]q
}
`, rName)
}

func testAccACLConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_memorydb_acl" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccACLConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_memorydb_acl" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
