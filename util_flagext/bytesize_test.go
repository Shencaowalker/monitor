package util_flagext

import (
	"testing"

	"encoding/json"

	"gopkg.in/yaml.v2"
)

func Test_ByteSize(t *testing.T) {
	for _, tc := range []struct {
		in  string
		err bool
		out int
	}{
		{
			in:  "abc",
			err: true,
		},
		{
			in:  "",
			err: false,
			out: 0,
		},
		{
			in:  "0",
			err: false,
			out: 0,
		},
		{
			in:  "1b",
			err: false,
			out: 1,
		},
		{
			in:  "100kb",
			err: false,
			out: 100 << 10,
		},
		{
			in:  "100 KB",
			err: false,
			out: 100 << 10,
		},
		{
			// ensure lowercase works
			in:  "50mb",
			err: false,
			out: 50 << 20,
		},
		{
			// ensure mixed capitalization works
			in:  "50Mb",
			err: false,
			out: 50 << 20,
		},
		{
			in:  "256GB",
			err: false,
			out: 256 << 30,
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			var bs ByteSize
			bs.Set(tc.in)
			// _ := bs.Set(tc.in)
			// if tc.err {
			// 	require.NotNil(t, err)
			// } else {
			// 	require.Nil(t, err)
			// 	require.Equal(t, tc.out, bs.Get().(int))

			// }

		})
	}
}

func Test_ByteSizeYAML(t *testing.T) {
	for _, tc := range []struct {
		in  string
		err bool
		out ByteSize
	}{
		{
			in:  "256GB",
			out: ByteSize(256 << 30),
		},
		{
			in:  "abc",
			err: true,
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			var out ByteSize
			yaml.Unmarshal([]byte(tc.in), &out)
			// err := yaml.Unmarshal([]byte(tc.in), &out)
			// if tc.err {
			// 	require.NotNil(t, err)
			// } else {
			// 	require.Nil(t, err)
			// 	require.Equal(t, tc.out, out)
			// }
		})
	}
}

func Test_ByteSizeJSON(t *testing.T) {
	for _, tc := range []struct {
		in  string
		err bool
		out ByteSize
	}{
		{
			in:  `{ "bytes": "256GB" }`,
			out: ByteSize(256 << 30),
		},
		{
			// JSON shouldn't allow to set integer as value for ByteSize field.
			in:  `{ "bytes": 2.62144e+07 }`,
			err: true,
		},
		{
			in:  `{ "bytes": "abc" }`,
			err: true,
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			var out struct {
				Bytes ByteSize `json:"bytes"`
			}
			json.Unmarshal([]byte(tc.in), &out)
			// err := json.Unmarshal([]byte(tc.in), &out)
			// if tc.err {
			// 	require.NotNil(t, err)
			// } else {
			// 	require.Nil(t, err)
			// 	require.Equal(t, tc.out, out.Bytes)
			// }
		})
	}
}
