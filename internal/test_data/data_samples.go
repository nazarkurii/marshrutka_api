package testdata

import (
	"maryan_api/config"
	"maryan_api/internal/entity"
	"math/rand"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

// SET foreign_key_checks = 0;
// SET @schema = DATABASE();

// SET SESSION group_concat_max_len = 1000000;

// SELECT CONCAT('DROP TABLE ', GROUP_CONCAT(CONCAT('`', table_name, '`'))) INTO @drop_sql
// FROM information_schema.tables
// WHERE table_schema = @schema;

// PREPARE stmt FROM @drop_sql;
// EXECUTE stmt;
// DEALLOCATE PREPARE stmt;
// SET foreign_key_checks = 1;

// nazar@debian:~$ mysql -h marshrutka-marshrutka.i.aivencloud.com   -u avnadmin   -p -D defaultdb  -P 27657   --ssl-ca=/home/nazar/nazar/marshrutka/api/certificates/ca.pem

func CreateTestData(db *gorm.DB) {
	err := db.Create(config.CreateTestData()).Error
	if err != nil {
		panic(err)
	}

	config.LoadCountriesConfig(db)

	drivers := entity.TestDrivers()
	err = db.Create(drivers).Error
	if err != nil {
		panic(err)
	}

	buses := entity.TestBuses()

	for i, _ := range buses {
		buses[i].LeadDriverID = uuid.NullUUID{drivers[i*2].ID, true}
		buses[i].AssistantDriverID = uuid.NullUUID{drivers[(i*2)+1].ID, true}
	}

	err = db.Create(buses).Error
	if err != nil {
		panic(err)
	}

	countries, masterCountry := config.GetConfig()

	var line = 100
	var busIndex int
	for _, country := range countries {

		departureTime := time.Now()
		var trips = make([]entity.Trip, 50)
		for i := 0; i < 50; i++ {

			tripiID := uuid.New()
			outboundConnectionID := uuid.New()
			returnConnectionID := uuid.New()
			trips[i] = entity.Trip{
				ID: tripiID,
				OutboundConnection: entity.Connection{
					ID:                   outboundConnectionID,
					Line:                 line,
					Price:                (rand.Intn(250) + 100) * 100,
					DepartureCountryID:   masterCountry.ID,
					DestinationCountryID: country.ID,
					DepartureTime:        departureTime.UTC(),
					ArrivalTime:          departureTime.Add(time.Hour*15 + time.Hour*time.Duration(rand.Intn(20))).UTC(),
					BusID:                buses[busIndex].ID,
					Updates: []entity.ConnectionUpdate{{
						ConnectionID: outboundConnectionID,
						Status:       entity.RegisteredConnectionStatus,
					}},
					Type:              entity.ComertialConnectionType,
					SellBefore:        departureTime.Add(-time.Hour * 24).UTC(),
					LuggageVolumeLeft: uint(buses[busIndex].LuggageVolume),
					MaxLength:         int(buses[busIndex].MaxLength),
					MaxHeight:         int(buses[busIndex].MaxHeight),
					MaxWidth:          int(buses[busIndex].MaxWidth),
				},
				ReturnConnection: entity.Connection{
					ID:                   returnConnectionID,
					Line:                 line,
					Price:                (rand.Intn(250) + 100) * 100,
					DepartureCountryID:   country.ID,
					DestinationCountryID: masterCountry.ID,
					DepartureTime:        departureTime.Add(time.Hour * 60).UTC(),
					ArrivalTime:          departureTime.Add(time.Hour*60 + time.Hour*15 + time.Hour*time.Duration(rand.Intn(20))).UTC(),
					BusID:                buses[busIndex].ID,
					Updates: []entity.ConnectionUpdate{{
						ConnectionID: returnConnectionID,
						Status:       entity.RegisteredConnectionStatus,
					}},
					Type:              entity.ComertialConnectionType,
					SellBefore:        departureTime.Add(time.Hour * 36).UTC().UTC(),
					LuggageVolumeLeft: uint(buses[busIndex].LuggageVolume),
				},
				Updates: []entity.TripUpdate{{
					TripID: tripiID,
					Status: entity.TripStatusRegistered,
				}},
			}

			departureTime = departureTime.Add(time.Hour * 24 * 10)
		}
		err = db.Create(trips).Error
		if err != nil {
			panic(err)
		}

		line += 100
		busIndex++
		departureTime = time.Now()
	}
}
