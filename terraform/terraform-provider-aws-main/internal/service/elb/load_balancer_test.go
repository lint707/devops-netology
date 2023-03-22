package elb_test

import ( // nosemgrep:ci.aws-sdk-go-multiple-service-imports
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelb "github.com/hashicorp/terraform-provider-aws/internal/service/elb"
)

func TestAccELBLoadBalancer_basic(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testAccCheckLoadBalancerAttributes(&conf),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8000",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
					resource.TestCheckResourceAttr(resourceName, "cross_zone_load_balancing", "true"),
					resource.TestCheckResourceAttr(resourceName, "desync_mitigation_mode", "defensive"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccELBLoadBalancer_disappears(t *testing.T) {
	var loadBalancer elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &loadBalancer),
					testAccCheckLoadBalancerDisappears(&loadBalancer),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccELBLoadBalancer_fullCharacterRange(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_fullRangeOfCharacters(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_AccessLogs_enabled(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"
	rName := fmt.Sprintf("tf-test-access-logs-%d", sdkacctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_accessLogs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
				),
			},

			{
				Config: testAccLoadBalancerConfig_accessLogsOn(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "access_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_logs.0.bucket", rName),
					resource.TestCheckResourceAttr(resourceName, "access_logs.0.interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "access_logs.0.enabled", "true"),
				),
			},

			{
				Config: testAccLoadBalancerConfig_accessLogs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "access_logs.#", "0"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_AccessLogs_disabled(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"
	rName := fmt.Sprintf("tf-test-access-logs-%d", sdkacctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_accessLogs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
				),
			},
			{
				Config: testAccLoadBalancerConfig_accessLogsDisabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "access_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_logs.0.bucket", rName),
					resource.TestCheckResourceAttr(resourceName, "access_logs.0.interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "access_logs.0.enabled", "false"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_accessLogs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(
						resourceName, "access_logs.#", "0"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_namePrefix(t *testing.T) {
	var conf elb.LoadBalancerDescription
	nameRegex := regexp.MustCompile("^test-")
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestMatchResourceAttr(resourceName, "name", nameRegex),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_generatedName(t *testing.T) {
	var conf elb.LoadBalancerDescription
	generatedNameRegexp := regexp.MustCompile("^tf-lb-")
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestMatchResourceAttr(resourceName, "name", generatedNameRegexp),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_generatesNameForZeroValue(t *testing.T) {
	var conf elb.LoadBalancerDescription
	generatedNameRegexp := regexp.MustCompile("^tf-lb-")
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_zeroValueName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestMatchResourceAttr(resourceName, "name", generatedNameRegexp),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_availabilityZones(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "3"),
				),
			},

			{
				Config: testAccLoadBalancerConfig_availabilityZonesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "2"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_tags(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_tags1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testAccCheckLoadBalancerAttributes(&conf),
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
				Config: testAccLoadBalancerConfig_tags2("key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testAccCheckLoadBalancerAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_tags1("key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testAccCheckLoadBalancerAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_ListenerSSLCertificateID_iamServerCertificate(t *testing.T) {
	var conf elb.LoadBalancerDescription
	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")
	rName := fmt.Sprintf("tf-acctest-%s", sdkacctest.RandString(10))
	resourceName := "aws_elb.test"

	testCheck := func(*terraform.State) error {
		if len(conf.ListenerDescriptions) != 1 {
			return fmt.Errorf(
				"TestAccELBLoadBalancer_ListenerSSLCertificateID_iamServerCertificate expected 1 listener, got %d",
				len(conf.ListenerDescriptions))
		}
		return nil
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccLoadBalancerConfig_listenerIAMServerCertificate(rName, certificate, key, "tcp"),
				ExpectError: regexp.MustCompile(`ssl_certificate_id may be set only when protocol is 'https' or 'ssl'`),
			},
			{
				Config: testAccLoadBalancerConfig_listenerIAMServerCertificate(rName, certificate, key, "https"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testCheck,
				),
			},
			{
				Config:      testAccLoadBalancerConfig_listenerIAMServerCertificateAddInvalidListener(rName, certificate, key),
				ExpectError: regexp.MustCompile(`ssl_certificate_id may be set only when protocol is 'https' or 'ssl'`),
			},
		},
	})
}

func TestAccELBLoadBalancer_Swap_subnets(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_subnets,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "2"),
				),
			},

			{
				Config: testAccLoadBalancerConfig_subnetSwap,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists("aws_elb.test", &conf),
					resource.TestCheckResourceAttr("aws_elb.test", "subnets.#", "2"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_instanceAttaching(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	testCheckInstanceAttached := func(count int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(conf.Instances) != count {
				return fmt.Errorf("instance count does not match")
			}
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testAccCheckLoadBalancerAttributes(&conf),
				),
			},

			{
				Config: testAccLoadBalancerConfig_newInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					testCheckInstanceAttached(1),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_listener(t *testing.T) {
	var conf elb.LoadBalancerDescription
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8000",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
				),
			},
			{
				Config: testAccLoadBalancerConfig_listenerMultipleListeners,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8000",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "22",
						"instance_protocol": "tcp",
						"lb_port":           "22",
						"lb_protocol":       "tcp",
					}),
				),
			},
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8000",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
				),
			},
			{
				Config: testAccLoadBalancerConfig_listenerUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8080",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
				),
			},
			{
				PreConfig: func() {
					// Simulate out of band listener removal
					conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn
					input := &elb.DeleteLoadBalancerListenersInput{
						LoadBalancerName:  conf.LoadBalancerName,
						LoadBalancerPorts: []*int64{aws.Int64(80)},
					}
					if _, err := conn.DeleteLoadBalancerListeners(input); err != nil {
						t.Fatalf("Error deleting listener: %s", err)
					}
				},
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8000",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
				),
			},
			{
				PreConfig: func() {
					// Simulate out of band listener addition
					conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn
					input := &elb.CreateLoadBalancerListenersInput{
						LoadBalancerName: conf.LoadBalancerName,
						Listeners: []*elb.Listener{
							{
								InstancePort:     aws.Int64(22),
								InstanceProtocol: aws.String("tcp"),
								LoadBalancerPort: aws.Int64(22),
								Protocol:         aws.String("tcp"),
							},
						},
					}
					if _, err := conn.CreateLoadBalancerListeners(input); err != nil {
						t.Fatalf("Error creating listener: %s", err)
					}
				},
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "listener.*", map[string]string{
						"instance_port":     "8000",
						"instance_protocol": "http",
						"lb_port":           "80",
						"lb_protocol":       "http",
					}),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_healthCheck(t *testing.T) {
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_healthCheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "health_check.0.healthy_threshold", "5"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_healthCheckUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "health_check.0.healthy_threshold", "10"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_timeout(t *testing.T) {
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_idleTimeout,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "idle_timeout", "200"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_idleTimeoutUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "idle_timeout", "400"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_connectionDraining(t *testing.T) {
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_connectionDraining,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "connection_draining", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_draining_timeout", "400"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_connectionDrainingUpdateTimeout,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "connection_draining", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_draining_timeout", "600"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_connectionDrainingUpdateDisable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "connection_draining", "false"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_securityGroups(t *testing.T) {
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// ELBs get a default security group
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_securityGroups,
				Check: resource.ComposeTestCheckFunc(
					// Count should still be one as we swap in a custom security group
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_desyncMitigationMode(t *testing.T) {
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_desyncMitigationMode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desync_mitigation_mode", "strictest"),
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

func TestAccELBLoadBalancer_desyncMitigationMode_update(t *testing.T) {
	resourceName := "aws_elb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_desyncMitigationModeUpdateDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desync_mitigation_mode", "defensive"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoadBalancerConfig_desyncMitigationModeUpdateMonitor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desync_mitigation_mode", "monitor"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoadBalancerConfig_desyncMitigationModeUpdateDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desync_mitigation_mode", "defensive"),
				),
			},
		},
	})
}

