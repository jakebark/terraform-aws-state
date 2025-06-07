Amazon S3 bucket (and optional DDB table) for managing terraform state.

## Inputs

```hcl
module "state" {
  source  = "jakebark/state/aws"
  version = "1.2.0"
}
```

### Optional Inputs

```hcl
module "state" {
  ...
  name                  = "tf-state-${data.aws_caller_identity.current.account_id}"
  ddb_table             = false
  kms_key               = aws_kms_key.this.arn
  access_logging_bucket = aws_s3_bucket.this.id
}
```

## Outputs

- `module.s3.bucket`

## Related Resources

- [jakebark/state/aws](https://registry.terraform.io/modules/jakebark/state/aws/latest)
