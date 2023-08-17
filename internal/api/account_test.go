package api

import (
	dbsql "database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/sql"
	mockdb "github.com/zvash/accounting-system/internal/sql/mock"
	"github.com/zvash/accounting-system/internal/util"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, resp.StatusCode, http.StatusOK)
				requireBodyMatchAccount(t, resp, account)
			},
		}, {
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sql.Account{}, sql.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, resp.StatusCode, http.StatusNotFound)
			},
		}, {
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sql.Account{}, dbsql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		}, {
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server := newTestServer(t, store)
			resp, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}

}

func randomAccount() sql.Account {
	return sql.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, resp *http.Response, account sql.Account) {
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var gotAccount sql.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)
}
