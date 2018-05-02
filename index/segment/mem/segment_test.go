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
	re "regexp"
	"testing"

	"github.com/m3db/m3ninx/doc"
	"github.com/m3db/m3ninx/index"

	"github.com/stretchr/testify/require"
)

func TestSegmentInsert(t *testing.T) {
	tests := []struct {
		name  string
		input doc.Document
	}{
		{
			name: "document without an ID",
			input: doc.Document{
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("apple"),
						Value: []byte("red"),
					},
				},
			},
		},
		// {
		// 	name: "document with an ID",
		// 	input: doc.Document{
		// 		ID: []byte("123"),
		// 		Fields: []doc.Field{
		// 			doc.Field{
		// 				Name:  []byte("apple"),
		// 				Value: []byte("red"),
		// 			},
		// 		},
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			segment, err := NewSegment(0, NewOptions())
			require.NoError(t, err)

			id, err := segment.Insert(test.input)
			require.NoError(t, err)

			r, err := segment.Reader()
			require.NoError(t, err)
			println(fmt.Sprintf("here"))

			testDocument(t, test.input, r)

			// The ID must be searchable.
			pl, err := r.MatchTerm(doc.IDReservedFieldName, id)
			require.NoError(t, err)

			iter, err := r.Docs(pl)
			require.NoError(t, err)

			require.True(t, iter.Next())
			actual := iter.Current()

			require.True(t, compareDocs(test.input, actual))

			require.NoError(t, iter.Close())
			require.NoError(t, r.Close())
			require.NoError(t, segment.Close())
		})
	}
}

func TestSegmentInsertDuplicateID(t *testing.T) {
	var (
		id    = []byte("123")
		first = doc.Document{
			ID: id,
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("apple"),
					Value: []byte("red"),
				},
			},
		}
		second = doc.Document{
			ID: id,
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("apple"),
					Value: []byte("red"),
				},
				doc.Field{
					Name:  []byte("variety"),
					Value: []byte("fuji"),
				},
			},
		}
	)

	segment, err := NewSegment(0, NewOptions())
	require.NoError(t, err)

	_, err = segment.Insert(first)
	require.NoError(t, err)

	r, err := segment.Reader()
	require.NoError(t, err)

	pl, err := r.MatchTerm(doc.IDReservedFieldName, id)
	require.NoError(t, err)

	iter, err := r.Docs(pl)
	require.NoError(t, err)

	require.True(t, iter.Next())
	actual := iter.Current()

	// Only the first document should be indexed.
	require.True(t, compareDocs(first, actual))
	require.False(t, compareDocs(second, actual))

	require.NoError(t, iter.Close())
	require.NoError(t, r.Close())
	require.NoError(t, segment.Close())
}

func TestSegmentInsertBatch(t *testing.T) {
	tests := []struct {
		name  string
		input index.Batch
	}{
		{
			name: "valid batch",
			input: index.NewBatch(
				[]doc.Document{
					doc.Document{
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("apple"),
							},
							doc.Field{
								Name:  []byte("color"),
								Value: []byte("red"),
							},
						},
					},
					doc.Document{
						ID: []byte("831992"),
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("banana"),
							},
							doc.Field{
								Name:  []byte("color"),
								Value: []byte("yellow"),
							},
						},
					},
				},
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			segment, err := NewSegment(0, NewOptions())
			require.NoError(t, err)

			err = segment.InsertBatch(test.input)
			require.NoError(t, err)

			r, err := segment.Reader()
			require.NoError(t, err)

			for _, doc := range test.input.Docs {
				testDocument(t, doc, r)
			}

			require.NoError(t, r.Close())
			require.NoError(t, segment.Close())
		})
	}
}

func TestSegmentInsertBatchError(t *testing.T) {
	tests := []struct {
		name  string
		input index.Batch
	}{
		{
			name: "invalid document",
			input: index.NewBatch(
				[]doc.Document{
					doc.Document{
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("apple"),
							},
							doc.Field{
								Name:  []byte("color\xff"),
								Value: []byte("red"),
							},
						},
					},
					doc.Document{
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("banana"),
							},
							doc.Field{
								Name:  []byte("color"),
								Value: []byte("yellow"),
							},
						},
					},
				},
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			segment, err := NewSegment(0, NewOptions())
			require.NoError(t, err)

			err = segment.InsertBatch(test.input)
			require.Error(t, err)
			require.False(t, index.IsBatchPartialError(err))
		})
	}
}

