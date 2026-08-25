// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventWebhookResource(t *testing.T) {
	resourceName := "sendgrid_event_webhook.test"

	url := fmt.Sprintf("https://test-acc-%s.com", acctest.RandString(16))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEventWebhookResourceConfig(url, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEventWebhookResourceConfig(url, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			// Enable signature verification: public_key must be populated
			{
				Config: testAccEventWebhookResourceConfig(url, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "signed", "true"),
					resource.TestCheckResourceAttrWith(resourceName, "public_key", checkNonEmpty),
				),
			},
			// Regression test: updating an unrelated field (enabled) while
			// `signed` stays true must not clear public_key. Update() only
			// populated public_key when `signed` itself changed, defaulting
			// to "" otherwise and clobbering the previously known key.
			{
				Config: testAccEventWebhookResourceConfig(url, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "signed", "true"),
					resource.TestCheckResourceAttrWith(resourceName, "public_key", checkNonEmpty),
				),
			},
		},
	})
}

func checkNonEmpty(value string) error {
	if value == "" {
		return fmt.Errorf("expected a non-empty value")
	}
	return nil
}

func testAccEventWebhookResourceConfig(url string, enabled bool, signed bool) string {
	return fmt.Sprintf(`
resource "sendgrid_event_webhook" "test" {
  url     = "%s"
  enabled = %t
  signed  = %t
}
`, url, enabled, signed)
}
