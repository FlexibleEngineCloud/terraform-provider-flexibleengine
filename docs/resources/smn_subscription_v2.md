---
subcategory: "Simple Message Notification (SMN)"
---

# flexibleengine\_smn\_subscription\_v2

Manages a V2 subscription resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name		  = "topic_1"
  display_name    = "The display name of topic_1"
}

resource "flexibleengine_smn_subscription_v2" "subscription_1" {
  topic_urn       = "${flexibleengine_smn_topic_v2.topic_1.id}"
  endpoint        = "mailtest@gmail.com"
  protocol        = "email"
  remark          = "O&M"
}

resource "flexibleengine_smn_subscription_v2" "subscription_2" {
  topic_urn       = "${flexibleengine_smn_topic_v2.topic_1.id}"
  endpoint        = "13600000000"
  protocol        = "sms"
  remark          = "O&M"
}
```

## Argument Reference

The following arguments are supported:

* `topic_urn` - (Required) Resource identifier of a topic, which is unique.

* `endpoint` - (Required) Message endpoint.
     For an HTTP subscription, the endpoint starts with http\://.
     For an HTTPS subscription, the endpoint starts with https\://.
     For an email subscription, the endpoint is a mail address.
     For an SMS message subscription, the endpoint is a phone number.

* `protocol` - (Required) Protocol of the message endpoint. Currently, email,
     sms, http, and https are supported.

* `remark` - (Optional) Remark information. The remarks must be a UTF-8-coded
     character string containing 128 bytes.

* `subscription_urn` - (Optional) Resource identifier of a subscription, which
     is unique.

* `owner` - (Optional) Project ID of the topic creator.

* `status` - (Optional) Subscription status.
     0 indicates that the subscription is not confirmed.
     1 indicates that the subscription is confirmed.
     3 indicates that the subscription is canceled.


## Attributes Reference

The following attributes are exported:

* `topic_urn` - See Argument Reference above.
* `endpoint` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `remark` - See Argument Reference above.
* `subscription_urn` - See Argument Reference above.
* `owner` - See Argument Reference above.
* `status` - See Argument Reference above.
