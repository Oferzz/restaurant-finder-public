output "eks_name" {
  value = aws_eks_cluster.this.name
}

output "openid_provider_arn" {
  value = aws_iam_openid_connect_provider.this[0].arn
}

output "openid_provider_url" {
  value = aws_iam_openid_connect_provider.this[0].url
}

output "nodes_iam_role_name" {
  value = aws_iam_role.nodes.name
}
