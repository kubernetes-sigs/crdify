package results

import (
	"errors"
	"fmt"
	"strings"
)

type Result struct {
	Error      error
	Subresults []*Result
}

func ErrorFromResult(res *Result, depth int) error {
	if res == nil {
		return nil
	}

	if res.Error == nil && len(res.Subresults) == 0 {
		return nil
	}

	var out strings.Builder
	nestedErrors := []error{}
	for _, subresult := range res.Subresults {
		nestedErrors = append(nestedErrors, ErrorFromResult(subresult, depth+1))
	}

	if res.Error != nil {
		out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", depth*2), res.Error.Error()))
	}

	for _, nestedErr := range nestedErrors {
		if nestedErr != nil {
			out.WriteString(nestedErr.Error())
		}
	}
	return errors.New(out.String())
}
