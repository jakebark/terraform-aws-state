Amazon S3 bucket and DDB table for managing terraform state.

## Inputs

```hcl
module "state" {
  source = "github.com/jakebark/state"
  name   = "name"
}
```

## Outputs

- `module.s3.bucket`
- `module.s3.table` 
