package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func GetCompanyRepository(ruc string) (CompanyDTO, error) {
	db := pgsqlClient
	var c CompanyDTO
	err := db.QueryRow(CTX, `
		SELECT content
		FROM sunat_rucs
		WHERE ruc=$1
		LIMIT 1;`, ruc).Scan(&c.Data)
	if err != nil {
		fmt.Println("Error:", err)
		return CompanyDTO{}, err
	}

	rows, err1 := db.Query(CTX, `
	SELECT type_id, content
	FROM sunat_ruc_extras
	WHERE ruc=$1
	LIMIT 10;`, ruc)

	if err1 != nil {
		return c, nil
	}
	var typeId int
	var content string
	var locales []string = make([]string, 0)
	for rows.Next() {
		err := rows.Scan(&typeId, &content)
		if err != nil {
			continue
		}
		if typeId == SUNAT_BUEN_CONTRIBUYENTE || typeId == 1 {
			c.EsBuenContribuyente = true
		}
		if typeId == SUNAT_AGENTE_RETENCION || typeId == 2 {
			c.EsAgenteRetencion = true
		}
		if typeId == SUNAT_LOCAL_ANEXO || typeId == 3 {
			locales = append(locales, content)
		}
	}
	c.Locales = locales
	return c, nil
}

func GetCompanyFullRepository(ruc string) (CompanyDTO, error) {
	var c CompanyDTO
	err := pgsqlClient.QueryRow(CTX, `
		SELECT content FROM sunat_ruc_extenses
		WHERE ruc=$1
		LIMIT 1;`, ruc).Scan(&c.Data)
	if err != nil {
		return CompanyDTO{}, err
	}
	return c, nil
}

func GetFromCacheRepository(key string, value interface{}) error {
	cached, _ := redisClient.Get(CTX, key).Result()
	if cached == "" || cached == "404" {
		return errors.New("not cached")
	}
	return json.Unmarshal([]byte(cached), &value)
}
func SetToCacheRepostory(key string, value interface{}, expiration time.Duration) error {
	t := fmt.Sprintf("%T", value)
	if t == "string" {
		redisClient.Set(CTX, key, value, expiration).Err()
	} else {
		cache, _ := json.Marshal(value)
		redisClient.Set(CTX, key, string(cache), expiration).Err()
	}
	return nil
}
