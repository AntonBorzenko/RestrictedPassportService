package services

import (
	"encoding/json"
	"fmt"
	"github.com/AntonBorzenko/RestrictedPassportService/config"
	"github.com/AntonBorzenko/RestrictedPassportService/utils"
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	"github.com/AntonBorzenko/RestrictedPassportService/utils/passports"
	"github.com/AntonBorzenko/RestrictedPassportService/utils/uint_set"
	"log"
	"net/http"
)

type PassportService struct {
	Set uint_set.UintSet
}

func NewPassportService() *PassportService {
	set := uint_set.NewSqliteSet(config.Cfg.DBFile, true, false)
	return &PassportService{set}
}

func respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}

func respondError(w http.ResponseWriter, err string) {
	respond(w, map[string]interface{}{ "success": false, "error": err })
}

func okResponse(w http.ResponseWriter, r *http.Request) {
	respond(w, map[string]interface{} { "success": true })
}

func (service *PassportService) checkPassport(w http.ResponseWriter, r *http.Request) {
	passport, ok := r.URL.Query()["passport"]
	if !ok || len(passport) < 1 {
		respondError(w, "Url Param 'passport' is missing")
		return
	}

	passportNumber, err := passports.ConvertPassportToUint64(passport[0])
	if err != nil {
		respondError(w, "Cannot parse 'passport' Param")
		return
	}

	result, err := service.Set.Has(passportNumber)
	if err != nil {
		respondError(w, "Cannot parse 'passport' Param")
		return
	}

	respond(w, map[string]interface{} { "success": true, "passportBanned": result })
}

func (service *PassportService) StartApi() {
	if !utils.FileExists(config.Cfg.DBFile) {
		service.UpdateDb()
	}
	fmt.Printf("starting server at port %v\n", config.Cfg.Port)

	serviceUriPrefix := "/passport-api"
	http.HandleFunc(serviceUriPrefix, okResponse)
	http.HandleFunc(serviceUriPrefix + "/check", service.checkPassport)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Cfg.Port), nil))
}

func (service *PassportService) UpdateDb() {
	passportGenerator := passports.GetPassportsGenerator(config.Cfg.DBFileUrl, config.Cfg.DBBatchSize)
	set := service.Set.(*uint_set.SqliteSet)
	e.Check(set.InsertMultiple(passportGenerator, false))
	set.CreateIndex()
}
