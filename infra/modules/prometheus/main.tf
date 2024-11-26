# terraform {
#   required_version = ">= 1.0"
# }

# // Helm Provider

# data "aws_eks_cluster" "eks" {
#     name = var.eks_name
# }

# data "aws_eks_cluster_auth" "eks" {
#     name = var.eks_name
# }

# provider "helm" {
#   kubernetes {
#     host                   = data.aws_eks_cluster.eks.endpoint
#     cluster_ca_certificate = base64decode(data.aws_eks_cluster.eks.certificate_authority.0.data)
#     token                  = data.aws_eks_cluster_auth.eks.token
#   }
# }

# // Prometheus

# resource "helm_release" "kube-prometheus" {
#   name       = "kube-prometheus-stack"
#   namespace  = "monitoring"
#   create_namespace = true
#   version    = var.kube_prometheus_helm_version
#   repository = "https://prometheus-community.github.io/helm-charts"
#   chart      = "kube-prometheus-stack"
# }