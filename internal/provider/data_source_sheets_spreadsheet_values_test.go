package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSheetsSpreadsheetValues_basic(t *testing.T) {
	// example spreadsheet from Google's docs: https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	const (
		spreadsheetID = "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
		readRange     = "Class Data!A2:E"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy:
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSheetsSpreadsheetValues(spreadsheetID, readRange),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("A1", "Alexandra"),
					resource.TestCheckOutput("E12", "Art"),
				),
			},
		},
	})
}

func testAccDataSourceSheetsSpreadsheetValues(spreadsheetID, readRange string) string {
	return fmt.Sprintf(`
data "gsuite_sheets_spreadsheet_values" "test" {
	spreadsheet_id = "%s"
	range = "%s"
}

output "A1" {
	# this syntax doesn't seem to work, so using element
	# value = "${data.gsuite_sheets_spreadsheet_values.test.values[0][0]}"

	value = "${element(data.gsuite_sheets_spreadsheet_values.test.values[0], 0)}"
}

output "E12" {
	# this syntax doesn't seem to work, so using element
	# value = "${data.gsuite_sheets_spreadsheet_values.test.values[11][4]}"

	value = "${element(data.gsuite_sheets_spreadsheet_values.test.values[11], 4)}"
}
`, spreadsheetID, readRange)
}
