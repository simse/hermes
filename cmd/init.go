package cmd

import (
	"fmt"
	"os"

	"github.com/simse/hermes/internal/constants"

	"github.com/simse/hermes/internal/cdn"
	"github.com/simse/hermes/internal/console"
	"github.com/simse/hermes/internal/edge"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
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
	bucket.InitS3()

	awsConnectionSpinner.PersistWith(console.Check, " Connection to AWS successful\n")

	//role, _ := edge.CreateExecutionRole("test")

	//fmt.Println(role)

	//time.Sleep(time.Second * 5)

	//edge.CreateLambdaFunction("test", "arn:aws:iam::616657489041:role/test", "nodejs12.x", "hermesOriginRequest.zip")

	// time.Sleep(time.Second * 5)

	//handler, _ := edge.PublishLambdaFunction("test")

	// bucket.Create("simons-big-test", "us-east-1", false)
	// identity, _ := cdn.CreateOAI("hermes")
	// bucket.AddOAIPermissions("simons-big-test", *identity.S3CanonicalUserId)

	//associations := make(map[string]string)
	//associations["origin-request"] = *handler.FunctionArn

	//cdn.CreateDistribution("novelhardware.com", "novelhardware.com", *identity.Id, "PriceClass_All", []string{"novelhardware.com"}, associations)

	//fmt.Println(deploy.ScanDir("./assets/default-site"))

	// return nil

	// Ask for domain
	domain := ""
	domainPrompt := &survey.Input{
		Message: "Which domain would you like to use?",
	}
	err := survey.AskOne(domainPrompt, &domain)
	if err == terminal.InterruptErr {
		fmt.Print("\n")
		return nil
	}

	// Check certificate
	domainCheckSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Checking certificate status for domain...")
	domainCheckSpinner.Start()

	certificateCheck, certificateCheckError := certificate.Exists(domain)
	if certificateCheckError != nil {
		panic(certificateCheckError)
	}

	if certificateCheck {
		domainCheckSpinner.PersistWith(console.Check, " Domain has certificate")
	} else {
		domainCheckSpinner.PersistWith(console.Cross, " Domain does not have certificate :(")

		fmt.Print("\nLet's fix it\n")

		return nil
	}

	// Check if there's a conflicting CloudFront distribution
	cfConflictSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Checking for conflicting CloudFront distributions...")
	cfConflictSpinner.Start()

	if cdn.CheckConflictCNAME(domain) {
		cfConflictSpinner.PersistWith(console.Cross, " There's a CloudFront distribution using this domain. Please remove it before proceeding.")

		return nil
	}

	cfConflictSpinner.PersistWith(console.Check, " No conflicting CloudFront distributions")

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
		shouldPickNewName := false

		// Bucket exists but is owner by another user
		if bucketExistsReason == bucket.ErrBucketExistsForeign {
			fmt.Println("\nA bucket with the name:" + domain + ", already exists.\n")

			shouldPickNewName = true

			// Bucket exists, and is owned by current user
		} else {
			fmt.Println("This bucket already exists in your account")
			fmt.Print("\n")

			overwriteBucketPrompt := &survey.Confirm{
				Message: "Would you like to overwrite the bucket?",
			}
			survey.AskOne(overwriteBucketPrompt, &actionSheet.DeleteBucket)

			if !actionSheet.DeleteBucket {
				shouldPickNewName = true
				fmt.Print("\n")
			} else {
				fmt.Println("cool")
			}
		}

		// Pick alternative bucket name
		if shouldPickNewName {
			// Generate alternative bucket names
			alternativeNamesSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Generating alternative names...")
			alternativeNamesSpinner.Start()

			alternativeNames := bucket.GenerateAlternativeNames(domain)
			alternativeNamesSpinner.PersistWith(console.Check, " Generated alternative bucket names")

			// alternativeBucketName := ""
			alternativeBucketNamePrompt := &survey.Select{
				Message: "Please pick another bucket name: ",
				Options: alternativeNames,
			}

			survey.AskOne(alternativeBucketNamePrompt, &actionSheet.BucketName)
		}

	} else {
		actionSheet.BucketName = domain
	}

	fmt.Println("")

	// Ask about region
	// alternativeBucketName := ""
	bucketRegion := &survey.Select{
		Message: "Please pick a bucket region: ",
		Options: constants.AWSRegionsList,
		Help:    "Pick region closest to YOU or the primary deploy server",
	}
	survey.AskOne(bucketRegion, &actionSheet.BucketRegion)

	if actionSheet.BucketRegion != "us-east-1" {
		session.InitSecondarySession(actionSheet.BucketRegion)
	}

	fmt.Println("")

	bucketPublicPrompt := &survey.Confirm{
		Message: "Would you like to make the bucket public?",
		Help:    "It's recommended to say no, so your files are only available through your domains.",
	}
	survey.AskOne(bucketPublicPrompt, &actionSheet.BucketPublic)

	// Ask about OAI
	actionSheet.NewOriginAccessIdentity = true
	// TODO: Offer to use existing OAI

	// Confirm action sheet
	console.WhiteUnderline.Print("\n\nConfirm setup")
	fmt.Println("\nBefore continuing please verify that all details are correct.")
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
		Help:    "Press Ctrl+C to cancel and no changes will be made.",
	}
	err = survey.AskOne(confirmActionSheetPrompt, &confirmActionSheet)
	if err == terminal.InterruptErr {
		fmt.Print("\n")
		return nil
	}

	if !confirmActionSheet {
		return nil
	}

	// Do action sheet

	// Create bucket
	// createBucketSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " Creating bucket...")
	// createBucketSpinner.Start()

	// err = bucket.Create(actionSheet.BucketName, actionSheet.BucketRegion, actionSheet.BucketPublic)
	// if err != nil {
	// 	panic(err)
	// }

	// createBucketSpinner.PersistWith(console.Check, " Created bucket")

	// Create OAI (if neccesary)
	// Deploy lambda functions
	// Create CloudFront distribution
	// Create default deploy
	// Inform about domain changes

	// manifest := deploy.Manifest{
	// 	InitVersion:   about.Version,
	// 	DeployVersion: about.Version,
	// 	Domain:        domain,
	// }

	// fmt.Println(manifest)

	return nil
}