// Unit test for listeners hash
func TestLoadBalancerListenerHash(t *testing.T) {
	cases := map[string]struct {
		Left  map[string]interface{}
		Right map[string]interface{}
		Match bool
	}{
		"protocols are case insensitive": {
			map[string]interface{}{
				"instance_port":     80,
				"instance_protocol": "TCP",
				"lb_port":           80,
				"lb_protocol":       "TCP",
			},
			map[string]interface{}{
				"instance_port":     80,
				"instance_protocol": "Tcp",
				"lb_port":           80,
				"lb_protocol":       "tcP",
			},
			true,
		},
	}

	for tn, tc := range cases {
		leftHash := tfelb.ListenerHash(tc.Left)
		rightHash := tfelb.ListenerHash(tc.Right)
		if leftHash == rightHash != tc.Match {
			t.Fatalf("%s: expected match: %t, but did not get it", tn, tc.Match)
		}
	}
}

func TestValidLoadBalancerNameCannotBeginWithHyphen(t *testing.T) {
	var n = "-Testing123"
	_, errors := tfelb.ValidName(n, "SampleKey")

	if len(errors) != 1 {
		t.Fatalf("Expected the ELB Name to trigger a validation error")
	}
}

func TestValidLoadBalancerNameCanBeAnEmptyString(t *testing.T) {
	var n = ""
	_, errors := tfelb.ValidName(n, "SampleKey")

	if len(errors) != 0 {
		t.Fatalf("Expected the ELB Name to pass validation")
	}
}

