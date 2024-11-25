provider "aws" {
  region = var.region
}

// VPC

module "vpc" {
  source = "./modules/vpc"

  azs             = ["us-east-1a", "us-east-1b"]
  private_subnets = ["10.0.0.0/19", "10.0.32.0/19"]
  public_subnets  = ["10.0.64.0/19", "10.0.96.0/19"]

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
    "kubernetes.io/cluster/dev-demo"  = "owned"
  }

  public_subnet_tags = {
    "kubernetes.io/role/elb"         = 1
    "kubernetes.io/cluster/dev-demo" = "owned"
  }
}

// EKS cluster

module "eks" {
  source = "./modules/eks"

  eks_version = "1.30"
  eks_name    = "demo"
  subnet_ids  = module.vpc.private_subnet_ids

  node_groups = {
    general = {
      capacity_type  = "ON_DEMAND"
      instance_types = ["t3a.medium"]
      scaling_config = {
        desired_size = 3
        max_size     = 5
        min_size     = 0
      }
    }
  }
}

// DynamoDB

module "restaurants_table" {
  source         = "./modules/dynamodb"
  table_name     = "restaurants"
  billing_mode   = "PROVISIONED"
  hash_key       = "restaurant_id"
  hash_key_type  = "S"
  read_capacity  = 10
  write_capacity = 5
}

module "audit_logs_table" {
  source         = "./modules/dynamodb"
  table_name     = "audit_logs"
  billing_mode   = "PROVISIONED"
  hash_key       = "timestamp"
  hash_key_type  = "S"
  read_capacity  = 10
  write_capacity = 5
}

// EBS-CSI driver

module "ebs-csi" {
  source = "./modules/ebs-csi"

  eks_name = module.eks.eks_name
}

// ECR Repository

module "ecr" {
  source = "./modules/ecr"
}

// Prometheus

module "prometheus"{
  source = "./modules/prometheus"

  kube_prometheus_helm_version = "66.2.1"
  eks_name = module.eks.eks_name
}
