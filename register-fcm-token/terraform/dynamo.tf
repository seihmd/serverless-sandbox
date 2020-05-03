resource "aws_dynamodb_table" "fcm_dynamo_table" {
  name     = "FcmTable"
  hash_key = "user_id"
  range_key = "os"
  write_capacity = 1
  read_capacity = 1

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "os"
    type = "S"
  }
}