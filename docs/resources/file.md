---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opencga_file Resource - terraform-provider-opencga"
subcategory: ""
description: |-
  
---

# opencga_file (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `path` (String) Directory path, this does not have to be the absolute path if a root is configured. e.g. sample/, /genomes/sample
- `uri` (String) File absolute path (URI), e.g. /genomes/sample/A00001.cram

### Optional

- `study` (String) The `id` of the study this file is associated with.

### Read-Only

- `id` (String) The ID of this resource.


