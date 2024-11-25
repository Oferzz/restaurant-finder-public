resource "aws_ecr_repository" "restaurant_api" {
  name                 = "restaurant-api"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = {
    Project     = "restaurant-api"
  }
}
