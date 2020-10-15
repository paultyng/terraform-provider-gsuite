package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSheetsSpreadsheetValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSheetsSpreadsheetValuesRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"spreadsheet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value_render": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "UNFORMATTED_VALUE",
				ValidateFunc: validation.StringInSlice([]string{"FORMATTED_VALUE", "UNFORMATTED_VALUE", "FORMULA"}, false),
			},
			"date_time_render": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SERIAL_NUMBER",
				ValidateFunc: validation.StringInSlice([]string{"SERIAL_NUMBER", "FORMATTED_STRING"}, false),
			},

			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func dataSourceSheetsSpreadsheetValuesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config)

	spreadsheetID := d.Get("spreadsheet_id").(string)
	readRange := d.Get("range").(string)
	valueRender := d.Get("value_render").(string)
	dateTimeRender := d.Get("date_time_render").(string)

	req := c.Sheets.Spreadsheets.Values.Get(spreadsheetID, readRange)
	req.Context(ctx)
	req.ValueRenderOption(valueRender)
	req.DateTimeRenderOption(dateTimeRender)

	resp, err := req.Do()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", spreadsheetID, readRange))

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		d.Set("values", [][]interface{}{})
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "No values found in spreadsheet range",
			},
		}
	}

	// normalize the sizing of the returned value
	// TODO: maybe this should be calculated off the range?
	maxCols := 0
	for _, row := range resp.Values {
		if cols := len(row); cols > maxCols {
			maxCols = cols
		}
	}

	values := []interface{}{}
	for _, row := range resp.Values {
		for len(row) < maxCols {
			row = append(row, "")
		}
		values = append(values, row)
	}
	d.Set("values", values)
	return nil
}
