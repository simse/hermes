package cmd

import (
	"fmt"
	"os"

	"github.com/simse/hermes/internal/cdn"
	"github.com/simse/hermes/internal/console"
	"github.com/simse/hermes/internal/edge"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/simse/hermes/internal/bucket"
	"github.com/simse/hermes/internal/certificate"
	"github.com/simse/hermes/internal/session"
	"github.com/urfave/cli/v2"
)

// ActionSheet describes the hermes stack options to perform
type ActionSheet struct {
	BucketRegion                 string
	BucketName                   string
	BucketPublic                 bool
	DeleteBucket                 bool
	Domain                       string
	NewOriginAccessIdentity      bool
	ExistingOriginAccessIdentity string
}

// InitCommand is the main command for setting up a hermes stack
func InitCommand(c *cli.Context) error {
	fmt.Print("Hermes: v0.1.3\n")

	// Connect to AWS
	awsConnectionSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Connecting to AWS...")
	awsConnectionSpinner.Start()

	session.InitSession()
	certificate.InitACM()
	cdn.InitCloudFront()
	edge.InitLambda()

	awsConnectionSpinner.PersistWith(spin.Spinner{Frames: []string{"✓"}}, " Connection to AWS successful\n")

	// Ask for domain
	domain := ""
	domainPrompt := &survey.Input{
		Message: "Which domain would you like to use?",
	}
	survey.AskOne(domainPrompt, &domain)

	// Check certificate
	domainCheckSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Checking certificate status for domain...")
	domainCheckSpinner.Start()

	certificateCheck, certificateCheckError := certificate.Exists(domain)
	if certificateCheckError != nil {
		panic(certificateCheckError)
	}

	if certificateCheck {
		domainCheckSpinner.PersistWith(spin.Spinner{Frames: []string{"✓"}}, " Domain has certificate")
	} else {
		domainCheckSpinner.PersistWith(spin.Spinner{Frames: []string{"X"}}, " Domain does not have certificate :(")

		fmt.Print("\nLet's fix it\n")

		return nil
	}

	// Check if there's a conflicting CloudFront distribution
	cfConflictSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Checking for conflicting CloudFront distributions...")
	cfConflictSpinner.Start()

	if cdn.CheckConflictCNAME(domain) {
		cfConflictSpinner.PersistWith(spin.Spinner{Frames: []string{"X"}}, " There's a CloudFront distribution using this domain. Please remove it before proceeding.")

		return nil
	}

	cfConflictSpinner.PersistWith(spin.Spinner{Frames: []string{"✓"}}, " No conflicting CloudFront distributions")

	/*
		// Inquire about CloudFront distribution
		whiteUnderline.Println("\n\nCloudFront")
	*/

	actionSheet := ActionSheet{
		Domain: domain,
	}

	// Inquire about S3 bucket
	console.WhiteUnderline.Println("\n\nS3 bucket")
	fmt.Println("hermes requires an S3 bucket to store your website and its own config")
	fmt.Println("\nYou can use an existing bucket or create a new one.")

	useExistingBucket := false
	existingBucketPrompt := &survey.Confirm{
		Message: "Would you like to use an existing bucket?",
	}
	survey.AskOne(existingBucketPrompt, &useExistingBucket)

	// Connect to S3
	bucket.InitS3()

	if useExistingBucket {
		fmt.Println("This hasn't been implemented yet, sorry")
		return nil
	}

	bucketExists, bucketExistsReason := bucket.Exists(domain)

	// Check for bucket conflict
	if bucketExists {
		if bucketExistsReason == bucket.ErrBucketExistsForeign {
			fmt.Println("\nA bucket with the name: , already exists.")

			overwriteBucketPrompt := &survey.Confirm{
				Message: "Would you like to overwrite the bucket?",
			}
			survey.AskOne(overwriteBucketPrompt, &actionSheet.DeleteBucket)

			if actionSheet.DeleteBucket {
				actionSheet.BucketName = domain
			} else {
				// alternativeBucketName := ""
				alternativeBucketNamePrompt := &survey.Select{
					Message: "Please pick another bucket name: ",
					Options: []string{"hello", "hello2"},
				}

				survey.AskOne(alternativeBucketNamePrompt, &actionSheet.BucketName)
			}
		} else {
			fmt.Println("This bucket is owned by you, please delete it and try again.")
		}

	} else {
		actionSheet.BucketName = domain
	}

	fmt.Println("")

	// Ask about region
	// alternativeBucketName := ""
	alternativeBucketNamePrompt := &survey.Select{
		Message: "Please pick a bucket region: ",
		Options: []string{
			"us-east-1",
			"us-east-2",
			"us-west-1",
			"us-west-2",
			"af-south-1",
			"ap-east-1",
			"ap-south-1",
			"ap-northeast-2",
			"ap-southeast-1",
			"ap-northeast-1",
			"ap-southeast-2",
			"ca-central-1",
			"eu-central-1",
			"eu-west-1",
			"eu-west-2",
			"eu-south-1",
			"eu-west-3",
			"eu-north-1",
			"me-south-1",
			"sa-east-1",
		},
		Help: "Pick region closest to YOU or the primary deploy server",
	}
	survey.AskOne(alternativeBucketNamePrompt, &actionSheet.BucketRegion)

	fmt.Println("")

	bucketPublicPrompt := &survey.Confirm{
		Message: "Would you like to make the bucket public?",
		Help:    "By default all files will only be accesible through your domain (CloudFront)",
	}
	survey.AskOne(bucketPublicPrompt, &actionSheet.BucketPublic)

	// fmt.Println(actionSheet)

	// Ask about OAI
	actionSheet.NewOriginAccessIdentity = true

	// Confirm action sheet
	console.WhiteUnderline.Print("\n\nConfirm setup")
	fmt.Println("\nhermes has not yet created any resouces. Before continuing please verify that all details are correct.")
	fmt.Print("\n")

	console.ShowLegend()

	console.ShowUsing("domain", domain)

	if actionSheet.DeleteBucket {
		console.ShowDelete("s3", "s3://"+actionSheet.BucketName)

		console.ShowCreate("s3", "s3://"+actionSheet.BucketName)
	} else {
		console.ShowCreate("s3", "s3://"+actionSheet.BucketName)
	}

	if actionSheet.NewOriginAccessIdentity {
		console.ShowCreate("cloudfront", "new Origin Access Identity")

		console.ShowAttach("cloudfront", "new Origin Access Identity", "s3", "s3://"+actionSheet.BucketName)
	} else {
		console.ShowAttach("cloudfront", "Origin Access Identity: sjkfsdjfdsjfhdsj", "s3", "s3://"+actionSheet.BucketName)
	}

	console.ShowCreate("lambda@edge", "hermesOriginResponse")
	console.ShowCreate("lambda@edge", "hermesOriginRequest")
	console.ShowAttach("lambda@edge", "hermesOriginResponse", "cloudfront", "")
	console.ShowAttach("lambda@edge", "hermesOriginRequest", "cloudfront", "")

	fmt.Print("\n")

	confirmActionSheet := false
	confirmActionSheetPrompt := &survey.Confirm{
		Message: "Shall we proceed?",
	}
	survey.AskOne(confirmActionSheetPrompt, &confirmActionSheet)

	// Create bucket
	// Create OAI (if neccesary)
	// Deploy lambda functions
	// Create CloudFront distribution
	// Create default deploy
	// Inform about domain changes

	return nil
}
