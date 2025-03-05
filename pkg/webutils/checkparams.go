package webutils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
)

func CheckParamToInt(context echo.Context, paramName string) (int, error) {
	paramValue, err := CheckParam(context, paramName)
	if err != nil {
		return -1, err
	}

	paramInt, err := strconv.Atoi(paramValue)
	if err != nil {
		return -1, err
	}
	return paramInt, nil
}

func CheckParam(context echo.Context, paramName string) (string, error) {
	paramValue := context.Param(paramName)
	if paramValue == "" {
		return "", fmt.Errorf("missing url param: %s", paramName)
	}
	return paramValue, nil
}
