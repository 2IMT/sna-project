package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Environment struct {
	BotToken string
	DbHost   string
	DbPort   string
	DbName   string
	DbUser   string
	DbPass   string
}

var (
	int64Type   = reflect.TypeOf(0)
	float64Type = reflect.TypeOf(0.0)
	stringType  = reflect.TypeOf("")
)

func LoadEnvironment(filenames ...string) (Environment, error) {
	var result Environment

	n := reflect.TypeOf(result).NumField()
	for i := 0; i < n; i++ {
		name := reflect.TypeOf(result).Field(i).Name
		nameSnakeCase := camelToSnake(name)

		variable, exists := os.LookupEnv(nameSnakeCase)
		if !exists {
			errorString := fmt.Sprintf(
				"variable %s requested by Environment.%s is not set",
				nameSnakeCase, name)
			return Environment{}, errors.New(errorString)
		}

		switch reflect.TypeOf(result).Field(i).Type {
		case int64Type:
			if value, err := strconv.ParseInt(variable, 10, 64); err == nil {
				reflect.ValueOf(&result).Elem().Field(i).SetInt(value)
			} else {
				errorString := fmt.Sprintf(
					"failed to parse int64 from %s=%s requested by "+
						"Environment.%s",
					nameSnakeCase, variable, name)
				return Environment{}, errors.New(errorString)
			}

		case float64Type:
			if value, err := strconv.ParseFloat(variable, 64); err == nil {
				reflect.ValueOf(&result).Elem().Field(i).SetFloat(value)
			} else {
				errorString := fmt.Sprintf(
					"failed to parse float64 from %s=%s requested by "+
						"Environment.%s",
					nameSnakeCase, variable, name)
				return Environment{}, errors.New(errorString)
			}

		case stringType:
			reflect.ValueOf(&result).Elem().Field(i).SetString(variable)

		default:
			errorString := fmt.Sprintf(
				"Environment.%s has unknown type. Allowed types are int64, "+
					"float64, string",
				name)
			return Environment{}, errors.New(errorString)
		}
	}

	return result, nil
}

func camelToSnake(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToUpper(snake)
}
