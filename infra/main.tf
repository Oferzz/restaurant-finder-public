provider "aws" {
  region = var.region
}

data "external" "aws_account_id" {
  program = [
    "sh", "-c",
    "aws sts get-caller-identity --query 'Account' --output text | jq -n --arg account $(cat) '{result: $account}'"
  ]
}

// VPC

module "vpc" {
  source = "./modules/vpc"

  azs             = ["us-east-1a", "us-east-1b"]
  private_subnets = ["10.0.0.0/19", "10.0.32.0/19"]
  public_subnets  = ["10.0.64.0/19", "10.0.96.0/19"]

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
    "kubernetes.io/cluster/restaurant-finder-cluster"  = "owned"
  }

  public_subnet_tags = {
    "kubernetes.io/role/elb"         = 1
    "kubernetes.io/cluster/restaurant-finder-cluster" = "owned"
  }
}

// EKS cluster

module "eks" {
  source = "./modules/eks"

  eks_version = "1.30"
  eks_name    = "restaurant-finder-cluster"
  subnet_ids  = module.vpc.private_subnet_ids
  account_id = data.external.aws_account_id.result["result"]
  enable_irsa = true

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

// ECR Repository

module "ecr" {
  source = "./modules/ecr"
}
