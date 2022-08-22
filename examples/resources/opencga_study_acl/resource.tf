resource "opencga_study_acl" "study_acl" {
  study    = opencga_study.a_cohort.alias
  member   = "user1"
  template = "admin"
}
