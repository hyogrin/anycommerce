provider "aws" {
  region = local.region
}

provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

  exec {
    api_version = "client.authentication.k8s.io/v1alpha1"
    command     = "aws"
    # This requires the awscli to be installed locally where Terraform is executed
    args = ["eks", "get-token", "--cluster-name", module.eks.cluster_id]
  }
}

locals {
  name            = "anycommerce"
  cluster_version = "1.22"
  region          = "ap-northeast-2"

  tags = {
    Cluster    = local.name
  }
}

data "aws_route53_zone" "apne2" {
  name = local.domain_name
}

module "acm_app" {
  source  = "terraform-aws-modules/acm/aws"
  version = "~> 3.0"

  domain_name               = local.domain_name
  zone_id                   = data.aws_route53_zone.apne2.id
  subject_alternative_names = ["*.${local.domain_name}"]
}

data "aws_caller_identity" "current" {}

################################################################################
# EKS Module
################################################################################

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 18.0"

  cluster_name                    = local.name
  cluster_version                 = local.cluster_version
  cluster_endpoint_private_access = true
  cluster_endpoint_public_access  = true

  enable_irsa     = true

  cluster_addons = {
    coredns = {
      resolve_conflicts = "OVERWRITE"
    }
    kube-proxy = {}
    vpc-cni = {
      resolve_conflicts        = "OVERWRITE"
    }
  }

  cluster_encryption_config = [{
    provider_key_arn = aws_kms_key.eks.arn
    resources        = ["secrets"]
  }]

  cluster_tags = {
    # This should not affect the name of the cluster primary security group
    # Ref: https://github.com/terraform-aws-modules/terraform-aws-eks/pull/2006
    # Ref: https://github.com/terraform-aws-modules/terraform-aws-eks/pull/2008
    Name = local.name
  }

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  manage_aws_auth_configmap = true
  # aws_auth_users    = var.map_users
  aws_auth_roles    = var.map_roles
  # aws_auth_accounts = var.map_accounts

  # Extend cluster security group rules
  cluster_security_group_additional_rules = {
    egress_nodes_ephemeral_ports_tcp = {
      description                = "To node 1025-65535"
      protocol                   = "tcp"
      from_port                  = 1025
      to_port                    = 65535
      type                       = "egress"
      source_node_security_group = true
    }
  }

  # Extend node-to-node security group rules
  node_security_group_additional_rules = {
    ingress_self_all = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "ingress"
      self        = true
    }
    egress_all = {
      description      = "Node all egress"
      protocol         = "-1"
      from_port        = 0
      to_port          = 0
      type             = "egress"
      cidr_blocks      = ["0.0.0.0/0"]
    }
    ingress_allow_access_from_control_plane = {
      type                          = "ingress"
      protocol                      = "tcp"
      from_port                     = 9443
      to_port                       = 9443
      source_cluster_security_group = true
      description                   = "Allow access from control plane to webhook port of AWS load balancer controller"
    }
    ingress_allow_webhook_access_from_control_plane = {
      type                          = "ingress"
      protocol                      = "tcp"
      from_port                     = 8443
      to_port                       = 8443
      source_cluster_security_group = true
      description                   = "Allow access from control plane to webhook port of custom webhook controller"
    }
  }

  eks_managed_node_group_defaults = {
    ami_type       = "AL2_x86_64"
    instance_types = ["m5.large"]

    # We are using the IRSA created below for permissions
    # However, we have to deploy with the policy attached FIRST (when creating a fresh cluster)
    # and then turn this off after the cluster/node group is created. Without this initial policy,
    # the VPC CNI fails to assign IPs and nodes cannot join the cluster
    # See https://github.com/aws/containers-roadmap/issues/1666 for more context
    iam_role_attach_cni_policy = true
  }

  eks_managed_node_groups = {
    graviton = {
      name = "graviton"

      ami_type = "AL2_ARM_64"

      min_size     = 1
      max_size     = 8
      desired_size = 2

      instance_types = ["m6g.large"]

      iam_role_additional_policies = [
        "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
      ]

      labels = {
        "node-role" = "graviton"
      }
    }

    amd64 = {
      name = "amd64"

      min_size     = 1
      max_size     = 8
      desired_size = 2

      iam_role_additional_policies = [
        "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
      ]

      labels = {
        "node-role" = "amd64"
      }
    }
  }

  tags = local.tags
}

