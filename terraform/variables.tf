// Terraform variables.
variable "aws_region" {
  type = string
}

# ========== API Gateway ==========
variable "api_gateway_name" {
  type    = string
  default = "pr09-api-gateway"
}

variable "api_gateway_role" {
  type    = string
  default = "pr09-api-gateway-role"
}

variable "api_gateway_stage" {
  type    = string
  default = "dev"
}

variable "authorizer_name" {
  type    = string
  default = "pr09-authorizer"
}

# ========== Lambdas ==========
variable "router_lambda_name" {
  type    = string
  default = "pr09-router-lambda"
}

variable "router_lambda_role" {
  type    = string
  default = "pr09-router-lambda-role"
}

variable "router_lambda_inline_policy" {
  type    = string
  default = "pr09-router-lambda-inline-policy"
}

variable "link_shortener_lambda_name" {
  type    = string
  default = "pr09-link-shortener-lambda"
}

variable "link_shortener_lambda_role" {
  type    = string
  default = "pr09-link-shortener-lambda-role"
}
variable "authorizer_lambda_inline_policy" {
  type    = string
  default = "pr09-authorizer-lambda-inline-policy"
}

variable "authorizer_lambda_name" {
  type    = string
  default = "pr09-authorizer-lambda"
}

variable "authorizer_lambda_role" {
  type    = string
  default = "pr09-authorizer-lambda-role"
}

variable "automation_lambda_name" {
  type    = string
  default = "pr09-automation-lambda"
}

# variable "lambda_iam_role" {
#   type    = string
#   default = "pr09-automation-lambda-role"
# }

variable "lambda_ecr_repo" {
  type    = string
  default = "pr09-lambda-ecr-repo"
}

# ========== Secrets Manager ==========
variable "token_secret_name" {
  type    = string
  default = "pr09-ai-assistant-api-key"
}

variable "token_secret_value" {
  type = string
}

variable "router_lambda_secrets_manager_policy" {
  type    = string
  default = "pr09-router-lambda-secrets-manager-policy"
}

# ========== Dub API ==========
variable "dub_api_key" {
  type = string
}
