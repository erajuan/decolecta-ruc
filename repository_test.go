package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func setupCompanyDB(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	previous := pgsqlClient
	pgsqlClient = mock
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		mock.Close()
		pgsqlClient = previous
	})
	return mock
}

func expectBasicCompany(mock pgxmock.PgxPoolIface) {
	mock.ExpectQuery(`SELECT content FROM sunat_rucs WHERE ruc=\$1 LIMIT 1;`).WithArgs(testRUC).
		WillReturnRows(pgxmock.NewRows([]string{"content"}).AddRow(basicRecord))
	mock.ExpectQuery(`SELECT type_id, content FROM sunat_ruc_extras WHERE ruc=\$1 LIMIT 40;`).WithArgs(testRUC).
		WillReturnRows(pgxmock.NewRows([]string{"type_id", "content"}).
			AddRow(SUNAT_BUEN_CONTRIBUYENTE, "").AddRow(SUNAT_AGENTE_RETENCION, "").
			AddRow(SUNAT_LOCAL_ANEXO, "150131;-;-;AV.;LIMA;123;-;4;-;-;-"))
}

func expectAdvancedCompany(mock pgxmock.PgxPoolIface) {
	mock.ExpectQuery(`SELECT content FROM sunat_ruc_extenses WHERE ruc=\$1 LIMIT 1;`).WithArgs(testRUC).
		WillReturnRows(pgxmock.NewRows([]string{"content"}).AddRow(advancedRecord))
}

func TestDatabaseBackedEndpointsHappyPath(t *testing.T) {
	for _, tc := range []struct {
		path   string
		full   bool
		format string
	}{
		{"/ruc?numero=" + testRUC, false, "native"},
		{"/ruc/full?numero=" + testRUC, true, "native"},
		{"/apisV1/ruc?numero=" + testRUC, false, "v1"},
		{"/apisV2/ruc?numero=" + testRUC, false, "v2"},
		{"/apisV2/ruc/full?numero=" + testRUC, true, "v2"},
		{"/ruc/" + testRUC, false, "pro5"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			server := setupCompanyCache(t)
			server.FlushAll()
			mock := setupCompanyDB(t)
			expectBasicCompany(mock)
			if tc.full {
				expectAdvancedCompany(mock)
			}
			response, err := NewApp().Test(httptest.NewRequest("GET", tc.path, nil))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != 200 {
				t.Fatalf("status = %d", response.StatusCode)
			}
			var got map[string]interface{}
			if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if tc.format == "pro5" {
				if got["success"] != true {
					t.Fatalf("response = %v", got)
				}
				got = got["data"].(map[string]interface{})
			}
			if got["estado"] != "ACTIVO" || got["distrito"] != "SAN ISIDRO" {
				t.Fatalf("unexpected company: %v", got)
			}
			if tc.format == "v2" || tc.format == "native" {
				retention, contributor, annexes := "EsAgenteRetencion", "EsBuenContribuyente", "localesAnexos"
				if tc.format == "native" {
					retention, contributor, annexes = "es_agente_retencion", "es_buen_contribuyente", "locales_anexos"
				}
				if got[retention] != true || got[contributor] != true {
					t.Errorf("missing flags: %v", got)
				}
				addresses, ok := got[annexes].([]interface{})
				if !ok || len(addresses) != 1 {
					t.Fatalf("unexpected annexes: %v", got[annexes])
				}
				address := addresses[0].(map[string]interface{})
				if address["direccion"] != "AV. LIMA NRO 123 INT. 4" || address["distrito"] != "SAN ISIDRO" {
					t.Errorf("unexpected annex: %v", address)
				}
			}
			if tc.full {
				field := "actividadEconomica"
				if tc.format == "native" {
					field = "actividad_economica"
				}
				if got[field] != "SERVICIOS" {
					t.Errorf("%s = %v", field, got[field])
				}
			}
		})
	}
}
