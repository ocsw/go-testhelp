/*
Copyright 2021 Danielle Zephyr Malament

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testhelp

import (
	"errors"
	"testing"
)

// Tests Panics(), PanicsGet(), PanicsStr(), and PanicsRE()
func TestPanicsX4(t *testing.T) {
	var didPanic bool
	var pContainsStr bool
	var pMatchesRE bool
	var pVal interface{}

	pErr := errors.New("ppp123")

	tests := []struct {
		name string
		f    func()
		// input to the function being tested
		inputStr string
		inputRE  string
		// want from the function being tested
		wantPanics   bool
		wantContains bool
		wantMatches  bool
		wantPVal     interface{}
	}{
		// panic, with correct string pVal
		{"p correct str", func() { panic("ppp123") }, "ppp", "p{3}[0-9]{3}", true, true, true, "ppp123"},

		// panic, with a string pVal and empty comparison strings
		{"p str, empty", func() { panic("ppp123") }, "", "", true, true, true, "ppp123"},

		// panic, with wrong string pVal
		{"p wrong str", func() { panic("ppp123") }, "1234", "p{3}[0-9]{4}", true, false, false, "ppp123"},

		// panic, with correct error pVal treated as string
		{"p correct err", func() { panic(pErr) }, "ppp", "p{3}[0-9]{3}", true, true, true, pErr},

		// panic, with an error pVal treated as a string, and empty comparison strings
		{"p err, empty", func() { panic(pErr) }, "", "", true, true, true, pErr},

		// panic, with wrong error pVal treated as string
		{"p wrong err", func() { panic(pErr) }, "1234", "p{3}[0-9]{4}", true, false, false, pErr},

		// panic, with non-string pVal
		{"p non-str", func() { panic(27.5) }, "ppp", "p{3}[0-9]{3}", true, false, false, 27.5},

		// panic, with non-string pVal and empty comparison strings
		{"p non-str, empty", func() { panic(27.5) }, "", "", true, false, false, 27.5},

		// non-panic
		{"np", func() {}, "ppp", "p{3}[0-9]{3}", false, false, false, nil},
	}
	for _, test := range tests {
		// test Panics()
		didPanic = Panics(test.f)
		if didPanic != test.wantPanics {
			if test.wantPanics == true {
				t.Errorf("Panics(): Expected function to panic in test '%s'", test.name)
			} else {
				t.Errorf("Panics(): Expected function not to panic in test '%s'", test.name)
			}
		}

		// test PanicsGet()
		didPanic, pVal = PanicsGet(test.f)
		if didPanic != test.wantPanics {
			if test.wantPanics == true {
				t.Errorf("PanicsGet(): Expected function to panic in test '%s'", test.name)
			} else {
				t.Errorf("PanicsGet(): Expected function not to panic in test '%s'", test.name)
			}
		}
		if pVal != test.wantPVal {
			t.Errorf("PanicsGet(): Incorrect panic value: expected\n%#+v\ngot\n%#+v\nin test '%s'",
				test.wantPVal, pVal, test.name)
		}

		// test PanicsStr()
		didPanic, pContainsStr, pVal = PanicsStr(test.f, test.inputStr)
		if didPanic != test.wantPanics {
			if test.wantPanics == true {
				t.Errorf("PanicsStr(): Expected function to panic in test '%s'", test.name)
			} else {
				t.Errorf("PanicsStr(): Expected function not to panic in test '%s'", test.name)
			}
		}
		if pContainsStr != test.wantContains {
			if test.wantContains == true {
				t.Errorf("PanicsStr(): Expected panic value to contain string in test '%s'", test.name)
			} else {
				t.Errorf("PanicsStr(): Expected panic value not to contain string in test '%s'", test.name)
			}
		}
		if pVal != test.wantPVal {
			t.Errorf("PanicsStr(): Incorrect panic value: expected string containing\n%#+v\ngot\n%#+v\nin test '%s'",
				test.wantPVal, pVal, test.name)
		}

		// test PanicsRE()
		didPanic, pMatchesRE, pVal = PanicsRE(test.f, test.inputRE)
		if didPanic != test.wantPanics {
			if test.wantPanics == true {
				t.Errorf("PanicsRE(): Expected function to panic in test '%s'", test.name)
			} else {
				t.Errorf("PanicsRE(): Expected function not to panic in test '%s'", test.name)
			}
		}
		if pMatchesRE != test.wantMatches {
			if test.wantMatches == true {
				t.Errorf("PanicsRE(): Expected panic value to match regexp in test '%s'", test.name)
			} else {
				t.Errorf("PanicsRE(): Expected panic value not to match regexp in test '%s'", test.name)
			}
		}
		if pVal != test.wantPVal {
			t.Errorf("PanicsRE(): Incorrect panic value: expected string matching\n%#+v\ngot\n%#+v\nin test '%s'",
				test.wantPVal, pVal, test.name)
		}
	}
}

func TestPanicsREPanicsWithBadRE(t *testing.T) {
	var didPanic bool
	var pContainsStr bool
	var pVal interface{}

	badRE := "[a-z" // no closing ]
	// want this from PanicsStr(), while testing a func containing PanicsRE()
	wantStr := "Regexp could not be compiled"

	tests := []struct {
		name string
		f    func()
	}{
		{"string pVal", func() { PanicsRE(func() { panic("ppp") }, badRE) }},
		{"non-string pVal", func() { PanicsRE(func() { panic(27) }, badRE) }},
	}
	for _, test := range tests {
		// It's a little suspect to use PanicsStr() here, but PanicsStr() and PanicsRE() don't reference each other, and
		// we have other tests for PanicsStr()
		didPanic, pContainsStr, pVal = PanicsStr(test.f, wantStr)
		if !didPanic {
			t.Fatalf("Expected PanicsRE() itself to panic in test '%s'", test.name)
		} else if !pContainsStr {
			t.Fatalf("Incorrect panic value from PanicsRE() itself: expected string containing\n%#+v\ngot\n%#+v\nin "+
				"test '%s'", wantStr, pVal, test.name)
		}
	}
}

func TestPanicsVal(t *testing.T) {
	var didPanic bool
	var pEquals bool
	var pVal interface{}

	tests := []struct {
		name string
		f    func()
		// input to PanicsVal()
		inputVal interface{}
		// want from PanicsVal()
		wantPanics bool
		wantEquals bool
		wantPVal   interface{}
	}{
		// panic, with correct string pVal
		{"p correct str", func() { panic("ppp123") }, "ppp123", true, true, "ppp123"},

		// panic, with wrong string pVal
		{"p wrong str", func() { panic("ppp123") }, "ppp234", true, false, "ppp123"},

		// panic, with correct float pVal
		{"p correct float", func() { panic(27.5) }, 27.5, true, true, 27.5},

		// panic, with wrong float pVal
		{"p wrong float", func() { panic(27.5) }, 42.1, true, false, 27.5},

		// panic, with float pVal and string inputVal
		{"p float/str", func() { panic(27.5) }, "27.5", true, false, 27.5},

		// panic, with float pVal and int inputVal
		{"p float/int", func() { panic(27.0) }, 27, true, false, 27.0},

		// non-panic
		{"np", func() {}, "ppp", false, false, nil},
	}
	for _, test := range tests {
		didPanic, pEquals, pVal = PanicsVal(test.f, test.inputVal)
		if didPanic != test.wantPanics {
			if test.wantPanics == true {
				t.Errorf("PanicsVal(): Expected function to panic in test '%s'", test.name)
			} else {
				t.Errorf("PanicsVal(): Expected function not to panic in test '%s'", test.name)
			}
		}
		if pEquals != test.wantEquals {
			if test.wantEquals == true {
				t.Errorf("PanicsVal(): Expected panic value to equal input value in test '%s'", test.name)
			} else {
				t.Errorf("PanicsVal(): Expected panic value not to equal input value in test '%s'", test.name)
			}
		}
		if pVal != test.wantPVal {
			t.Errorf("PanicsVal(): Incorrect panic value: expected\n%#+v\ngot\n%#+v\nin test '%s'",
				test.wantPVal, pVal, test.name)
		}
	}
}

func TestPanicsValPanicsWithUncomparableType(t *testing.T) {
	var didPanic bool
	var pContainsStr bool
	var pVal interface{}

	// want this from PanicsStr(), while testing a func containing PanicsVal()
	wantStr := "runtime error: comparing uncomparable type"

	tests := []struct {
		name string
		f    func()
	}{
		{
			"p string slices, correct", func() {
				PanicsVal(func() { panic([]string{"a", "b"}) }, []string{"a", "b"})
			},
		},
		{
			"p string slices, wrong", func() {
				PanicsVal(func() { panic([]string{"a", "b"}) }, []string{"c", "d"})
			},
		},
	}
	for _, test := range tests {
		// It's a little suspect to use PanicsStr() here, but PanicsStr() and PanicsVal() don't reference each other,
		// and we have other tests for PanicsStr()
		didPanic, pContainsStr, pVal = PanicsStr(test.f, wantStr)
		if !didPanic {
			t.Fatalf("Expected PanicsVal() itself to panic in test '%s'", test.name)
		} else if !pContainsStr {
			t.Fatalf("Incorrect panic value from PanicsVal() itself: expected string containing\n%#+v\ngot\n%#+v\nin "+
				"test '%s'", wantStr, pVal, test.name)
		}
	}
}

// Tests NotPanics() and NotPanicsGet()
func TestNotPanicsX2(t *testing.T) {
	var didNotPanic bool
	var pVal interface{}

	tests := []struct {
		name          string
		f             func()
		wantNotPanics bool
		wantPVal      interface{}
	}{
		{"np", func() {}, true, nil},
		{"p string", func() { panic("ppp") }, false, "ppp"},
		{"p non-string", func() { panic(27) }, false, 27},
	}
	for _, test := range tests {
		// test NotPanics()
		didNotPanic = NotPanics(test.f)
		if didNotPanic != test.wantNotPanics {
			if test.wantNotPanics == true {
				t.Errorf("NotPanics(): Expected function not to panic in test '%s'", test.name)
			} else {
				t.Errorf("NotPanics(): Expected function to panic in test '%s'", test.name)
			}
		}

		// test NotPanicsGet()
		didNotPanic, pVal = NotPanicsGet(test.f)
		if didNotPanic != test.wantNotPanics {
			if test.wantNotPanics == true {
				t.Errorf("NotPanicsGet(): Expected function not to panic in test '%s'", test.name)
			} else {
				t.Errorf("NotPanicsGet(): Expected function to panic in test '%s'", test.name)
			}
		}
		if pVal != test.wantPVal {
			t.Errorf("NotPanicsGet(): Incorrect panic value: expected\n%#+v\ngot\n%#+v\nin test '%s'",
				test.wantPVal, pVal, test.name)
		}
	}
}

type PanicStrRETest struct {
	Name    string
	F       func()
	WantStr string
	WantRE  string
}

// Tests PanicsLoop(), PanicsStrLoop(), and PanicsRELoop()
func TestPanicsLoopX3(t *testing.T) {
	var noPanic []string
	var noCM []string // CM = Contains/Matches
	var plainTable []PanicTest
	var strTable []PanicStrTest
	var reTable []PanicRETest

	pErr := errors.New("ppp110")
	notPanicFunc := func(testName string) { noPanic = append(noPanic, testName) }
	notCMFunc := func(testName string) { noCM = append(noCM, testName) }

	tests := []struct {
		name        string
		pTable      []PanicStrRETest
		wantNoPanic []string
		wantNoCM    []string
	}{
		{
			"p, cm; p, cm",
			[]PanicStrRETest{
				// Handle the error -> string case while we're at it, for PanicsStrLoop() and PanicsRELoop()
				{"p, cm; p, cm: 1", func() { panic(pErr) }, "ppp", "p{3}[0-9]{3}"},
				{"p, cm; p, cm: 2", func() { panic("ppp111") }, "ppp", "p{3}[0-9]{3}"},
			},
			[]string{},
			[]string{},
		},
		{
			"p, cm; p, ncm",
			[]PanicStrRETest{
				{"p, cm; p, ncm: 1", func() { panic("ppp120") }, "ppp", "p{3}[0-9]{3}"},
				{"p, cm; p, ncm: 2", func() { panic("ppp121") }, "zzz", "z{3}[0-9]{3}"},
			},
			[]string{},
			[]string{"p, cm; p, ncm: 2"},
		},
		{
			"p, cm; np",
			[]PanicStrRETest{
				{"p, cm; np: 1", func() { panic("ppp130") }, "ppp", "p{3}[0-9]{3}"},
				{"p, cm; np: 2", func() {}, "ppp", "p{3}[0-9]{3}"},
			},
			[]string{"p, cm; np: 2"},
			[]string{},
		},

		{
			"p, ncm; p, cm",
			[]PanicStrRETest{
				{"p, ncm; p, cm: 1", func() { panic("ppp210") }, "ccc", "c{3}[0-9]{3}"},
				{"p, ncm; p, cm: 2", func() { panic("ppp211") }, "ppp", "p{3}[0-9]{3}"},
			},
			[]string{},
			[]string{"p, ncm; p, cm: 1"},
		},
		{
			"p, ncm; p, ncm",
			[]PanicStrRETest{
				{"p, ncm; p, ncm: 1", func() { panic("ppp220") }, "ccc", "c{3}[0-9]{3}"},
				{"p, ncm; p, ncm: 2", func() { panic("ppp221") }, "zzz", "z{3}[0-9]{3}"},
			},
			[]string{},
			[]string{"p, ncm; p, ncm: 1", "p, ncm; p, ncm: 2"},
		},
		{
			"p, ncm; np",
			[]PanicStrRETest{
				{"p, ncm; np: 1", func() { panic("ppp230") }, "ccc", "c{3}[0-9]{3}"},
				{"p, ncm; np: 2", func() {}, "ppp", "p{3}[0-9]{3}"},
			},
			[]string{"p, ncm; np: 2"},
			[]string{"p, ncm; np: 1"},
		},

		{
			"np; p, cm",
			[]PanicStrRETest{
				{"np; p, cm: 1", func() {}, "ppp", "p{3}[0-9]{3}"},
				{"np; p, cm: 2", func() { panic("ppp311") }, "ppp", "p{3}[0-9]{3}"},
			},
			[]string{"np; p, cm: 1"},
			[]string{},
		},
		{
			"np; p, ncm",
			[]PanicStrRETest{
				{"np; p, ncm: 1", func() {}, "ppp", "p{3}[0-9]{3}"},
				{"np; p, ncm: 2", func() { panic("ppp321") }, "zzz", "z{3}[0-9]{3}"},
			},
			[]string{"np; p, ncm: 1"},
			[]string{"np; p, ncm: 2"},
		},
		{
			"np; np",
			[]PanicStrRETest{
				{"np; np: 1", func() {}, "ppp", "p{3}[0-9]{3}"},
				{"np; np: 2", func() {}, "ppp", "p{3}[0-9]{3}"},
			},
			[]string{"np; np: 1", "np; np: 2"},
			[]string{},
		},
	}
	for _, test := range tests {
		// test PanicsLoop()
		noPanic = nil
		plainTable = []PanicTest{}
		for _, tableEntry := range test.pTable {
			plainTable = append(plainTable, PanicTest{tableEntry.Name, tableEntry.F})
		}
		PanicsLoop(plainTable, notPanicFunc)
		if len(noPanic) != len(test.wantNoPanic) {
			t.Errorf("PanicsLoop(): Wrong number of panic-test failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoPanic), len(noPanic), test.name, test.wantNoPanic, noPanic)
		} else {
			for i := 0; i < len(noPanic); i++ {
				if noPanic[i] != test.wantNoPanic[i] {
					t.Errorf("PanicsLoop(): Wrong panic-test failure: expected '%s', got '%s'",
						test.wantNoPanic[i], noPanic[i])
				}
			}
		}

		// test PanicsStrLoop()
		noPanic = nil
		noCM = nil
		strTable = []PanicStrTest{}
		for _, tableEntry := range test.pTable {
			strTable = append(strTable, PanicStrTest{tableEntry.Name, tableEntry.F, tableEntry.WantStr})
		}
		PanicsStrLoop(strTable, notPanicFunc, notCMFunc)
		if len(noPanic) != len(test.wantNoPanic) {
			t.Errorf("PanicsStrLoop(): Wrong number of panic-test failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoPanic), len(noPanic), test.name, test.wantNoPanic, noPanic)
		} else {
			for i := 0; i < len(noPanic); i++ {
				if noPanic[i] != test.wantNoPanic[i] {
					t.Errorf("PanicsStrLoop(): Wrong panic-test failure: expected '%s', got '%s'",
						test.wantNoPanic[i], noPanic[i])
				}
			}
		}
		if len(noCM) != len(test.wantNoCM) {
			t.Errorf("PanicsStrLoop(): Wrong number of panic-contains failures: expected %d, got %d in test table "+
				"'%s'\nExpected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoCM), len(noCM), test.name, test.wantNoCM, noCM)
		} else {
			for i := 0; i < len(noCM); i++ {
				if noCM[i] != test.wantNoCM[i] {
					t.Errorf("PanicsStrLoop(): Wrong panic-contains failure: expected '%s', got '%s'",
						test.wantNoCM[i], noCM[i])
				}
			}
		}

		// test PanicsRELoop()
		noPanic = nil
		noCM = nil
		reTable = []PanicRETest{}
		for _, tableEntry := range test.pTable {
			reTable = append(reTable, PanicRETest{tableEntry.Name, tableEntry.F, tableEntry.WantRE})
		}
		PanicsRELoop(reTable, notPanicFunc, notCMFunc)
		if len(noPanic) != len(test.wantNoPanic) {
			t.Errorf("PanicsRELoop(): Wrong number of panic-test failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoPanic), len(noPanic), test.name, test.wantNoPanic, noPanic)
		} else {
			for i := 0; i < len(noPanic); i++ {
				if noPanic[i] != test.wantNoPanic[i] {
					t.Errorf("PanicsRELoop(): Wrong panic-test failure: expected '%s', got '%s'",
						test.wantNoPanic[i], noPanic[i])
				}
			}
		}
		if len(noCM) != len(test.wantNoCM) {
			t.Errorf("PanicsRELoop(): Wrong number of panic-matches failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoCM), len(noCM), test.name, test.wantNoCM, noCM)
		} else {
			for i := 0; i < len(noCM); i++ {
				if noCM[i] != test.wantNoCM[i] {
					t.Errorf("PanicsRELoop(): Wrong panic-matches failure: expected '%s', got '%s'",
						test.wantNoCM[i], noCM[i])
				}
			}
		}
	}
}

func TestPanicsRELoopPanicsWithBadRE(t *testing.T) {
	// for PanicsStr(), while testing a func containing PanicsRELoop()
	var didPanic bool
	var pContainsStr bool
	var pVal interface{}
	wantStr := "Regexp could not be compiled"

	// for the PanicsRELoop() being run by PanicsStr()
	var noPanic []string
	var noMatches []string
	badRE1 := "[a-z" // no closing ]
	badRE2 := "[0-9" // no closing ]
	notPanicFunc := func(testName string) { noPanic = append(noPanic, testName) }
	notMatchesFunc := func(testName string) { noMatches = append(noMatches, testName) }

	tests := []struct {
		name          string
		pTable        []PanicRETest
		wantNoMatches []string
	}{
		{
			"ok, not ok",
			[]PanicRETest{
				// ok but wrong
				{"ok, not ok: 1", func() { panic("ppp111") }, "c{3}[0-9]{3}"},
				{"ok, not ok: 2", func() { panic("ppp112") }, badRE2},
			},
			// first test within PanicsRELoop() proceeds normally, second one panics
			[]string{"ok, not ok: 1"},
		},
		{
			"not ok, ok",
			[]PanicRETest{
				{"not ok, ok: 1", func() { panic("ppp221") }, badRE1},
				// ok but wrong
				{"not ok, ok: 2", func() { panic("ppp222") }, "z{3}[0-9]{3}"},
			},
			[]string{},
		},
		{
			"not ok, not ok",
			[]PanicRETest{
				{"not ok, not ok: 1", func() { panic("ppp331") }, badRE1},
				{"not ok, not ok: 2", func() { panic("ppp332") }, badRE2},
			},
			[]string{},
		},
	}
	for _, test := range tests {
		noPanic = nil
		noMatches = nil

		// It's a little suspect to use PanicsStr() here, but PanicsStr() and PanicsRELoop() don't reference each
		// other, and we have other tests for PanicsStr()
		didPanic, pContainsStr, pVal = PanicsStr(func() {
			// nolint: scopelint
			PanicsRELoop(test.pTable, notPanicFunc, notMatchesFunc)
		}, wantStr)
		if !didPanic {
			t.Errorf("Expected PanicsRELoop() itself to panic in test table '%s'", test.name)
		} else if !pContainsStr {
			t.Errorf("Incorrect panic value from PanicsRELoop() itself: expected string containing\n"+
				"%#+v\ngot\n%#+v\nin test table '%s'", wantStr, pVal, test.name)
		}

		// Now test the results of PanicsRELoop() itself
		if len(noPanic) != 0 {
			t.Errorf("PanicsRELoop(): Unexpected panic-test failure(s): expected none, got %d in test table '%s'",
				len(noPanic), test.name)
		}
		if len(noMatches) != len(test.wantNoMatches) {
			t.Errorf("PanicsRELoop(): Wrong number of panic-matches failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoMatches), len(noMatches), test.name, test.wantNoMatches, noMatches)
		} else {
			for i := 0; i < len(noMatches); i++ {
				if noMatches[i] != test.wantNoMatches[i] {
					t.Errorf("PanicsRELoop(): Wrong panic-matches failure: expected '%s', got '%s'",
						test.wantNoMatches[i], noMatches[i])
				}
			}
		}
	}
}

func TestPanicsValLoop(t *testing.T) {
	var noPanic []string
	var noEquals []string

	notPanicFunc := func(testName string) { noPanic = append(noPanic, testName) }
	notEqualsFunc := func(testName string) { noEquals = append(noEquals, testName) }

	tests := []struct {
		name         string
		pTable       []PanicValTest
		wantNoPanic  []string
		wantNoEquals []string
	}{
		{
			"p, eq; p, eq",
			[]PanicValTest{
				{"p, eq; p, eq: 1", func() { panic("ppp110") }, "ppp110"},
				{"p, eq; p, eq: 2", func() { panic("ppp111") }, "ppp111"},
			},
			[]string{},
			[]string{},
		},
		{
			"p, eq; p, neq",
			[]PanicValTest{
				// Non-strings (ints), equal and not equal
				{"p, eq; p, neq: 1", func() { panic(120) }, 120},
				{"p, eq; p, neq: 2", func() { panic(121) }, 129},
			},
			[]string{},
			[]string{"p, eq; p, neq: 2"},
		},
		{
			"p, eq; np",
			[]PanicValTest{
				{"p, eq; np: 1", func() { panic("ppp130") }, "ppp130"},
				{"p, eq; np: 2", func() {}, "ppp131"},
			},
			[]string{"p, eq; np: 2"},
			[]string{},
		},

		{
			"p, neq; p, eq",
			[]PanicValTest{
				{"p, neq; p, eq: 1", func() { panic("ppp210") }, "ccc210"},
				{"p, neq; p, eq: 2", func() { panic("ppp211") }, "ppp211"},
			},
			[]string{},
			[]string{"p, neq; p, eq: 1"},
		},
		{
			"p, neq; p, neq",
			[]PanicValTest{
				// String vs. int, float vs. int
				{"p, neq; p, neq: 1", func() { panic("220") }, 220},
				{"p, neq; p, neq: 2", func() { panic(221.0) }, 221},
			},
			[]string{},
			[]string{"p, neq; p, neq: 1", "p, neq; p, neq: 2"},
		},
		{
			"p, neq; np",
			[]PanicValTest{
				{"p, neq; np: 1", func() { panic("ppp230") }, "ccc230"},
				{"p, neq; np: 2", func() {}, "ppp231"},
			},
			[]string{"p, neq; np: 2"},
			[]string{"p, neq; np: 1"},
		},

		{
			"np; p, eq",
			[]PanicValTest{
				{"np; p, eq: 1", func() {}, "ppp310"},
				{"np; p, eq: 2", func() { panic("ppp311") }, "ppp311"},
			},
			[]string{"np; p, eq: 1"},
			[]string{},
		},
		{
			"np; p, neq",
			[]PanicValTest{
				{"np; p, neq: 1", func() {}, "ppp320"},
				{"np; p, neq: 2", func() { panic("ppp321") }, "zzz321"},
			},
			[]string{"np; p, neq: 1"},
			[]string{"np; p, neq: 2"},
		},
		{
			"np; np",
			[]PanicValTest{
				{"np; np: 1", func() {}, "ppp330"},
				{"np; np: 2", func() {}, "ppp331"},
			},
			[]string{"np; np: 1", "np; np: 2"},
			[]string{},
		},
	}
	for _, test := range tests {
		noPanic = nil
		noEquals = nil
		PanicsValLoop(test.pTable, notPanicFunc, notEqualsFunc)
		if len(noPanic) != len(test.wantNoPanic) {
			t.Errorf("PanicsValLoop(): Wrong number of panic-test failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoPanic), len(noPanic), test.name, test.wantNoPanic, noPanic)
		} else {
			for i := 0; i < len(noPanic); i++ {
				if noPanic[i] != test.wantNoPanic[i] {
					t.Errorf("PanicsValLoop(): Wrong panic-test failure: expected '%s', got '%s'",
						test.wantNoPanic[i], noPanic[i])
				}
			}
		}
		if len(noEquals) != len(test.wantNoEquals) {
			t.Errorf("PanicsValLoop(): Wrong number of panic-equals failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoEquals), len(noEquals), test.name, test.wantNoEquals, noEquals)
		} else {
			for i := 0; i < len(noEquals); i++ {
				if noEquals[i] != test.wantNoEquals[i] {
					t.Errorf("PanicsValLoop(): Wrong panic-equals failure: expected '%s', got '%s'",
						test.wantNoEquals[i], noEquals[i])
				}
			}
		}
	}
}

func TestPanicsValLoopPanicsWithUncomparableType(t *testing.T) {
	// for PanicsStr(), while testing a func containing PanicsValLoop()
	var didPanic bool
	var pContainsStr bool
	var pVal interface{}
	wantStr := "runtime error: comparing uncomparable type"

	// for the PanicsValLoop() being run by PanicsStr()
	var noPanic []string
	var noEquals []string
	notPanicFunc := func(testName string) { noPanic = append(noPanic, testName) }
	notEqualsFunc := func(testName string) { noEquals = append(noEquals, testName) }

	tests := []struct {
		name         string
		pTable       []PanicValTest
		wantNoEquals []string
	}{
		{
			"ok, not ok",
			[]PanicValTest{
				// ok but wrong
				{"ok, not ok: 1", func() { panic("ppp111") }, "zzz111"},
				{"ok, not ok: 2", func() { panic([]string{"a", "b"}) }, []string{"a", "b"}},
			},
			// first test within PanicsValLoop() proceeds normally, second one panics
			[]string{"ok, not ok: 1"},
		},
		{
			"not ok, ok",
			[]PanicValTest{
				{"not ok, ok: 1", func() { panic([]string{"a", "b"}) }, []string{"a", "b"}},
				// ok but wrong
				{"not ok, ok: 2", func() { panic("ppp222") }, "zzz222"},
			},
			[]string{},
		},
		{
			"not ok, not ok",
			[]PanicValTest{
				// one not ok but correct, one not ok and wrong
				{"not ok, not ok: 1", func() { panic([]string{"a", "b"}) }, []string{"a", "b"}},
				{"not ok, not ok: 2", func() { panic([]string{"a", "b"}) }, []string{"c", "d"}},
			},
			[]string{},
		},
	}
	for _, test := range tests {
		noPanic = nil
		noEquals = nil

		// It's a little suspect to use PanicsStr() here, but PanicsStr() and PanicsValLoop() don't reference each
		// other, and we have other tests for PanicsStr()
		didPanic, pContainsStr, pVal = PanicsStr(func() {
			// nolint: scopelint
			PanicsValLoop(test.pTable, notPanicFunc, notEqualsFunc)
		}, wantStr)
		if !didPanic {
			t.Errorf("Expected PanicsValLoop() itself to panic in test table '%s'", test.name)
		} else if !pContainsStr {
			t.Errorf("Incorrect panic value from PanicsValLoop() itself: expected string containing\n"+
				"%#+v\ngot\n%#+v\nin test table '%s'", wantStr, pVal, test.name)
		}

		// Now test the results of PanicsValLoop() itself
		if len(noPanic) != 0 {
			t.Errorf("PanicsValLoop(): Unexpected panic-test failure(s): expected none, got %d in test table '%s'",
				len(noPanic), test.name)
		}
		if len(noEquals) != len(test.wantNoEquals) {
			t.Errorf("PanicsValLoop(): Wrong number of panic-equals failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantNoEquals), len(noEquals), test.name, test.wantNoEquals, noEquals)
		} else {
			for i := 0; i < len(noEquals); i++ {
				if noEquals[i] != test.wantNoEquals[i] {
					t.Errorf("PanicsValLoop(): Wrong panic-equals failure: expected '%s', got '%s'",
						test.wantNoEquals[i], noEquals[i])
				}
			}
		}
	}
}

func TestNotPanicsLoop(t *testing.T) {
	var failed []string

	elseFunc := func(testName string) { failed = append(failed, testName) }

	tests := []struct {
		name       string
		pTable     []PanicTest
		wantFailed []string
	}{
		{
			"neither panics",
			[]PanicTest{
				{"neither panics: 1", func() {}},
				{"neither panics: 2", func() {}},
			},
			[]string{},
		},
		{
			"first panics",
			[]PanicTest{
				{"first panics: 1", func() { panic("fp1") }},
				{"first panics: 2", func() {}},
			},
			[]string{"first panics: 1"},
		},
		{
			"second panics",
			[]PanicTest{
				{"second panics: 1", func() {}},
				{"second panics: 2", func() { panic("sp2") }},
			},
			[]string{"second panics: 2"},
		},
		{
			"both panic",
			[]PanicTest{
				{"both panic: 1", func() { panic("bp1") }},
				{"both panic: 2", func() { panic("bp2") }},
			},
			[]string{"both panic: 1", "both panic: 2"},
		},
	}
	for _, test := range tests {
		failed = nil
		NotPanicsLoop(test.pTable, elseFunc)
		if len(failed) != len(test.wantFailed) {
			t.Errorf("Wrong number of not-panic-test failures: expected %d, got %d in test table '%s'\n"+
				"Expected failures:\n%#+v\nGot:\n%#+v",
				len(test.wantFailed), len(failed), test.name, test.wantFailed, failed)
		} else {
			for i := 0; i < len(failed); i++ {
				if failed[i] != test.wantFailed[i] {
					t.Errorf("Wrong not-panic-test failure: expected '%s', got '%s'", test.wantFailed[i], failed[i])
				}
			}
		}
	}
}
