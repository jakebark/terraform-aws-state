output "bucket" {
  value = aws_s3_bucket.this
}

output "table" {
  value = aws_dynamodb_table.this
}
