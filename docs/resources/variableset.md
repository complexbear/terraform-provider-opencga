---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opencga_variableset Resource - terraform-provider-opencga"
subcategory: ""
description: |-
  
---

# opencga_variableset (Resource)



## Example Usage

```terraform
resource "opencga_variableset" "new_var_set" {
  study       = "NS"
  name        = "New Var Set 2"
  description = "Another new variable set"
  unique      = true
  variables   = file("sample.json")
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Description, can be left blank
- `name` (String) Variable Set name
- `unique` (Boolean) True if there can only be 1 instance of this attached to a record item. False to allow for multiple instances.
- `variables` (String) Json content representing the variables in this variable set. Json definitions can be read directly from the GelReportModels repo.

### Optional

- `check_description` (Boolean) If true the description content will be checked against the state
- `study` (String) The study that this variable set belongs to

### Read-Only

- `id` (String) The ID of this resource.


