// Terraform variables.
variable "lambda_iam_role" {
  type    = string
  default = "pr09-automation-lambda-role"
}

variable "router_lambda_role" {
  type    = string
  default = "pr09-router-lambda-role"
}

variable "lambda_inline_policy" {
  type    = string
  default = "pr09-router-lambda-inline-policy"
}

variable "router_lambda_name" {
  type    = string
  default = "pr09-router-lambda"
}

variable "automation_lambda_name" {
  type    = string
  default = "pr09-automation-lambda"
}

variable "api_gateway_name" {
  type    = string
  default = "pr09-api-gateway"
}
