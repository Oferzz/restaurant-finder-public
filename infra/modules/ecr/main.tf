resource "aws_ecr_repository" "restaurant_finder" {
  name                 = "restaurant-finder"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = {
    Project     = "restaurant-finder"
  }
}
