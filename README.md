# state
S3 bucket and DDB table for managing terraform state.

## Inputs

```hcl
module "state" {
  source = "https://github.com/jakebark/terraform-aws-state"
  name   = "bucket and table name"
}
```

## Outputs

- `module.s3.bucket`
- `module.s3.table` 