func TestSegmentInsertBatchPartialError(t *testing.T) {
	tests := []struct {
		name  string
		input index.Batch
	}{
		{
			name: "invalid document",
			input: index.NewBatch(
				[]doc.Document{
					doc.Document{
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("apple"),
							},
							doc.Field{
								Name:  []byte("color\xff"),
								Value: []byte("red"),
							},
						},
					},
					doc.Document{

						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("banana"),
							},
							doc.Field{
								Name:  []byte("color"),
								Value: []byte("yellow"),
							},
						},
					},
				},
				index.AllowPartialUpdates(),
			),
		},
		{
			name: "duplicate ID",
			input: index.NewBatch(
				[]doc.Document{
					doc.Document{
						ID: []byte("831992"),
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("apple"),
							},
							doc.Field{
								Name:  []byte("color"),
								Value: []byte("red"),
							},
						},
					},
					doc.Document{
						ID: []byte("831992"),
						Fields: []doc.Field{
							doc.Field{
								Name:  []byte("fruit"),
								Value: []byte("banana"),
							},
							doc.Field{
								Name:  []byte("color"),
								Value: []byte("yellow"),
							},
						},
					},
				},
				index.AllowPartialUpdates(),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			segment, err := NewSegment(0, NewOptions())
			require.NoError(t, err)

			err = segment.InsertBatch(test.input)
			require.Error(t, err)
			require.True(t, index.IsBatchPartialError(err))

			batchErr := err.(*index.BatchPartialError)
			idxs := batchErr.Indices()
			failedDocs := make(map[int]struct{}, len(idxs))
			for _, idx := range idxs {
				failedDocs[idx] = struct{}{}
			}

			r, err := segment.Reader()
			require.NoError(t, err)

			for i, doc := range test.input.Docs {
				_, ok := failedDocs[i]
				if ok {
					// Don't test documents which were not indexed.
					continue
				}
				testDocument(t, doc, r)
			}

			require.NoError(t, r.Close())
			require.NoError(t, segment.Close())
		})
	}
}

func TestSegmentInsertBatchPartialErrorInvalidDoc(t *testing.T) {
	b1 := index.NewBatch(
		[]doc.Document{
			doc.Document{
				ID: []byte("abc"),
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("fruit"),
						Value: []byte("apple"),
					},
					doc.Field{
						Name:  []byte("color\xff"),
						Value: []byte("red"),
					},
				},
			},
			doc.Document{
				ID: []byte("abc"),
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("fruit"),
						Value: []byte("banana"),
					},
					doc.Field{
						Name:  []byte("color"),
						Value: []byte("yellow"),
					},
				},
			},
		},
		index.AllowPartialUpdates(),
	)
	segment, err := NewSegment(0, NewOptions())
	require.NoError(t, err)

	err = segment.InsertBatch(b1)
	require.Error(t, err)
	require.True(t, index.IsBatchPartialError(err))
	be := err.(*index.BatchPartialError)
	require.Len(t, be.Indices(), 1)
	require.Equal(t, be.Indices()[0], 0)

	r, err := segment.Reader()
	require.NoError(t, err)
	iter, err := r.AllDocs()
	require.NoError(t, err)
	require.True(t, iter.Next())
	require.Equal(t, b1.Docs[1], iter.Current())
	require.False(t, iter.Next())
	require.NoError(t, iter.Close())
	require.NoError(t, r.Close())
	require.NoError(t, segment.Close())
}

func TestSegmentInsertBatchPartialErrorAlreadyIndexing(t *testing.T) {
	b1 := index.NewBatch(
		[]doc.Document{
			doc.Document{
				ID: []byte("abc"),
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("fruit"),
						Value: []byte("apple"),
					},
					doc.Field{
						Name:  []byte("color"),
						Value: []byte("red"),
					},
				},
			},
		},
		index.AllowPartialUpdates())

	b2 := index.NewBatch(
		[]doc.Document{
			doc.Document{
				ID: []byte("abc"),
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("fruit"),
						Value: []byte("apple"),
					},
					doc.Field{
						Name:  []byte("color"),
						Value: []byte("red"),
					},
				},
			},
			doc.Document{
				ID: []byte("cdef"),
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("color"),
						Value: []byte("blue"),
					},
				},
			},
			doc.Document{
				ID: []byte("cdef"),
				Fields: []doc.Field{
					doc.Field{
						Name:  []byte("color"),
						Value: []byte("blue"),
					},
				},
			},
		},
		index.AllowPartialUpdates())

	segment, err := NewSegment(0, NewOptions())
	require.NoError(t, err)

	err = segment.InsertBatch(b1)
	require.NoError(t, err)

	err = segment.InsertBatch(b2)
	require.Error(t, err)
	require.True(t, index.IsBatchPartialError(err))
	ind := err.(*index.BatchPartialError).Indices()
	require.Len(t, ind, 1)
	require.Equal(t, 2, ind[0])
}

