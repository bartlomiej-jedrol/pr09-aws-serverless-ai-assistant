// Terraform variables.
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
variable "generic_lambda_role" {
  type    = string
  default = "pr09-lambda-role"
}

variable "router_lambda_name" {
  type    = string
  default = "pr09-router-lambda"
}

variable "authorizer_lambda_name" {
  type    = string
  default = "pr09-authorizer-lambda"
}

variable "automation_lambda_name" {
  type    = string
  default = "pr09-automation-lambda"
}

# variable "lambda_iam_role" {
#   type    = string
#   default = "pr09-automation-lambda-role"
# }

variable "router_lambda_inline_policy" {
  type    = string
  default = "pr09-router-lambda-inline-policy"
}

variable "lambda_ecr_repo" {
  type    = string
  default = "pr09-lambda-ecr-repo"
}