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
	"fmt"
	"regexp"
	"strings"
)

// A PanicTest encapsulates a function that is intended to panic, along with a name for it in diagnostic messages.
type PanicTest struct {
	Name string
	F    func()
}

// A PanicStrTest encapsulates a function that is intended to panic, along with a name for it in diagnostic messages,
// plus a string that should be contained in the panic value.
type PanicStrTest struct {
	Name    string
	F       func()
	WantStr string
}

// A PanicRETest encapsulates a function that is intended to panic, along with a name for it in diagnostic messages,
// plus a string representing a regular expression that should match the panic value.
type PanicRETest struct {
	Name   string
	F      func()
	WantRE string
}

// A PanicValTest encapsulates a function that is intended to panic, along with a name for it in diagnostic messages,
// plus a value that should equal the panic value.
type PanicValTest struct {
	Name    string
	F       func()
	WantVal interface{}
}

// Panics tests if the given function panics, and returns a boolean that is true if it does.
//
// It is strongly suggested to test the actual panic value with PanicsGet, PanicsStr, PanicsRE, or PanicsVal
// instead of using this function.
func Panics(f func()) (didPanic bool) {
	defer func() {
		didPanic = recover() != nil
	}()
	f()
	return false // overridden by the deferred function; here for the compiler
}

// NotPanics tests if the given function does not panic, and returns a boolean that is true if it does not.
//
// (This function is the opposite of Panics, and is only included to help make the semantics of tests clearer.)
func NotPanics(f func()) (didNotPanic bool) {
	return !Panics(f)
}

// PanicsGet tests if the given function panics, and returns a boolean that is true if it does.  If the function does
// panic, the panic value itself is also returned.  (Specifically, this is the return value from recover, which is nil
// if the function did not panic.)
//
// Note that PanicsStr, PanicsRE, and PanicsVal provide ways to test the panic value that are generally more convenient
// than this function.
func PanicsGet(f func()) (didPanic bool, pVal interface{}) {
	defer func() {
		pVal = recover()
		didPanic = pVal != nil
	}()
	f()
	return false, nil // overridden by the deferred function; here for the compiler
}

// NotPanicsGet tests if the given function does not panic, and returns a boolean that is true if it does not.  If the
// function does panic, the panic value itself is also returned.  (Specifically, this is the return value from
// recover, which is nil if the function did not panic.)
//
// (This function is the opposite of PanicsGet, and is only included to help make the semantics of tests clearer.)
func NotPanicsGet(f func()) (didNotPanic bool, pVal interface{}) {
	didPanic, _, v := PanicsStr(f, "")
	return !didPanic, v
}

// PanicsStr tests if the given function panics, and returns a boolean that is true if it does.  It also takes a string,
// to allow checking the contents of the panic; if the function does panic, and the panic can be cast to a string
// containing wantStr, pContainsStr will be true.  If the panic can be cast to an error value, the error's Error string
// will be used for the check.  The panic value itself is also returned.  (Specifically, this is the return value from
// recover, which is nil if the function did not panic.)
//
// For example, supposing you have a function that should panic with a nil input, but might also panic for some other
// reason.  You want to test what it does with a nil input, but what if it turns out to be panicking for another reason
// and doesn't actually panic on a nil input?  Use PanicStr like so to find out:
//
//	func NotNilTask(strs []string) {
//		stuff, err := allocateABunchOfStuff()
//		if err != nil {
//			panic("Out of Stuff")
//		}
//		if strs == nil {
//			panic("Can't do anything - input was nil")
//		}
//		// Do things
//		// ...
//	}
//
//	func TestNotNilTaskPanicsWithNil(t *testing.T) {
//		wantStr := "input was nil"
//		didPanic, pContainsStr, pVal := testhelp.PanicsStr(func() {
//			NotNilTask(nil)
//		}, wantStr)
//		if !didPanic {
//			t.Fatalf("Expected doing a NotNilTask with a nil input to panic")
//		} else if !pContainsStr {
//			t.Fatalf("Incorrect panic value: expected a string containing\n\"%s\"\ngot\n%#+v", wantStr, pVal)
//		}
//	}
//
// The contents check can be bypassed by setting wantStr to "", which is contained by every string.  In this case,
// pContainsStr will always be true (assuming the panic can be cast to a string), and you will still get the panic
// value.
func PanicsStr(f func(), wantStr string) (didPanic bool, pContainsStr bool, pVal interface{}) {
	defer func() {
		pVal = recover()
		didPanic = pVal != nil
		pStr, ok := pVal.(string)
		if !ok {
			var pErr error // pre-allocated so we can reuse ok
			pErr, ok = pVal.(error)
			if !ok {
				pContainsStr = false
			} else {
				pStr = pErr.Error()
			}
		}
		if ok { // one of the type assertions succeeded
			pContainsStr = strings.Contains(pStr, wantStr)
		}
	}()
	f()
	return false, false, nil // overridden by the deferred function; here for the compiler
}

