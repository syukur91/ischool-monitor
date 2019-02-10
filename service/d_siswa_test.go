package service

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/query"
)

func TestCreateSiswa(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		nama           string
		idKelas        int
		idWaliKelas    int
		tingkat        int
		alamat         string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful add siswa",
			nama:         "Cendana",
			idKelas:      1,
			idWaliKelas:  2,
			alamat:       "Jalan Cendana",
			tingkat:      3,
		},
		{
			scenarioName:   "Error add siswa id kelas is not set",
			nama:           "Cendano",
			idKelas:        0,
			idWaliKelas:    2,
			alamat:         "Jalan Cendana",
			tingkat:        3,
			expectedErrMsg: "Siswa id kelas is not set",
		},
		{
			scenarioName:   "Error add siswa id wali kelas is not set",
			nama:           "Cendani",
			idKelas:        1,
			idWaliKelas:    0,
			alamat:         "Jalan Cendana",
			tingkat:        3,
			expectedErrMsg: "Siswa id wali kelas is not set",
		},
		{
			scenarioName:   "Error add siswa alamat is not set",
			nama:           "Cendani",
			idKelas:        1,
			idWaliKelas:    2,
			alamat:         "",
			tingkat:        3,
			expectedErrMsg: "Siswa alamat is not set",
		},
		{
			scenarioName:   "Error add siswa tingkat is not set",
			nama:           "Cendani",
			idKelas:        1,
			idWaliKelas:    2,
			alamat:         "Jalan Cendana",
			tingkat:        0,
			expectedErrMsg: "Siswa tingkat is not set",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			_, err := siswaService.CreateSiswa(&schema.CreateSiswaRequest{
				Nama:        v.nama,
				IDKelas:     v.idKelas,
				IDWaliKelas: v.idWaliKelas,
				Alamat:      v.alamat,
				Tingkat:     v.tingkat,
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

func TestListPaginationSiswa(t *testing.T) {
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
			expectedFirstNama: "Sakura",
		},
		{
			scenarioName:      "get third 1",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 4, PageSize: 2},
			expectedLength:    1,
			expectedTotal:     5,
			expectedFirstNama: "Lee",
		},
	}

	// insert 4 more data
	siswaService.CreateSiswa(&schema.CreateSiswaRequest{
		Nama:        "Naruto",
		IDKelas:     1,
		IDWaliKelas: 2,
		Alamat:      "Jalan Naruto",
		Tingkat:     3,
	})

	siswaService.CreateSiswa(&schema.CreateSiswaRequest{
		Nama:        "Sakura",
		IDKelas:     1,
		IDWaliKelas: 2,
		Alamat:      "Jalan Naruto",
		Tingkat:     3,
	})

	siswaService.CreateSiswa(&schema.CreateSiswaRequest{
		Nama:        "Sasuke",
		IDKelas:     2,
		IDWaliKelas: 3,
		Alamat:      "Jalan Sasuke",
		Tingkat:     3,
	})

	siswaService.CreateSiswa(&schema.CreateSiswaRequest{
		Nama:        "Lee",
		IDKelas:     2,
		IDWaliKelas: 3,
		Alamat:      "Jalan Lee",
		Tingkat:     3,
	})

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			siswas, total, err := siswaService.ListSiswas(&v.query)
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

				if len(siswas) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(siswas))
					return
				}

				if siswas[0].Nama != v.expectedFirstNama {
					t.Errorf("expect name %s, but got %s", v.expectedFirstNama, siswas[0].Nama)
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

func TestListFilterSiswa(t *testing.T) {
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
							Value:    "Sa",
						},
					},
				},
			},
			expectedLength: 2,
			expectedTotal:  2,
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
							Value:    "Na",
						},
					},
				},
			},
			expectedLength: 1,
			expectedTotal:  1,
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			siswass, total, err := siswaService.ListSiswas(&v.query)
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

				if len(siswass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(siswass))
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

func TestListSortSiswa(t *testing.T) {
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
							Value:    "Sa",
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
			expectedLength:      2,
			expectedTotal:       2,
			expectedSortedNames: []string{"Sakura", "Sasuke"},
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			siswas, total, err := siswaService.ListSiswas(&v.query)
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

				if len(siswas) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(siswas))
					return
				}

				if siswas[0].Nama != v.expectedSortedNames[0] {
					t.Errorf("expect name %s and index 0, but got %s", v.expectedSortedNames[0], siswas[0].Nama)
					return
				}

				if siswas[1].Nama != v.expectedSortedNames[1] {
					t.Errorf("expect name %s and index 2, but got %s", v.expectedSortedNames[1], siswas[1].Nama)
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

func TestGetSiswa(t *testing.T) {
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
			expectedErrMsg: "Siswa with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			getSiswaResponse, err := siswaService.GetSiswa(v.id)
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
				if getSiswaResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.nama != getSiswaResponse.Nama {
					t.Errorf("expect name %s, but got %s", v.nama, getSiswaResponse.Nama)
					return
				}
			}

		})
	}
}

func TestUpdateSiswa(t *testing.T) {
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
			expectedErrMsg: "Siswa with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			updatedSiswaResponse, err := siswaService.UpdateSiswa(v.id, &schema.UpdateSiswaRequest{
				Nama:   v.nama,
				Alamat: v.alamat,
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
				if updatedSiswaResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.expectedNama != updatedSiswaResponse.Nama {
					t.Errorf("expect nama %s, but got %s", v.expectedNama, updatedSiswaResponse.Nama)
					return
				}

				if v.expectedAlamat != updatedSiswaResponse.Alamat {
					t.Errorf("expect tingkat %s, but got %s", v.expectedAlamat, updatedSiswaResponse.Alamat)
					return
				}
			}
		})
	}
}

func TestDeleteSiswa(t *testing.T) {
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
			scenarioName:   "Failure delete: siswa with id not exists",
			id:             "10",
			expectedErrMsg: "Siswa with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			err := siswaService.DeleteSiswa(v.id)
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
