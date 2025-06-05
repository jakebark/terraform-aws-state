variable "name" {
  type    = string
  default = null
}

variable "kms_key" {
  type    = string
  default = null
}

variable "access_logging_bucket" {
  type    = string
  default = null
}

variable "ddb_table" {
  type    = bool
  default = false
}
