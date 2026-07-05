locals {
  services = ["api-gateway","user-service","product-service","order-service","payment-service","notification-service"]
}
resource "aws_ecr_repository" "services" {
  for_each = toset(local.services)
  name = "${var.project_name}/${each.key}"; image_tag_mutability = "MUTABLE"
  image_scanning_configuration { scan_on_push = true }
  tags = var.tags
}
resource "aws_ecr_lifecycle_policy" "keep10" {
  for_each = aws_ecr_repository.services; repository = each.value.name
  policy = jsonencode({ rules = [{ rulePriority=1, selection={tagStatus="any",countType="imageCountMoreThan",countNumber=10}, action={type="expire"} }] })
}
variable "project_name" {}; variable "tags" { type = map(string) }
output "repository_urls" { value = { for name,repo in aws_ecr_repository.services : name => repo.repository_url } }
