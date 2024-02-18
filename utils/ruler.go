package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validatorKeywords = []string{
	`(alert).*(\()`,
	`(prompt).*(\()`,
	`(eval).*(\()`,
	`(window).*(\[)`,
	`<script`,
	`</script`,
	`<x`,
	`<X`,
	`<http`,
	`(function).*(\()`,
	`<iframe`,
	`(href).*(=)`,
	`<br>`,
	"alert`",
	`(find).*(\()`,
	`(top).*(\[)`,
	`(vibrate).*(\()`,
	`<object`,
	`<embed`,
	`<img`,
	`<layer`,
	`<style`,
	`<meta`,
	`=".*"`,
	`<html`,
	`(echo).*(\()`,
	`(confirm).*(\()`,
	`(write).*(\()`,
	`</svg`,
	`<div`,
	`</image`,
	`form>`,
	`(vectors).*(\()`,
	`<body`,
	`(url).*(\()`,
	`math>`,
	`-->`,
	`<!--`,
	`<!attlist`,
	`<label`,
	`<%`,
	`xmp>`,
	`template>`,
	`<!doctype`,
	`=confirm`,
	`=cmd`,
}

func ValidateVar(validate *validator.Validate, data interface{}, validation string) error {
	containChar := ""
	validate.RegisterValidation("unsafe", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field().String()

		allowedRegex := regexp.MustCompile(`[\p{Han}\p{Hiragana}\p{Katakana}\p{Hangul}\p{So}]`)

		for _, keyword := range validatorKeywords {
			if matched, _ := regexp.MatchString(keyword, fieldValue); matched && !allowedRegex.MatchString(string(keyword)) {
				containChar = fmt.Sprint("keyword ", string(keyword))
				return false
			}
		}

		containedChars := ""
		for _, char := range fieldValue {
			if strings.ContainsAny(string(char), "[](){}<>=") && !allowedRegex.MatchString(string(char)) {
				containedChars = fmt.Sprint(containedChars, string(char))
			}
		}

		if containedChars != "" {
			containChar = fmt.Sprint("characters ", string(containedChars))
			return false
		}

		return true
	})

	if err := validate.Var(data, validation); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		errString := ""
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "gte" {
				errString += fmt.Sprintf("Minimum is %s\n", err.Param())
			} else if err.Tag() == "lte" {
				errString += fmt.Sprintf("Maximum is %s\n", err.Param())
			} else if err.Tag() == "min" {
				if err.Kind() == reflect.TypeOf(err.Field()).Kind() {
					errString += fmt.Sprintf("Minimum %s length is %s\n", err.Field(), err.Param())
				} else {
					errString += fmt.Sprintf("Minimum %s is %s\n", err.Field(), err.Param())
				}
			} else if err.Tag() == "max" {
				if err.Kind() == reflect.TypeOf(err.Field()).Kind() {
					errString += fmt.Sprintf("Maximum %s length is %s\n", err.Field(), err.Param())
				} else {
					errString += fmt.Sprintf("Maximum %s is %s\n", err.Field(), err.Param())
				}
			} else if err.Tag() == "required" {
				errString += fmt.Sprintf("is %s\n", err.Tag())
			} else if err.Tag() == "html_encoded" {
				errString += "contains unencoded characters\n"
			} else if err.Tag() == "unsafe" {
				errString += fmt.Sprintf("contains unsafe %s\n", containChar)
			} else if err.Tag() == "gt" {
				errString += "has invalid value\n"
			} else if err.Tag() == "http_url" {
				errString += fmt.Sprintf("must contains %s\n", err.Tag())
			} else if err.Tag() == "alphanum" {
				errString += fmt.Sprintf("only accepts %s", err.Tag())
			} else if err.Tag() == "number" {
				errString += fmt.Sprintf("only accepts %s", err.Tag())
			} else {
				errString += fmt.Sprintf("Invalid input %s\n", err.Tag())
			}
		}
		return fmt.Errorf(strings.TrimSpace(errString))

	}

	return nil
}

func ValidateStruct(validate *validator.Validate, model interface{}) error {

	containChar := ""
	validate.RegisterValidation("unsafe", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field().String()

		allowedRegex := regexp.MustCompile(`[\p{Han}\p{Hiragana}\p{Katakana}\p{Hangul}\p{So}]`)

		for _, keyword := range validatorKeywords {
			if matched, _ := regexp.MatchString(keyword, fieldValue); matched && !allowedRegex.MatchString(string(keyword)) {
				containChar = fmt.Sprint("keyword ", string(keyword))
				return false
			}
		}

		containedChars := ""
		for _, char := range fieldValue {
			if strings.ContainsAny(string(char), "[](){}<>=") && !allowedRegex.MatchString(string(char)) {
				containedChars = fmt.Sprint(containedChars, string(char))
			}
		}

		if containedChars != "" {
			containChar = fmt.Sprint("characters ", string(containedChars))
			return false
		}

		return true
	})

	err2 := validate.Struct(model)
	if err2 != nil {
		if _, ok := err2.(*validator.InvalidValidationError); ok {
			return err2
		}

		errString := ""
		for _, err := range err2.(validator.ValidationErrors) {
			if err.Tag() == "gte" {
				errString += fmt.Sprintf("Minimum %s is %s\n", err.Field(), err.Param())
			} else if err.Tag() == "lte" {
				errString += fmt.Sprintf("Maximum %s is %s\n", err.Field(), err.Param())
			} else if err.Tag() == "min" {
				if err.Kind() == reflect.TypeOf(err.Field()).Kind() {
					errString += fmt.Sprintf("Minimum %s length is %s\n", err.Field(), err.Param())
				} else {
					errString += fmt.Sprintf("Minimum %s is %s\n", err.Field(), err.Param())
				}
			} else if err.Tag() == "max" {
				if err.Kind() == reflect.TypeOf(err.Field()).Kind() {
					errString += fmt.Sprintf("Maximum %s length is %s\n", err.Field(), err.Param())
				} else {
					errString += fmt.Sprintf("Maximum %s is %s\n", err.Field(), err.Param())
				}
			} else if err.Tag() == "required" {
				errString += fmt.Sprintf("%s is %s\n", err.Field(), err.Tag())
			} else if err.Tag() == "url" {
				errString += fmt.Sprintf("%s must contains %s\n", err.Field(), err.Tag())
			} else if err.Tag() == "html_encoded" {
				errString += fmt.Sprintf("%s contains unencoded characters\n", err.Field())
			} else if err.Tag() == "unsafe" {
				errString += fmt.Sprintf("%s contains unsafe %s\n", err.Field(), containChar)
			} else if err.Tag() == "gt" {
				errString += fmt.Sprintf("%s has invalid value\n", err.Field())
			} else if err.Tag() == "http_url" {
				errString += fmt.Sprintf("%s must contains %s\n", err.Field(), err.Tag())
			} else if err.Tag() == "alphanum" {
				errString += fmt.Sprintf("%s only accepts %s", err.Field(), err.Tag())
			} else if err.Tag() == "number" {
				errString += fmt.Sprintf("%s only accepts %s", err.Field(), err.Tag())
			} else {
				errString += fmt.Sprintf("Invalid input for %s: %s\n", err.Field(), err.Tag())
			}
		}

		return fmt.Errorf(strings.TrimSpace(errString))
	}

	return nil
}
