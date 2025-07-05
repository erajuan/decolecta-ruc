package main

import (
	"fmt"
	"time"
)

func getCompanyKey(ruc string, group string) string {
	return fmt.Sprintf("%s%s%s", "r", ruc, group)
}
func GetCompanyService(ruc string) (Company, error) {
	key := getCompanyKey(ruc, "1")
	var cached string
	err := GetFromCacheRepository(key, &cached)
	c := Company{
		NumeroDocumento: ruc,
	}
	if err == nil {
		c.Extender(cached)
		return c, err
	}

	data, err1 := GetCompanyRepository(ruc)
	if err1 != nil {
		return Company{}, err
	}
	c.Extender(data.Data)
	c.EsAgenteRetencion = data.EsAgenteRetencion
	c.EsBuenContribuyente = data.EsBuenContribuyente
	for _, local := range data.Locales {
		c.LocalesAnexos = append(c.LocalesAnexos, GetLocalAnexoAddress(local))
	}
	SetToCacheRepostory(key, data.Data, time.Hour*24)
	return c, nil
}

func GetCompanyAdvanceService(ruc string) (CompanyAdvance, error) {
	key := getCompanyKey(ruc, "1")
	var cached1 string
	err := GetFromCacheRepository(key, &cached1)
	c := Company{
		NumeroDocumento: ruc,
	}
	if err == nil {
		c.Extender(cached1)
	} else {
		data1, err1 := GetCompanyRepository(ruc)
		if err1 != nil {
			return CompanyAdvance{}, err
		}
		c.EsAgenteRetencion = data1.EsAgenteRetencion
		c.Extender(data1.Data)
		c.EsAgenteRetencion = data1.EsAgenteRetencion
		c.EsBuenContribuyente = data1.EsBuenContribuyente
		for _, local := range data1.Locales {
			c.LocalesAnexos = append(c.LocalesAnexos, GetLocalAnexoAddress(local))
		}
	}
	ca := CompanyAdvance{
		Company: c,
	}
	key2 := getCompanyKey(ruc, "2")
	var cached2 string
	err = GetFromCacheRepository(key2, &cached2)
	if err == nil {
		ca.AdvanceExtender(cached2)
		return ca, nil
	}
	data2, err2 := GetCompanyFullRepository(ruc)
	if err2 != nil {
		return ca, nil
	}
	ca.AdvanceExtender(data2.Data)
	SetToCacheRepostory(key2, data2.Data, time.Hour)
	return ca, nil
}
