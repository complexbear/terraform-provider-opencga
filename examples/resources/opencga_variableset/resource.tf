resource "opencga_variableset" "new_var_set" {
  study       = "NS"
  name        = "New Var Set 2"
  description = "Another new variable set"
  unique      = true
  variables   = file("sample.json")
}
