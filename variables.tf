variable "access_logging_bucket" {
  description = "s3 server access logging bucket name"
  type        = string
  default     = null
}

variable "enable_access_logging" {
  description = "enable s3 server access logging"
  type        = bool
  default     = false
}

variable "ddb_table" {
  description = "enable DynamoDB table"
  type        = bool
  default     = false
}

variable "kms_key" {
  description = "AWS KMS key arn"
  type        = string
  default     = null
}

variable "name" {
  description = "bucket name"
  type        = string
  default     = null
}

