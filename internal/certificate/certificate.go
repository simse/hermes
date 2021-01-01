package certificate

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/simse/hermes/internal/session"
)

// ACM contains a AWS Certificate Manager service for interaction
var ACM *acm.ACM

// InitACM create an ACM service
func InitACM() {
	svc := acm.New(session.Session)
	ACM = svc
}

// Exists checks if a certificate exists to cover a domain
func Exists(domain string) (bool, error) {
	domainExists := false

	// List all issued certificates
	listCertInput := acm.ListCertificatesInput{
		CertificateStatuses: []*string{aws.String("ISSUED")},
	}
	listCertOutput, err := ACM.ListCertificates(&listCertInput)
	if err != nil {
		panic(err)
	}

	// Check if any certificate covers the domain name
	for _, certificate := range listCertOutput.CertificateSummaryList {
		if *certificate.DomainName == domain {
			return true, nil
		}
	}

	// Check all certificates for aliases
	for _, certificate := range listCertOutput.CertificateSummaryList {
		// Get all certificate information for a cert given ARN
		describeCertInput := acm.DescribeCertificateInput{
			CertificateArn: certificate.CertificateArn,
		}

		describeCertOutput, describeCertErr := ACM.DescribeCertificate(&describeCertInput)
		if describeCertErr != nil {
			panic(describeCertErr)
		}

		// Check certificates
		for _, alias := range describeCertOutput.Certificate.SubjectAlternativeNames {
			// If alias matches domain, that domain is covered by a certificate
			if *alias == domain {
				return true, nil
			}
		}
	}

	return domainExists, nil
}

// Get returns a certficate given name
func Get(domain string) (*acm.CertificateDetail, error) {
	// List all issued certificates
	listCertInput := acm.ListCertificatesInput{
		CertificateStatuses: []*string{aws.String("ISSUED")},
	}
	listCertOutput, err := ACM.ListCertificates(&listCertInput)
	if err != nil {
		panic(err)
	}

	// Check if any certificate covers the domain name
	for _, certificate := range listCertOutput.CertificateSummaryList {
		if *certificate.DomainName == domain {
			// Get all certificate information for a cert given ARN
			describeCertInput := acm.DescribeCertificateInput{
				CertificateArn: certificate.CertificateArn,
			}

			describeCertOutput, describeCertErr := ACM.DescribeCertificate(&describeCertInput)
			if describeCertErr != nil {
				panic(describeCertErr)
			}

			return describeCertOutput.Certificate, nil
		}
	}

	// Check all certificates for aliases
	for _, certificate := range listCertOutput.CertificateSummaryList {
		// Get all certificate information for a cert given ARN
		describeCertInput := acm.DescribeCertificateInput{
			CertificateArn: certificate.CertificateArn,
		}

		describeCertOutput, describeCertErr := ACM.DescribeCertificate(&describeCertInput)
		if describeCertErr != nil {
			panic(describeCertErr)
		}

		// Check certificates
		for _, alias := range describeCertOutput.Certificate.SubjectAlternativeNames {
			// If alias matches domain, that domain is covered by a certificate
			if *alias == domain {
				return describeCertOutput.Certificate, nil
			}
		}
	}

	return &acm.CertificateDetail{}, nil
}
