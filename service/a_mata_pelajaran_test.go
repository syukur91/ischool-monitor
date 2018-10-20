package service

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/query"

	"github.com/arifsetiawan/go-common/env"
)

var mataPelajaranService *Mata_PelajaranService

func TestMain(m *testing.M) {
	if len(os.Getenv("TEST_DB_CONNECTION_STR")) == 0 {
		log.Fatalln("Database connection string is not set. Set TEST_DB_CONNECTION_STR in environment")
	}

	db, err := sqlx.Connect(env.Getenv("TEST_DB_DRIVER", "postgres"), os.Getenv("TEST_DB_CONNECTION_STR"))
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	mataPelajaranService = NewMata_PelajaranService(db)

	code := m.Run()
	os.Exit(code)
}

func TestCreateMataPelajaran(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		nama           string
		kode           string
		tingkat        int
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful add",
			nama:         "PPKN",
			kode:         "KD_PPKN_1",
			tingkat:      1,
		},
		{
			scenarioName:   "Failure add: mata pelajaran kode is not set",
			nama:           "PENJASKES",
			kode:           "",
			tingkat:        1,
			expectedErrMsg: "Mata_Pelajaran kode is not set",
		},
		{
			scenarioName:   "Failure add: mata pelajaran with same kode exist",
			nama:           "PPKN",
			kode:           "KD_PPKN_1",
			tingkat:        1,
			expectedErrMsg: "Mata_Pelajaran with same kode already exists. Use different kode",
		},
		{
			scenarioName:   "Failure add: mata pelajaran kode is not set",
			nama:           "PPKN",
			kode:           "",
			tingkat:        1,
			expectedErrMsg: "Mata_Pelajaran kode is not set",
		},
		{
			scenarioName: "Successful add: another mata pelajaran",
			nama:         "PENJASKES",
			kode:         "KD_PENJASKES_1",
			tingkat:      1,
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			_, err := mataPelajaranService.CreateMata_Pelajaran(&schema.CreateMata_PelajaranRequest{
				Nama:    v.nama,
				Kode:    v.kode,
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

func TestListPaginationMataPelajaran(t *testing.T) {
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
			expectedFirstKode: "KD_PPKN_1",
		},
		{
			scenarioName:      "get second 2",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 2, PageSize: 2},
			expectedLength:    2,
			expectedTotal:     5,
			expectedFirstKode: "KD_PPKN_2",
		},
		{
			scenarioName:      "get third 1",
			query:             query.GridParams{Take: 2, Page: 2, Skip: 4, PageSize: 2},
			expectedLength:    1,
			expectedTotal:     5,
			expectedFirstKode: "KD_PENJASKES_2",
		},
	}

	// insert 4 more data
	mataPelajaranService.CreateMata_Pelajaran(&schema.CreateMata_PelajaranRequest{
		Nama:    "PPKN",
		Kode:    "KD_PPKN_2",
		Tingkat: 2,
	})

	mataPelajaranService.CreateMata_Pelajaran(&schema.CreateMata_PelajaranRequest{
		Nama:    "PPKN",
		Kode:    "KD_PPKN_3",
		Tingkat: 3,
	})

	mataPelajaranService.CreateMata_Pelajaran(&schema.CreateMata_PelajaranRequest{
		Nama:    "PENJASKES",
		Kode:    "KD_PENJASKES_1",
		Tingkat: 1,
	})

	mataPelajaranService.CreateMata_Pelajaran(&schema.CreateMata_PelajaranRequest{
		Nama:    "PENJASKES",
		Kode:    "KD_PENJASKES_2",
		Tingkat: 2,
	})

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			mataPelajarans, total, err := mataPelajaranService.ListMata_Pelajarans(&v.query)
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

				if len(mataPelajarans) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(mataPelajarans))
					return
				}

				if mataPelajarans[0].Kode != v.expectedFirstKode {
					t.Errorf("expect name %s, but got %s", v.expectedFirstKode, mataPelajarans[0].Kode)
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

func TestListFilterMataPelajaran(t *testing.T) {
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
							Field:    "kode",
							Operator: "contains",
							Value:    "PPKN",
						},
					},
				},
			},
			expectedLength: 3,
			expectedTotal:  3,
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
							Field:    "kode",
							Operator: "contains",
							Value:    "NextWhat",
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
			mataPelajarans, total, err := mataPelajaranService.ListMata_Pelajarans(&v.query)
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

				if len(mataPelajarans) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(mataPelajarans))
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

