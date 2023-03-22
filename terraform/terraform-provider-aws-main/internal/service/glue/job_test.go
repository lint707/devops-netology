package glue_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/glue"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfglue "github.com/hashicorp/terraform-provider-aws/internal/service/glue"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccGlueJob_basic(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"
	roleResourceName := "aws_iam_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_required(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "glue", fmt.Sprintf("job/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "execution_class", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.%", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", roleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "2880"),
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

func TestAccGlueJob_disappears(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_required(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					acctest.CheckResourceDisappears(acctest.Provider, tfglue.ResourceJob(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccGlueJob_basicStreaming(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"
	roleResourceName := "aws_iam_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_requiredStreaming(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "glue", fmt.Sprintf("job/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.name", "gluestreaming"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", roleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "0"),
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
func TestAccGlueJob_command(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_command(rName, "testscriptlocation1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation1"),
				),
			},
			{
				Config: testAccJobConfig_command(rName, "testscriptlocation2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation2"),
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

func TestAccGlueJob_defaultArguments(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_defaultArguments(rName, "job-bookmark-disable", "python"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.--job-bookmark-option", "job-bookmark-disable"),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.--job-language", "python"),
				),
			},
			{
				Config: testAccJobConfig_defaultArguments(rName, "job-bookmark-enable", "scala"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.--job-bookmark-option", "job-bookmark-enable"),
					resource.TestCheckResourceAttr(resourceName, "default_arguments.--job-language", "scala"),
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

func TestAccGlueJob_nonOverridableArguments(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_nonOverridableArguments(rName, "job-bookmark-disable", "python"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.--job-bookmark-option", "job-bookmark-disable"),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.--job-language", "python"),
				),
			},
			{
				Config: testAccJobConfig_nonOverridableArguments(rName, "job-bookmark-enable", "scala"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.--job-bookmark-option", "job-bookmark-enable"),
					resource.TestCheckResourceAttr(resourceName, "non_overridable_arguments.--job-language", "scala"),
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

func TestAccGlueJob_description(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_description(rName, "First Description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "description", "First Description"),
				),
			},
			{
				Config: testAccJobConfig_description(rName, "Second Description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "description", "Second Description"),
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

func TestAccGlueJob_glueVersion(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_versionMaxCapacity(rName, "0.9"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "glue_version", "0.9"),
				),
			},
			{
				Config: testAccJobConfig_versionMaxCapacity(rName, "1.0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "glue_version", "1.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobConfig_versionNumberOfWorkers(rName, "2.0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "glue_version", "2.0"),
				),
			},
		},
	})
}

func TestAccGlueJob_executionClass(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_executionClass(rName, "FLEX"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "execution_class", "FLEX"),
				),
			},
			{
				Config: testAccJobConfig_executionClass(rName, "STANDARD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "execution_class", "STANDARD"),
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
func TestAccGlueJob_executionProperty(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccJobConfig_executionProperty(rName, 0),
				ExpectError: regexp.MustCompile(`expected execution_property.0.max_concurrent_runs to be at least`),
			},
			{
				Config: testAccJobConfig_executionProperty(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "execution_property.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "execution_property.0.max_concurrent_runs", "1"),
				),
			},
			{
				Config: testAccJobConfig_executionProperty(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "execution_property.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "execution_property.0.max_concurrent_runs", "2"),
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

func TestAccGlueJob_maxRetries(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccJobConfig_maxRetries(rName, 11),
				ExpectError: regexp.MustCompile(`expected max_retries to be in the range`),
			},
			{
				Config: testAccJobConfig_maxRetries(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "0"),
				),
			},
			{
				Config: testAccJobConfig_maxRetries(rName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "10"),
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

func TestAccGlueJob_notificationProperty(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccJobConfig_notificationProperty(rName, 0),
				ExpectError: regexp.MustCompile(`expected notification_property.0.notify_delay_after to be at least`),
			},
			{
				Config: testAccJobConfig_notificationProperty(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "notification_property.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_property.0.notify_delay_after", "1"),
				),
			},
			{
				Config: testAccJobConfig_notificationProperty(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "notification_property.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_property.0.notify_delay_after", "2"),
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

func TestAccGlueJob_tags(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccJobConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccGlueJob_streamingTimeout(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_timeout(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),
				),
			},
			{
				Config: testAccJobConfig_timeout(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "timeout", "2"),
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
func TestAccGlueJob_timeout(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_timeout(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),
				),
			},
			{
				Config: testAccJobConfig_timeout(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "timeout", "2"),
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

func TestAccGlueJob_security(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_securityConfiguration(rName, "default_encryption"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "security_configuration", "default_encryption"),
				),
			},
			{
				Config: testAccJobConfig_securityConfiguration(rName, "custom_encryption2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "security_configuration", "custom_encryption2"),
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

func TestAccGlueJob_workerType(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_workerType(rName, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "worker_type", "Standard"),
				),
			},
			{
				Config: testAccJobConfig_workerType(rName, "G.1X"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "worker_type", "G.1X"),
				),
			},
			{
				Config: testAccJobConfig_workerType(rName, "G.2X"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "worker_type", "G.2X"),
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

func TestAccGlueJob_pythonShell(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_pythonShell(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "0.0625"),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "command.0.name", "pythonshell"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobConfig_pythonShellVersion(rName, "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "command.0.python_version", "2"),
					resource.TestCheckResourceAttr(resourceName, "command.0.name", "pythonshell"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobConfig_pythonShellVersion(rName, "3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "command.0.python_version", "3"),
					resource.TestCheckResourceAttr(resourceName, "command.0.name", "pythonshell"),
				),
			},
			{
				Config: testAccJobConfig_pythonShellVersion(rName, "3.9"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "command.0.python_version", "3.9"),
					resource.TestCheckResourceAttr(resourceName, "command.0.name", "pythonshell"),
				),
			},
		},
	})
}

func TestAccGlueJob_maxCapacity(t *testing.T) {
	var job glue.Job
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_job.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobConfig_maxCapacity(rName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "10"),
					resource.TestCheckResourceAttr(resourceName, "command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "command.0.script_location", "testscriptlocation"),
					resource.TestCheckResourceAttr(resourceName, "command.0.name", "glueetl"),
				),
			},
			{
				Config: testAccJobConfig_maxCapacity(rName, 15),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "15"),
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

func testAccCheckJobExists(n string, v *glue.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Glue Job ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn

		output, err := tfglue.FindJobByName(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckJobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_glue_job" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn

		_, err := tfglue.FindJobByName(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Glue Job %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccJobConfig_base(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_iam_policy" "AWSGlueServiceRole" {
  arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSGlueServiceRole"
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "glue.${data.aws_partition.current.dns_suffix}"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "test" {
  policy_arn = data.aws_iam_policy.AWSGlueServiceRole.arn
  role       = aws_iam_role.test.name
}
`, rName)
}

func testAccJobConfig_command(rName, scriptLocation string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = %[2]q
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, scriptLocation))
}

func testAccJobConfig_defaultArguments(rName, jobBookmarkOption, jobLanguage string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  default_arguments = {
    "--job-bookmark-option" = %[2]q
    "--job-language"        = %[3]q
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, jobBookmarkOption, jobLanguage))
}

func testAccJobConfig_nonOverridableArguments(rName, jobBookmarkOption, jobLanguage string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  non_overridable_arguments = {
    "--job-bookmark-option" = %[2]q
    "--job-language"        = %[3]q
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, jobBookmarkOption, jobLanguage))
}

func testAccJobConfig_description(rName, description string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  description  = %[1]q
  max_capacity = 10
  name         = %[2]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, description, rName))
}

func testAccJobConfig_versionMaxCapacity(rName, glueVersion string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  glue_version = %[1]q
  max_capacity = 10
  name         = %[2]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, glueVersion, rName))
}

func testAccJobConfig_versionNumberOfWorkers(rName, glueVersion string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  glue_version      = %[1]q
  name              = %[2]q
  number_of_workers = 2
  role_arn          = aws_iam_role.test.arn
  worker_type       = "Standard"

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, glueVersion, rName))
}

func testAccJobConfig_executionClass(rName, executionClass string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  execution_class   = %[2]q
  name              = %[1]q
  number_of_workers = 2
  role_arn          = aws_iam_role.test.arn
  worker_type       = "G.1X"
  glue_version      = "3.0"

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, executionClass))
}

func testAccJobConfig_executionProperty(rName string, maxConcurrentRuns int) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  execution_property {
    max_concurrent_runs = %[2]d
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, maxConcurrentRuns))
}

func testAccJobConfig_maxRetries(rName string, maxRetries int) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  max_retries  = %[1]d
  name         = %[2]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, maxRetries, rName))
}

func testAccJobConfig_notificationProperty(rName string, notifyDelayAfter int) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  notification_property {
    notify_delay_after = %[2]d
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, notifyDelayAfter))
}

func testAccJobConfig_required(rName string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccJobConfig_requiredStreaming(rName string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn

  command {
    name            = "gluestreaming"
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccJobConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  name              = %[1]q
  number_of_workers = 2
  role_arn          = aws_iam_role.test.arn
  worker_type       = "Standard"

  command {
    script_location = "testscriptlocation"
  }

  tags = {
    %[2]q = %[3]q
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, tagKey1, tagValue1))
}

func testAccJobConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  name              = %[1]q
  number_of_workers = 2
  role_arn          = aws_iam_role.test.arn
  worker_type       = "Standard"

  command {
    script_location = "testscriptlocation"
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccJobConfig_timeout(rName string, timeout int) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity = 10
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn
  timeout      = %[2]d

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, timeout))
}

func testAccJobConfig_securityConfiguration(rName string, securityConfiguration string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  max_capacity           = 10
  name                   = %[1]q
  role_arn               = aws_iam_role.test.arn
  security_configuration = %[2]q

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, securityConfiguration))
}

func testAccJobConfig_workerType(rName string, workerType string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  name              = %[1]q
  role_arn          = aws_iam_role.test.arn
  worker_type       = %[2]q
  number_of_workers = 10

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, workerType))
}

func testAccJobConfig_pythonShell(rName string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn
  max_capacity = 0.0625

  command {
    name            = "pythonshell"
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccJobConfig_pythonShellVersion(rName string, pythonVersion string) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn
  max_capacity = 0.0625

  command {
    name            = "pythonshell"
    script_location = "testscriptlocation"
    python_version  = %[2]q
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, pythonVersion))
}

func testAccJobConfig_maxCapacity(rName string, maxCapacity float64) string {
	return acctest.ConfigCompose(testAccJobConfig_base(rName), fmt.Sprintf(`
resource "aws_glue_job" "test" {
  name         = %[1]q
  role_arn     = aws_iam_role.test.arn
  max_capacity = %[2]g

  command {
    script_location = "testscriptlocation"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, maxCapacity))
}
