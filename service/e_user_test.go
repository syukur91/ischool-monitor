package service

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/query"
)

func TestCreateUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		nama           string
		alamat         string
		password       string
		telepon        string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful add siswa",
			nama:         "Cendana",
			alamat:       "Jalan Cendana",
			password:     "adasdhakasdaj",
			telepon:      "02919192",
		},
		{
			scenarioName:   "Error add user alamat is not set",
			nama:           "Cendano",
			alamat:         "",
			password:       "adasdhakasdaj",
			telepon:        "02919192",
			expectedErrMsg: "User alamat is not set",
		},
		{
			scenarioName:   "Error add user password is not set",
			nama:           "Cendani",
			alamat:         "Jalan Cendani",
			password:       "",
			telepon:        "02919192",
			expectedErrMsg: "User password is not set",
		},
		{
			scenarioName:   "Error add user telepon is not set",
			nama:           "Cendani",
			alamat:         "Jalan Cendani",
			password:       "dasjdasdj",
			telepon:        "",
			expectedErrMsg: "User telepon is not set",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			_, err := userService.CreateUser(&schema.CreateUserRequest{
				Nama:     v.nama,
				Alamat:   v.alamat,
				Password: v.password,
				Telepon:  v.telepon,
			})
			//t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}
		})
	}
}

func TestListPaginationUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName      string
		query             query.GridParams
		expectedLength    int
		expectedTotal     int
		expectedErrMsg    string
		expectedFirstNama string
	}{
		{
			scenarioName:      "get first 2",
			query:             query.GridParams{Take: 2, Page: 1, Skip: 0, PageSize: 2},
			expectedLength:    2,
			expectedTotal:     5,
			expectedFirstNama: "Cendana",
		},
		{
			scenarioName:      "get second 2",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 2, PageSize: 2},
			expectedLength:    2,
			expectedTotal:     5,
			expectedFirstNama: "UserB",
		},
		{
			scenarioName:      "get third 1",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 4, PageSize: 2},
			expectedLength:    1,
			expectedTotal:     5,
			expectedFirstNama: "UserD",
		},
	}

	// insert 4 more data
	userService.CreateUser(&schema.CreateUserRequest{
		Nama:     "UserA",
		Alamat:   "AlamatA",
		Password: "weak",
		Telepon:  "089618465310",
	})

	userService.CreateUser(&schema.CreateUserRequest{
		Nama:     "UserB",
		Alamat:   "AlamatB",
		Password: "weak",
		Telepon:  "089618465310",
	})

	userService.CreateUser(&schema.CreateUserRequest{
		Nama:     "UserC",
		Alamat:   "AlamatC",
		Password: "weak",
		Telepon:  "089618465310",
	})

	userService.CreateUser(&schema.CreateUserRequest{
		Nama:     "UserD",
		Alamat:   "AlamatD",
		Password: "weak",
		Telepon:  "089618465310",
	})

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			users, total, err := userService.ListUsers(&v.query)
			// t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}

			if errMsg == "" {

				if len(users) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(users))
					return
				}

				if users[0].Nama != v.expectedFirstNama {
					t.Errorf("expect name %s, but got %s", v.expectedFirstNama, users[0].Nama)
					return
				}

				if total != v.expectedTotal {
					t.Errorf("expect len %d, but got %d", v.expectedTotal, total)
					return
				}
			}
		})
	}
}

func TestListFilterUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		query          query.GridParams
		expectedLength int
		expectedTotal  int
		expectedErrMsg string
	}{
		{
			scenarioName: "filter nama",
			query: query.GridParams{
				Take:      10,
				Page:      1,
				Skip:      0,
				PageSize:  10,
				HasFilter: true,
				Filter: query.GridFilterMain{
					Logic: "and",
					Filters: []query.GridFilter{
						query.GridFilter{
							Field:    "nama",
							Operator: "contains",
							Value:    "Cend",
						},
					},
				},
			},
			expectedLength: 1,
			expectedTotal:  1,
		},
		{
			scenarioName: "get na",
			query: query.GridParams{
				Take:      10,
				Page:      1,
				Skip:      0,
				PageSize:  10,
				HasFilter: true,
				Filter: query.GridFilterMain{
					Logic: "and",
					Filters: []query.GridFilter{
						query.GridFilter{
							Field:    "nama",
							Operator: "contains",
							Value:    "User",
						},
					},
				},
			},
			expectedLength: 4,
			expectedTotal:  4,
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			users, total, err := userService.ListUsers(&v.query)
			// t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}

			if errMsg == "" {

				if len(users) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(users))
					return
				}

				if total != v.expectedTotal {
					t.Errorf("expect len %d, but got %d", v.expectedTotal, total)
					return
				}
			}
		})
	}
}

func TestListSortUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName        string
		query               query.GridParams
		expectedLength      int
		expectedTotal       int
		expectedSortedNames []string
		expectedErrMsg      string
	}{
		{
			scenarioName: "sort Next",
			query: query.GridParams{
				Take:      10,
				Page:      1,
				Skip:      0,
				PageSize:  10,
				HasFilter: true,
				Filter: query.GridFilterMain{
					Logic: "and",
					Filters: []query.GridFilter{
						query.GridFilter{
							Field:    "nama",
							Operator: "contains",
							Value:    "User",
						},
					},
				},
				HasSort: true,
				Sort: []query.GridSort{
					query.GridSort{
						Field: "nama",
						Dir:   "asc",
					},
				},
			},
			expectedLength:      4,
			expectedTotal:       4,
			expectedSortedNames: []string{"UserA", "UserB", "UserC", "UserD"},
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			users, total, err := userService.ListUsers(&v.query)
			//t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}

			if errMsg == "" {

				if len(users) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(users))
					return
				}

				if users[0].Nama != v.expectedSortedNames[0] {
					t.Errorf("expect name %s and index 0, but got %s", v.expectedSortedNames[0], users[0].Nama)
					return
				}

				if users[1].Nama != v.expectedSortedNames[1] {
					t.Errorf("expect name %s and index 2, but got %s", v.expectedSortedNames[1], users[1].Nama)
					return
				}

				if total != v.expectedTotal {
					t.Errorf("expect len %d, but got %d", v.expectedTotal, total)
					return
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		nama           string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful get by id",
			id:           "1",
			nama:         "Cendana",
		},
		{
			scenarioName:   "Failure get: siswa with id not exists",
			id:             "10",
			nama:           "NextWhat",
			expectedErrMsg: "User with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			getUserResponse, err := userService.GetUser(v.id)
			//t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
				//ae, _ := err.(*apierror.APIError)
				//t.Logf("%+v", ae.Err)
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}

			// If error is empty, check for response
			if errMsg == "" {
				if getUserResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.nama != getUserResponse.Nama {
					t.Errorf("expect name %s, but got %s", v.nama, getUserResponse.Nama)
					return
				}
			}

		})
	}
}

func TestUpdateUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		nama           string
		alamat         string
		telpon         string
		password       string
		expectedErrMsg string
		expectedNama   string
		expectedAlamat string
	}{
		{
			scenarioName:   "Successful name update by id",
			id:             "1",
			nama:           "Cendanas",
			expectedNama:   "Cendanas",
			expectedAlamat: "Jalan Cendana",
		},
		{
			scenarioName:   "Successful alamat update by id",
			id:             "1",
			expectedNama:   "Cendanas",
			alamat:         "Jalan Nanos",
			expectedAlamat: "Jalan Nanos",
		},
		{
			scenarioName:   "Failure update: siswa with id not exists",
			id:             "10",
			expectedErrMsg: "User with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			updatedUserResponse, err := userService.UpdateUser(v.id, &schema.UpdateUserRequest{
				Nama:     v.nama,
				Alamat:   v.alamat,
				Password: v.password,
				Telepon:  v.alamat,
			})
			//t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}

			if errMsg == "" {
				if updatedUserResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.expectedNama != updatedUserResponse.Nama {
					t.Errorf("expect nama %s, but got %s", v.expectedNama, updatedUserResponse.Nama)
					return
				}

				if v.expectedAlamat != updatedUserResponse.Alamat {
					t.Errorf("expect tingkat %s, but got %s", v.expectedAlamat, updatedUserResponse.Alamat)
					return
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful delete by id",
			id:           "4",
		},
		{
			scenarioName:   "Failure delete: user with id not exists",
			id:             "10",
			expectedErrMsg: "User with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			err := userService.DeleteUser(v.id)
			//t.Logf("%+v, %+v", v, err)

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if v.expectedErrMsg != errMsg {
				t.Errorf("expect error %s, but got %s", v.expectedErrMsg, errMsg)
				return
			}
		})
	}
}
