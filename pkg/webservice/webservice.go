package webservice

import (
	"encoding/json"
	"net/http"
	"strconv"

	"s0counter/global"
)

func init() {
	InitWebService()
}

func InitWebService() (err error) {
	for pattern, f := range map[string]func(http.ResponseWriter, *http.Request){
		"version":     httpGetVersion,
		"currentdata": httpReadCurrentData,
	} {
		if ok, set := global.Config.Webserver.Webservices[pattern]; ok && set {
			http.HandleFunc("/"+pattern, f)
		}
	}

	port := ":" + strconv.Itoa(global.Config.Webserver.Port)
	go http.ListenAndServe(port, nil)
	return
}

// httpGetVersion prints the SW Version
func httpGetVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(global.VERSION)); err != nil {
		errorLog.Println(err)
		return
	}
}

// httpReadCurrentData supplies the data of al meters
func httpReadCurrentData(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(global.AllMeters, "", "  ")
	if err != nil {
		errorLog.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		errorLog.Println(err)
		return
	}
}
