// Code generated by mockery v2.32.4. DO NOT EDIT.

package gamemanager

import (
	context "context"

	gamemanager "github.com/jj-style/chain-react/src/server/game_manager"
	mock "github.com/stretchr/testify/mock"
)

// MockGameManager is an autogenerated mock type for the GameManager type
type MockGameManager struct {
	mock.Mock
}

type MockGameManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGameManager) EXPECT() *MockGameManager_Expecter {
	return &MockGameManager_Expecter{mock: &_m.Mock}
}

// CreateGame provides a mock function with given fields: _a0
func (_m *MockGameManager) CreateGame(_a0 context.Context) (*gamemanager.Game, error) {
	ret := _m.Called(_a0)

	var r0 *gamemanager.Game
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*gamemanager.Game, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *gamemanager.Game); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gamemanager.Game)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGameManager_CreateGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateGame'
type MockGameManager_CreateGame_Call struct {
	*mock.Call
}

// CreateGame is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *MockGameManager_Expecter) CreateGame(_a0 interface{}) *MockGameManager_CreateGame_Call {
	return &MockGameManager_CreateGame_Call{Call: _e.mock.On("CreateGame", _a0)}
}

func (_c *MockGameManager_CreateGame_Call) Run(run func(_a0 context.Context)) *MockGameManager_CreateGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGameManager_CreateGame_Call) Return(_a0 *gamemanager.Game, _a1 error) *MockGameManager_CreateGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGameManager_CreateGame_Call) RunAndReturn(run func(context.Context) (*gamemanager.Game, error)) *MockGameManager_CreateGame_Call {
	_c.Call.Return(run)
	return _c
}

// GetGame provides a mock function with given fields: _a0
func (_m *MockGameManager) GetGame(_a0 context.Context) (*gamemanager.Game, error) {
	ret := _m.Called(_a0)

	var r0 *gamemanager.Game
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*gamemanager.Game, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *gamemanager.Game); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gamemanager.Game)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGameManager_GetGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGame'
type MockGameManager_GetGame_Call struct {
	*mock.Call
}

// GetGame is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *MockGameManager_Expecter) GetGame(_a0 interface{}) *MockGameManager_GetGame_Call {
	return &MockGameManager_GetGame_Call{Call: _e.mock.On("GetGame", _a0)}
}

func (_c *MockGameManager_GetGame_Call) Run(run func(_a0 context.Context)) *MockGameManager_GetGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGameManager_GetGame_Call) Return(_a0 *gamemanager.Game, _a1 error) *MockGameManager_GetGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGameManager_GetGame_Call) RunAndReturn(run func(context.Context) (*gamemanager.Game, error)) *MockGameManager_GetGame_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGameManager creates a new instance of MockGameManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGameManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGameManager {
	mock := &MockGameManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}