resource "aws_iam_policy" "node_additional" {
  name        = "${local.name}-additional"
  description = "Example usage of node additional policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
          "Effect": "Allow",
          "Action": [
              "dynamodb:*",
              "personalize:*",
              "route53:*"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "iam:CreateServiceLinkedRole"
          ],
          "Resource": "*",
          "Condition": {
              "StringEquals": {
                  "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "ec2:DescribeAccountAttributes",
              "ec2:DescribeAddresses",
              "ec2:DescribeAvailabilityZones",
              "ec2:DescribeInternetGateways",
              "ec2:DescribeVpcs",
              "ec2:DescribeVpcPeeringConnections",
              "ec2:DescribeSubnets",
              "ec2:DescribeSecurityGroups",
              "ec2:DescribeInstances",
              "ec2:DescribeNetworkInterfaces",
              "ec2:DescribeTags",
              "ec2:GetCoipPoolUsage",
              "ec2:DescribeCoipPools",
              "elasticloadbalancing:DescribeLoadBalancers",
              "elasticloadbalancing:DescribeLoadBalancerAttributes",
              "elasticloadbalancing:DescribeListeners",
              "elasticloadbalancing:DescribeListenerCertificates",
              "elasticloadbalancing:DescribeSSLPolicies",
              "elasticloadbalancing:DescribeRules",
              "elasticloadbalancing:DescribeTargetGroups",
              "elasticloadbalancing:DescribeTargetGroupAttributes",
              "elasticloadbalancing:DescribeTargetHealth",
              "elasticloadbalancing:DescribeTags"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "cognito-idp:DescribeUserPoolClient",
              "acm:ListCertificates",
              "acm:DescribeCertificate",
              "iam:ListServerCertificates",
              "iam:GetServerCertificate",
              "waf-regional:GetWebACL",
              "waf-regional:GetWebACLForResource",
              "waf-regional:AssociateWebACL",
              "waf-regional:DisassociateWebACL",
              "wafv2:GetWebACL",
              "wafv2:GetWebACLForResource",
              "wafv2:AssociateWebACL",
              "wafv2:DisassociateWebACL",
              "shield:GetSubscriptionState",
              "shield:DescribeProtection",
              "shield:CreateProtection",
              "shield:DeleteProtection"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "ec2:AuthorizeSecurityGroupIngress",
              "ec2:RevokeSecurityGroupIngress"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "ec2:CreateSecurityGroup"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "ec2:CreateTags"
          ],
          "Resource": "arn:aws:ec2:*:*:security-group/*",
          "Condition": {
              "StringEquals": {
                  "ec2:CreateAction": "CreateSecurityGroup"
              },
              "Null": {
                  "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "ec2:CreateTags",
              "ec2:DeleteTags"
          ],
          "Resource": "arn:aws:ec2:*:*:security-group/*",
          "Condition": {
              "Null": {
                  "aws:RequestTag/elbv2.k8s.aws/cluster": "true",
                  "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "ec2:AuthorizeSecurityGroupIngress",
              "ec2:RevokeSecurityGroupIngress",
              "ec2:DeleteSecurityGroup"
          ],
          "Resource": "*",
          "Condition": {
              "Null": {
                  "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:CreateLoadBalancer",
              "elasticloadbalancing:CreateTargetGroup"
          ],
          "Resource": "*",
          "Condition": {
              "Null": {
                  "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:CreateListener",
              "elasticloadbalancing:DeleteListener",
              "elasticloadbalancing:CreateRule",
              "elasticloadbalancing:DeleteRule"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:AddTags",
              "elasticloadbalancing:RemoveTags"
          ],
          "Resource": [
              "arn:aws:elasticloadbalancing:*:*:targetgroup/*/*",
              "arn:aws:elasticloadbalancing:*:*:loadbalancer/net/*/*",
              "arn:aws:elasticloadbalancing:*:*:loadbalancer/app/*/*"
          ],
          "Condition": {
              "Null": {
                  "aws:RequestTag/elbv2.k8s.aws/cluster": "true",
                  "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:AddTags",
              "elasticloadbalancing:RemoveTags"
          ],
          "Resource": [
              "arn:aws:elasticloadbalancing:*:*:listener/net/*/*/*",
              "arn:aws:elasticloadbalancing:*:*:listener/app/*/*/*",
              "arn:aws:elasticloadbalancing:*:*:listener-rule/net/*/*/*",
              "arn:aws:elasticloadbalancing:*:*:listener-rule/app/*/*/*"
          ]
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:ModifyLoadBalancerAttributes",
              "elasticloadbalancing:SetIpAddressType",
              "elasticloadbalancing:SetSecurityGroups",
              "elasticloadbalancing:SetSubnets",
              "elasticloadbalancing:DeleteLoadBalancer",
              "elasticloadbalancing:ModifyTargetGroup",
              "elasticloadbalancing:ModifyTargetGroupAttributes",
              "elasticloadbalancing:DeleteTargetGroup"
          ],
          "Resource": "*",
          "Condition": {
              "Null": {
                  "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
              }
          }
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:RegisterTargets",
              "elasticloadbalancing:DeregisterTargets"
          ],
          "Resource": "arn:aws:elasticloadbalancing:*:*:targetgroup/*/*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "elasticloadbalancing:SetWebAcl",
              "elasticloadbalancing:ModifyListener",
              "elasticloadbalancing:AddListenerCertificates",
              "elasticloadbalancing:RemoveListenerCertificates",
              "elasticloadbalancing:ModifyRule"
          ],
          "Resource": "*"
      }
    ]
  })

  tags = local.tags
}


