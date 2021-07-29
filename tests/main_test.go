package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

// DoSomething is a method on MyMockedObject that implements some interface
// and just records the activity, and returns what the Mock object tells it to.
//
// In the real object, this method would do something useful, but since this
// is a mocked object - we're just going to stub it out.
//
// NOTE: This method is not being tested here, code that uses this object is.
func (m *MyMockedObject) DoSomething(number int) (bool, error) {
	args := m.Called(number)
	fmt.Println(args)
	return true, nil
}

/*
	Actual test functions
*/

// TestSomething is an example of how to use our test object to
// make assertions about some target code we are testing.
func TestSomething(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("DoSomething", 123).Return(false, nil)
	testObj.DoSomething(124)
}
