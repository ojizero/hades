// Package chainables provides some experiments and thoughts I've had
// around working with Go codebases. The main objective is to
// attempt and reduce redundancy in code by incorporating
// and hiding away repetitive code blocks.
//
package chainables

import (
	"reflect"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

// With attempts to replicate Elixir's `with` clause in Go
// as means to reduce the boilerplate of error handling.
//
// However it is much more simplistic/dumbified version of
// that found in Elixir!
//
// Here be dragons as this is just an thought experiment
// being done!
//
func With(fns ...interface{}) error {
	fvs := mustBeValidFuncs(fns)
	prevOut := []reflect.Value{}
	for _, fv := range fvs {
		prevOut = fv.Call(prevOut)
		if len(prevOut) == 0 {
			continue
		}
		hasErr := prevOut[len(prevOut)-1].Type().Implements(errType)
		if hasErr {
			if err, ok := prevOut[len(prevOut)-1].Interface().(error); ok && err != nil {
				return err
			}
			prevOut = prevOut[0 : len(prevOut)-1]
		}
	}
	return nil
}

func mustBeValidFuncs(fns []interface{}) []reflect.Value {
	fvs := []reflect.Value{}
	prevOuts := []reflect.Type{}
	for _, fn := range fns {
		v := reflect.ValueOf(fn)
		fvs = append(fvs, v)
		t := v.Type()
		if t.Kind() != reflect.Func {
			panic("non function passed")
		}
		nIn := t.NumIn()
		nOut := t.NumOut()
		if nIn != len(prevOuts) && nIn != len(prevOuts)-1 {
			panic("arguments of functions don't align with outputs of previous ones")
		}
		for i := 0; i < nIn; i += 1 {
			if t.In(i).Kind() != prevOuts[i].Kind() {
				panic("arguments of functions don't align with types of previous ones")
			}
		}
		prevOuts = []reflect.Type{}
		for i := 0; i < nOut; i += 1 {
			prevOuts = append(prevOuts, t.Out(i))
		}
	}
	return fvs
}
