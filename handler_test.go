package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

const testRUC = "20601030013"
const basicRecord = "REXTIE S.A.C.>1>1>150131>CAL.>LAS CAMELIAS>URB.>JARDIN>256>701A>->->->-"
const advancedRecord = "SAC,SERVICIOS,10,C,M,0,150131"

// These tests share the application's globals and must not run in parallel.
func setupCompanyCache(t *testing.T) *miniredis.Miniredis {
	t.Helper()
	server := miniredis.RunT(t)
	previousClient, previousDepartments := redisClient, DEPARTAMENTOS
	redisClient = redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { redisClient.Close(); redisClient = previousClient; DEPARTAMENTOS = previousDepartments })
	data, err := os.ReadFile("ubigeos.json")
	if err != nil {
		t.Fatal(err)
	}
	DEPARTAMENTOS = map[string]Departamento{}
	if err := json.Unmarshal(data, &DEPARTAMENTOS); err != nil {
		t.Fatal(err)
	}
	for key, value := range map[string]string{"r" + testRUC + "1": basicRecord, "r" + testRUC + "2": advancedRecord, "r101234567811": basicRecord} {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if err := server.Set(key, string(encoded)); err != nil {
			t.Fatal(err)
		}
	}
	return server
}

func expectedV2() map[string]interface{} {
	return map[string]interface{}{
		"razonSocial": "REXTIE S.A.C.", "tipoDocumento": "6", "numeroDocumento": testRUC,
		"estado": "ACTIVO", "condicion": "HABIDO", "direccion": "CAL. LAS CAMELIAS NRO 256 INT. 701A URB. JARDIN ",
		"ubigeo": "150131", "viaTipo": "CAL.", "viaNombre": "LAS CAMELIAS", "zonaCodigo": "URB.", "zonaTipo": "JARDIN",
		"numero": "256", "interior": "701A", "lote": "-", "dpto": "-", "manzana": "-", "kilometro": "-",
		"distrito": "SAN ISIDRO", "provincia": "LIMA", "departamento": "LIMA",
		"EsAgenteRetencion": false, "EsBuenContribuyente": false, "localesAnexos": nil,
	}
}

func TestEndpointsHappyPath(t *testing.T) {
	setupCompanyCache(t)
	app := NewApp()
	for _, tc := range []struct {
		name, path, nameKey, numberKey string
		full, pro5                     bool
	}{
		{"basic", "/ruc?numero=" + testRUC, "razon_social", "numero_documento", false, false},
		{"basic DNI", "/ruc?dni=12345678", "razon_social", "numero_documento", false, false},
		{"full", "/ruc/full?numero=" + testRUC, "razon_social", "numero_documento", true, false},
		{"pro5", "/ruc/" + testRUC, "nombre_o_razon_social", "ruc", false, true},
		{"V1", "/apisV1/ruc?numero=" + testRUC, "nombre", "numeroDocumento", false, false},
		{"V2", "/apisV2/ruc?numero=" + testRUC, "razonSocial", "numeroDocumento", false, false},
		{"full V2", "/apisV2/ruc/full?numero=" + testRUC, "razonSocial", "numeroDocumento", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response, err := app.Test(httptest.NewRequest("GET", tc.path, nil))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != 200 {
				t.Fatalf("status = %d", response.StatusCode)
			}
			if !strings.Contains(response.Header.Get("Content-Type"), "application/json") {
				t.Fatal("expected JSON response")
			}
			var got map[string]interface{}
			if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if tc.pro5 {
				if got["success"] != true {
					t.Fatalf("response = %#v", got)
				}
				data, ok := got["data"].(map[string]interface{})
				if !ok {
					t.Fatal("missing data")
				}
				got = data
			}
			number := testRUC
			if tc.name == "basic DNI" {
				number = "10123456781"
			}
			for key, want := range map[string]interface{}{tc.nameKey: "REXTIE S.A.C.", tc.numberKey: number, "estado": "ACTIVO", "condicion": "HABIDO", "distrito": "SAN ISIDRO"} {
				if got[key] != want {
					t.Errorf("%s = %v, want %v", key, got[key], want)
				}
			}
			if tc.full && got["tipo"] != "SAC" {
				t.Errorf("tipo = %v", got["tipo"])
			}
			if tc.name == "V2" || tc.name == "full V2" {
				want := expectedV2()
				if tc.full {
					for key, value := range map[string]interface{}{"tipo": "SAC", "actividadEconomica": "SERVICIOS", "numeroTrabajadores": "10", "tipoFacturacion": "COMPUTARIZADO", "tipoContabilidad": "MANUAL", "comercioExterior": "SIN ACTIVIDAD"} {
						want[key] = value
					}
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("JSON contract mismatch\ngot: %#v\nwant: %#v", got, want)
				}
			}
		})
	}
	t.Run("home", func(t *testing.T) {
		response, err := app.Test(httptest.NewRequest("GET", "/", nil))
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 200 || !strings.Contains(string(body), "Welcome to the SUNAT RUC API") {
			t.Fatalf("unexpected home response: %d %s", response.StatusCode, body)
		}
	})
}

func TestCacheRepositoryRoundTrip(t *testing.T) {
	server := setupCompanyCache(t)
	want := map[string]string{"name": "REXTIE S.A.C."}
	if err := SetToCacheRepostory("test-company", want, time.Hour); err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if err := GetFromCacheRepository("test-company", &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if ttl := server.TTL("test-company"); ttl != time.Hour {
		t.Fatalf("TTL = %v", ttl)
	}
}