func TestSegmentReaderMatchExact(t *testing.T) {
	docs := []doc.Document{
		doc.Document{
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("fruit"),
					Value: []byte("apple"),
				},
				doc.Field{
					Name:  []byte("color"),
					Value: []byte("red"),
				},
			},
		},
		doc.Document{
			ID: []byte("83"),
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("fruit"),
					Value: []byte("banana"),
				},
				doc.Field{
					Name:  []byte("color"),
					Value: []byte("yellow"),
				},
			},
		},
		doc.Document{
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("fruit"),
					Value: []byte("apple"),
				},
				doc.Field{
					Name:  []byte("color"),
					Value: []byte("green"),
				},
			},
		},
	}

	segment, err := NewSegment(0, NewOptions())
	require.NoError(t, err)

	for _, doc := range docs {
		_, err = segment.Insert(doc)
		require.NoError(t, err)
	}

	r, err := segment.Reader()
	require.NoError(t, err)

	pl, err := r.MatchTerm([]byte("fruit"), []byte("apple"))
	require.NoError(t, err)

	iter, err := r.Docs(pl)
	require.NoError(t, err)

	actualDocs := make([]doc.Document, 0)
	for iter.Next() {
		actualDocs = append(actualDocs, iter.Current())
	}

	require.NoError(t, iter.Err())
	require.NoError(t, iter.Close())

	expectedDocs := []doc.Document{docs[0], docs[2]}
	require.Equal(t, len(expectedDocs), len(actualDocs))
	for i := range actualDocs {
		require.True(t, compareDocs(expectedDocs[i], actualDocs[i]))
	}

	require.NoError(t, r.Close())
	require.NoError(t, segment.Close())
}

func TestSegmentReaderMatchRegex(t *testing.T) {
	docs := []doc.Document{
		doc.Document{
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("fruit"),
					Value: []byte("banana"),
				},
				doc.Field{
					Name:  []byte("color"),
					Value: []byte("yellow"),
				},
			},
		},
		doc.Document{
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("fruit"),
					Value: []byte("apple"),
				},
				doc.Field{
					Name:  []byte("color"),
					Value: []byte("red"),
				},
			},
		},
		doc.Document{
			ID: []byte("42"),
			Fields: []doc.Field{
				doc.Field{
					Name:  []byte("fruit"),
					Value: []byte("pineapple"),
				},
				doc.Field{
					Name:  []byte("color"),
					Value: []byte("yellow"),
				},
			},
		},
	}

	segment, err := NewSegment(0, NewOptions())
	require.NoError(t, err)

	for _, doc := range docs {
		_, err = segment.Insert(doc)
		require.NoError(t, err)
	}

	r, err := segment.Reader()
	require.NoError(t, err)

	field, regexp := []byte("fruit"), []byte(".*ple")
	compiled := re.MustCompile(string(regexp))
	pl, err := r.MatchRegexp(field, regexp, compiled)
	require.NoError(t, err)

	iter, err := r.Docs(pl)
	require.NoError(t, err)

	actualDocs := make([]doc.Document, 0)
	for iter.Next() {
		actualDocs = append(actualDocs, iter.Current())
	}

	require.NoError(t, iter.Err())
	require.NoError(t, iter.Close())

	expectedDocs := []doc.Document{docs[1], docs[2]}
	require.Equal(t, len(expectedDocs), len(actualDocs))
	for i := range actualDocs {
		require.True(t, compareDocs(expectedDocs[i], actualDocs[i]))
	}

	require.NoError(t, r.Close())
	require.NoError(t, segment.Close())
}

func testDocument(t *testing.T, d doc.Document, r index.Reader) {
	for _, f := range d.Fields {
		name, value := f.Name, f.Value
		pl, err := r.MatchTerm(name, value)
		require.NoError(t, err)

		iter, err := r.Docs(pl)
		require.NoError(t, err)

		require.True(t, iter.Next())
		actual := iter.Current()

		// The document must have an ID.
		hasID := actual.ID != nil
		require.True(t, hasID)

		require.True(t, compareDocs(d, actual))

		require.False(t, iter.Next())
		require.NoError(t, iter.Err())
		require.NoError(t, iter.Close())
	}
}

// compareDocs returns whether two documents are equal. If the actual doc contains
// an ID but the expected doc does not then the ID is excluded from the comparison
// since it was auto-generated.
func compareDocs(expected, actual doc.Document) bool {
	if actual.HasID() && !expected.HasID() {
		actual.ID = nil
	}
	return expected.Equal(actual)
}
