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

package searcher

import (
	"testing"

	"github.com/m3db/m3ninx/search"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestValidateSearchers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	firstTestSearcher1 := search.NewMockSearcher(mockCtrl)
	firstTestSearcher2 := search.NewMockSearcher(mockCtrl)
	firstTestSeachers := search.Searchers{firstTestSearcher1, firstTestSearcher2}

	secondTestSearcher1 := search.NewMockSearcher(mockCtrl)
	secondTestSearcher2 := search.NewMockSearcher(mockCtrl)
	secondTestSeachers := search.Searchers{secondTestSearcher1, secondTestSearcher2}

	gomock.InOrder(
		// All searchers have the same length in the first test.
		firstTestSearcher1.EXPECT().Len().Return(3),
		firstTestSearcher2.EXPECT().Len().Return(3),

		// The searchers do not all have the same length in the second test.
		secondTestSearcher1.EXPECT().Len().Return(3),
		secondTestSearcher2.EXPECT().Len().Return(4),
		secondTestSearcher1.EXPECT().Close().Return(nil),
		secondTestSearcher2.EXPECT().Close().Return(nil),
	)

	tests := []struct {
		name      string
		ss        search.Searchers
		expectErr bool
	}{
		{
			name:      "all searchers have the same length",
			ss:        firstTestSeachers,
			expectErr: false,
		},
		{
			name:      "all searchers do not have the same length",
			ss:        secondTestSeachers,
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateSearchers(test.ss)
			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
