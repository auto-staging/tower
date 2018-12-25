package controller

import "regexp"

func validateIAMRoleARN(arn string) bool {
	regex := regexp.MustCompile(`arn:aws:iam::\d{12}:role/?[a-zA-Z_0-9+=,.@\-_/]+`)
	return regex.MatchString(arn)
}
