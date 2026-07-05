variable "project_name"       { default = "ecommerce" }
variable "environment"        { default = "dev" }
variable "region"             { default = "us-east-1" }
variable "cluster_name"       { default = "ecommerce-eks-dev" }
variable "node_instance_type" { default = "t3.medium" }
variable "db_password"        { sensitive = true }
