package app

import (
	"log"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedVirtualApp struct {
	mock.Mock
}

func (app *MockedVirtualApp) HandleRequest(ch *channel.Channel, ty string, data interface{}) error {
	args := app.Called(ch, ty, data)
	return args.Error(0)
}

func (app *MockedVirtualApp) Type() string {
	return "mock"
}

func TestAppManager(t *testing.T) {
	setup := func() (store.Store, *AppManager, *MockedVirtualApp) {
		logger := log.New(os.Stdout, "test-virtual-app", log.Flags())
		s := store.NewMemStore(testactors.Alice.PrivateKey)
		appMgr := NewAppManager(logger, s)
		mApp := &MockedVirtualApp{}
		return s, appMgr, mApp
	}

	t.Run("Dispatch request to registered virtual app", func(t *testing.T) {
		s, appMgr, mApp := setup()
		appMgr.RegisterApp(mApp)
		c := td.Objectives.Directfund.GenericDFO().C
		s.SetChannel(c)
		mApp.On("HandleRequest", c, "ping", nil).Return(nil)
		req := &AppRequest{
			AppType:     "mock",
			RequestType: "ping",
			ChannelId:   c.ChannelId(),
		}
		err := appMgr.HandleRequest(req)
		require.NoError(t, err)
		mApp.AssertCalled(t, "HandleRequest", c, "ping", nil)
	})

	t.Run("Doesn't dispatch request to registered virtual app", func(t *testing.T) {
		_, appMgr, mApp := setup()
		err := appMgr.HandleRequest(&AppRequest{
			AppType: "mock",
		})
		require.Error(t, ErrAppNotRegistered, err)
		mApp.AssertNotCalled(t, "HandleRequest", mock.Anything, mock.Anything)
	})

}
