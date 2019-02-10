package service

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/query"
)

func TestCreateWaliKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		nama           string
		alamat         string
		telpon         string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful add wali kelas",
			nama:         "Nano",
			alamat:       "Jalan Cendana",
			telpon:       "0228973218",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			_, err := waliKelasService.CreateWali_Kelas(&schema.CreateWali_KelasRequest{
				Nama:   v.nama,
				Alamat: v.alamat,
				Telpon: v.telpon,
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

func TestListPaginationWaliKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName      string
		query             query.GridParams
		expectedLength    int
		expectedTotal     int
		expectedErrMsg    string
		expectedFirstKode string
	}{
		{
			scenarioName:      "get first 2",
			query:             query.GridParams{Take: 2, Page: 1, Skip: 0, PageSize: 2},
			expectedLength:    2,
			expectedTotal:     5,
			expectedFirstKode: "Nano",
		},
		{
			scenarioName:      "get second 2",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 2, PageSize: 2},
			expectedLength:    2,
			expectedTotal:     5,
			expectedFirstKode: "Nani",
		},
		{
			scenarioName:      "get third 1",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 4, PageSize: 2},
			expectedLength:    1,
			expectedTotal:     5,
			expectedFirstKode: "Noni",
		},
	}

	// insert 4 more data
	waliKelasService.CreateWali_Kelas(&schema.CreateWali_KelasRequest{
		Nama:   "Nana",
		Alamat: "Jalan Nana",
		Telpon: "08916162625",
	})

	waliKelasService.CreateWali_Kelas(&schema.CreateWali_KelasRequest{
		Nama:   "Nani",
		Alamat: "Jalan Nani",
		Telpon: "08972552516",
	})

	waliKelasService.CreateWali_Kelas(&schema.CreateWali_KelasRequest{
		Nama:   "Nini",
		Alamat: "Jalan Nini",
		Telpon: "089725525564",
	})

	waliKelasService.CreateWali_Kelas(&schema.CreateWali_KelasRequest{
		Nama:   "Noni",
		Alamat: "Jalan Noni",
		Telpon: "089725525234",
	})

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			waliKelass, total, err := waliKelasService.ListWali_Kelass(&v.query)
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

				if len(waliKelass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(waliKelass))
					return
				}

				if waliKelass[0].Nama != v.expectedFirstKode {
					t.Errorf("expect name %s, but got %s", v.expectedFirstKode, waliKelass[0].Nama)
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

func TestListFilterWaliKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		query          query.GridParams
		expectedLength int
		expectedTotal  int
		expectedErrMsg string
	}{
		{
			scenarioName: "filter kode",
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
							Value:    "N",
						},
					},
				},
			},
			expectedLength: 5,
			expectedTotal:  5,
		},
		{
			scenarioName: "get nextwhat",
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
							Value:    "Na",
						},
					},
				},
			},
			expectedLength: 3,
			expectedTotal:  3,
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			waliKelass, total, err := waliKelasService.ListWali_Kelass(&v.query)
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

				if len(waliKelass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(waliKelass))
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

func TestListSortWaliKelas(t *testing.T) {
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
							Value:    "Na",
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
			expectedLength:      3,
			expectedTotal:       3,
			expectedSortedNames: []string{"Nana", "Nani", "Nano"},
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			waliKelass, total, err := waliKelasService.ListWali_Kelass(&v.query)
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

				if len(waliKelass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(waliKelass))
					return
				}

				if waliKelass[0].Nama != v.expectedSortedNames[0] {
					t.Errorf("expect name %s and index 0, but got %s", v.expectedSortedNames[0], waliKelass[0].Nama)
					return
				}

				if waliKelass[2].Nama != v.expectedSortedNames[2] {
					t.Errorf("expect name %s and index 2, but got %s", v.expectedSortedNames[2], waliKelass[2].Nama)
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

func TestGetWaliKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		nama           string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful get by id",
			id:           "1",
			nama:         "Nano",
		},
		{
			scenarioName:   "Failure get: wali kelas with id not exists",
			id:             "10",
			nama:           "NextWhat",
			expectedErrMsg: "Wali_Kelas with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			getWaliKelasResponse, err := waliKelasService.GetWali_Kelas(v.id)
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
				if getWaliKelasResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.nama != getWaliKelasResponse.Nama {
					t.Errorf("expect name %s, but got %s", v.nama, getWaliKelasResponse.Nama)
					return
				}
			}

		})
	}
}

func TestUpdateWaliKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		nama           string
		alamat         string
		telpon         string
		expectedErrMsg string
		expectedNama   string
		expectedAlamat string
	}{
		{
			scenarioName:   "Successful name update by id",
			id:             "1",
			nama:           "Nanos",
			expectedNama:   "Nanos",
			expectedAlamat: "Jalan Cendana",
		},
		{
			scenarioName:   "Successful alamat update by id",
			id:             "1",
			expectedNama:   "Nanos",
			alamat:         "Jalan Nanos",
			expectedAlamat: "Jalan Nanos",
		},
		{
			scenarioName:   "Failure update: wali kelas with id not exists",
			id:             "10",
			expectedErrMsg: "Wali_Kelas with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			updatedWaliKelasResponse, err := waliKelasService.UpdateWali_Kelas(v.id, &schema.UpdateWali_KelasRequest{
				Nama:   v.nama,
				Alamat: v.alamat,
				Telpon: v.telpon,
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
				if updatedWaliKelasResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.expectedNama != updatedWaliKelasResponse.Nama {
					t.Errorf("expect nama %s, but got %s", v.expectedNama, updatedWaliKelasResponse.Nama)
					return
				}

				if v.expectedAlamat != updatedWaliKelasResponse.Alamat {
					t.Errorf("expect tingkat %s, but got %s", v.expectedAlamat, updatedWaliKelasResponse.Alamat)
					return
				}
			}
		})
	}
}

func TestDeleteWaliKelas(t *testing.T) {
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
			scenarioName:   "Failure delete: wali kelas with id not exists",
			id:             "10",
			expectedErrMsg: "Wali_Kelas with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			err := waliKelasService.DeleteWali_Kelas(v.id)
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
