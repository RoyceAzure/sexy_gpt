package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateStringLen(str string, n int) error {
	if len(str) < n {
		return fmt.Errorf("string length must greater than %d", n)
	}
	return nil
}

func ValidateEmptyString(str string) error {
	if str == "" {
		return fmt.Errorf("string is empty")
	}
	return nil
}

func ValidateStrongPas(str string) []error {
	var errs []error
	if err := ValidateEmptyString(str); err != nil {
		errs = append(errs, err)
		return errs
	}
	var isDigit, islower, isUpper bool
	for _, ru := range str {
		if unicode.IsDigit(ru) {
			isDigit = true
		}
		if unicode.IsLower(ru) {
			islower = true
		}
		if unicode.IsUpper(ru) {
			isUpper = true
		}
		if isDigit && islower && isUpper {
			return nil
		}
	}

	if err := ValidateStringLen(str, 6); err != nil {
		errs = append(errs, err)
	}
	if !isDigit {
		errs = append(errs, fmt.Errorf("string must contain digit"))
	}
	if !islower {
		errs = append(errs, fmt.Errorf("string must contain lowercase letters"))
	}
	if !isUpper {
		errs = append(errs, fmt.Errorf("string must contain uppercase letters"))
	}
	return errs
}

func ValidateUUID(str string) error {
	_, err := uuid.Parse(str)
	return err
}

func ValidateEmailFormat(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	res, _ := regexp.MatchString(pattern, email)
	if !res {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func FieldViolations(field string, errs []error) *errdetails.BadRequest_FieldViolation {
	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString(err.Error())
		sb.WriteString(",")
	}
	temp := sb.String()
	return FieldViolation(field, fmt.Errorf(temp[:len(temp)-1]))
}
