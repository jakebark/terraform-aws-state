Amazon S3 bucket and DDB table for managing terraform state.

## Inputs

```hcl
module "state" {
  source = "github.com/jakebark/state"
  name   = "name"
}
```

### Optional Inputs

```hcl
module "state" {
  ...
  kms_key               = aws_kms_key.this.arn
  access_logging_bucket = aws_s3_bucket.this.id
}
```

## Outputs

- `module.s3.bucket`
- `module.s3.table` 
