package main

import (
	"fmt"
	"os"
	"regexp"
)

//CustomEnv - Get true custom_env for gitlab pipelines
func CustomEnv(a, b string) (customenv string, err error) {
	//uat??
	_, exist := os.LookupEnv(a)
	if !exist {
		err := fmt.Errorf("Error:\n Environment variable not found: %s", b)
		return "", err
	}
	_, exist = os.LookupEnv(b)
	if !exist {
		err := fmt.Errorf("Error:\n Environment variable not found: %s", b)
		return "", err
	}

	if len(a) > 13 && a[13:len(a)-1] == "0" {
		pattern := regexp.MustCompile("0")
		a = pattern.ReplaceAllString(a, "")
	}

	c := "test" + os.Getenv(a)[13:] + os.Getenv(b)[4:]
	return c, nil
}
