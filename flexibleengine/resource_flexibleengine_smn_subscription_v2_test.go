package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/smn/v2/subscriptions"
)

func TestAccSMNV2Subscription_basic(t *testing.T) {
	var subscription2 subscriptions.SubscriptionGet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMNSubscriptionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2SubscriptionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2SubscriptionExists("flexibleengine_smn_subscription_v2.subscription_2", &subscription2),
					resource.TestCheckResourceAttr(
						"flexibleengine_smn_subscription_v2.subscription_2", "endpoint",
						"13600000000"),
				),
			},
		},
	})
}

func testAccCheckSMNSubscriptionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	smnClient, err := config.SmnV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine smn: %s", err)
	}
	var subscription *subscriptions.SubscriptionGet
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_smn_subscription_v2" {
			continue
		}
		foundList, err := subscriptions.List(smnClient).Extract()
		if err != nil {
			return err
		}
		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription != nil {
			return fmt.Errorf("subscription still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2SubscriptionExists(n string, subscription *subscriptions.SubscriptionGet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		smnClient, err := config.SmnV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine smn client: %s", err)
		}

		foundList, err := subscriptions.List(smnClient).Extract()
		if err != nil {
			return err
		}
		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription == nil {
			return fmt.Errorf("subscription not found")
		}

		return nil
	}
}

var TestAccSMNV2SubscriptionConfig_basic = `
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name		  = "topic_1"
  display_name    = "The display name of topic_1"
}

#resource "flexibleengine_smn_subscription_v2" "subscription_1" {
#  topic_urn       = "${flexibleengine_smn_topic_v2.topic_1.id}"
#  endpoint        = "mailtest@gmail.com"
#  protocol        = "email"
#  remark          = "O&M"
#}

resource "flexibleengine_smn_subscription_v2" "subscription_2" {
  topic_urn       = "${flexibleengine_smn_topic_v2.topic_1.id}"
  endpoint        = "13600000000"
  protocol        = "sms"
  remark          = ""
}
`
