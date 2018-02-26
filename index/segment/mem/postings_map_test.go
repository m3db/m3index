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

package mem

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostingsMap(t *testing.T) {
	opts := NewOptions()
	pm := newPostingsMap(opts)

	require.NoError(t, pm.addID([]byte("foo"), 1))
	require.NoError(t, pm.addID([]byte("bar"), 2))
	require.NoError(t, pm.addID([]byte("baz"), 3))

	pl := pm.get([]byte("foo"))
	require.Equal(t, 1, pl.Len())
	require.True(t, pl.Contains(1))

	re := regexp.MustCompile("ba.*")
	pls := pm.getRegex(re)
	require.Equal(t, 2, len(pls))

	clone := pls[0].Clone()
	clone.Union(pls[1])
	require.True(t, clone.Contains(2))
	require.True(t, clone.Contains(3))
}
