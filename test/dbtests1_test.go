// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package dbprototest

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/cloud/datastore"
)

type testProto1Equivalent struct {
	//UnusedField               string
	SomeSingleString string      `datastore:"single_string"`
	SomeMultiString  []string    `datastore:"multi_string"`
	SomeSingleTime   time.Time   `datastore:"single_time"`
	SomeMultiInt     []int64     `datastore:"multi_int"`
	SomeMultiTimes   []time.Time `datastore:"multi_time"`
}

func TestDatastore_CanonicalToCodeGen(t *testing.T) {
	foo := &testProto1Equivalent{
		//UnusedField: "foo",
		SomeSingleString: "one",
		SomeMultiString:  []string{"one", "two", "three"},
		SomeSingleTime:   time.Unix(100, 0).UTC(),
		SomeMultiInt:     []int64{100, 200, 300},
		SomeMultiTimes:   []time.Time{time.Unix(100, 0), time.Unix(200, 0)},
	}

	props, err := datastore.SaveStruct(foo)
	require.NoError(t, err, "should not error on using upstream datastore yo")
	m := &TestProto1{}
	err = m.Load(props)
	require.NoError(t, err, "should not error reading a correct canonical code gen value")

	assert.Equal(t, foo.SomeSingleString, m.SomeSingleString, "single string passing must work")
	assert.Equal(t, foo.SomeMultiString, m.SomeMultiString, "multi string passing must work")
	outTstamp, _ := ptypes.Timestamp(m.SomeSingleTime)
	assert.Equal(t, foo.SomeSingleTime, outTstamp, "single time passing must work")
	assert.Equal(t, len(foo.SomeMultiTimes), len(m.SomeMultiTimes), "multi time passing must have the same number of stuff")
	assert.EqualValues(t, []int32{100, 200, 300}, m.SomeMultiInt, "multi int parsing must work")
}

func TestDatastore_CodeGenToCanonical(t *testing.T) {
	input := &TestProto1{
		//UnusedField: "foo",
		SomeSingleString: "one",
		SomeMultiString:  []string{"one", "two", "three"},
		SomeSingleTime:   &timestamp.Timestamp{Seconds: 100},
		SomeMultiInt:     []int32{100, 200, 300},
		SomeMultiTimes: []*timestamp.Timestamp{
			&timestamp.Timestamp{Seconds: 100},
			&timestamp.Timestamp{Seconds: 200},
		},
	}

	props, err := input.Save()
	require.NoError(t, err, "should not error on using upstream datastore yo")
	output := &testProto1Equivalent{}
	err = datastore.LoadStruct(output, props)
	require.NoError(t, err, "should not error reading a correct canonical code gen value")

	assert.Equal(t, input.SomeSingleString, output.SomeSingleString, "single string passing must work")
	assert.Equal(t, input.SomeMultiString, output.SomeMultiString, "multi string passing must work")
	expectedTime, _ := ptypes.Timestamp(input.SomeSingleTime)
	assert.Equal(t, expectedTime, output.SomeSingleTime, "single time passing must work")
	assert.Equal(t, len(input.SomeMultiTimes), len(output.SomeMultiTimes), "multi time passing must have the same number of stuff")
	assert.EqualValues(t, []int64{100, 200, 300}, output.SomeMultiInt, "multi int parsing must work")
}
