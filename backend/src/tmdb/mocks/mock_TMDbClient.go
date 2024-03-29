// Code generated by mockery v2.32.4. DO NOT EDIT.

package tmdb

import (
	mock "github.com/stretchr/testify/mock"

	tmdb "github.com/jj-style/go-tmdb"
)

// MockTMDbClient is an autogenerated mock type for the TMDbClient type
type MockTMDbClient struct {
	mock.Mock
}

type MockTMDbClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTMDbClient) EXPECT() *MockTMDbClient_Expecter {
	return &MockTMDbClient_Expecter{mock: &_m.Mock}
}

// GetPersonInfo provides a mock function with given fields: id, options
func (_m *MockTMDbClient) GetPersonInfo(id int, options map[string]string) (*tmdb.Person, error) {
	ret := _m.Called(id, options)

	var r0 *tmdb.Person
	var r1 error
	if rf, ok := ret.Get(0).(func(int, map[string]string) (*tmdb.Person, error)); ok {
		return rf(id, options)
	}
	if rf, ok := ret.Get(0).(func(int, map[string]string) *tmdb.Person); ok {
		r0 = rf(id, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tmdb.Person)
		}
	}

	if rf, ok := ret.Get(1).(func(int, map[string]string) error); ok {
		r1 = rf(id, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTMDbClient_GetPersonInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonInfo'
type MockTMDbClient_GetPersonInfo_Call struct {
	*mock.Call
}

// GetPersonInfo is a helper method to define mock.On call
//   - id int
//   - options map[string]string
func (_e *MockTMDbClient_Expecter) GetPersonInfo(id interface{}, options interface{}) *MockTMDbClient_GetPersonInfo_Call {
	return &MockTMDbClient_GetPersonInfo_Call{Call: _e.mock.On("GetPersonInfo", id, options)}
}

func (_c *MockTMDbClient_GetPersonInfo_Call) Run(run func(id int, options map[string]string)) *MockTMDbClient_GetPersonInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(map[string]string))
	})
	return _c
}

func (_c *MockTMDbClient_GetPersonInfo_Call) Return(_a0 *tmdb.Person, _a1 error) *MockTMDbClient_GetPersonInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTMDbClient_GetPersonInfo_Call) RunAndReturn(run func(int, map[string]string) (*tmdb.Person, error)) *MockTMDbClient_GetPersonInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonLatest provides a mock function with given fields:
func (_m *MockTMDbClient) GetPersonLatest() (*tmdb.PersonLatest, error) {
	ret := _m.Called()

	var r0 *tmdb.PersonLatest
	var r1 error
	if rf, ok := ret.Get(0).(func() (*tmdb.PersonLatest, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *tmdb.PersonLatest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tmdb.PersonLatest)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTMDbClient_GetPersonLatest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonLatest'
type MockTMDbClient_GetPersonLatest_Call struct {
	*mock.Call
}

// GetPersonLatest is a helper method to define mock.On call
func (_e *MockTMDbClient_Expecter) GetPersonLatest() *MockTMDbClient_GetPersonLatest_Call {
	return &MockTMDbClient_GetPersonLatest_Call{Call: _e.mock.On("GetPersonLatest")}
}

func (_c *MockTMDbClient_GetPersonLatest_Call) Run(run func()) *MockTMDbClient_GetPersonLatest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTMDbClient_GetPersonLatest_Call) Return(_a0 *tmdb.PersonLatest, _a1 error) *MockTMDbClient_GetPersonLatest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTMDbClient_GetPersonLatest_Call) RunAndReturn(run func() (*tmdb.PersonLatest, error)) *MockTMDbClient_GetPersonLatest_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonMovieCredits provides a mock function with given fields: id, options
func (_m *MockTMDbClient) GetPersonMovieCredits(id int, options map[string]string) (*tmdb.PersonMovieCredits, error) {
	ret := _m.Called(id, options)

	var r0 *tmdb.PersonMovieCredits
	var r1 error
	if rf, ok := ret.Get(0).(func(int, map[string]string) (*tmdb.PersonMovieCredits, error)); ok {
		return rf(id, options)
	}
	if rf, ok := ret.Get(0).(func(int, map[string]string) *tmdb.PersonMovieCredits); ok {
		r0 = rf(id, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tmdb.PersonMovieCredits)
		}
	}

	if rf, ok := ret.Get(1).(func(int, map[string]string) error); ok {
		r1 = rf(id, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTMDbClient_GetPersonMovieCredits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonMovieCredits'
type MockTMDbClient_GetPersonMovieCredits_Call struct {
	*mock.Call
}

// GetPersonMovieCredits is a helper method to define mock.On call
//   - id int
//   - options map[string]string
func (_e *MockTMDbClient_Expecter) GetPersonMovieCredits(id interface{}, options interface{}) *MockTMDbClient_GetPersonMovieCredits_Call {
	return &MockTMDbClient_GetPersonMovieCredits_Call{Call: _e.mock.On("GetPersonMovieCredits", id, options)}
}

func (_c *MockTMDbClient_GetPersonMovieCredits_Call) Run(run func(id int, options map[string]string)) *MockTMDbClient_GetPersonMovieCredits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(map[string]string))
	})
	return _c
}

func (_c *MockTMDbClient_GetPersonMovieCredits_Call) Return(_a0 *tmdb.PersonMovieCredits, _a1 error) *MockTMDbClient_GetPersonMovieCredits_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTMDbClient_GetPersonMovieCredits_Call) RunAndReturn(run func(int, map[string]string) (*tmdb.PersonMovieCredits, error)) *MockTMDbClient_GetPersonMovieCredits_Call {
	_c.Call.Return(run)
	return _c
}

// SearchPerson provides a mock function with given fields: name, options
func (_m *MockTMDbClient) SearchPerson(name string, options map[string]string) (*tmdb.PersonSearchResults, error) {
	ret := _m.Called(name, options)

	var r0 *tmdb.PersonSearchResults
	var r1 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) (*tmdb.PersonSearchResults, error)); ok {
		return rf(name, options)
	}
	if rf, ok := ret.Get(0).(func(string, map[string]string) *tmdb.PersonSearchResults); ok {
		r0 = rf(name, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tmdb.PersonSearchResults)
		}
	}

	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(name, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTMDbClient_SearchPerson_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchPerson'
type MockTMDbClient_SearchPerson_Call struct {
	*mock.Call
}

// SearchPerson is a helper method to define mock.On call
//   - name string
//   - options map[string]string
func (_e *MockTMDbClient_Expecter) SearchPerson(name interface{}, options interface{}) *MockTMDbClient_SearchPerson_Call {
	return &MockTMDbClient_SearchPerson_Call{Call: _e.mock.On("SearchPerson", name, options)}
}

func (_c *MockTMDbClient_SearchPerson_Call) Run(run func(name string, options map[string]string)) *MockTMDbClient_SearchPerson_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(map[string]string))
	})
	return _c
}

func (_c *MockTMDbClient_SearchPerson_Call) Return(_a0 *tmdb.PersonSearchResults, _a1 error) *MockTMDbClient_SearchPerson_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTMDbClient_SearchPerson_Call) RunAndReturn(run func(string, map[string]string) (*tmdb.PersonSearchResults, error)) *MockTMDbClient_SearchPerson_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTMDbClient creates a new instance of MockTMDbClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTMDbClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTMDbClient {
	mock := &MockTMDbClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
