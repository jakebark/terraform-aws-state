variable "name" {
  type = string
}

variable "kms_key" {
  type    = string
  default = null
}

variable "access_logging_bucket" {
  type    = string
  default = null
}
