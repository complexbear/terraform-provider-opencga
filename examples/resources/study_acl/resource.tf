resource "opencga_study_acl" "lss_study_acl" {
  study    = opencga_study.a_study.alias
  member   = "user1"
  template = "admin"
}
