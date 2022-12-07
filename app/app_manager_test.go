package app

import (
	"log"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedVirtualApp struct {
	mock.Mock
}

func (app *MockedVirtualApp) HandleRequest(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	ty string,
	data interface{},
) error {
	args := app.Called(ch, ty, data)
	return args.Error(0)
}

func (app *MockedVirtualApp) Id() string {
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
		err := s.SetChannel(c)
		require.NoError(t, err)
		mApp.On("HandleRequest", c, "ping", nil).Return(nil)
		req := &types.AppRequest{
			AppId:       "mock",
			RequestType: "ping",
			ChannelId:   c.ChannelId(),
		}
		err = appMgr.HandleRequest(req)
		require.NoError(t, err)
		mApp.AssertCalled(t, "HandleRequest", c, "ping", nil)
	})

	t.Run("Doesn't dispatch request to registered virtual app", func(t *testing.T) {
		_, appMgr, mApp := setup()
		err := appMgr.HandleRequest(&types.AppRequest{
			AppId: "mock",
		})
		require.Error(t, ErrAppNotRegistered, err)
		mApp.AssertNotCalled(t, "HandleRequest", mock.Anything, mock.Anything)
	})

}
