module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6.0"

  identifier           = "${var.project_name}-postgres"
  engine               = "postgres"
  engine_version       = "16.2"
  instance_class       = var.instance_class
  allocated_storage    = 20
  storage_encrypted    = true
  db_name              = "ecommerce"
  username             = "dbadmin"
  password             = var.db_password
  multi_az             = true
  deletion_protection  = true
  skip_final_snapshot  = false
  backup_retention_period = 7

  vpc_security_group_ids = [var.security_group_id]
  db_subnet_group_name   = aws_db_subnet_group.this.name
  tags = var.tags
}

resource "aws_db_subnet_group" "this" {
  name       = "${var.project_name}-db-subnet"
  subnet_ids = var.subnet_ids
  tags       = var.tags
}

variable "project_name"      {}
variable "instance_class"    { default = "db.t3.medium" }
variable "db_password"       { sensitive = true }
variable "subnet_ids"        { type = list(string) }
variable "security_group_id" {}
variable "tags"              { type = map(string) }

output "endpoint" { value = module.rds.db_instance_endpoint; sensitive = true }
