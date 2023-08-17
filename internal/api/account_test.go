package api

import (
	dbsql "database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/sql"
	mockdb "github.com/zvash/accounting-system/internal/sql/mock"
	"github.com/zvash/accounting-system/internal/token"
	"github.com/zvash/accounting-system/internal/util"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccountById(gomock.Any(), gomock.Eq(sql.GetUserAccountByIdParams{
						ID:    account.ID,
						Owner: user.Username,
					})).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				requireBodyMatchAccount(t, resp, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccountById(gomock.Any(), gomock.Eq(sql.GetUserAccountByIdParams{
						ID:    account.ID,
						Owner: user.Username,
					})).
					Times(1).
					Return(sql.Account{}, sql.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, resp.StatusCode, http.StatusNotFound)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccountById(gomock.Any(), gomock.Eq(sql.GetUserAccountByIdParams{
						ID:    account.ID,
						Owner: user.Username,
					})).
					Times(1).
					Return(sql.Account{}, dbsql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccountById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccountById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccountById(gomock.Any(), gomock.Eq(sql.GetUserAccountByIdParams{
						ID:    account.ID,
						Owner: "unauthorized_user",
					})).
					Times(1).
					Return(sql.Account{}, sql.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
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
			tc.setupAuth(t, request, server.tokenMaker)
			resp, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}

}

func randomAccount(owner string) sql.Account {
	return sql.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
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
