provider "aws" {
  alias = "cloudfront"
  region = "us-east-1" # CloudFront expects ACM resources in us-east-1 region only

  # Make it faster by skipping something
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_credentials_validation = true

  # skip_requesting_account_id should be disabled to generate valid ARN in apigatewayv2_api_execution_arn
  skip_requesting_account_id = false
}

locals {
  domain_name = "${var.domain_name}"
  subdomain   = "shop"
}

module "cloudfront" {
  providers = {
    aws = aws.cloudfront
  }

  source = "terraform-aws-modules/cloudfront/aws"

  aliases = ["${local.subdomain}.${local.domain_name}"]

  comment             = "Anycommerce CloudFront"
  enabled             = true
  is_ipv6_enabled     = true
  price_class         = "PriceClass_All"
  retain_on_delete    = false
  wait_for_deployment = false
  web_acl_id          = aws_wafv2_web_acl.waf_cf.arn

  create_origin_access_identity = true
  origin_access_identities = {
    s3_bucket_assets = "Anycommerce CloudFront can access"
  }

  logging_config = {
    bucket = module.log_bucket.s3_bucket_bucket_domain_name
    prefix = "cloudfront"
  }

  default_root_object = "index.html"

  origin = {
    app = {
      domain_name = "app.${local.domain_name}"
      custom_origin_config = {
        http_port              = 80
        https_port             = 443
        origin_protocol_policy = "match-viewer"
        origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
      }

      custom_header = [
        {
          name  = "X-Forwarded-Scheme"
          value = "https"
        },
        {
          name  = "X-Frame-Options"
          value = "SAMEORIGIN"
        },
        {
          name  = "X-Custom-Header"
          value = "awesome-sanghwa-demo-123"
        },
      ]
    }

    s3_assets = {
      domain_name = module.s3_assets.s3_bucket_bucket_regional_domain_name
      s3_origin_config = {
        origin_access_identity = "s3_bucket_assets"
      }
      origin_shield = {
        enabled              = true
        origin_shield_region = var.aws_region
      }
    }
  }

  default_cache_behavior = {
    target_origin_id           = "app"
    viewer_protocol_policy     = "allow-all"

    allowed_methods = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    compress        = false
    query_string    = true
    headers         = ["Authorization"]

    min_ttl                = 0
    default_ttl            = 0
    max_ttl                = 0
  }

  ordered_cache_behavior = [
    {
      path_pattern           = "index.html"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.html"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.css"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.ico"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
        {
      path_pattern           = "*.jpg"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.js"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.map"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.png"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.svg"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.ttf"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    },
    {
      path_pattern           = "*.webmanifest"
      target_origin_id       = "s3_assets"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true

      cache_policy_id = data.aws_cloudfront_cache_policy.cache_policy_optimized.id
      use_forwarded_values = false
    }
  ]

  viewer_certificate = {
    acm_certificate_arn = module.acm.acm_certificate_arn
    ssl_support_method  = "sni-only"
  }

  # geo_restriction = {
  #   restriction_type = "whitelist"
  #   locations        = ["KR"]
  # }
}

######
# ACM
######

module "acm" {
  providers = {
    aws = aws.cloudfront
  } 
  source  = "terraform-aws-modules/acm/aws"
  version = "~> 3.0"

  domain_name               = local.domain_name
  zone_id                   = data.aws_route53_zone.this.id
  subject_alternative_names = ["${local.subdomain}.${local.domain_name}"]
}

#############
# S3 buckets
#############

data "aws_canonical_user_id" "current" {}

module "s3_assets" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 2.0"

  bucket        = "anycommerce-assets-${random_pet.this.id}"
  force_destroy = true

  website = {
    index_document = "index.html"
  }
}

module "log_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 2.0"

  bucket = "logs-${random_pet.this.id}"
  acl    = null
  grant = [{
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
    id          = data.aws_canonical_user_id.current.id
    }, {
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
    id          = "c4c1ede66af53448b93c283ce9448c4ba468c9432aa01d700d3878632f77d2d0"
    # Ref. https://github.com/terraform-providers/terraform-provider-aws/issues/12512
    # Ref. https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html
  }]
  force_destroy = true
}

###########################
# Origin Access Identities
###########################
data "aws_iam_policy_document" "s3_policy" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${module.s3_assets.s3_bucket_arn}/*"]

    principals {
      type        = "AWS"
      identifiers = module.cloudfront.cloudfront_origin_access_identity_iam_arns
    }
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = module.s3_assets.s3_bucket_id
  policy = data.aws_iam_policy_document.s3_policy.json
}

##########
# policy
##########

data "aws_cloudfront_cache_policy" "cache_policy_amplify" {
    provider = aws.cloudfront
    name = "Managed-Amplify"
}

data "aws_cloudfront_cache_policy" "cache_policy_optimized" {
    provider = aws.cloudfront
    name = "Managed-CachingOptimized"
}

##########
# Route53
##########

data "aws_route53_zone" "this" {
  name = local.domain_name
}

module "records" {
  source  = "terraform-aws-modules/route53/aws//modules/records"
  version = "2.0.0" # @todo: revert to "~> 2.0" once 2.1.0 is fixed properly

  zone_id = data.aws_route53_zone.this.zone_id

  records = [
    {
      name = local.subdomain
      type = "A"
      alias = {
        name    = module.cloudfront.cloudfront_distribution_domain_name
        zone_id = module.cloudfront.cloudfront_distribution_hosted_zone_id
      }
    },
  ]
}

########
# Extra
########

resource "random_pet" "this" {
  length = 2
}