func TestValidLoadBalancerNameCannotBeLongerThan32Characters(t *testing.T) {
	var n = "Testing123dddddddddddddddddddvvvv"
	_, errors := tfelb.ValidName(n, "SampleKey")

	if len(errors) != 1 {
		t.Fatalf("Expected the ELB Name to trigger a validation error")
	}
}

func TestValidLoadBalancerNameCannotHaveSpecialCharacters(t *testing.T) {
	var n = "Testing123%%"
	_, errors := tfelb.ValidName(n, "SampleKey")

	if len(errors) != 1 {
		t.Fatalf("Expected the ELB Name to trigger a validation error")
	}
}

func TestValidLoadBalancerNameCannotEndWithHyphen(t *testing.T) {
	var n = "Testing123-"
	_, errors := tfelb.ValidName(n, "SampleKey")

	if len(errors) != 1 {
		t.Fatalf("Expected the ELB Name to trigger a validation error")
	}
}

func TestValidLoadBalancerAccessLogsInterval(t *testing.T) {
	type testCases struct {
		Value    int
		ErrCount int
	}

	invalidCases := []testCases{
		{
			Value:    0,
			ErrCount: 1,
		},
		{
			Value:    10,
			ErrCount: 1,
		},
		{
			Value:    -1,
			ErrCount: 1,
		},
	}

	for _, tc := range invalidCases {
		_, errors := tfelb.ValidAccessLogsInterval(tc.Value, "interval")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error.", tc.Value)
		}
	}

}

