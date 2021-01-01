package edge

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/markbates/pkger"
	"github.com/simse/hermes/internal/session"
)

// Lambda stores the current connection to AWS lambda
var Lambda *lambda.Lambda

// InitLambda creates a lambda service
func InitLambda() {
	Lambda = lambda.New(session.Session)
}

// CreateLambdaFunction creates a lambda function
func CreateLambdaFunction(name, role, runtime, path string) error {
	funcFile, openError := pkger.Open("/assets/lambda/" + path)
	if openError != nil {
		panic(openError)
	}

	funcZip, _ := ioutil.ReadAll(funcFile)

	lambdaCode := lambda.FunctionCode{
		ZipFile: funcZip,
	}

	tags := make(map[string]*string)
	tags["X-Created-By"] = aws.String("hermes")

	createFunctionInput := lambda.CreateFunctionInput{
		Code:         &lambdaCode,
		FunctionName: aws.String(name),
		Role:         aws.String(role),
		Runtime:      aws.String(runtime),
		Handler:      aws.String("index.handler"),
		Tags:         tags,
	}

	createFunctionOutput, createFunctionErr := Lambda.CreateFunction(&createFunctionInput)
	fmt.Println(createFunctionOutput)

	if createFunctionErr != nil {
		panic(createFunctionErr)
	}

	getFunctionInput := lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	}

	waitError := Lambda.WaitUntilFunctionExists(&getFunctionInput)
	if waitError != nil {
		panic(waitError)
	}

	return nil
}

// PublishLambdaFunction publishes a lambda function
func PublishLambdaFunction(name string) (*lambda.FunctionConfiguration, error) {
	publishVersionInput := lambda.PublishVersionInput{
		FunctionName: aws.String(name),
	}

	publishVersionOutput, err := Lambda.PublishVersion(&publishVersionInput)
	if err != nil {
		panic(err)
	}

	return publishVersionOutput, nil
}

// CreateExecutionRole creates a role so the function works with lambda@edge
func CreateExecutionRole(name string) (string, error) {
	svc := iam.New(session.Session)

	policyFile, openError := pkger.Open("/assets/lambda/rolePolicyDocument.json")
	if openError != nil {
		panic(openError)
	}

	policyDocument, _ := ioutil.ReadAll(policyFile)

	var tags []*iam.Tag
	tags = append(tags, &iam.Tag{
		Key:   aws.String("X-Created-By"),
		Value: aws.String("hermes"),
	})

	createRoleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(policyDocument)),
		RoleName:                 aws.String(name),
		Tags:                     tags,
	}

	createRoleOutput, err := svc.CreateRole(&createRoleInput)
	if err != nil {
		panic(err)
	}

	attachRolePolicyInput := iam.AttachRolePolicyInput{
		RoleName:  createRoleOutput.Role.RoleName,
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	}
	_, attachRoleErr := svc.AttachRolePolicy(&attachRolePolicyInput)
	if attachRoleErr != nil {
		panic(err)
	}

	getRoleInput := iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	svc.WaitUntilRoleExists(&getRoleInput)

	return *createRoleOutput.Role.Arn, nil
}

/*
// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func zipFiles(files []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	//defer newZipFile.Close()

	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return buf, err
		}
	}

	return buf, nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
*/
