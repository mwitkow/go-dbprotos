// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package dbprototest

import (
	"time"
	"testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/cloud/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/golang/protobuf/ptypes"
)

type testProto1Equivalent struct {
	//UnusedField               string
	SomeSingleString string    `datastore:"single_string"`
	SomeMultiString  []string  `datastore:"multi_string"`
	SomeSingleTime   time.Time `datastore:"single_time"`
	SomeMultiInt   []int64 `datastore:"multi_int"`
	SomeMultiTime []time.Time `datastore:"multi_time"`

}

func TestDatastore_CanonicalToCodeGen(t *testing.T) {
	foo := &testProto1Equivalent{
		//UnusedField: "foo",
		SomeSingleString: "one",
		SomeMultiString: []string{"one", "two", "three"},
		SomeSingleTime: time.Unix(100, 0).UTC(),
		SomeMultiInt: []int64{100, 200, 300},
		SomeMultiTime: []time.Time{time.Unix(100, 0), time.Unix(200, 0)},
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
	assert.Equal(t, len(foo.SomeMultiTime), len(m.SomeMultiTimes), "multi time passing must have the same number of stuff")
	assert.EqualValues(t, []int32{100, 200, 300}, m.SomeMultiInt, "multi int parsing must work")
}