func TestValidLoadBalancerHealthCheckTarget(t *testing.T) {
	type testCase struct {
		Value    string
		ErrCount int
	}

	randomRunes := func(n int) string {
		rand.Seed(time.Now().UTC().UnixNano())

		// A complete set of modern Katakana characters.
		runes := []rune("アイウエオ" +
			"カキクケコガギグゲゴサシスセソザジズゼゾ" +
			"タチツテトダヂヅデドナニヌネノハヒフヘホ" +
			"バビブベボパピプペポマミムメモヤユヨラリ" +
			"ルレロワヰヱヲン")

		s := make([]rune, n)
		for i := range s {
			s[i] = runes[rand.Intn(len(runes))]
		}
		return string(s)
	}

	validCases := []testCase{
		{
			Value:    "TCP:1234",
			ErrCount: 0,
		},
		{
			Value:    "http:80/test",
			ErrCount: 0,
		},
		{
			Value:    fmt.Sprintf("HTTP:8080/%s", randomRunes(5)),
			ErrCount: 0,
		},
		{
			Value:    "SSL:8080",
			ErrCount: 0,
		},
	}

	for _, tc := range validCases {
		_, errors := tfelb.ValidHeathCheckTarget(tc.Value, "target")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q not to trigger a validation error.", tc.Value)
		}
	}

	invalidCases := []testCase{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "TCP:",
			ErrCount: 1,
		},
		{
			Value:    "TCP:1234/",
			ErrCount: 1,
		},
		{
			Value:    "SSL:8080/",
			ErrCount: 1,
		},
		{
			Value:    "HTTP:8080",
			ErrCount: 1,
		},
		{
			Value:    "incorrect-value",
			ErrCount: 1,
		},
		{
			Value:    "TCP:123456",
			ErrCount: 1,
		},
		{
			Value:    "incorrect:80/",
			ErrCount: 1,
		},
		{
			Value: fmt.Sprintf("HTTP:8080/%s%s",
				sdkacctest.RandStringFromCharSet(512, sdkacctest.CharSetAlpha), randomRunes(512)),
			ErrCount: 1,
		},
	}

	for _, tc := range invalidCases {
		_, errors := tfelb.ValidHeathCheckTarget(tc.Value, "target")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error.", tc.Value)
		}
	}
}

func testAccCheckLoadBalancerDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_elb" {
			continue
		}

		describe, err := conn.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
			LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
		})

		if err == nil {
			if len(describe.LoadBalancerDescriptions) != 0 &&
				*describe.LoadBalancerDescriptions[0].LoadBalancerName == rs.Primary.ID {
				return fmt.Errorf("ELB still exists")
			}
		}

		// Verify the error
		providerErr, ok := err.(awserr.Error)
		if !ok {
			return err
		}

		if providerErr.Code() != elb.ErrCodeAccessPointNotFoundException {
			return fmt.Errorf("Unexpected error: %s", err)
		}
	}

	return nil
}

func testAccCheckLoadBalancerDisappears(loadBalancer *elb.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn

		input := elb.DeleteLoadBalancerInput{
			LoadBalancerName: loadBalancer.LoadBalancerName,
		}
		_, err := conn.DeleteLoadBalancer(&input)

		return err
	}
}

func testAccCheckLoadBalancerAttributes(conf *elb.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		l := elb.Listener{
			InstancePort:     aws.Int64(8000),
			InstanceProtocol: aws.String("HTTP"),
			LoadBalancerPort: aws.Int64(80),
			Protocol:         aws.String("HTTP"),
		}

		if !reflect.DeepEqual(conf.ListenerDescriptions[0].Listener, &l) {
			return fmt.Errorf(
				"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
				conf.ListenerDescriptions[0].Listener,
				l)
		}

		if *conf.DNSName == "" {
			return fmt.Errorf("empty dns_name")
		}

		return nil
	}
}

func testAccCheckLoadBalancerExists(n string, res *elb.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ELB ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn

		describe, err := conn.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
			LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if len(describe.LoadBalancerDescriptions) != 1 ||
			*describe.LoadBalancerDescriptions[0].LoadBalancerName != rs.Primary.ID {
			return fmt.Errorf("ELB not found")
		}

		*res = *describe.LoadBalancerDescriptions[0]

		// Confirm source_security_group_id for ELBs in a VPC
		// 	See https://github.com/hashicorp/terraform/pull/3780
		if res.VPCId != nil {
			sgid := rs.Primary.Attributes["source_security_group_id"]
			if sgid == "" {
				return fmt.Errorf("Expected to find source_security_group_id for ELB, but was empty")
			}
		}

		return nil
	}
}

const testAccLoadBalancerConfig_basic = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  cross_zone_load_balancing = true
}
`

func testAccLoadBalancerConfig_tags1(tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  tags = {
    %[1]q = %[2]q
  }

  cross_zone_load_balancing = true
}
`, tagKey1, tagValue1)
}

