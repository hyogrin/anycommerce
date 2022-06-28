variable "aws_region" {
  default = "ap-northeast-2"
}

variable "aws_account_id" {
  default = "123456789012"
}

variable "k8s_version" {
  default = "1.22"
}

variable "name" {
  default = "anycommerce"
}

variable "domain_name" {
  default = "dev.sanghwa.people.aws.dev"
}

variable "repos" {
  type = list
  default = ["carts", "order", "product", "recommendation", "user"]
}

variable "map_roles" {
  description = "Additional IAM roles to add to the aws-auth configmap."
  type = list(object({
    rolearn  = string
    username = string
    groups   = list(string)
  }))

  default = [
    {
      rolearn  = "arn:aws:iam::123456:role/Admin"
      username = "Admin"
      groups   = ["system:masters"]
    },
  ]
}

variable "map_users" {
  description = "Additional IAM users to add to the aws-auth configmap."
  type = list(object({
    userarn  = string
    username = string
    groups   = list(string)
  }))

  default = []
}
