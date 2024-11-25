terraform {
  backend "s3" {
    bucket         = "tr-state-3926106"
    key            = "terraform/state/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
  }
}
