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

package pilosa

import (
	"fmt"
	"testing"

	"github.com/m3db/m3ninx/postings"

	"github.com/stretchr/testify/require"
)

func TestRoaringPostingsListEmpty(t *testing.T) {
	d := NewPostingsList()
	require.True(t, d.IsEmpty())
	require.Equal(t, 0, d.Len())
}

func TestRoaringPostingsListMax(t *testing.T) {
	d := NewPostingsList()
	d.Insert(42)
	d.Insert(78)
	d.Insert(103)

	max, err := d.Max()
	require.NoError(t, err)
	require.Equal(t, postings.ID(103), max)

	d = NewPostingsList()
	_, err = d.Max()
	require.Error(t, err)
}

func TestRoaringPostingsListInsert(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())
	// Idempotency of inserts.
	d.Insert(1)
	require.Equal(t, 1, d.Len())
	require.True(t, d.Contains(1))
}

func TestRoaringPostingsListClone(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())

	c := d.Clone()
	require.True(t, c.Contains(1))
	require.Equal(t, 1, c.Len())

	// Ensure only clone is uniquely backed.
	c.Insert(2)
	require.True(t, c.Contains(2))
	require.Equal(t, 2, c.Len())
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())
}

func TestRoaringPostingsListIntersect(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())

	c := d.Clone()
	require.True(t, c.Contains(1))

	d.Insert(2)
	c.Insert(3)

	require.NoError(t, d.Intersect(c))
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())
	require.True(t, c.Contains(1))
	require.True(t, c.Contains(3))
	require.Equal(t, 2, c.Len())
}

func TestRoaringPostingsListDifference(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())

	c := d.Clone()
	require.True(t, c.Contains(1))

	d.Insert(2)
	d.Insert(3)
	require.NoError(t, d.Difference(c))

	require.False(t, d.Contains(1))
	require.True(t, c.Contains(1))
	require.Equal(t, 2, d.Len())
	require.Equal(t, 1, c.Len())
	require.True(t, d.Contains(3))
	require.True(t, d.Contains(2))
}

func TestRoaringPostingsListUnion(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())

	c := d.Clone()
	require.True(t, c.Contains(1))
	d.Insert(2)
	c.Insert(3)

	require.NoError(t, d.Union(c))
	require.True(t, d.Contains(1))
	require.True(t, d.Contains(2))
	require.True(t, d.Contains(3))
	require.Equal(t, 3, d.Len())
	require.True(t, c.Contains(1))
	require.True(t, c.Contains(3))
	require.Equal(t, 2, c.Len())
}

func TestRoaringPostingsListRemoveRange(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	d.Insert(2)
	d.Insert(7)
	d.Insert(9)

	d.RemoveRange(2, 8)
	require.Equal(t, 2, d.Len())
	require.True(t, d.Contains(1))
	require.False(t, d.Contains(2))
	require.False(t, d.Contains(7))
	require.True(t, d.Contains(9))
}

func TestRoaringPostingsListReset(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	require.True(t, d.Contains(1))
	require.Equal(t, 1, d.Len())
	d.Reset()
	require.True(t, d.IsEmpty())
	require.Equal(t, 0, d.Len())
}

func TestRoaringPostingsListIter(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	d.Insert(2)
	require.Equal(t, 2, d.Len())

	it := d.Iterator()
	defer it.Close()
	found := map[postings.ID]bool{
		1: false,
		2: false,
	}
	for it.Next() {
		found[it.Current()] = true
	}

	for id, ok := range found {
		require.True(t, ok, id)
	}
}

func TestRoaringPostingsListIterInsertAfter(t *testing.T) {
	d := NewPostingsList()
	d.Insert(1)
	d.Insert(2)
	require.Equal(t, 2, d.Len())

	it := d.Iterator()
	defer it.Close()
	numElems := 0
	d.Insert(3)
	require.Equal(t, 3, d.Len())
	found := map[postings.ID]bool{
		1: false,
		2: false,
	}
	for it.Next() {
		found[it.Current()] = true
		numElems++
	}

	for id, ok := range found {
		require.True(t, ok, fmt.Sprintf("%v", id))
	}
	require.Equal(t, 2, numElems)
}