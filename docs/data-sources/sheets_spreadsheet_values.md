---
page_title: "gsuite_sheets_spreadsheet_values Data Source - terraform-provider-gsuite"
subcategory: ""
description: |-
  
---

# Data Source `gsuite_sheets_spreadsheet_values`





## Schema

### Required

- **range** (String, Required)
- **spreadsheet_id** (String, Required)

### Optional

- **date_time_render** (String, Optional)
- **id** (String, Optional) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **value_render** (String, Optional)

### Read-only

- **values** (List of List of String, Read-only)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **read** (String, Optional)


