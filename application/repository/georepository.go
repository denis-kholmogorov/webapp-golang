package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"web/application/domain"
	"web/application/errorhandler"
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

func (r GeoRepository) FindAll() []domain.Country {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	vars, err := txn.Query(ctx, getAllCountries)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("GeoRepository:FindAll() Error query %s", err)})
	}
	captchaList := domain.CountriesList{}
	err = json.Unmarshal(vars.Json, &captchaList)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("GeoRepository:FindAll() Error Unmarshal %s", err)})
	}
	return captchaList.List

}

func (r GeoRepository) FindCitiesByCountryId(id string) []domain.City {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$countryId"] = id
	vars, err := txn.QueryWithVars(ctx, getAllCities, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("GeoRepository:FindCitiesByCountryId() Error query %s", err)})
	}
	captchaList := domain.CountriesList{}
	err = json.Unmarshal(vars.Json, &captchaList)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("GeoRepository:FindCitiesByCountryId() Error Unmarshal %s", err)})
	}
	return captchaList.List[0].Cities
}

var getAllCountries = `{ countriesList (func: type(Country)) {
	id : uid
	title: countryTitle
    cities{
		id:uid
		title: cityTitle
		}
	}
}`

var getAllCities = `query Cities($countryId: string) { countriesList (func: uid($countryId)) {
	cities{
		id: uid
		title: cityTitle
    	}
	}
}`
