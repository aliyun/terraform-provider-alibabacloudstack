package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// below copy/pasta from https://github.com/hashicorp/terraform-plugin-sdk/v2/blob/master/helper/validation/validation.go
// alibabacloudstack vendor contains very old version of Terraform which lacks this functions

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func intBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		return
	}
}
func validateAllowedSplitStringValue(ss []string, splitStr string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		existed := false
		tsList := strings.Split(value, splitStr)

		for _, ts := range tsList {
			existed = false
			for _, s := range ss {
				if ts == s {
					existed = true
					break
				}
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid string value should in %#v, got %q",
				k, ss, value))
		}
		return

	}
}
func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
		return
	}

	mark, _ := strconv.Atoi(strings.Split(ipnet.String(), "/")[1])
	if mark < 16 || mark > 29 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a network CIDR which mark between 16 and 29",
			k))
	}

	return
}
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}
func validateForwardPort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "any" {
		valueConv, err := strconv.Atoi(value)
		if err != nil || valueConv < 1 || valueConv > 65535 {
			errors = append(errors, fmt.Errorf("%q must be a valid port between 1 and 65535 or any ", k))
		}
	}
	return
}
func validateStringConvertInt64() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		if value, ok := v.(string); ok {
			_, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				errors = append(errors, fmt.Errorf(
					"%q should be convert to int64, got %q", k, value))
			}
		} else {
			errors = append(errors, fmt.Errorf(
				"%q should be convert to string, got %q", k, value))
		}

		return
	}
}

func validateOssBucketDateTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as date YYYY-MM-DD Format", value))
	}
	return
}

func validateDBConnectionPort(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			errors = append(errors, err)
		}
		if port < 3001 || len(value) > 3999 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 3001 and larger than 3999.", k))
		}
	}
	return
}

func validateOnsGroupId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !(strings.HasPrefix(value, "GID-") || strings.HasPrefix(value, "GID_")) {
		errors = append(errors, fmt.Errorf("%q is invalid, it must start with 'GID-' or 'GID_'", k))
	}
	if reg := regexp.MustCompile(`^[\w\-]{7,64}$`); !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q length is limited to 7-64 and only characters such as letters, digits, '_' and '-' are allowed", k))
	}
	return
}
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

func validateRR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("RR is invalid, it can not starts or ends with '-'"))
	}

	if len(value) > 253 {
		errors = append(errors, fmt.Errorf("RR can not longer than 253 characters."))
	}

	for _, part := range strings.Split(value, ".") {
		if len(part) > 63 {
			errors = append(errors, fmt.Errorf("Each part of RR split with . can not longer than 63 characters."))
			return
		}
	}
	return
}

// Validate length(2~128) and prefix of the name.
func validateNormalName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}
	return
}