func TestListSortMataPelajaran(t *testing.T) {
	testScenarios := []struct {
		scenarioName        string
		query               query.GridParams
		expectedLength      int
		expectedTotal       int
		expectedSortedKodes []string
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
							Field:    "kode",
							Operator: "contains",
							Value:    "PPKN",
						},
					},
				},
				HasSort: true,
				Sort: []query.GridSort{
					query.GridSort{
						Field: "kode",
						Dir:   "asc",
					},
				},
			},
			expectedLength:      3,
			expectedTotal:       3,
			expectedSortedKodes: []string{"KD_PPKN_1", "KD_PPKN_2", "KD_PPKN_3"},
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			mataPelajarans, total, err := mataPelajaranService.ListMata_Pelajarans(&v.query)
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

				if len(mataPelajarans) != v.expectedLength {
					t.Errorf("expect len %d, but got %d", v.expectedLength, len(mataPelajarans))
					return
				}

				if mataPelajarans[0].Kode != v.expectedSortedKodes[0] {
					t.Errorf("expect name %s and index 0, but got %s", v.expectedSortedKodes[0], mataPelajarans[0].Kode)
					return
				}

				if mataPelajarans[2].Kode != v.expectedSortedKodes[2] {
					t.Errorf("expect name %s and index 2, but got %s", v.expectedSortedKodes[2], mataPelajarans[2].Kode)
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

func TestGetMataPelajaran(t *testing.T) {
	testScenarios := []struct {
		scenarioName   string
		id             string
		nama           string
		expectedErrMsg string
	}{
		{
			scenarioName: "Successful get by id",
			id:           "1",
			nama:         "PPKN",
		},
		{
			scenarioName:   "Failure get: mata pelajaran with id not exists",
			id:             "10",
			nama:           "NextWhat",
			expectedErrMsg: "Mata_Pelajaran with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			getMataPelajaranResponse, err := mataPelajaranService.GetMata_Pelajaran(v.id)
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
				if getMataPelajaranResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.nama != getMataPelajaranResponse.Nama {
					t.Errorf("expect name %s, but got %s", v.nama, getMataPelajaranResponse.Nama)
					return
				}
			}

		})
	}
}

func TestUpdateMataPelajaran(t *testing.T) {
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
			nama:            "PPKN2",
			expectedNama:    "PPKN2",
			expectedTingkat: 1,
		},
		{
			scenarioName:    "Successful tingkat update by id",
			id:              "1",
			expectedNama:    "PPKN2",
			tingkat:         2,
			expectedTingkat: 2,
		},
		{
			scenarioName:   "Failure update: mata pelajaran with id not exists",
			id:             "10",
			expectedErrMsg: "Mata_Pelajaran with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			updatedMataPelajaranResponse, err := mataPelajaranService.UpdateMata_Pelajaran(v.id, &schema.UpdateMata_PelajaranRequest{
				Nama:    v.nama,
				Kode:    v.kode,
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
				if updatedMataPelajaranResponse == nil {
					t.Errorf("expect response, but got nil")
					return
				}

				if v.expectedNama != updatedMataPelajaranResponse.Nama {
					t.Errorf("expect nama %s, but got %s", v.expectedNama, updatedMataPelajaranResponse.Nama)
					return
				}

				if v.expectedTingkat != updatedMataPelajaranResponse.Tingkat {
					t.Errorf("expect tingkat %d, but got %d", v.expectedTingkat, updatedMataPelajaranResponse.Tingkat)
					return
				}
			}
		})
	}
}

func TestDeleteMataPelajaran(t *testing.T) {
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
			scenarioName:   "Failure delete: mata pelajaran with id not exists",
			id:             "10",
			expectedErrMsg: "Mata_Pelajaran with id: 10 is not exists",
		},
	}

	for _, v := range testScenarios {
		t.Run(v.scenarioName, func(t *testing.T) {
			err := mataPelajaranService.DeleteMata_Pelajaran(v.id)
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
