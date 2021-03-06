package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"

	Exchange "madaoQT/exchange"
	Mongo "madaoQT/mongo"
)

type ExchangeController struct {
	Ctx iris.Context
	// [ Your fields here ]
	// Request lifecycle data
	// Models
	// Database
	// Global properties
	Sessions *sessions.Sessions `iris:"persistence"`
}

type ExchangeInfo struct {
	Name   string `json:"name"`
	API    string `json:"api"`
	Secret string `json:"secret"`
}

func (e *ExchangeController) authen() (bool, iris.Map) {
	if DEBUG {
		return true, iris.Map{}
	}
	{
		s := e.Sessions.Start(e.Ctx)
		username := s.Get("name")
		if username == nil || username == "" {
			result := iris.Map{
				"result": false,
				"error":  errorCodeInvalidSession,
			}
			return false, result
		}
		return true, iris.Map{}
	}

}

// Get: /exchange/list
func (e *ExchangeController) GetList() iris.Map {
	var exchangeList []string

	for _, exchange := range Exchange.ExchangeNameList {
		exchangeList = append(exchangeList, exchange)
	}

	return iris.Map{
		"result":    true,
		"exchanges": exchangeList,
	}
}

func (e *ExchangeController) PostAddkey() iris.Map {

	var errMsg string
	var err error
	// var encryptedAPI, encryptedSecret []byte

	// s := e.Sessions.Start(e.Ctx)
	// username := s.Get("name")
	// password := s.Get("password")

	exchangesDB := Mongo.ExchangeDB{}

	info := ExchangeInfo{}

	// if username == nil || password == nil {
	// 	errMsg = "Invalid Session"
	// 	goto _ERROR
	// }

	err = e.Ctx.ReadJSON(&info)
	if err != nil {
		errMsg = err.Error()
		goto _ERROR
	}

	if err = exchangesDB.Connect(); err != nil {
		errMsg = err.Error()
		goto _ERROR
	}

	// err, encryptedAPI = Utils.GCM_encrypt(password.(string), username.(string), info.APIKey)
	// if err != nil {
	// 	errMsg = err.Error()
	// 	goto _ERROR
	// }

	// err, encryptedSecret = Utils.GCM_encrypt(password.(string), username.(string), info.SecretKey)
	// if err != nil {
	// 	errMsg = err.Error()
	// 	goto _ERROR
	// }

	if err = exchangesDB.Insert(&Mongo.ExchangeInfo{
		Name: info.Name,
		// API:    string(encryptedAPI),
		// Secret: string(encryptedSecret),
		// User:   username.(string),
		API:    []byte(info.API),
		Secret: []byte(info.Secret),
	}); err != nil {
		errMsg = "Fail to insert record into database"
		goto _ERROR
	}

	return iris.Map{
		"result": true,
	}

_ERROR:
	return iris.Map{
		"result": false,
		"error":  errMsg,
	}
}
