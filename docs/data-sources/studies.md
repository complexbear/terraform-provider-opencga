---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opencga_studies Data Source - terraform-provider-opencga"
subcategory: ""
description: |-
  Use this data source to list several existing Studies for use in other resources
---

# opencga_studies (Data Source)

Use this data source to list several existing Studies for use in other resources

## Example Usage

```terraform
data "opencga_studies" "all" {
  project = "1000000001"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `alias_filter` (String) A study alias to limit the search
- `id_filter` (Number) A study id to limit the search
- `project` (String) A project id or alias to limit the search

### Read-Only

- `id` (String) The ID of this resource.
- `studies` (List of Object) (see [below for nested schema](#nestedatt--studies))

<a id="nestedatt--studies"></a>
### Nested Schema for `studies`

Read-Only:

- `alias` (String)
- `description` (String)
- `id` (Number)
- `name` (String)


