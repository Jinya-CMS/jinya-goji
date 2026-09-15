package locate

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"jinya.de/goji/configuration"
	"jinya.de/goji/database"
)

type Location struct {
	City    string `json:"city" db:"city"`
	Country string `json:"country" db:"country"`
}

func PopulateDatabase() error {
	databaseUrl := configuration.GetString("jinya", "ip_database_url", "https://download.db-ip.com/free/dbip-city-lite-{YEAR}-{MONTH}.csv.gz")

	yesterday := time.Now().AddDate(0, 0, -1)
	databaseUrl = strings.ReplaceAll(databaseUrl, "{YEAR}", yesterday.Format("2006"))
	databaseUrl = strings.ReplaceAll(databaseUrl, "{MONTH}", yesterday.Format("01"))

	resp, err := http.Get(databaseUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download database: %s", resp.Status)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	csvReader := csv.NewReader(gzReader)
	csvReader.ReuseRecord = true

	dbMap := database.GetDbMap()
	_, err = dbMap.Exec(`
delete
from ip_address`)
	if err != nil {
		return err
	}

	batchSize := 1000
	counter := 0

	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(record) < 6 {
			continue
		}

		ipFromStr := record[0]
		ipToStr := record[1]
		countryCode := record[3]
		city := record[5]

		ipFrom := net.ParseIP(ipFromStr)

		var addressType string

		if ipFrom.To4() != nil {
			addressType = "ipv4"
		} else {
			addressType = "ipv6"
		}

		_, err = tx.Exec(`
insert into ip_address (address_type, ip_start, ip_end, country, city)
values ($1, $2, $3, $4, $5)`, addressType, ipFromStr, ipToStr, countryCode, city)
		if err != nil {
			return err
		}

		counter++
		if counter%batchSize == 0 {
			err = tx.Commit()
			if err != nil {
				return err
			}
			tx, err = dbMap.Begin()
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func LocalizeIp(ipStr string) (Location, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return Location{Country: "-1", City: "-1"}, nil
	}

	addr, err := database.GetDbMap().SelectOneType[Location](`
select country, city
from ip_address
where ip_start <= cast($1 as inet)
  and ip_end >= cast($1 as inet)
limit 1
`, ipStr)

	if err != nil {
		return Location{Country: "ZZ", City: "ZZ"}, nil
	}

	return Location{
		Country: addr.Country,
		City:    addr.City,
	}, nil
}
