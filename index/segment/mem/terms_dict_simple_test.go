// Copyright (c) 2017 Uber Technologies, Inc.
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

package mem

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/m3db/m3ninx/doc"
	"github.com/m3db/m3ninx/postings"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/suite"
)

var (
	testRandomSeed         int64 = 42
	testMinSuccessfulTests       = 1000

	sampleRegexes = []interface{}{
		`a`,
		`a.`,
		`a.b`,
		`ab`,
		`a.b.c`,
		`abc`,
		`a|^`,
		`a|b`,
		`(a)`,
		`(a)|b`,
		`a*`,
		`a+`,
		`a?`,
		`a{2}`,
		`a{2,3}`,
		`a{2,}`,
		`a*?`,
		`a+?`,
		`a??`,
		`a{2}?`,
		`a{2,3}?`,
		`a{2,}?`,
	}
)

type newSimpleTermsDictFn func() *simpleTermsDict

type simpleTermsDictionaryTestSuite struct {
	suite.Suite

	fn        newSimpleTermsDictFn
	termsDict *simpleTermsDict
}

func (t *simpleTermsDictionaryTestSuite) SetupTest() {
	t.termsDict = t.fn()
}

func (t *simpleTermsDictionaryTestSuite) TestInsert() {
	props := getProperties()
	props.Property(
		"The dictionary should supporting inserting fields",
		prop.ForAll(
			func(f doc.Field, id postings.ID) (bool, error) {
				err := t.termsDict.Insert(f, id)
				if err != nil {
					return false, fmt.Errorf("unexpected error inserting %v into terms dictionary: %v", f, err)
				}

				return true, nil
			},
			genField(),
			genDocID(),
		))

	props.TestingRun(t.T())
}

func (t *simpleTermsDictionaryTestSuite) TestMatchExact() {
	props := getProperties()
	props.Property(
		"The dictionary should support exact match queries",
		prop.ForAll(
			func(f doc.Field, id postings.ID) (bool, error) {
				err := t.termsDict.Insert(f, id)
				if err != nil {
					return false, fmt.Errorf("unexpected error inserting %v into terms dictionary: %v", f, err)
				}

				pl, err := t.termsDict.MatchExact(f.Name, []byte(f.Value))
				if err != nil {
					return false, fmt.Errorf("unexpexted error retrieving postings list: %v", err)
				}
				if pl == nil {
					return false, fmt.Errorf("postings list of documents matching query should not be nil")
				}
				if !pl.Contains(id) {
					return false, fmt.Errorf("id of new document '%v' is not in postings list of matching documents", id)
				}

				return true, nil
			},
			genField(),
			genDocID(),
		))

	props.TestingRun(t.T())
}

func (t *simpleTermsDictionaryTestSuite) TestMatchExactNoResults() {
	props := getProperties()
	props.Property(
		"Exact match queries which return no results are valid",
		prop.ForAll(
			func(f doc.Field) (bool, error) {
				pl, err := t.termsDict.MatchExact(f.Name, []byte(f.Value))
				if err != nil {
					return false, fmt.Errorf("unexpexted error retrieving postings list: %v", err)
				}
				if pl == nil {
					return false, fmt.Errorf("postings list returned should not be nil")
				}
				if pl.Len() != 0 {
					return false, fmt.Errorf("postings list contains unexpected IDs")
				}

				return true, nil
			},
			genField(),
		))

	props.TestingRun(t.T())
}

func (t *simpleTermsDictionaryTestSuite) TestMatchRegex() {
	props := getProperties()
	props.Property(
		"The dictionary should support regular expression queries",
		prop.ForAll(
			func(input fieldAndRegex, id postings.ID) (bool, error) {
				var (
					f       = input.field
					pattern = input.pattern
					re      = input.re
				)
				err := t.termsDict.Insert(f, id)
				if err != nil {
					return false, fmt.Errorf("unexpected error inserting %v into terms dictionary: %v", f, err)
				}

				pl, err := t.termsDict.MatchRegex(f.Name, []byte(pattern), re)
				if err != nil {
					return false, fmt.Errorf("unexpexted error retrieving postings list: %v", err)
				}
				if pl == nil {
					return false, fmt.Errorf("postings list of documents matching query should not be nil")
				}
				if !pl.Contains(id) {
					return false, fmt.Errorf("id of new document '%v' is not in list of matching documents", id)
				}

				return true, nil
			},
			genFieldAndRegex(),
			genDocID(),
		))

	props.TestingRun(t.T())
}

func (t *simpleTermsDictionaryTestSuite) TestMatchRegexNoResults() {
	props := getProperties()
	props.Property(
		"Regular expression queries which no results are valid",
		prop.ForAll(
			func(input fieldAndRegex, id postings.ID) (bool, error) {
				var (
					f       = input.field
					pattern = input.pattern
					re      = input.re
				)
				pl, err := t.termsDict.MatchRegex(f.Name, []byte(pattern), re)
				if err != nil {
					return false, fmt.Errorf("unexpexted error retrieving postings list: %v", err)
				}
				if pl == nil {
					return false, fmt.Errorf("postings list returned should not be nil")
				}
				if pl.Len() != 0 {
					return false, fmt.Errorf("postings list contains unexpected IDs")
				}

				return true, nil
			},
			genFieldAndRegex(),
			genDocID(),
		))

	props.TestingRun(t.T())
}

func TestSimpleTermsDictionary(t *testing.T) {
	opts := NewOptions()
	suite.Run(t, &simpleTermsDictionaryTestSuite{
		fn: func() *simpleTermsDict {
			return newSimpleTermsDict(opts).(*simpleTermsDict)
		},
	})
}

func getProperties() *gopter.Properties {
	params := gopter.DefaultTestParameters()
	params.Rng.Seed(testRandomSeed)
	params.MinSuccessfulTests = testMinSuccessfulTests
	return gopter.NewProperties(params)
}

func genField() gopter.Gen {
	return gopter.CombineGens(
		gen.AnyString(),
		gen.AnyString(),
	).Map(func(values []interface{}) doc.Field {
		var (
			name  = values[0].(string)
			value = values[1].(string)
		)
		f := doc.Field{
			Name:  []byte(name),
			Value: []byte(value),
		}
		return f
	})
}

func genDocID() gopter.Gen {
	return gen.UInt32().
		Map(func(value uint32) postings.ID {
			return postings.ID(value)
		})
}

type fieldAndRegex struct {
	field   doc.Field
	pattern string
	re      *regexp.Regexp
}

func genFieldAndRegex() gopter.Gen {
	return gen.OneConstOf(sampleRegexes...).
		FlatMap(func(value interface{}) gopter.Gen {
			pattern := value.(string)
			return fieldFromRegex(pattern)
		}, reflect.TypeOf(fieldAndRegex{}))
}

func fieldFromRegex(pattern string) gopter.Gen {
	return gopter.CombineGens(
		gen.AnyString(),
		gen.RegexMatch(pattern),
	).Map(func(values []interface{}) fieldAndRegex {
		var (
			name  = values[0].(string)
			value = values[1].(string)
		)
		f := doc.Field{
			Name:  []byte(name),
			Value: []byte(value),
		}
		return fieldAndRegex{
			field:   f,
			pattern: pattern,
			re:      regexp.MustCompile(pattern),
		}
	})
}
