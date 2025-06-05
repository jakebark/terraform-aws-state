locals {
  name = "tf-state-${data.aws_caller_identity.current.account_id}"
}