// PanicsRE tests if the given function panics, and returns a boolean that is true if it does.  It also takes a string,
// to allow checking the contents of the panic; if the function does panic, and the panic can be cast to a string
// matching the regular expression given by wantRE, pMatchesRE will be true.  If the panic can be cast to an error
// value, the error's Error string will be used for the check.  The panic value itself is also returned.
// (Specifically, this is the return value from recover, which is nil if the function did not panic.)
//
// See PanicsStr for a plain-string-flavored version of how to use this function.
//
// The contents check can be bypassed by setting wantRE to "", which matches any string.  In this case,
// pMatchesRE will always be true (assuming the panic can be cast to a string), and you will still get the panic value.
//
// PanicsRE itself panics if wantRE does not represent a valid regular expression.
func PanicsRE(f func(), wantRE string) (didPanic bool, pMatchesRE bool, pVal interface{}) {
	// Compile so that we can fail immediately if the RE is invalid
	re, err := regexp.Compile(wantRE)
	if err != nil {
		panic(fmt.Sprintf("Regexp could not be compiled: %s", err))
	}

	defer func() {
		pVal = recover()
		didPanic = pVal != nil
		pStr, ok := pVal.(string)
		if !ok {
			var pErr error // pre-allocated so we can reuse ok
			pErr, ok = pVal.(error)
			if !ok {
				pMatchesRE = false
			} else {
				pStr = pErr.Error()
			}
		}
		if ok { // one of the type assertions succeeded
			pMatchesRE = re.MatchString(pStr)
		}
	}()
	f()
	return false, false, nil // overridden by the deferred function; here for the compiler
}

// PanicsVal tests if the given function panics, and returns a boolean that is true if it does.  It also takes a value,
// to allow checking the contents of the panic; if the function does panic, and the panic value equals wantVal, pEquals
// will be true.  The panic value itself is also returned.  (Specifically, this is the return value from recover, which
// is nil if the function did not panic.)
//
// See PanicsStr for a string-flavored version of how to use this function.
//
// PanicsVal itself panics if pVal and wantVal are of the same type, but it's not a type that Go can compare with ==.
func PanicsVal(f func(), wantVal interface{}) (didPanic bool, pEquals bool, pVal interface{}) {
	defer func() {
		pVal = recover()
		didPanic = pVal != nil
		pEquals = pVal == wantVal
	}()
	f()
	return false, false, nil // overridden by the deferred function; here for the compiler
}

// PanicsLoop runs through a slice of panic tests.  For any test function that does not panic, elseFunc is called with
// the name from the test's struct.
//
// It is strongly suggested to test the actual panic values with PanicsGetLoop, PanicsStrLoop, PanicsRELoop, or
// PanicsValLoop instead of using this function.
func PanicsLoop(tests []PanicTest, elseFunc func(testName string)) {
	for _, test := range tests {
		if !Panics(test.F) {
			elseFunc(test.Name)
		}
	}
}

// PanicsGetLoop runs through a slice of panic tests.  For any test function that does not panic, elseFunc is called
// with the name from the test's struct.  For any test function that does panic, valFunc is called with the panic value.
//
// Note that PanicsStrLoop, PanicsRELoop, and PanicsValLoop provide ways to test the panic values that are generally
// more convenient than this function.
func PanicsGetLoop(tests []PanicTest, elseFunc func(testName string), valFunc func(pVal interface{})) {
	for _, test := range tests {
		didPanic, pVal := PanicsGet(test.F)
		if !didPanic {
			elseFunc(test.Name)
		} else {
			valFunc(pVal)
		}
	}
}

