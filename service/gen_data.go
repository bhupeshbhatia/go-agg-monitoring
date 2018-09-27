package service

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/TerrexTech/uuuid"
	"github.com/pkg/errors"

	"github.com/bhupeshbhatia/go-agg-monitoring/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ModifyMetricData struct {
	Mon       model.Metric
	Ethylene  float64
	CarbonDi  float64
	TempIn    float64
	TempOut   float64
	Humidity  float64
	Timestamp int64
	Randnum   int
}

func random(min, max float64) float64 {
	return float64(rand.Intn(int(max)-int(min)) + int(min))
}

func generateRandomValue(num1, num2 float64) float64 {
	// rand.Seed(time.Now().Unix())
	return random(num1, num2)
}

func generateNewUUID() uuuid.UUID {
	uuid, err := uuuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "Unable to generate UUID")
		log.Println(err)
	}
	return uuid
}

var locationName = []string{"A101", "B201", "O301", "M401", "S501", "T601", "L701", "P801", "G901", "SW1001"}

func GenSensorData() model.Metric {
	randLocation := generateRandomValue(1, 10)
	randDateArr := generateRandomValue(1, 7) //in hours

	metric := model.Metric{
		ItemID:       generateNewUUID(),
		RsCustomerID: generateNewUUID(),
		DeviceID:     generateNewUUID(),
		Ethylene:     generateRandomValue(4, 100),
		CarbonDi:     generateRandomValue(400, 1200),
		TempOut:      generateRandomValue(-10, 35),
		TempIn:       generateRandomValue(22, 27),
		Humidity:     generateRandomValue(30, 60),
		Timestamp:    time.Now().Add(time.Duration(randDateArr) * time.Hour).Unix(),
		Location:     locationName[int(randLocation)],
	}

	return metric
}

func TestIfDataGenerated() {
	metric := []model.Metric{}
	for i := 0; i < 100; i++ {
		metric = append(metric, GenSensorData())
	}

	jSensorData, err := json.Marshal(&metric)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(jSensorData))
}
