package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DataExchangeFormatTestStackProps struct {
	awscdk.StackProps
}

func NewDataExchangeFormatTestStackStack(scope constructs.Construct, id string, props *DataExchangeFormatTestStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := awsec2.NewVpc(stack, jsii.String("DataExchangeFormatTestVpc"), nil)

	// Create a new security group
	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("DataExchangeFormatTestSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:         vpc,
		Description: jsii.String("Security group for the Data Exchange Format Test"),
	})

	// Add an inbound rule to allow TCP traffic on port 80 from any IP
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP traffic"), jsii.Bool(false))

	awsec2.NewInstance(stack, jsii.String("DataExchangeFormatTestServer"), &awsec2.InstanceProps{
		Vpc:          vpc,
		InstanceType: awsec2.NewInstanceType(jsii.String("t2.micro")),
		MachineImage: awsec2.NewAmazonLinuxImage(&awsec2.AmazonLinuxImageProps{}),
	})

	awsec2.NewInstance(stack, jsii.String("DataExchangeFormatTestGoClient"), &awsec2.InstanceProps{
		Vpc:          vpc,
		InstanceType: awsec2.NewInstanceType(jsii.String("t2.micro")),
		MachineImage: awsec2.NewAmazonLinuxImage(&awsec2.AmazonLinuxImageProps{}),
	})

	awsec2.NewInstance(stack, jsii.String("DataExchangeFormatTestJsClient"), &awsec2.InstanceProps{
		Vpc:          vpc,
		InstanceType: awsec2.NewInstanceType(jsii.String("t2.micro")),
		MachineImage: awsec2.NewAmazonLinuxImage(&awsec2.AmazonLinuxImageProps{}),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewDataExchangeFormatTestStackStack(app, "DataExchangeFormatTestStack", &DataExchangeFormatTestStackProps{
		awscdk.StackProps{
			Env: env(),
			// Eventually I want to be able to reuse an S3 bucket
			// Synthesizer: awscdk.NewDefaultStackSynthesizer(&awscdk.DefaultStackSynthesizerProps{
			// 	FileAssetsBucketName: jsii.String("eric-experiments-k9xjhbjk3r4s"),
			// }),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
