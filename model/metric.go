package model

import (
	"encoding/json"
	"log"

	"github.com/TerrexTech/uuuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/pkg/errors"
)

type Metric struct {
	ID           objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ItemID       uuuid.UUID        `bson:"item_id,omitempty" json:"item_id,omitempty"`
	RsCustomerID uuuid.UUID        `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	DeviceID     uuuid.UUID        `bson:"device_id,omitempty" json:"device_id,omitempty"`
	Ethylene     float64           `bson:"ethylene,omitempty" json:"ethylene,omitempty"`
	CarbonDi     float64           `bson:"carbon_di,omitempty" json:"carbon_di,omitempty"`
	TempOut      float64           `bson:"temp_out,omitempty" json:"temp_out,omitempty"`
	TempIn       float64           `bson:"temp_in,omitempty" json:"temp_in,omitempty"`
	Humidity     float64           `bson:"humidity,omitempty" json:"humidity,omitempty"`
	Timestamp    int64             `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Location     string            `bson:"location,omitempty" json:"location,omitempty"`
}

type MarshalMetric struct {
	ID           objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ItemID       string            `bson:"item_id,omitempty" json:"item_id,omitempty"`
	RsCustomerID string            `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	DeviceID     string            `bson:"device_id,omitempty" json:"device_id,omitempty"`
	Ethylene     float64           `bson:"ethylene,omitempty" json:"ethylene,omitempty"`
	CarbonDi     float64           `bson:"carbon_di,omitempty" json:"carbon_di,omitempty"`
	TempOut      float64           `bson:"temp_out,omitempty" json:"temp_out,omitempty"`
	TempIn       float64           `bson:"temp_in,omitempty" json:"temp_in,omitempty"`
	Humidity     float64           `bson:"humidity,omitempty" json:"humidity,omitempty"`
	Timestamp    int64             `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Location     string            `bson:"location,omitempty" json:"location,omitempty"`
}

func (m *Metric) MarshalJSON() ([]byte, error) {
	mon := &MarshalMetric{
		Ethylene:  m.Ethylene,
		CarbonDi:  m.CarbonDi,
		TempOut:   m.TempOut,
		TempIn:    m.TempIn,
		Humidity:  m.Humidity,
		Timestamp: m.Timestamp,
		Location:  m.Location,
	}

	if m.ItemID.String() != (uuuid.UUID{}).String() {
		mon.ItemID = m.ItemID.String()
	}
	if m.DeviceID.String() != (uuuid.UUID{}).String() {
		mon.DeviceID = m.DeviceID.String()
	}
	if m.RsCustomerID.String() != (uuuid.UUID{}).String() {
		mon.RsCustomerID = m.RsCustomerID.String()
	}

	e, _ := json.Marshal(mon)
	log.Println(string(e))

	return json.Marshal(mon)
}

func (m *Metric) UnmarshalJSON(in []byte) error {
	mars := make(map[string]interface{})
	err := bson.Unmarshal(in, mars)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	m.ID = mars["_id"].(objectid.ObjectID)

	m.ItemID, err = uuuid.FromString(mars["item_id"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing ItemID for inventory")
		return err
	}

	m.DeviceID, err = uuuid.FromString(mars["device_id"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing DeviceID for inventory")
		return err
	}

	m.RsCustomerID, err = uuuid.FromString(mars["rs_customer_id"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing DeviceID for inventory")
		return err
	}

	if mars["location"] != nil {
		m.Location = mars["location"].(string)
	}

	if mars["timestamp"] != nil {
		m.Timestamp = mars["timestamp"].(int64)

	}

	if mars["ethylene"] != nil {
		m.Ethylene = mars["ethylene"].(float64)
	}

	return nil
}
