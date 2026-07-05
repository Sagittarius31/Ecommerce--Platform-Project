module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = var.cluster_name
  cluster_version = "1.30"
  vpc_id          = var.vpc_id
  subnet_ids      = var.subnet_ids
  enable_irsa     = true
  cluster_endpoint_public_access = true
  enable_cluster_creator_admin_permissions = true

  eks_managed_node_groups = {
    general = {
      instance_types = [var.node_instance_type]
      min_size       = 2
      max_size       = 10
      desired_size   = 3
    }
  }

  tags = var.tags
}

variable "cluster_name"        {}
variable "vpc_id"              {}
variable "subnet_ids"          { type = list(string) }
variable "node_instance_type"  { default = "t3.medium" }
variable "tags"                { type = map(string) }

output "cluster_name"     { value = module.eks.cluster_name }
output "cluster_endpoint" { value = module.eks.cluster_endpoint; sensitive = true }
output "oidc_provider_arn" { value = module.eks.oidc_provider_arn }