// NotPanicsLoop runs through a slice of panic tests.  For any test function that panics, elseFunc is called with the
// name from the test's struct.
//
// It is strongly suggested to test the actual panic values with NotPanicsGetLoop, PanicsStrLoop, PanicsRELoop, or
// PanicsValLoop instead of using this function.
func NotPanicsLoop(tests []PanicTest, elseFunc func(testName string)) {
	for _, test := range tests {
		if Panics(test.F) {
			elseFunc(test.Name)
		}
	}
}

// NotPanicsGetLoop runs through a slice of panic tests.  For any test function that panics, elseFunc is called with the
// name from the test's struct and the panic value.
//
// Note that PanicsStrLoop, PanicsRELoop, and PanicsValLoop provide ways to test the panic values that are generally
// more convenient than this function.
func NotPanicsGetLoop(tests []PanicTest, elseFunc func(testName string, pVal interface{})) {
	for _, test := range tests {
		didPanic, pVal := PanicsGet(test.F)
		if didPanic {
			elseFunc(test.Name, pVal)
		}
	}
}

// PanicsStrLoop runs through a slice of panic tests, including checking the panic values to make sure they contain
// specific strings.  For any test function that does not panic, notPanicFunc is called with the name from the test's
// struct.  For any test function that does panic, but for which the panic value cannot be cast to a string or error
// containing the test's WantStr, notContainsFunc is called with test information and the panic value.  If wantStrAll
// is not nil, it is used in place of the tests' wantStrs.  See also PanicsStr.
//
// See NotContainsFuncErrorFactory and NotContainsFuncFatalFactory for good starting points for notContainsFunc.
func PanicsStrLoop(tests []PanicStrTest, wantStrAll *string, notPanicFunc func(testName string),
	notContainsFunc func(testName string, wantStr string, pVal interface{}),
) {
	var realWantStr string
	var didPanic, pContainsStr bool
	var pVal interface{}

	for _, test := range tests {
		if wantStrAll != nil {
			realWantStr = *wantStrAll
		} else {
			realWantStr = test.WantStr
		}
		didPanic, pContainsStr, pVal = PanicsStr(test.F, realWantStr)
		if !didPanic {
			notPanicFunc(test.Name)
		} else if !pContainsStr {
			notContainsFunc(test.Name, realWantStr, pVal)
		}
	}
}

// PanicsRELoop runs through a slice of panic tests, including checking the panic values to make sure they match
// specific regular expressions.  For any test function that does not panic, notPanicFunc is called with the name from
// the test's struct.  For any test function that does panic, but for which the panic value cannot be cast to a string
// or error matching the test's WantRE, notMatchesFunc is called with test information and the panic value.  If
// wantREAll is not nil, it is used in place of the tests' wantREs.  See also PanicsRE.
//
// See NotMatchesFuncErrorFactory and NotMatchesFuncFatalFactory for good starting points for notMatchesFunc.
//
// PanicsRELoop itself panics when attempting to run any test for which WantRE does not represent a valid regular
// expression.
func PanicsRELoop(tests []PanicRETest, wantREAll *string, notPanicFunc func(testName string),
	notMatchesFunc func(testName string, wantRE string, pVal interface{}),
) {
	var realWantRE string
	var didPanic, pMatchesRE bool
	var pVal interface{}

	for _, test := range tests {
		if wantREAll != nil {
			realWantRE = *wantREAll
		} else {
			realWantRE = test.WantRE
		}
		didPanic, pMatchesRE, pVal = PanicsRE(test.F, realWantRE)
		if !didPanic {
			notPanicFunc(test.Name)
		} else if !pMatchesRE {
			notMatchesFunc(test.Name, realWantRE, pVal)
		}
	}
}

