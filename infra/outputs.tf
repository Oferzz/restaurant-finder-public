output "vpc_id" {
  description = "VPC id."
  value       = module.vpc.vpc_id
}

# output "private_subnet_ids" {
#   value = module.vpc.aws_subnet.private[*].id
# }

# output "public_subnet_ids" {
#   value = module.vpc.aws_subnet.public[*].id
# }
