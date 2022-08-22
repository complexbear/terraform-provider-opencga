resource "opencga_study" "a_cohort" {
  project = data.opencga_project.a_project.id

  name        = "Germline Study"
  alias       = "GS"
  description = "Example study"
}
