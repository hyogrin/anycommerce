module "s3_applicaion_logs" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 2.0"

  bucket        = "anycommerce-logs-${random_pet.this.id}"
  force_destroy = true
}

resource "aws_kinesis_firehose_delivery_stream" "anycommerce_stream" {
  name        = "anycommerce-stream"
  destination = "s3"

  s3_configuration {
    role_arn            = aws_iam_role.firehose_role.arn
    bucket_arn          = module.s3_applicaion_logs.s3_bucket_arn
    buffer_size         = 5
    buffer_interval     = 300
    compression_format  = "GZIP"
    prefix              = "raw/year=!{timestamp:YYYY}/month=!{timestamp:MM}/day=!{timestamp:dd}/hour=!{timestamp:HH}/"
    error_output_prefix = "!{firehose:error-output-type}/!{timestamp:yyyy/MM/dd/HH}/"
  }

  # elasticsearch_configuration {
  #   domain_arn = "arn:aws:es:ap-northeast-2:${var.aws_account_id}:domain/anycommerce-logging"
  #   role_arn   = aws_iam_role.firehose_role.arn
  #   index_name = "kinesis"

  #   s3_backup_mode = "AllDocuments"
  #   buffering_interval = 60
    # processing_configuration {
    #   enabled = "true"

    #   processors {
    #     type = "Lambda"

    #     parameters {
    #       parameter_name  = "LambdaArn"
    #       parameter_value = "${aws_lambda_function.lambda_processor.arn}:$LATEST"
    #     }
    #   }
    # }
  # }
}


resource "aws_iam_role" "firehose_role" {
  name = "firehose-service-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "firehose-elasticsearch" {
  role   = aws_iam_role.firehose_role.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "es:*"
      ],
      "Resource": [
        "arn:aws:es:ap-northeast-2:${var.aws_account_id}:domain/anycommerce-logging",
        "arn:aws:es:ap-northeast-2:${var.aws_account_id}:domain/anycommerce-logging/*"
      ]
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
EOF
}
