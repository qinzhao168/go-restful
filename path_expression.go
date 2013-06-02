// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package restful

import (
	"bytes"
	"regexp"
	"strings"
)

type PathExpression struct {
	LiteralCount int
	VarCount     int
	Matcher      *regexp.Regexp
	Source       string
}

func NewPathExpression(path string) (*PathExpression, error) {
	expression, literalCount, varCount := templateToRegularExpression(path)
	compiled, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return &PathExpression{literalCount, varCount, compiled, expression}, nil
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-370003.7.3
func templateToRegularExpression(template string) (expression string, literalCount int, varCount int) {
	var buffer bytes.Buffer
	buffer.WriteString("^")
	tokens := strings.Split(template, "/")
	for _, each := range tokens {
		if each == "" {
			continue
		}
		buffer.WriteString("/")
		if strings.HasPrefix(each, "{") {
			// ignore var spec
			varCount += 1
			buffer.WriteString("([^/]+?)")
		} else {
			literalCount += len(each)
			encoded := each // TODO URI encode
			buffer.WriteString(regexp.QuoteMeta(encoded))
		}
	}
	return strings.TrimRight(buffer.String(), "/") + "(/.*)?$", literalCount, varCount
}
