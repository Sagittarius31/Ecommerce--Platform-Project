module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"
  name = "${var.project_name}-vpc"
  cidr = "10.0.0.0/16"
  azs             = ["${var.region}a","${var.region}b","${var.region}c"]
  private_subnets = ["10.0.1.0/24","10.0.2.0/24","10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24","10.0.102.0/24","10.0.103.0/24"]
  enable_nat_gateway = true; single_nat_gateway = false; enable_dns_hostnames = true
  private_subnet_tags = { "kubernetes.io/role/internal-elb" = "1" }
  public_subnet_tags  = { "kubernetes.io/role/elb" = "1" }
  tags = var.tags
}
variable "project_name" {}; variable "region" {}; variable "cluster_name" {}; variable "tags" { type = map(string) }
output "vpc_id"          { value = module.vpc.vpc_id }
output "private_subnets" { value = module.vpc.private_subnets }
output "public_subnets"  { value = module.vpc.public_subnets }
