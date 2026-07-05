terraform {
  required_version = ">= 1.7"
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.0" }
  }
  backend "s3" {
    bucket         = "YOUR-TERRAFORM-STATE-BUCKET"
    key            = "ecommerce/dev/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-state-lock"
    encrypt        = true
  }
}
provider "aws" {
  region = var.region
  default_tags { tags = local.tags }
}
locals {
  tags = { Project = var.project_name, Environment = var.environment, ManagedBy = "terraform" }
}
module "vpc" {
  source = "./modules/vpc"
  project_name = var.project_name; region = var.region; cluster_name = var.cluster_name; tags = local.tags
}
module "eks" {
  source = "./modules/eks"
  cluster_name = var.cluster_name; vpc_id = module.vpc.vpc_id; subnet_ids = module.vpc.private_subnets; tags = local.tags
}
module "ecr" {
  source = "./modules/ecr"
  project_name = var.project_name; tags = local.tags
}
