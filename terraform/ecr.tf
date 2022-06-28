resource "aws_ecr_repository" "foo" {
  for_each = toset(var.repos)
  name                 = "anycommerce/${each.key}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}