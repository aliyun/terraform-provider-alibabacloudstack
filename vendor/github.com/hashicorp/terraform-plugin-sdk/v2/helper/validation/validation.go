package validation

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"net"
	"regexp"
)

func ValidateRegexp(v interface{}, k string) (ws []string, errors []error) {
	if _, err := regexp.Compile(v.(string)); err != nil {
		errors = append(errors, fmt.Errorf("%q: %s", k, err))
	}
	return
}

func SingleIP() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		ip := net.ParseIP(v)
		if ip == nil {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid IP, got: %s", k, v))
		}
		return
	}
}
// ValidateJsonString is a SchemaValidateFunc which tests to make sure the
// supplied string is valid JSON.
func ValidateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := structure.NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}