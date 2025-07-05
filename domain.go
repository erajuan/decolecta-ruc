package main

import (
	"regexp"
	"strings"
)

// Provincia de un departamento con sus distritos
type Provincia struct {
	Nombre    string            `json:"nombre"`
	Distritos map[string]string `json:"distritos"`
}

// Departamento con sus provincias
type Departamento struct {
	Nombre     string               `json:"nombre"`
	Provincias map[string]Provincia `json:"provincias"`
}

// DEPARTAMENTOS del peru
var DEPARTAMENTOS = map[string]Departamento{}

// Company
type Company struct {
	RazonSocial         string           `json:"razonSocial"`
	NumeroDocumento     string           `json:"numeroDocumento"`
	Estado              string           `json:"estado"`
	Condicion           string           `json:"condicion"`
	Direccion           string           `json:"direccion"`
	Ubigeo              string           `json:"ubigeo"`
	ViaTipo             string           `json:"viaTipo"`
	ViaNombre           string           `json:"viaNombre"`
	ZonaCodigo          string           `json:"zonaCodigo"`
	ZonaTipo            string           `json:"zonaTipo"`
	Numero              string           `json:"numero"`
	Interior            string           `json:"interior"`
	Lote                string           `json:"lote"`
	Dpto                string           `json:"dpto"`
	Manzana             string           `json:"manzana"`
	Kilometro           string           `json:"kilometro"`
	Distrito            string           `json:"distrito"`
	Provincia           string           `json:"provincia"`
	Departamento        string           `json:"departamento"`
	EsAgenteRetencion   bool             `json:"EsAgenteRetencion"`
	EsBuenContribuyente bool             `json:"EsBuenContribuyente"`
	LocalesAnexos       []CompanyAddress `json:"localesAnexos"`
}

// CompanyAdvance
type CompanyAdvance struct {
	Company
	Tipo               string `json:"tipo"`
	ActividadEconomica string `json:"actividadEconomica"`
	NumeroTrabajadores string `json:"numeroTrabajadores"`
	TipoFacturacion    string `json:"tipoFacturacion"`
	TipoContabilidad   string `json:"tipoContabilidad"`
	ComercioExterior   string `json:"comercioExterior"`
}

type CompanyDTO struct {
	Data                string
	Extras              []SunatExtrasDTO
	EsAgenteRetencion   bool
	EsBuenContribuyente bool
	Locales             []string
}

type CompanyAddress struct {
	Direccion    string `json:"direccion"`
	Ubigeo       string `json:"ubigeo"`
	Departamento string `json:"departamento"`
	Provincia    string `json:"provincia"`
	Distrito     string `json:"distrito"`
}

type SunatExtrasDTO struct {
	TypeId   int
	Position int
	Data     string
}

type CompanyWithholdingAgent struct {
	Ruc               string `json:"ruc"`
	EsAgenteRetencion bool   `json:"EsAgenteRetencion"`
	StartAt           string `json:"desde"`
	Resolution        string `json:"resolucion"`
}

func (c *Company) Extender(cached string) {
	data := strings.Split(cached, ";;")
	c.RazonSocial = data[0]

	switch data[1] {
	case "1":
		c.Estado = "ACTIVO"
	case "0":
		c.Estado = "NO ACTIVO"
	default:
		c.Estado = data[1]
	}

	switch data[2] {
	case "1":
		c.Condicion = "HABIDO"
	case "0":
		c.Condicion = "NO HABIDO"
	default:
		c.Condicion = data[2]
	}

	if len(data) == 3 {
		c.Ubigeo = "-"
		c.Direccion = "-"
		c.ViaTipo = "-"
		c.ViaNombre = "-"
		c.ZonaCodigo = "-"
		c.ZonaTipo = "-"
		c.Numero = "-"
		c.Interior = "-"
		c.Lote = "-"
		c.Dpto = "-"
		c.Manzana = "-"
		c.Kilometro = "-"
	} else {
		re1, _ := regexp.Compile("^[-]+$")
		c.Ubigeo = data[3]
		if re1.FindStringIndex(data[4]) == nil {
			c.ViaTipo = data[4]
		} else {
			c.ViaTipo = "-"
		}
		c.ViaNombre = data[5]

		if re1.FindStringIndex(data[6]) == nil {
			c.ZonaCodigo = data[6]
		} else {
			c.ZonaCodigo = "-"
		}

		if re1.FindStringIndex(data[7]) == nil {
			c.ZonaTipo = data[7]
		} else {
			c.ZonaTipo = "-"
		}
		c.Numero = data[8]
		c.Interior = data[9]
		c.Lote = data[10]
		c.Dpto = data[11]
		c.Manzana = data[12]
		c.Kilometro = data[13]
		c.Completar()
	}
}

