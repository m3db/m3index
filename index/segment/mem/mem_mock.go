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

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3ninx/index/segment/mem/types.go

package mem

import (
	"reflect"
	"regexp"

	"github.com/m3db/m3ninx/doc"
	"github.com/m3db/m3ninx/postings"

	"github.com/golang/mock/gomock"
)

// MocktermsDict is a mock of termsDict interface
type MocktermsDict struct {
	ctrl     *gomock.Controller
	recorder *MocktermsDictMockRecorder
}

// MocktermsDictMockRecorder is the mock recorder for MocktermsDict
type MocktermsDictMockRecorder struct {
	mock *MocktermsDict
}

// NewMocktermsDict creates a new mock instance
func NewMocktermsDict(ctrl *gomock.Controller) *MocktermsDict {
	mock := &MocktermsDict{ctrl: ctrl}
	mock.recorder = &MocktermsDictMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MocktermsDict) EXPECT() *MocktermsDictMockRecorder {
	return _m.recorder
}

// Insert mocks base method
func (_m *MocktermsDict) Insert(field doc.Field, id postings.ID) error {
	ret := _m.ctrl.Call(_m, "Insert", field, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (_mr *MocktermsDictMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Insert", reflect.TypeOf((*MocktermsDict)(nil).Insert), arg0, arg1)
}

// MatchExact mocks base method
func (_m *MocktermsDict) MatchExact(name []byte, value []byte) (postings.List, error) {
	ret := _m.ctrl.Call(_m, "MatchExact", name, value)
	ret0, _ := ret[0].(postings.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MatchExact indicates an expected call of MatchExact
func (_mr *MocktermsDictMockRecorder) MatchExact(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "MatchExact", reflect.TypeOf((*MocktermsDict)(nil).MatchExact), arg0, arg1)
}

// MatchRegex mocks base method
func (_m *MocktermsDict) MatchRegex(name []byte, pattern []byte, re *regexp.Regexp) (postings.List, error) {
	ret := _m.ctrl.Call(_m, "MatchRegex", name, pattern, re)
	ret0, _ := ret[0].(postings.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MatchRegex indicates an expected call of MatchRegex
func (_mr *MocktermsDictMockRecorder) MatchRegex(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "MatchRegex", reflect.TypeOf((*MocktermsDict)(nil).MatchRegex), arg0, arg1, arg2)
}

// MockreadableSegment is a mock of readableSegment interface
type MockreadableSegment struct {
	ctrl     *gomock.Controller
	recorder *MockreadableSegmentMockRecorder
}

// MockreadableSegmentMockRecorder is the mock recorder for MockreadableSegment
type MockreadableSegmentMockRecorder struct {
	mock *MockreadableSegment
}

// NewMockreadableSegment creates a new mock instance
func NewMockreadableSegment(ctrl *gomock.Controller) *MockreadableSegment {
	mock := &MockreadableSegment{ctrl: ctrl}
	mock.recorder = &MockreadableSegmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockreadableSegment) EXPECT() *MockreadableSegmentMockRecorder {
	return _m.recorder
}

// matchExact mocks base method
func (_m *MockreadableSegment) matchExact(name []byte, value []byte) (postings.List, error) {
	ret := _m.ctrl.Call(_m, "matchExact", name, value)
	ret0, _ := ret[0].(postings.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// matchExact indicates an expected call of matchExact
func (_mr *MockreadableSegmentMockRecorder) matchExact(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "matchExact", reflect.TypeOf((*MockreadableSegment)(nil).matchExact), arg0, arg1)
}

// matchRegex mocks base method
func (_m *MockreadableSegment) matchRegex(name []byte, pattern []byte, re *regexp.Regexp) (postings.List, error) {
	ret := _m.ctrl.Call(_m, "matchRegex", name, pattern, re)
	ret0, _ := ret[0].(postings.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// matchRegex indicates an expected call of matchRegex
func (_mr *MockreadableSegmentMockRecorder) matchRegex(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "matchRegex", reflect.TypeOf((*MockreadableSegment)(nil).matchRegex), arg0, arg1, arg2)
}

// getDoc mocks base method
func (_m *MockreadableSegment) getDoc(id postings.ID) (doc.Document, error) {
	ret := _m.ctrl.Call(_m, "getDoc", id)
	ret0, _ := ret[0].(doc.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getDoc indicates an expected call of getDoc
func (_mr *MockreadableSegmentMockRecorder) getDoc(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "getDoc", reflect.TypeOf((*MockreadableSegment)(nil).getDoc), arg0)
}