// PanicsValLoop runs through a slice of panic tests, including checking the panic values.  For any test function that
// does not panic, notPanicFunc is called with the name from the test's struct.  For any test function that does panic,
// but for which the panic value does not equal the test's WantVal, notEqualsFunc is called with test information and
// the panic value.  If wantValAll is not nil, it is used in place of the tests' wantVals.  See also PanicsVal.
//
// See NotEqualsFuncErrorFactory and NotEqualsFuncFatalFactory for good starting points for notEqualsFunc.
//
// PanicsValLoop itself panics when attempting to run any test for which the panic value and the test's WantVal are of
// the same type, but it's not a type that Go can compare with ==.
func PanicsValLoop(tests []PanicValTest, wantValAll *interface{}, notPanicFunc func(testName string),
	notEqualsFunc func(testName string, wantVal interface{}, pVal interface{}),
) {
	var realWantVal interface{}
	var didPanic, pEquals bool
	var pVal interface{}

	for _, test := range tests {
		if wantValAll != nil {
			realWantVal = *wantValAll
		} else {
			realWantVal = test.WantVal
		}
		didPanic, pEquals, pVal = PanicsVal(test.F, realWantVal)
		if !didPanic {
			notPanicFunc(test.Name)
		} else if !pEquals {
			notEqualsFunc(test.Name, realWantVal, pVal)
		}
	}
}

// TestingT is a stub interface intended to be satisfied by a *testing.T.  It is here to help test factory functions
// such as NotContainsFuncErrorFactory.
type TestingT interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// NotContainsFuncErrorFactory returns a function suitable for passing to PanicsStrLoop as a notContainsFunc.  The
// returned function is a closure over a *testing.T which uses it to call Errorf with a generic informative message.
func NotContainsFuncErrorFactory(t TestingT) func(testName string, wantStr string, pVal interface{}) {
	return func(testName string, wantStr string, pVal interface{}) {
		t.Errorf("Incorrect panic value: expected a string containing\n\"%s\"\ngot\n%#+v\nin test '%s'",
			wantStr, pVal, testName)
	}
}

// NotContainsFuncFatalFactory returns a function suitable for passing to PanicsStrLoop as a notContainsFunc.  The
// returned function is a closure over a *testing.T which uses it to call Fatalf with a generic informative message.
func NotContainsFuncFatalFactory(t TestingT) func(testName string, wantStr string, pVal interface{}) {
	return func(testName string, wantStr string, pVal interface{}) {
		t.Fatalf("Incorrect panic value: expected a string containing\n\"%s\"\ngot\n%#+v\nin test '%s'",
			wantStr, pVal, testName)
	}
}

// NotMatchesFuncErrorFactory returns a function suitable for passing to PanicsRELoop as a notMatchesFunc.  The
// returned function is a closure over a *testing.T which uses it to call Errorf with a generic informative message.
func NotMatchesFuncErrorFactory(t TestingT) func(testName string, wantRE string, pVal interface{}) {
	return func(testName string, wantRE string, pVal interface{}) {
		t.Errorf("Incorrect panic value: expected a string matching\n\"%s\"\ngot\n%#+v\nin test '%s'",
			wantRE, pVal, testName)
	}
}

// NotMatchesFuncFatalFactory returns a function suitable for passing to PanicsRELoop as a notMatchesFunc.  The
// returned function is a closure over a *testing.T which uses it to call Fatalf with a generic informative message.
func NotMatchesFuncFatalFactory(t TestingT) func(testName string, wantRE string, pVal interface{}) {
	return func(testName string, wantRE string, pVal interface{}) {
		t.Fatalf("Incorrect panic value: expected a string matching\n\"%s\"\ngot\n%#+v\nin test '%s'",
			wantRE, pVal, testName)
	}
}

// NotEqualsFuncErrorFactory returns a function suitable for passing to PanicsValLoop as a notEqualsFunc.  The
// returned function is a closure over a *testing.T which uses it to call Errorf with a generic informative message.
func NotEqualsFuncErrorFactory(t TestingT) func(testName string, wantVal interface{}, pVal interface{}) {
	return func(testName string, wantVal interface{}, pVal interface{}) {
		t.Errorf("Incorrect panic value: expected\n%#+v\ngot\n%#+v\nin test '%s'",
			wantVal, pVal, testName)
	}
}

// NotEqualsFuncFatalFactory returns a function suitable for passing to PanicsValLoop as a notEqualsFunc.  The
// returned function is a closure over a *testing.T which uses it to call Fatalf with a generic informative message.
func NotEqualsFuncFatalFactory(t TestingT) func(testName string, wantVal interface{}, pVal interface{}) {
	return func(testName string, wantVal interface{}, pVal interface{}) {
		t.Fatalf("Incorrect panic value: expected\n%#+v\ngot\n%#+v\nin test '%s'",
			wantVal, pVal, testName)
	}
}
