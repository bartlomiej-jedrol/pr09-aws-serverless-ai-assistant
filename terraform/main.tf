terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket  = "bj-terraform-states"
    key     = "state-pr09-aws-serverless-api/terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

# Router Lambda Role
resource "aws_iam_role" "router_lambda_role" {
  name = var.router_lambda_role
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

# # Automation Lambda Role
# resource "aws_iam_role" "automation_lambda_role" {
#   name = var.automation_lambda_role
#   assume_role_policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Action = "sts:AssumeRole"
#         Effect = "Allow"
#         Principal = {
#           Service = "lambda.amazonaws.com"
#         }
#       },
#     ]
#   })

#   managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
# }

# Router Lambda
resource "aws_lambda_function" "router_lambda" {
  function_name = var.router_lambda_name
  handler       = "main"
  runtime       = "provided.al2023"
  filename      = "../go/build/main.zip"
  timeout       = 10

  role = aws_iam_role.router_lambda_role.arn
}

# # Automation Lambda
# resource "aws_lambda_function" "automation_lambda_function" {
#   function_name = var.automation_lambda_name
#   image_uri     = "${aws_ecr_repository.ecr-repo.repository_url}:latest" # todo
#   package_type  = "Image"
#   timeout       = 10

#   role = aws_iam_role.automation_lambda_role.arn
# }

# # Router Lambda Inline Policy to allow invoke automation lambda
# resource "aws_iam_role_policy" "router_lambda_inline_policy" {
#   name = var.lambda_inline_policy
#   role = aws_iam_role.lambda_iam_role.name
#   policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Effect = "Allow"
#         Action = [
#            "lambda:InvokeFunction",
#         ]
#         Resource = [aws_lambda_function.automation_lambda_function.arn]
#       },
#     ],
#   })
# }

# API Gateway
resource "aws_api_gateway_rest_api" "api_gateway" {
  name = var.api_gateway_name
}

# Root assistant endpoint
resource "aws_api_gateway_resource" "assistant" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  parent_id   = aws_api_gateway_rest_api.api_gateway.root_resource_id
  path_part   = "assistant"
}

# Catch all proxy
resource "aws_api_gateway_resource" "assistant_proxy" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  parent_id   = aws_api_gateway_resource.assistant.id
  path_part   = "{proxy+}"
}

# Catch all method
resource "aws_api_gateway_method" "router_any" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = aws_api_gateway_resource.assistant_proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}

# Catch all integration
resource "aws_api_gateway_integration" "router_lambda" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = aws_api_gateway_resource.assistant_proxy.id
  http_method = aws_api_gateway_method.router_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.router_lambda.invoke_arn
}

# API Gateway Deployment
resource "aws_api_gateway_deployment" "router_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  stage_name  = var.api_gateway_stage

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.api_gateway.body))
  }

  depends_on = [
    aws_api_gateway_resource.assistant_proxy,
    aws_api_gateway_method.router_any,
    aws_api_gateway_integration.router_lambda,
  ]
}

# API Gateway CloudWatch
resource "aws_cloudwatch_log_group" "cloud_watch_group" {
  name = "/aws/apigateway/${var.api_gateway_name}"
}

# Router Lambda CloudWatch
resource "aws_cloudwatch_log_group" "router_lambda_log_group" {
  name = "/aws/lambda/${var.router_lambda_name}"
}

# #ECR repository
resource "aws_ecr_repository" "ecr-repo" {
  name = var.lambda_ecr_repo
}
