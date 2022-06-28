resource "aws_codebuild_project" "x86" {
  for_each = toset(var.repos)

  badge_enabled  = false
  build_timeout  = 360
  name           = "${each.key}-x86"
  queued_timeout = 480
  service_role   = aws_iam_role.codebuild.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-x86_64-standard:2.0"
    image_pull_credentials_type = "CODEBUILD"
    privileged_mode             = true
    type                        = "LINUX_CONTAINER"

    environment_variable {
      name  = "AWS_DEFAULT_REGION"
      value = var.aws_region
    }

    environment_variable {
      name  = "AWS_ACCOUNT_ID"
      value = var.aws_account_id
    }

    environment_variable {
      name  = "IMAGE_REPO_NAME"
      value = "anycommerce/${each.key}"
    }

    environment_variable {
      name  = "IMAGE_TAG"
      value = "latest-amd64"
    }
  }

  logs_config {
    cloudwatch_logs {
      status = "ENABLED"
    }

    s3_logs {
      encryption_disabled = false
      status              = "DISABLED"
    }
  }

  source {
    location        = "https://git-codecommit.ap-northeast-2.amazonaws.com/v1/repos/${each.key}"
    type            = "CODECOMMIT"
  }

  source_version = "refs/heads/main"
}


resource "aws_codebuild_project" "arm64" {
  for_each = toset(var.repos)

  badge_enabled  = false
  build_timeout  = 360
  name           = "${each.key}-arm64"
  queued_timeout = 480
  service_role   = aws_iam_role.codebuild.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-aarch64-standard:2.0"
    image_pull_credentials_type = "CODEBUILD"
    privileged_mode             = true
    type                        = "ARM_CONTAINER"

    environment_variable {
      name  = "AWS_DEFAULT_REGION"
      value = var.aws_region
    }

    environment_variable {
      name  = "AWS_ACCOUNT_ID"
      value = var.aws_account_id
    }

    environment_variable {
      name  = "IMAGE_REPO_NAME"
      value = "anycommerce/${each.key}"
    }

    environment_variable {
      name  = "IMAGE_TAG"
      value = "latest-arm64v8"
    }
  }

  logs_config {
    cloudwatch_logs {
      status = "ENABLED"
    }

    s3_logs {
      encryption_disabled = false
      status              = "DISABLED"
    }
  }

  source {
    location        = "https://git-codecommit.ap-northeast-2.amazonaws.com/v1/repos/${each.key}"
    type            = "CODECOMMIT"
  }

  source_version = "refs/heads/main"
}

resource "aws_codebuild_project" "manifest" {
  for_each = toset(var.repos)

  badge_enabled  = false
  build_timeout  = 360
  name           = "${each.key}-manifest"
  queued_timeout = 480
  service_role   = aws_iam_role.codebuild.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-x86_64-standard:3.0"
    image_pull_credentials_type = "CODEBUILD"
    privileged_mode             = true
    type                        = "LINUX_CONTAINER"

    environment_variable {
      name  = "AWS_DEFAULT_REGION"
      value = var.aws_region
    }

    environment_variable {
      name  = "AWS_ACCOUNT_ID"
      value = var.aws_account_id
    }

    environment_variable {
      name  = "IMAGE_REPO_NAME"
      value = "anycommerce/${each.key}"
    }

    environment_variable {
      name  = "IMAGE_TAG"
      value = "latest"
    }
  }

  logs_config {
    cloudwatch_logs {
      status = "ENABLED"
    }

    s3_logs {
      encryption_disabled = false
      status              = "DISABLED"
    }
  }

  source {
    location        = "https://git-codecommit.ap-northeast-2.amazonaws.com/v1/repos/${each.key}"
    type            = "CODECOMMIT"
    buildspec       = "buildspec-manifest.yml"
  }

  source_version = "refs/heads/main"
}

// roles
resource "aws_iam_role" "codebuild" {
  name = "codebuild-service-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecr" {
  role       = aws_iam_role.codebuild.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
}

resource "aws_iam_role_policy" "codebuild_policy" {
  role = aws_iam_role.codebuild.name

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    },
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "codecommit:*",
        "cloudfront:*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateNetworkInterface",
        "ec2:DescribeDhcpOptions",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DeleteNetworkInterface",
        "ec2:DescribeSubnets",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeVpcs"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
POLICY
}

// This is for frontend

resource "aws_codebuild_project" "web" {
  badge_enabled  = false
  build_timeout  = 360
  name           = "web"
  queued_timeout = 480
  service_role   = aws_iam_role.codebuild.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:5.0"
    image_pull_credentials_type = "CODEBUILD"
    privileged_mode             = true
    type                        = "LINUX_CONTAINER"

    environment_variable {
      name  = "WEB_BUCKET_NAME"
      value = module.s3_assets.s3_bucket_id
    }

    environment_variable {
      name  = "CLOUDFRONT_DIST_ID"
      value = module.cloudfront.cloudfront_distribution_id
    }

    environment_variable {
      name  = "PRODUCTS_SERVICE_URL"
      value = "${local.subdomain}.${local.domain_name}"
    }

    environment_variable {
      name  = "USERS_SERVICE_URL"
      value = "${local.subdomain}.${local.domain_name}"
    }

    environment_variable {
      name  = "CARTS_SERVICE_URL"
      value = "${local.subdomain}.${local.domain_name}"
    }

    environment_variable {
      name  = "ORDERS_SERVICE_URL"
      value = "${local.subdomain}.${local.domain_name}"
    }

    environment_variable {
      name  = "RECOMMENDATIONS_SERVICE_URL"
      value = "${local.subdomain}.${local.domain_name}"
    }

    environment_variable {
      name  = "DEPLOYED_REGION"
      value = var.aws_region
    }

    environment_variable {
      name  = "IMAGE_ROOT_URL"
      value = "/images/"
    }
  }

  logs_config {
    cloudwatch_logs {
      status = "ENABLED"
    }

    s3_logs {
      encryption_disabled = false
      status              = "DISABLED"
    }
  }

  source {
    location        = "https://git-codecommit.ap-northeast-2.amazonaws.com/v1/repos/web"
    type            = "CODECOMMIT"
  }

  source_version = "refs/heads/main"
}