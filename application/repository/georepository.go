package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"log"
	"web/application/domain"
)

var geoRepo *GeoRepository
var isInitializedGeoRepo bool

type GeoRepository struct {
	conn *dgo.Dgraph
}

func GetGeoRepository() *GeoRepository {
	if !isInitializedGeoRepo {
		geoRepo = &GeoRepository{}
		geoRepo.conn = GetDGraphConn().connection
		isInitializedGeoRepo = true
	}
	return geoRepo
}

func (r GeoRepository) FindAll() ([]domain.Country, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	vars, err := txn.Query(ctx, getAllCountries)
	if err != nil {
		log.Printf("GeoRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("GeoRepository:FindAll() Error query %s", err)
	}
	captchaList := domain.CountriesList{}

	err = json.Unmarshal(vars.Json, &captchaList)
	return captchaList.List, nil

}

func (r GeoRepository) FindCitiesByCountryId(id string) ([]domain.City, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$countryId"] = id
	vars, err := txn.QueryWithVars(ctx, getAllCities, variables)

	captchaList := domain.CountriesList{}
	if err != nil {
		log.Printf("GeoRepository:FindCitiesByCountryId() Error query %s", err)
		return nil, fmt.Errorf("GeoRepository:FindCitiesByCountryId() Error query %s", err)
	}
	err = json.Unmarshal(vars.Json, &captchaList)
	if err != nil {
		log.Printf("GeoRepository:FindCitiesByCountryId() Error Unmarshal %s", err)
		return nil, fmt.Errorf("GeoRepository:FindCitiesByCountryId() Error Unmarshal %s", err)
	}
	return captchaList.List[0].Cities, nil
}

var getAllCountries = `{ countriesList (func: type(Country)) {
	id : uid
	title
	}
}`

var getAllCities = `query Cities($countryId: string) { countriesList (func: uid($countryId)) {
	cities{
		uid
		title
    	}
	}
}`