func testAccLoadBalancerConfig_tags2(tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }

  cross_zone_load_balancing = true
}
`, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccLoadBalancerConfig_fullRangeOfCharacters(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  name               = %[1]q
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`, rName)
}

const testAccLoadBalancerConfig_accessLogs = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`

func testAccLoadBalancerConfig_accessLogsOn(r string) string {
	return `
resource "aws_elb" "test" {
  # Must have bucket policy attached first
  depends_on = [aws_s3_bucket_policy.test]

  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  access_logs {
    interval = 5
    bucket   = aws_s3_bucket.accesslogs_bucket.bucket
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}
` + testAccLoadBalancerAccessLogsCommon(r)
}

func testAccLoadBalancerConfig_accessLogsDisabled(r string) string {
	return `
resource "aws_elb" "test" {
  # Must have bucket policy attached first
  depends_on = [aws_s3_bucket_policy.test]

  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  access_logs {
    interval = 5
    bucket   = aws_s3_bucket.accesslogs_bucket.bucket
    enabled  = false
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}
` + testAccLoadBalancerAccessLogsCommon(r)
}

func testAccLoadBalancerAccessLogsCommon(r string) string {
	return fmt.Sprintf(`
data "aws_elb_service_account" "current" {}

data "aws_partition" "current" {}

resource "aws_s3_bucket" "accesslogs_bucket" {
  bucket        = "%[1]s"
  force_destroy = true
}

resource "aws_s3_bucket_acl" "accesslogs_bucket_acl" {
  bucket = aws_s3_bucket.accesslogs_bucket.id
  acl    = "private"
}

resource "aws_s3_bucket_policy" "test" {
  bucket = aws_s3_bucket.accesslogs_bucket.id
  policy = <<EOF
{
  "Id": "Policy1446577137248",
  "Statement": [
    {
      "Action": "s3:PutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "${data.aws_elb_service_account.current.arn}"
      },
      "Resource": "${aws_s3_bucket.accesslogs_bucket.arn}/*",
      "Sid": "Stmt1446575236270"
    }
  ],
  "Version": "2012-10-17"
}
EOF
}
`, r)
}

const testAccLoadBalancerConfig_namePrefix = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  name_prefix        = "test-"
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`

const testAccLoadBalancerConfig_generatedName = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`

const testAccLoadBalancerConfig_zeroValueName = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  name               = ""
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}

# See https://github.com/hashicorp/terraform-provider-aws/issues/2498
output "lb_name" {
  value = aws_elb.test.name
}
`

const testAccLoadBalancerConfig_availabilityZonesUpdate = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`

const testAccLoadBalancerConfig_newInstance = `
data "aws_ami" "amzn-ami-minimal-hvm-ebs" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-minimal-hvm-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  instances = [aws_instance.test.id]
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.micro"
}
`

const testAccLoadBalancerConfig_healthCheck = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 5
    unhealthy_threshold = 5
    target              = "HTTP:8000/"
    interval            = 60
    timeout             = 30
  }
}
`

const testAccLoadBalancerConfig_healthCheckUpdate = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 10
    unhealthy_threshold = 5
    target              = "HTTP:8000/"
    interval            = 60
    timeout             = 30
  }
}
`

const testAccLoadBalancerConfig_listenerUpdate = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8080
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`

const testAccLoadBalancerConfig_listenerMultipleListeners = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  listener {
    instance_port     = 22
    instance_protocol = "tcp"
    lb_port           = 22
    lb_protocol       = "tcp"
  }
}
`

const testAccLoadBalancerConfig_idleTimeout = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  idle_timeout = 200
}
`

const testAccLoadBalancerConfig_idleTimeoutUpdate = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  idle_timeout = 400
}
`

const testAccLoadBalancerConfig_connectionDraining = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  connection_draining         = true
  connection_draining_timeout = 400
}
`

const testAccLoadBalancerConfig_connectionDrainingUpdateTimeout = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  connection_draining         = true
  connection_draining_timeout = 600
}
`

const testAccLoadBalancerConfig_connectionDrainingUpdateDisable = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  connection_draining = false
}
`

