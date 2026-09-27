package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFullV2EmptyAdvancedFields(t *testing.T) {
	setupCompanyCache(t)
	company := CompanyAdvance{Company: Company{NumeroDocumento: testRUC}}
	company.Extender(basicRecord)
	encoded, err := json.Marshal(company.ToApisV2())
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	want := expectedV2()
	for _, field := range []string{"tipo", "actividadEconomica", "numeroTrabajadores", "tipoFacturacion", "tipoContabilidad", "comercioExterior"} {
		want[field] = ""
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %#v", encoded, want)
	}
}

func TestV2FlagsAndAnnexes(t *testing.T) {
	annex := CompanyAddress{Direccion: "AV. LIMA 123", Ubigeo: "150131", Departamento: "LIMA", Provincia: "LIMA", Distrito: "SAN ISIDRO"}
	company := CompanyAdvance{Company: Company{EsAgenteRetencion: true, EsBuenContribuyente: true, LocalesAnexos: []CompanyAddress{annex}}}
	data, err := json.Marshal(company.ToApisV2())
	if err != nil {
		t.Fatal(err)
	}
	var got CompanyAdvanceApisV2
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !got.EsAgenteRetencion || !got.EsBuenContribuyente || !reflect.DeepEqual(got.LocalesAnexos, []CompanyAddress{annex}) {
		t.Fatalf("flags or annexes lost: %s", data)
	}
}

func TestAdvanceExtenderHappyPath(t *testing.T) {
	setupCompanyCache(t)
	company := CompanyAdvance{}
	company.AdvanceExtender("EIRL,0,0,M,C,0,150131")
	if company.Tipo != "EMPRESA INDIVIDUAL DE RESPONSABILIDAD LIMITADA" || company.ActividadEconomica != "NO DISPONIBLE" || company.NumeroTrabajadores != "NO DISPONIBLE" || company.TipoFacturacion != "MANUAL" || company.TipoContabilidad != "COMPUTARIZADO" || company.ComercioExterior != "SIN ACTIVIDAD" || company.Distrito != "SAN ISIDRO" {
		t.Fatalf("unexpected company: %+v", company)
	}
}

func TestCompanyWithoutAddress(t *testing.T) {
	company := Company{}
	company.Extender("REXTIE S.A.C.>1>1")
	if company.RazonSocial != "REXTIE S.A.C." || company.Estado != "ACTIVO" || company.Condicion != "HABIDO" || company.Direccion != "-" || company.Ubigeo != "-" {
		t.Fatalf("unexpected company: %+v", company)
	}
}

func TestGetLocalAnexoAddress(t *testing.T) {
	setupCompanyCache(t)
	got := GetLocalAnexoAddress("150131;-;-;AV.;LIMA;123;-;4;-;-;-")
	want := CompanyAddress{Direccion: "AV. LIMA NRO 123 INT. 4", Ubigeo: "150131", Departamento: "LIMA", Provincia: "LIMA", Distrito: "SAN ISIDRO"}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
