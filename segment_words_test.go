//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package segment

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestAdhocSegmentsWithType(t *testing.T) {

	tests := []struct {
		input         []byte
		output        [][]byte
		outputStrings []string
		outputTypes   []int
	}{
		{
			input: []byte("Now is the.\n End."),
			output: [][]byte{
				[]byte("Now"),
				[]byte(" "),
				[]byte("is"),
				[]byte(" "),
				[]byte("the"),
				[]byte("."),
				[]byte("\n"),
				[]byte(" "),
				[]byte("End"),
				[]byte("."),
			},
			outputStrings: []string{
				"Now",
				" ",
				"is",
				" ",
				"the",
				".",
				"\n",
				" ",
				"End",
				".",
			},
			outputTypes: []int{
				Letter,
				None,
				Letter,
				None,
				Letter,
				None,
				None,
				None,
				Letter,
				None,
			},
		},
		{
			input: []byte("3.5"),
			output: [][]byte{
				[]byte("3.5"),
			},
			outputStrings: []string{
				"3.5",
			},
			outputTypes: []int{
				Number,
			},
		},
		{
			input: []byte("cat3.5"),
			output: [][]byte{
				[]byte("cat3.5"),
			},
			outputStrings: []string{
				"cat3.5",
			},
			outputTypes: []int{
				Letter,
			},
		},
		{
			input: []byte("c"),
			output: [][]byte{
				[]byte("c"),
			},
			outputStrings: []string{
				"c",
			},
			outputTypes: []int{
				Letter,
			},
		},
		{
			input: []byte("こんにちは世界"),
			output: [][]byte{
				[]byte("こ"),
				[]byte("ん"),
				[]byte("に"),
				[]byte("ち"),
				[]byte("は"),
				[]byte("世"),
				[]byte("界"),
			},
			outputStrings: []string{
				"こ",
				"ん",
				"に",
				"ち",
				"は",
				"世",
				"界",
			},
			outputTypes: []int{
				Ideo,
				Ideo,
				Ideo,
				Ideo,
				Ideo,
				Ideo,
				Ideo,
			},
		},
		{
			input: []byte("你好世界"),
			output: [][]byte{
				[]byte("你"),
				[]byte("好"),
				[]byte("世"),
				[]byte("界"),
			},
			outputStrings: []string{
				"你",
				"好",
				"世",
				"界",
			},
			outputTypes: []int{
				Ideo,
				Ideo,
				Ideo,
				Ideo,
			},
		},
		{
			input: []byte("サッカ"),
			output: [][]byte{
				[]byte("サッカ"),
			},
			outputStrings: []string{
				"サッカ",
			},
			outputTypes: []int{
				Ideo,
			},
		},
	}

	for _, test := range tests {
		rv := make([][]byte, 0)
		rvstrings := make([]string, 0)
		rvtypes := make([]int, 0)
		segmenter := NewWordSegmenter(bytes.NewReader(test.input))
		// Set the split function for the scanning operation.
		for segmenter.Segment() {
			rv = append(rv, segmenter.Bytes())
			rvstrings = append(rvstrings, segmenter.Text())
			rvtypes = append(rvtypes, segmenter.Type())
		}
		if err := segmenter.Err(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(rv, test.output) {
			t.Fatalf("expected:\n%#v\ngot:\n%#v\nfor: '%s'", test.output, rv, test.input)
		}
		if !reflect.DeepEqual(rvstrings, test.outputStrings) {
			t.Fatalf("expected:\n%#v\ngot:\n%#v\nfor: '%s'", test.outputStrings, rvstrings, test.input)
		}
		if !reflect.DeepEqual(rvtypes, test.outputTypes) {
			t.Fatalf("expeced:\n%#v\ngot:\n%#v\nfor: '%s'", test.outputTypes, rvtypes, test.input)
		}
	}

}

func TestUnicodeSegments(t *testing.T) {

	for _, test := range unicodeWordTests {
		rv := make([][]byte, 0)
		scanner := bufio.NewScanner(bytes.NewReader(test.input))
		// Set the split function for the scanning operation.
		scanner.Split(SplitWords)
		for scanner.Scan() {
			rv = append(rv, scanner.Bytes())
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(rv, test.output) {
			t.Fatalf("expected:\n%#v\ngot:\n%#v\nfor: '%s'", test.output, rv, test.input)
		}
	}
}
