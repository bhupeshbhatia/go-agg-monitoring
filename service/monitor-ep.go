package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bhupeshbhatia/go-agg-monitoring/connectDB"
	"github.com/bhupeshbhatia/go-agg-monitoring/model"
	"github.com/pkg/errors"
)

type MetSearch struct {
	EndDate   int64  `bson:"end_date,omitempty" json:"end_date,omitempty"`
	StartDate int64  `bson:"start_date,omitempty" json:"start_date,omitempty"`
	SearchKey string `bson:"search_key,omitempty" json:"search_key,omitempty"`
	SearchVal string `bson:"search_val,omitempty" json:"search_val,omitempty"`
}

type MetDash struct {
	Eth   float64
	Dates int64 `bson:"dates,omitempty" json:"dates,omitempty"`
}

func LoadDataInMongo(w http.ResponseWriter, r *http.Request) {
	// DB connection
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		return
	}

	metric := []model.Metric{}
	for i := 0; i < 100; i++ {
		metric = append(metric, GenSensorData())
	}

	for _, val := range metric {
		log.Println(val.ItemID)
		insertResult, err := Db.Collection.InsertOne(val)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert event")
			log.Println(err)
			return
		}
		log.Println(insertResult)
	}

	_, err = json.Marshal(&metric)
	if err != nil {
		log.Println(err)
	}
}

//find results for timestamp field within a specified time range
func GetDataFromTime(req []byte) *[]model.Metric {
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		return nil
	}

	searchMet := MetSearch{}

	var findResults []interface{}

	err = json.Unmarshal(req, &searchMet)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal body - SearchBtwTimeRange")
		log.Println(err)
		return nil
	}

	findResults, err = Db.Collection.Find(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lt": &searchMet.EndDate,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		return nil
	}

	metric := []model.Metric{}

	for _, v := range findResults {
		result := v.(*model.Metric)
		metric = append(metric, *result)
	}

	return &metric
}

func PerDaySenVal(w http.ResponseWriter, r *http.Request) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric := GetDataFromTime(body) //Just need max time

	if metric == nil {
		log.Println(errors.New("Unable to get anything back from SearchBtwTimeRange function - EthPerDay"))
		return
	}

	// for _, v := range *metric {
	// 	soldWeight = v.SoldWeight + soldWeight
	// 	sweight = append(sweight, soldWeight)
	// }

	metSearch := []MetSearch{}
	err = json.Unmarshal(body, &metSearch)
	if err != nil {
		err = errors.Wrap(err, "Unable to Unmarshal timestamp from body - TwSaleWasteDonate")
		log.Println(err)
		return
	}

	var totalResult []byte

	dash := make(map[int]MetDash)
	for i, v := range metSearch {

		dash[i] = MetDash{
			Dates: v.StartDate,
		}

		log.Println(dash)
	}

	totalResult, err = json.Marshal(dash)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		return
	}
	w.Write(totalResult)
}

func SenTable(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		return
	}

	metric := GetDataFromTime(body) //Just need max time

	table := []model.Metric{}
	for _, v := range *metric {
		table = append(table, model.Metric{
			ItemID:       v.ItemID,
			RsCustomerID: v.RsCustomerID,
			DeviceID:     v.DeviceID,
			Ethylene:     v.Ethylene,
			CarbonDi:     v.CarbonDi,
			TempOut:      v.TempOut,
			TempIn:       v.TempIn,
			Humidity:     v.Humidity,
			Timestamp:    v.Timestamp,
			Location:     v.Location,
		})

		log.Println(table)
	}

	totalResult, err := json.Marshal(&table)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		return
	}
	w.Write(totalResult)
}