const testAccLoadBalancerConfig_securityGroups = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  security_groups = [aws_security_group.test.id]
}

resource "aws_security_group" "test" {
  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "tf_elb_sg_test"
  }
}
`

func testAccLoadBalancerConfig_listenerIAMServerCertificate(certName, certificate, key, lbProtocol string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_iam_server_certificate" "test_cert" {
  name             = "%[1]s"
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port      = 443
    instance_protocol  = "%[4]s"
    lb_port            = 443
    lb_protocol        = "%[4]s"
    ssl_certificate_id = aws_iam_server_certificate.test_cert.arn
  }
}
`, certName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key), lbProtocol)
}

func testAccLoadBalancerConfig_listenerIAMServerCertificateAddInvalidListener(certName, certificate, key string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_iam_server_certificate" "test_cert" {
  name             = "%[1]s"
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port      = 443
    instance_protocol  = "https"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = aws_iam_server_certificate.test_cert.arn
  }

  # lb_protocol tcp and ssl_certificate_id is not valid
  listener {
    instance_port      = 8443
    instance_protocol  = "tcp"
    lb_port            = 8443
    lb_protocol        = "tcp"
    ssl_certificate_id = aws_iam_server_certificate.test_cert.arn
  }
}
`, certName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}

const testAccLoadBalancerConfig_subnets = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "azelb" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name = "terraform-testacc-elb-subnets"
  }
}

resource "aws_subnet" "public_a_one" {
  vpc_id = aws_vpc.azelb.id

  cidr_block        = "10.1.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  tags = {
    Name = "tf-acc-elb-subnets-a-one"
  }
}

resource "aws_subnet" "public_b_one" {
  vpc_id = aws_vpc.azelb.id

  cidr_block        = "10.1.7.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "tf-acc-elb-subnets-b-one"
  }
}

resource "aws_subnet" "public_a_two" {
  vpc_id = aws_vpc.azelb.id

  cidr_block        = "10.1.2.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "tf-acc-elb-subnets-a-two"
  }
}

resource "aws_elb" "test" {
  name = "terraform-asg-deployment-example"

  subnets = [
    aws_subnet.public_a_one.id,
    aws_subnet.public_b_one.id,
  ]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  depends_on = [aws_internet_gateway.gw]
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.azelb.id

  tags = {
    Name = "main"
  }
}
`

const testAccLoadBalancerConfig_subnetSwap = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "azelb" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name = "terraform-testacc-elb-subnet-swap"
  }
}

resource "aws_subnet" "public_a_one" {
  vpc_id = aws_vpc.azelb.id

  cidr_block        = "10.1.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  tags = {
    Name = "tf-acc-elb-subnet-swap-a-one"
  }
}

resource "aws_subnet" "public_b_one" {
  vpc_id = aws_vpc.azelb.id

  cidr_block        = "10.1.7.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]
  tags = {
    Name = "tf-acc-elb-subnet-swap-b-one"
  }
}

resource "aws_subnet" "public_a_two" {
  vpc_id = aws_vpc.azelb.id

  cidr_block        = "10.1.2.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  tags = {
    Name = "tf-acc-elb-subnet-swap-a-two"
  }
}

resource "aws_elb" "test" {
  name = "terraform-asg-deployment-example"

  subnets = [
    aws_subnet.public_a_two.id,
    aws_subnet.public_b_one.id,
  ]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  depends_on = [aws_internet_gateway.gw]
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.azelb.id

  tags = {
    Name = "main"
  }
}
`

const testAccLoadBalancerConfig_desyncMitigationMode = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]
  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  desync_mitigation_mode = "strictest"
}
`

const testAccLoadBalancerConfig_desyncMitigationModeUpdateDefault = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
`

const testAccLoadBalancerConfig_desyncMitigationModeUpdateMonitor = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_elb" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  desync_mitigation_mode = "monitor"
}
`
