// Code generated by mockery v2.32.4. DO NOT EDIT.

package search

import mock "github.com/stretchr/testify/mock"

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

type MockRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepository) EXPECT() *MockRepository_Expecter {
	return &MockRepository_Expecter{mock: &_m.Mock}
}

// AddDocuments provides a mock function with given fields: docs, index
func (_m *MockRepository) AddDocuments(docs interface{}, index string) error {
	ret := _m.Called(docs, index)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string) error); ok {
		r0 = rf(docs, index)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_AddDocuments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddDocuments'
type MockRepository_AddDocuments_Call struct {
	*mock.Call
}

// AddDocuments is a helper method to define mock.On call
//   - docs interface{}
//   - index string
func (_e *MockRepository_Expecter) AddDocuments(docs interface{}, index interface{}) *MockRepository_AddDocuments_Call {
	return &MockRepository_AddDocuments_Call{Call: _e.mock.On("AddDocuments", docs, index)}
}

func (_c *MockRepository_AddDocuments_Call) Run(run func(docs interface{}, index string)) *MockRepository_AddDocuments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(string))
	})
	return _c
}

func (_c *MockRepository_AddDocuments_Call) Return(_a0 error) *MockRepository_AddDocuments_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_AddDocuments_Call) RunAndReturn(run func(interface{}, string) error) *MockRepository_AddDocuments_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
