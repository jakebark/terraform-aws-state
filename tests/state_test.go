package tests

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name                  string
	vars                  map[string]interface{}
	requiresKmsKey        bool
	requiresAccessLogging bool
	assertions            func(t *testing.T, opts *terraform.Options, awsRegion string)
}

func TestState(t *testing.T) {
	t.Parallel()

	awsRegion := "us-east-1"
	uniqueID := random.UniqueId()
	bucketName := fmt.Sprintf("tf-state-test-%s", uniqueID)

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"name": bucketName,
		},
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	actualBucketName := terraform.Output(t, terraformOptions, "bucket_name")
	assertS3Bucket(t, awsRegion, actualBucketName, false, "")
}

func assertS3Bucket(t *testing.T, awsRegion, bucketName string, expectLogging bool, kmsKeyArn string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	assert.NoError(t, err)
	s3Client := s3.New(sess)

	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(bucketName)})
	assert.NoError(t, err, "S3 bucket not found")

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
		assert.NotNil(t, logging.LoggingEnabled, "Access logging should be enabled")
	} else {
		assert.Nil(t, logging.LoggingEnabled, "Access logging should be disabled")
	}
}