# References to resources that do not exist yet when creating a cluster will cause a plan failure due to https://github.com/hashicorp/terraform/issues/4149
# There are two options users can take
# 1. Create the dependent resources before the cluster => `terraform apply -target <your policy or your security group> and then `terraform apply`
#   Note: this is the route users will have to take for adding additonal security groups to nodes since there isn't a separate "security group attachment" resource
# 2. For addtional IAM policies, users can attach the policies outside of the cluster definition as demonstrated below
resource "aws_iam_role_policy_attachment" "additional" {
  for_each = module.eks.eks_managed_node_groups

  policy_arn = aws_iam_policy.node_additional.arn
  role       = each.value.iam_role_name
}

################################################################################
# Supporting Resources
################################################################################

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  name = local.name
  cidr = "10.0.0.0/16"

  azs             = ["${local.region}a", "${local.region}b", "${local.region}c", "${local.region}d"]
  private_subnets      = ["10.0.0.0/20", "10.0.16.0/20", "10.0.32.0/20", "10.0.48.0/20"]
  public_subnets       = ["10.0.64.0/20", "10.0.80.0/20", "10.0.96.0/20", "10.0.112.0/20"]
  
  enable_ipv6                     = false
  assign_ipv6_address_on_creation = false
  create_egress_only_igw          = true

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  enable_flow_log                      = true
  create_flow_log_cloudwatch_iam_role  = true
  create_flow_log_cloudwatch_log_group = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.name}" = "shared"
    "kubernetes.io/role/elb"              = 1
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.name}" = "shared"
    "kubernetes.io/role/internal-elb"     = 1
  }

  tags = local.tags
}

resource "aws_kms_key" "eks" {
  description             = "EKS Secret Encryption Key"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = local.tags
}

resource "aws_kms_key" "ebs" {
  description             = "Customer managed key to encrypt EKS managed node group volumes"
  deletion_window_in_days = 7
  policy                  = data.aws_iam_policy_document.ebs.json
}

# This policy is required for the KMS key used for EKS root volumes, so the cluster is allowed to enc/dec/attach encrypted EBS volumes
data "aws_iam_policy_document" "ebs" {
  # Copy of default KMS policy that lets you manage it
  statement {
    sid       = "Enable IAM User Permissions"
    actions   = ["kms:*"]
    resources = ["*"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
  }

  # Required for EKS
  statement {
    sid = "Allow service-linked role use of the CMK"
    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling", # required for the ASG to manage encrypted volumes for nodes
        module.eks.cluster_iam_role_arn,                                                                                                            # required for the cluster / persistentvolume-controller to create encrypted PVCs
      ]
    }
  }

  statement {
    sid       = "Allow attachment of persistent resources"
    actions   = ["kms:CreateGrant"]
    resources = ["*"]

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling", # required for the ASG to manage encrypted volumes for nodes
        module.eks.cluster_iam_role_arn,                                                                                                            # required for the cluster / persistentvolume-controller to create encrypted PVCs
      ]
    }

    condition {
      test     = "Bool"
      variable = "kms:GrantIsForAWSResource"
      values   = ["true"]
    }
  }
}


data "aws_ami" "eks_default" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amazon-eks-node-${local.cluster_version}-v*"]
  }
}

data "aws_ami" "eks_default_arm" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amazon-eks-arm64-node-${local.cluster_version}-v*"]
  }
}

data "aws_ami" "eks_default_bottlerocket" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["bottlerocket-aws-k8s-${local.cluster_version}-x86_64-*"]
  }
}