func (c *CompanyAdvance) AdvanceExtender(cached string) {
	data := strings.Split(cached, ",")
	switch data[0] {
	case "EIRL":
		c.Tipo = "EMPRESA INDIVIDUAL DE RESPONSABILIDAD LIMITADA"
	case "0":
		c.Tipo = "PERSONA NATURAL SIN EMPRESA"
	default:
		c.Tipo = data[0]
	}
	if data[1] == "0" {
		c.ActividadEconomica = "NO DISPONIBLE"
	} else {
		c.ActividadEconomica = data[1]
	}
	if data[2] == "0" {
		c.NumeroTrabajadores = "NO DISPONIBLE"
	} else {
		c.NumeroTrabajadores = data[2]
	}
	switch data[3] {
	case "C":
		c.TipoFacturacion = "COMPUTARIZADO"
	case "M":
		c.TipoFacturacion = "MANUAL"
	default:
		c.TipoFacturacion = data[3]
	}
	switch data[4] {
	case "C":
		c.TipoContabilidad = "COMPUTARIZADO"

	case "M":
		c.TipoContabilidad = "MANUAL"

	default:
		c.TipoContabilidad = data[4]
	}
	if data[5] == "0" {
		c.ComercioExterior = "SIN ACTIVIDAD"
	} else {
		c.ComercioExterior = data[5]
	}
	c.Ubigeo = data[6]
	// Departamento-Provincia-Distrito
	dep := DEPARTAMENTOS[c.Ubigeo[0:2]]
	prov := dep.Provincias[c.Ubigeo[2:4]]
	c.Departamento = dep.Nombre
	c.Provincia = prov.Nombre
	c.Distrito = prov.Distritos[c.Ubigeo[4:6]]
}
func GetLocalAnexoAddress(data string) CompanyAddress {
	address := strings.Split(data, ";")
	longAddress := ""
	if address[0] == "-" {
		longAddress += address[1] + " " + address[2]
	}
	if address[3] != "-" {
		longAddress += " " + address[3]
	}
	if address[4] != "-" {
		longAddress += " " + address[4]
	}
	if address[5] != "-" {
		longAddress += " NRO " + address[5]
	}
	if address[6] != "-" {
		longAddress += " KM " + address[6]
	}
	if address[7] != "-" {
		longAddress += " INT. " + address[7]
	}
	if address[8] != "-" {
		longAddress += " LT " + address[8]
	}
	if address[9] != "-" {
		longAddress += " DEP. " + address[9]
	}
	if address[10] != "-" {
		longAddress += " MZ " + address[10]
	}
	companyAddress := CompanyAddress{
		Ubigeo:    strings.TrimSpace(address[0]),
		Direccion: strings.TrimSpace(longAddress),
	}
	companyAddress.Departamento, companyAddress.Provincia, companyAddress.Distrito = getUbigeo(address[0])
	return companyAddress
}
func getUbigeo(ubigeo string) (string, string, string) {
	dep := DEPARTAMENTOS[ubigeo[0:2]]
	prov := dep.Provincias[ubigeo[2:4]]
	return dep.Nombre, prov.Nombre, prov.Distritos[ubigeo[4:6]]
}

// Completar direccion
func (c *Company) Completar() {

	// Direccion
	if c.ViaTipo != "-" {
		c.Direccion += c.ViaTipo + " "
	}
	if c.ViaNombre != "-" {
		c.Direccion += c.ViaNombre + " "
	}
	// NUMERO
	if c.Numero != "-" {
		c.Direccion += "NRO " + c.Numero + " "
	}
	// INTERIOR
	if c.Interior != "-" {
		c.Direccion += "INT. " + c.Interior + " "
	}
	// MANZANA
	if c.Manzana != "-" {
		c.Direccion += "MZA. " + c.Manzana + " "
	}
	// LOTE
	if c.Lote != "-" {
		c.Direccion += "LOTE " + c.Lote + " "
	}
	// DPTO
	if c.Dpto != "-" {
		c.Direccion += "DEP. " + c.Dpto + " "
	}
	// ZONA
	if c.ZonaCodigo != "-" {
		c.Direccion += c.ZonaCodigo + " "
	}
	if c.ZonaTipo != "-" {
		c.Direccion += c.ZonaTipo + " "
	}
	// KM
	if c.Kilometro != "-" {
		c.Direccion += "KM. " + c.Kilometro + " "
	}
	// End Direccion
	// Departamento-Provincia-Distrito
	dep := DEPARTAMENTOS[c.Ubigeo[0:2]]
	prov := dep.Provincias[c.Ubigeo[2:4]]
	c.Departamento = dep.Nombre
	c.Provincia = prov.Nombre
	c.Distrito = prov.Distritos[c.Ubigeo[4:6]]
}

type CompanyPro5 struct {
	NombreRazonSocial   string   `json:"nombre_o_razon_social"`
	Ruc                 string   `json:"ruc"`
	Estado              string   `json:"estado"`
	Condicion           string   `json:"condicion"`
	Direccion           string   `json:"direccion"`
	DireccionCompleta   string   `json:"dirección_completa"`
	Distrito            string   `json:"distrito"`
	Provincia           string   `json:"provincia"`
	Departamento        string   `json:"departamento"`
	UbigeoSunat         string   `json:"ubigeo_sunat"`
	Ubigeo              []string `json:"ubigeo"`
	EsAgenteRetencion   bool     `json:"es_agente_de_retencion"`
	EsBuenContribuyente bool     `json:"es_buen_contribuyente"`
}

func (c Company) ToPro5() CompanyPro5 {
	ubigeo := []string{"", "", ""}
	if len(c.Ubigeo) == 6 {
		ubigeo[0] = c.Ubigeo[0:2]
		ubigeo[1] = c.Ubigeo[0:4]
		ubigeo[2] = c.Ubigeo
	}
	return CompanyPro5{
		NombreRazonSocial:   c.RazonSocial,
		Ruc:                 c.NumeroDocumento,
		Estado:              c.Estado,
		Condicion:           c.Condicion,
		Direccion:           c.Direccion,
		DireccionCompleta:   c.Direccion,
		UbigeoSunat:         c.Ubigeo,
		Ubigeo:              ubigeo,
		Distrito:            c.Distrito,
		Provincia:           c.Provincia,
		Departamento:        c.Departamento,
		EsAgenteRetencion:   false,
		EsBuenContribuyente: false,
	}
}
