terraform {
  required_version = "v0.12.20"
  experiments      = [
    variable_validation]
}

provider "aws" {
  version = "~> 2.0"
  region  = "ap-northeast-1"
}

variable "env" {
  type    = string
  default = "dev"

  validation {
    condition     = contains([
      "dev",
      "stg",
      "prod"], var.env)
    error_message = "The env variable is wrong."
  }
}

# https://www.terraform.io/docs/providers/aws/r/sns_platform_application.html
resource "aws_sns_platform_application" "fcm_application" {
  name                = "fcm_application"
  platform            = "GCM"
  platform_credential = data.aws_ssm_parameter.fcm_serverkey_value.value
  depends_on          = [
    data.aws_ssm_parameter.fcm_serverkey_value]
}

# aws ssm putparameter name '/fcm/serverkey' type SecureString value 'TODO' overwrite
resource "aws_ssm_parameter" "fcm_serverkey" {
  name        = "/fcm/${var.env}/serverkey"
  type        = "SecureString"
  value       = "UNINITIALIZED"
  description = "Firebase CloudMessaging"

  lifecycle {
    ignore_changes = [
      value]
  }
}

resource "aws_ssm_parameter" "fcm_dynamodb_table_name" {
  name  = "/fcm/${var.env}/dynamodb/table/name"
  type  = "String"
  value = aws_dynamodb_table.fcm_dynamo_table.name
}

resource "aws_ssm_parameter" "fcm_dynamodb_table_arn" {
  name  = "/fcm/${var.env}/dynamodb/table/arn"
  type  = "String"
  value = aws_dynamodb_table.fcm_dynamo_table.arn
}

resource "aws_ssm_parameter" "fcm_sns_arn" {
  name  = "/fcm/${var.env}/sns/arn"
  type  = "String"
  value = aws_sns_platform_application.fcm_application.arn
}

data "aws_ssm_parameter" "fcm_serverkey_value" {
  name       = "/fcm/${var.env}/serverkey"
  depends_on = [
    aws_ssm_parameter.fcm_serverkey]
}