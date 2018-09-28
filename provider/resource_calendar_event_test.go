package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccResourceCalendarEvent_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// TODO: CheckDestroy:
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCalendarEvent("60m"),
				// TODO: Check:
			},
			{
				Config:   testAccResourceCalendarEvent("1h"),
				PlanOnly: true,
				// TODO: Check:
			},
			{
				Config:   testAccResourceCalendarEvent("60m"),
				PlanOnly: true,
				// TODO: Check:
			},
		},
	})
}

func TestAccResourceCalendarEvent_acls(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// TODO: CheckDestroy:
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCalendarEvent_acls(),
				// TODO: Check:
			},
		},
	})
}

func testAccResourceCalendarEvent(reminder string) string {
	return fmt.Sprintf(`
resource "gsuite_calendar_event" "demo" {
	summary     = "Terraform Test Event"
	description = "Long-form description"
	location    = "Conference Room B"

	start = "2017-10-12T15:00:00-05:00"
	end   = "2017-10-12T17:00:00-05:00"

	attendee {
		email = "terraform-acctest@hashicorp.com"
	}

	attendee {
		email    = "terraform-acctest+alt@hashicorp.com"
		optional = true
	}

	reminder {
		method = "email"
		before = "%s"
	}

	reminder {
		method = "popup"
		before = "120m"
	}
}
`, reminder)
}

func testAccResourceCalendarEvent_acls() string {
	return `
resource "gsuite_calendar_event" "demo" {
	summary     = "My Open Event"
	description = "Anyone can do anything, anytime, anywhere."
	location    = "Wherever you want!"

	start = "2017-10-12T15:00:00-05:00"
	end   = "2017-10-12T17:00:00-05:00"

	attendee {
		email = "terraform-acctest@hashicorp.com"
	}

	guests_can_invite_others    = true
	guests_can_modify           = true
	guests_can_see_other_guests = true
}
`
}
