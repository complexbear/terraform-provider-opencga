terraform {
  required_providers {
    opencga = {
      source = "complexbear/opencga"
    }
  }
}

provider "opencga" {
  username = "user"
  password = "password"
  base_url = "https://localhost"
}