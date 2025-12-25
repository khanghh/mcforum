package validate

import (
	"errors"
	"regexp"
	"strings"

	"bbs-go/common/strs"
)

// IsUsername Validate username legality: must be 5-12 chars (digits, letters, _, -) and start with a letter.
func IsUsername(username string) error {
	if strs.IsBlank(username) {
		return errors.New("Please enter a username")
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return errors.New("Username must be 5-12 characters (digits, letters, _, -) and start with a letter")
	}
	matched, err = regexp.MatchString("^[a-zA-Z]", username)
	if err != nil || !matched {
		return errors.New("Username must be 5-12 characters (digits, letters, _, -) and start with a letter")
	}
	return nil
}

// IsEmail Validate whether the email is valid
func IsEmail(email string) (err error) {
	if strs.IsBlank(email) {
		err = errors.New("Invalid email format")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("Invalid email format")
	}
	return
}

// IsValidPassword Validate whether the password is valid
func IsValidPassword(password, rePassword string) error {
	if err := IsPassword(password); err != nil {
		return err
	}
	if password != rePassword {
		return errors.New("Passwords do not match")
	}
	return nil
}

func IsPassword(password string) error {
	if strs.IsBlank(password) {
		return errors.New("Please enter a password")
	}
	if strs.RuneLen(password) < 6 {
		return errors.New("Password is too simple")
	}
	if strs.RuneLen(password) > 1024 {
		return errors.New("Password length cannot exceed 128")
	}
	return nil
}

// IsURL Validate whether the URL is valid
func IsURL(url string) error {
	if strs.IsBlank(url) {
		return errors.New("Invalid URL format")
	}
	indexOfHttp := strings.Index(url, "http://")
	indexOfHttps := strings.Index(url, "https://")
	if indexOfHttp == 0 || indexOfHttps == 0 {
		return nil
	}
	return errors.New("Invalid URL format")
}
