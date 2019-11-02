package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vitaly-kashtalyan/go-smarthome/db"
	"github.com/vitaly-kashtalyan/go-smarthome/models"
	"github.com/vitaly-kashtalyan/go-smarthome/responses"
	"github.com/vitaly-kashtalyan/go-smarthome/routers"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load("app.config"); err != nil {
		log.Print("No app.config file found")
	}
	// Connects to the database and applying migrations
	db.GetDB().Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&models.Sensors{}, &models.SensorsHistory{}, &models.RelayStateHistory{})
}
func main() {
	scheduledJobs()

	// Runs the server
	r := routers.Setup()
	_ = r.Run(fmt.Sprintf(":%v", os.Getenv("PORT_ENV")))
}

func scheduledJobs() {
	nextTime := time.Now().Truncate(time.Minute)
	nextTime = nextTime.Add(time.Minute)
	time.Sleep(time.Until(nextTime))
	readDataSensors()
	parseRelayStatus()
	updateRelayState()
	go scheduledJobs()
}

func updateRelayState() {
	fmt.Println("run -> updateRelayState")
}

func parseRelayStatus() {
	relayStatus := responses.RelayStatus{}
	if err := getJSON("http://"+os.Getenv("HOST_RELAYS"), &relayStatus); err == nil {
		if relayStatus.Status == http.StatusOK {

			for _, relay := range relayStatus.Data {
				rsh := models.RelayStateHistory{}
				db.GetDB().Where(models.RelayStateHistory{RelayId: sql.NullInt32{Int32: relay.Id, Valid: true}}).
					Order("created_at desc").
					Limit(1).Find(&rsh)

				if rsh.ID == 0 || rsh.ID > 0 && rsh.State.Int32 != relay.State {
					var newRecord = models.RelayStateHistory{
						RelayId:   sql.NullInt32{Int32: relay.Id, Valid: true},
						State:     sql.NullInt32{Int32: relay.State, Valid: true},
						UpdatedAt: sql.NullTime{Time: time.Now(), Valid: false},
					}
					if err := db.GetDB().Create(&newRecord).Error; err != nil {
						log.Println("error creating relay history record: ", err)
					}

					if rsh.ID > 0 {
						db.GetDB().First(&models.RelayStateHistory{}, rsh.ID).UpdateColumns(models.RelayStateHistory{
							UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true}})
					}
				}
			}
		}
	} else {
		log.Println("error getting json object: ", err)
	}
}

func readDataSensors() {
	response := responses.ArduinoSensors{}
	if err := getJSON("http://"+os.Getenv("HOST_SENSORS"), &response); err == nil {
		tx := db.GetDB().Begin()
		for _, v := range response.Dht22 {
			if v.Status == http.StatusText(http.StatusOK) {
				if err := tx.Where(models.Sensors{Pin: v.Pin}).
					Assign(models.Sensors{Pin: v.Pin, Temperature: v.Temperature, Humidity: v.Humidity, UpdatedAt: time.Now()}).
					FirstOrCreate(&models.Sensors{}).Error; err != nil {
					tx.Rollback()
					log.Println(err)
				}
			}
		}
		for _, v := range response.Ds18b20 {
			if v.Status == http.StatusText(http.StatusOK) {
				if err := tx.Where(models.Sensors{Pin: v.Pin, DecSensor: v.Dec}).
					Assign(models.Sensors{Pin: v.Pin, Temperature: v.Temperature, DecSensor: v.Dec, UpdatedAt: time.Now()}).
					FirstOrCreate(&models.Sensors{}).Error; err != nil {
					tx.Rollback()
					log.Println(err)
				}
			}
		}
		tx.Commit()
	} else {
		log.Println("error getting json object: ", err)
	}
}

func getJSON(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot fetch URL %q: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http GET status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %v", err)
	}
	return nil
}
