package service

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/query"
)

func TestCreateKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		nama           string
		tingkat        int
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful add",
			nama:         "XII IPA 1",
			tingkat:      3,
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			_, err := kelasService.CreateKelas(&schema.CreateKelasRequest{
				Nama:    v.nama,
				Tingkat: v.tingkat,
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

func TestListPaginationKelas(t *testing.T) {
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
			expectedFirstKode: "XII IPA 1",
		},
		{
			scenarioName:      "get second 2",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 2, PageSize: 2},
			expectedLength:    2,
			expectedTotal:     5,
			expectedFirstKode: "XII IPA 3",
		},
		{
			scenarioName:      "get third 1",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 4, PageSize: 2},
			expectedLength:    1,
			expectedTotal:     5,
			expectedFirstKode: "XII IPA 5",
		},
	}

	// insert 4 more data
	kelasService.CreateKelas(&schema.CreateKelasRequest{
		Nama:    "XII IPA 2",
		Tingkat: 3,
	})

	kelasService.CreateKelas(&schema.CreateKelasRequest{
		Nama:    "XII IPA 3",
		Tingkat: 3,
	})

	kelasService.CreateKelas(&schema.CreateKelasRequest{
		Nama:    "XII IPA 4",
		Tingkat: 3,
	})

	kelasService.CreateKelas(&schema.CreateKelasRequest{
		Nama:    "XII IPA 5",
		Tingkat: 3,
	})

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			kelass, total, err := kelasService.ListKelass(&v.query)
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

				if len(kelass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(kelass))
					return
				}

				if kelass[0].Nama != v.expectedFirstKode {
					t.Errorf("expect name %s, but got %s", v.expectedFirstKode, kelass[0].Nama)
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

func TestListFilterKelas(t *testing.T) {
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
							Value:    "XII IPA",
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
							Value:    "XII IPS",
						},
					},
				},
			},
			expectedLength: 0,
			expectedTotal:  0,
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			kelass, total, err := kelasService.ListKelass(&v.query)
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

				if len(kelass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(kelass))
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

func TestListSortKelas(t *testing.T) {
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
							Value:    "XII IPA",
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
			expectedLength:      5,
			expectedTotal:       5,
			expectedSortedNames: []string{"XII IPA 1", "XII IPA 2", "XII IPA 3"},
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			kelass, total, err := kelasService.ListKelass(&v.query)
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

				if len(kelass) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(kelass))
					return
				}

				if kelass[0].Nama != v.expectedSortedNames[0] {
					t.Errorf("expect name %s and index 0, but got %s", v.expectedSortedNames[0], kelass[0].Nama)
					return
				}

				if kelass[2].Nama != v.expectedSortedNames[2] {
					t.Errorf("expect name %s and index 2, but got %s", v.expectedSortedNames[2], kelass[2].Nama)
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

func TestGetKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		nama           string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful get by id",
			id:           "1",
			nama:         "XII IPA 1",
		},
		{
			scenarioName:   "Failure get: kelas with id not exists",
			id:             "10",
			nama:           "NextWhat",
			expectedErrMsg: "Kelas with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			getKelasResponse, err := kelasService.GetKelas(v.id)
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
				if getKelasResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.nama != getKelasResponse.Nama {
					t.Errorf("expect name %s, but got %s", v.nama, getKelasResponse.Nama)
					return
				}
			}

		})
	}
}

func TestUpdateKelas(t *testing.T) {
	testScenarios := []struct {
		scenarioName    string
		id              string
		nama            string
		kode            string
		tingkat         int
		expectedErrMsg  string
		expectedNama    string
		expectedTingkat int
	}{
		{
			scenarioName:    "Successful name update by id",
			id:              "1",
			nama:            "XI IPA 1",
			expectedNama:    "XI IPA 1",
			expectedTingkat: 3,
		},
		{
			scenarioName:    "Successful tingkat update by id",
			id:              "1",
			expectedNama:    "XI IPA 1",
			tingkat:         2,
			expectedTingkat: 2,
		},
		{
			scenarioName:   "Failure update: kelas with id not exists",
			id:             "10",
			expectedErrMsg: "Kelas with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			updatedKelasResponse, err := kelasService.UpdateKelas(v.id, &schema.UpdateKelasRequest{
				Nama:    v.nama,
				Tingkat: v.tingkat,
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
				if updatedKelasResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.expectedNama != updatedKelasResponse.Nama {
					t.Errorf("expect nama %s, but got %s", v.expectedNama, updatedKelasResponse.Nama)
					return
				}

				if v.expectedTingkat != updatedKelasResponse.Tingkat {
					t.Errorf("expect tingkat %d, but got %d", v.expectedTingkat, updatedKelasResponse.Tingkat)
					return
				}
			}
		})
	}
}

func TestDeleteKelas(t *testing.T) {
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
			scenarioName:   "Failure delete: kelas with id not exists",
			id:             "10",
			expectedErrMsg: "Kelas with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			err := kelasService.DeleteKelas(v.id)
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
