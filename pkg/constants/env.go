package constants

import "os"

func ENVCredential() bool {
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		return true
	}
	return false
}
func ENVAWSProfile() bool {
	if os.Getenv("AWS_SHARED_CREDENTIALS_FILE") != "" {
		return true
	}
	return false
}
