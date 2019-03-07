package user_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/josejuanmontiel/golang/testing/mocks"
	"github.com/josejuanmontiel/golang/testing/user"
)

func TestUse(t *testing.T) {
	fmt.Println("TestUse")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)
	testUser := &user.User{Doer: mockDoer}

	// Expect Do to be called once with 123 and "Hello GoMock" as parameters, and return nil from the mocked call.
	mockDoer.EXPECT().DoSomething(gomock.Any(), "Hello GoMock").Return(nil).Times(1)

	err := testUser.Use()
	if err != nil {
		t.Fail()
	}
}

func TestUseReturnsErrorFromDo(t *testing.T) {
	fmt.Println("TestUseReturnsErrorFromDo")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dummyError := errors.New("dummy error")
	mockDoer := mocks.NewMockDoer(mockCtrl)
	testUser := &user.User{Doer: mockDoer}

	// Expect Do to be called once with 123 and "Hello GoMock" as parameters, and return dummyError from the mocked call.
	mockDoer.EXPECT().DoSomething(123, "Hello GoMock").Return(dummyError).Times(1)

	err := testUser.Use()

	if err != dummyError {
		t.Fail()
	}
}
