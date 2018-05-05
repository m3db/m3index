// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package doc

import (
	"bytes"
	"fmt"

	"github.com/golang/mock/gomock"
)

// DocumentMatcher matches a given document.
type DocumentMatcher interface {
	gomock.Matcher
}

// NewDocumentMatcher returns a new DocumentMatcher.
func NewDocumentMatcher(d Document) DocumentMatcher {
	return docMatcher{d}
}

type docMatcher struct {
	d Document
}

func (dm docMatcher) Matches(x interface{}) bool {
	other, ok := x.(Document)
	if !ok {
		return false
	}
	if !bytes.Equal(dm.d.ID, other.ID) {
		return false
	}
	if len(dm.d.Fields) != len(other.Fields) {
		return false
	}
	for i := range dm.d.Fields {
		if !bytes.Equal(dm.d.Fields[i].Name, other.Fields[i].Name) {
			return false
		}
		if !bytes.Equal(dm.d.Fields[i].Value, other.Fields[i].Value) {
			return false
		}
	}
	return true

}

func (dm docMatcher) String() string {
	return fmt.Sprintf("doc is %+v", dm.d)
}