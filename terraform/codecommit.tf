resource "aws_codecommit_repository" "app" {
  for_each = toset(var.repos)

  repository_name = "${each.key}"
  description     = "This is an ${each.key} repository"
  default_branch  = "main"
}

resource "aws_codecommit_repository" "web" {
  repository_name = "web"
  description     = "This is the web repository"
  default_branch  = "main"
}

resource "aws_codecommit_repository" "generator" {
  repository_name = "generator"
  description     = "This is the generator repository"
  default_branch  = "main"
}