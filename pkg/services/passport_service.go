package services

import (
	"encoding/json"
	"fmt"
	"github.com/AntonBorzenko/RestrictedPassportService/config"
	"github.com/AntonBorzenko/RestrictedPassportService/utils"
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	"github.com/AntonBorzenko/RestrictedPassportService/utils/net"
	"github.com/AntonBorzenko/RestrictedPassportService/utils/passports"
	"github.com/AntonBorzenko/RestrictedPassportService/utils/uint_set"
	"log"
	"net/http"
	"os"
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
	serviceUriPrefix := "/passportApi"
	http.HandleFunc(serviceUriPrefix, okResponse)
	http.HandleFunc(serviceUriPrefix + "/check", service.checkPassport)

	if config.Cfg.SslKey != "" && config.Cfg.SslCert != "" {
		log.Println("starting server at port 443")
		e.Check(http.ListenAndServeTLS(":443", config.Cfg.SslCert, config.Cfg.SslKey, nil))
	} else {
		log.Printf("starting server at port %v\n", config.Cfg.Port)
		e.Check(http.ListenAndServe(fmt.Sprintf(":%v", config.Cfg.Port), nil))
	}
}

func (service *PassportService) UpdateDb() {
	e.Check(passports.RemovePreviousFiles())

	tempFileName := e.CheckString(utils.CreateTempFile("passports_*.bz2"))
	defer os.Remove(tempFileName)

	log.Printf("Downloading file '%v' from url '%v'...\n", tempFileName, config.Cfg.DBFileUrl)
	e.Check(net.DownloadFile(tempFileName, config.Cfg.DBFileUrl))

	passportGenerator := passports.GetPassportsGenerator(tempFileName, config.Cfg.DBBatchSize)
	e.Check(service.Set.InsertMultiple(passportGenerator, false))

	if set, ok := service.Set.(*uint_set.SqliteSet); ok {
		set.CreateIndex()
	}
}
