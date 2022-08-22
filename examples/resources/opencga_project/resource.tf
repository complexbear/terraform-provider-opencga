resource "opencga_project" "new_project" {
  name        = "TestProj"
  alias       = "TP"
  description = "A test project"

  scientific_name = "Homo Sapiens"
  taxonomy_code   = 9606
  assembly        = "GRCh38"
}
