output "ecr_repository_uri" {
  value = aws_ecr_repository.restaurant_finder.repository_url
  description = "The URI of the ECR repository"
}
