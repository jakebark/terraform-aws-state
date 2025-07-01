package tests

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name                  string
	vars                  map[string]interface{}
	requiresKmsKey        bool
	requiresAccessLogging bool
	assertions            func(t *testing.T, opts *terraform.Options, awsRegion string)
}

func TestTerraformAwsState(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name: "TestDefaults",
			vars: map[string]interface{}{"ddb_table": false},
			assertions: func(t *testing.T, opts *terraform.Options, awsRegion string) {
				bucketName := terraform.Output(t, opts, "bucket_name")
				assertS3Bucket(t, awsRegion, bucketName, "", false)
			},
		},
		{
			name:           "TestKMS",
			vars:           map[string]interface{}{"ddb_table": false},
			requiresKmsKey: true,
			assertions: func(t *testing.T, opts *terraform.Options, awsRegion string) {
				bucketName := terraform.Output(t, opts, "bucket_name")
				kmsKeyArn := opts.Vars["kms_key"].(string)
				assertS3Bucket(t, awsRegion, bucketName, kmsKeyArn, false)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			awsRegion := "us-east-1"

			sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
			require.NoError(t, err)
			kmsClient := kms.New(sess)
			s3Client := s3.New(sess)

			vars := make(map[string]interface{})
			for k, v := range tc.vars {
				vars[k] = v
			}
			vars["name"] = fmt.Sprintf("tf-state-test-%s", random.UniqueId())

			if tc.requiresKmsKey {
				key, err := kmsClient.CreateKey(&kms.CreateKeyInput{})
				require.NoError(t, err)
				defer kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
					KeyId:               key.KeyMetadata.KeyId,
					PendingWindowInDays: aws.Int64(7),
				})
				vars["kms_key"] = *key.KeyMetadata.Arn
			}

			if tc.requiresAccessLogging {
				loggingBucketName := fmt.Sprintf("terratest-access-logs-%s", random.UniqueId())
				_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
					Bucket: aws.String(loggingBucketName),
				})
				require.NoError(t, err)
				defer s3Client.DeleteBucket(&s3.DeleteBucketInput{
					Bucket: aws.String(loggingBucketName),
				})
				vars["access_logging_bucket"] = loggingBucketName
			}

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../",
				Vars:         vars,
				EnvVars: map[string]string{
					"AWS_DEFAULT_REGION": awsRegion,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApply(t, terraformOptions)
			tc.assertions(t, terraformOptions, awsRegion)
		})
	}
}

func assertS3Bucket(t *testing.T, awsRegion, bucketName, kmsKeyArn string, expectLogging bool) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	assert.NoError(t, err)
	s3Client := s3.New(sess)

	encryption, err := s3Client.GetBucketEncryption(&s3.GetBucketEncryptionInput{Bucket: aws.String(bucketName)})
	assert.NoError(t, err)
	rule := encryption.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault
	if kmsKeyArn != "" {
		assert.Equal(t, "aws:kms", *rule.SSEAlgorithm)
		assert.Equal(t, kmsKeyArn, *rule.KMSMasterKeyID)
	} else {
		assert.Equal(t, "AES256", *rule.SSEAlgorithm)
	}

	logging, err := s3Client.GetBucketLogging(&s3.GetBucketLoggingInput{Bucket: aws.String(bucketName)})
	assert.NoError(t, err)
	if expectLogging {
		assert.NotNil(t, logging.LoggingEnabled)
	} else {
		assert.Nil(t, logging.LoggingEnabled)
	}
}
