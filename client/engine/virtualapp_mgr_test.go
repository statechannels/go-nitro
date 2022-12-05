package engine

import (
	"log"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/stretchr/testify/mock"
)

type MockedVirtualApp struct {
	mock.Mock
}

func (app *MockedVirtualApp) handleVirtualAppRequest(s *store.Store, request *VirtualAppRequest) error {
	args := app.Called(s, request)
	return args.Error(0)
}

func TestVirtualAppManager(t *testing.T) {
	setup := func() (*store.Store, *VirtualAppManager, *MockedVirtualApp) {
		logger := log.New(os.Stdout, "test-virtual-app", log.Flags())
		s := store.NewMemStore(testactors.Alice.PrivateKey)
		appMgr := NewVirtualAppManager(logger, &s)
		mApp := &MockedVirtualApp{}
		return &s, appMgr, mApp
	}

	t.Run("Dispatch request to registered virtual app", func(t *testing.T) {
		s, appMgr, mApp := setup()
		appMgr.RegisterVirtualApp("mock", mApp)
		mApp.On("handleVirtualAppRequest", mock.Anything, mock.Anything).Return(nil)
		req := &VirtualAppRequest{
			AppId: "mock",
		}
		appMgr.HandleRequest(req)
		mApp.AssertCalled(t, "handleVirtualAppRequest", s, req)
	})

	t.Run("Doesn't dispatch request to registered virtual app", func(t *testing.T) {
		_, appMgr, mApp := setup()
		appMgr.HandleRequest(&VirtualAppRequest{
			AppId: "mock",
		})
		mApp.AssertNotCalled(t, "handleVirtualAppRequest", mock.Anything, mock.Anything)
	})

}
