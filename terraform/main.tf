terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
  }
}

provider "aws" {
  region = var.aws_region
}

# ========== IAM ==========
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

# Authorizer Lambda Role
resource "aws_iam_role" "authorizer_lambda_role" {
  name = var.authorizer_lambda_role
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = [
            "lambda.amazonaws.com"
            # "apigateway.amazonaws.com"
          ]
        }
      },
    ]
  })

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

# Link Shortener Lambda Role
resource "aws_iam_role" "link_shortener_lambda_role" {
  name = var.link_shortener_lambda_role
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = [
            "lambda.amazonaws.com"
            # "apigateway.amazonaws.com"
          ]
        }
      },
    ]
  })

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

# Test Lambda Role
resource "aws_iam_role" "test_lambda_role" {
  name = "test-lambda-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = [
            "lambda.amazonaws.com"
          ]
        }
      },
    ]
  })

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

# User: arn:aws:sts::874397132032:assumed-role/pr09-router-lambda-role/pr09-router-lambda is not authorized to perform: lambda:InvokeFunction on resource: arn:aws:lambda:eu-central-1:874397132032:function:pr09-link-shortener-lambda-role because no identity-based policy allows the lambda:InvokeFunction action
#         1727951256451

# Router Lambda Inline Policy
resource "aws_iam_role_policy" "router_lambda_inline_policy" {
  name = var.router_lambda_inline_policy
  role = aws_iam_role.router_lambda_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction",
        ]
        Resource = [aws_lambda_function.link_shortener_lambda.arn]
      },
    ],
  })
}

# Authorizer Lambda Inline Policy
resource "aws_iam_role_policy" "authorizer_lambda_inline_policy" {
  name = var.authorizer_lambda_inline_policy
  role = aws_iam_role.authorizer_lambda_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
        ]
        Resource = [aws_secretsmanager_secret.token_secret.arn]
      },
    ],
  })
}

# ========== API Gateway ==========
resource "aws_api_gateway_rest_api" "api_gateway" {
  name = var.api_gateway_name
}

# API Gateway Role
resource "aws_iam_role" "api_gateway_role" {
  name = var.api_gateway_role
  path = "/"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      }
    ]
  })

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole", "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"]
}

# API Gateway Router Lambda permission
resource "aws_lambda_permission" "api_gateway_router_lambda" {
  statement_id  = "AllowAPIGatewayInvokeRouterLambda"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.router_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/* part allows invocation from any stage, method and resource path
  # within API Gateway.
  source_arn = "${aws_api_gateway_rest_api.api_gateway.execution_arn}/*/*"
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
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = aws_api_gateway_resource.assistant_proxy.id
  http_method = "ANY"
  # authorization = "NONE"
  authorization = "CUSTOM"
  authorizer_id = "ecpkve"
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
resource "aws_api_gateway_deployment" "api_gateway_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id

  # Triggers determine when a resource should be recreated
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.assistant.id,
      aws_api_gateway_resource.assistant_proxy.id,
      aws_api_gateway_method.router_any.id,
      aws_api_gateway_integration.router_lambda.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

# API Gateway Stage
resource "aws_api_gateway_stage" "api_stage" {
  deployment_id = aws_api_gateway_deployment.api_gateway_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  stage_name    = var.api_gateway_stage
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_log_group.arn
    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
    })
  }
}


# ========== Lambdas ==========
# Router Lambda
resource "aws_lambda_function" "router_lambda" {
  function_name = var.router_lambda_name
  handler       = "main"
  runtime       = "provided.al2023"
  filename      = "../go/build/router_lambda.zip"
  timeout       = 10

  role = aws_iam_role.router_lambda_role.arn

  environment {
    variables = {
      OPENAI_API_KEY = var.openai_api_key
    }
  }
}

# Authorizer Lambda
resource "aws_lambda_function" "authorizer_lambda" {
  function_name = var.authorizer_lambda_name
  handler       = "main"
  runtime       = "provided.al2023"
  filename      = "../go/build/authorizer_lambda.zip"
  timeout       = 10
  role          = aws_iam_role.authorizer_lambda_role.arn

  environment {
    variables = {
      SECRET_NAME = var.token_secret_name
    }
  }
}

# Link Shortener Lambda
resource "aws_lambda_function" "link_shortener_lambda" {
  function_name = var.link_shortener_lambda_name
  handler       = "main"
  runtime       = "provided.al2023"
  filename      = "../go/build/link_shortener_lambda.zip"
  timeout       = 10
  role          = aws_iam_role.link_shortener_lambda_role.arn

  environment {
    variables = {
      DUB_API_KEY = var.dub_api_key
    }
  }
}

# # Automation Lambda
# resource "aws_lambda_function" "automation_lambda_function" {
#   function_name = var.automation_lambda_name
#   image_uri     = "${aws_ecr_repository.ecr-repo.repository_url}:latest" # todo
#   package_type  = "Image"
#   timeout       = 10

#   role = aws_iam_role.automation_lambda_role.arn
# }

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

# ECR repository
resource "aws_ecr_repository" "ecr-repo" {
  name = var.lambda_ecr_repo
}

# ========== Secrets Manager ==========
# Secrets Manager secret
resource "aws_secretsmanager_secret" "token_secret" {
  name        = var.token_secret_name
  description = "Access token for API authorization"
}

# Secrets Manager secret version
resource "aws_secretsmanager_secret_version" "token_secret_version" {
  secret_id     = aws_secretsmanager_secret.token_secret.id
  secret_string = var.token_secret_value
}

# ========== CloudWatch ==========
# API Gateway CloudWatch
resource "aws_cloudwatch_log_group" "api_gateway_log_group" {
  name = "/aws/apigateway/${var.api_gateway_name}"
}

# Router Lambda CloudWatch
resource "aws_cloudwatch_log_group" "router_lambda_log_group" {
  name = "/aws/lambda/${var.router_lambda_name}"
}

# Authorizer Lambda CloudWatch
resource "aws_cloudwatch_log_group" "authorizer_lambda_log_group" {
  name = "/aws/lambda/${var.authorizer_lambda_name}"
}

# Link Shortener Lambda CloudWatch
resource "aws_cloudwatch_log_group" "link_shortener_lambda_log_group" {
  name = "/aws/lambda/${var.link_shortener_lambda_name}"
}
