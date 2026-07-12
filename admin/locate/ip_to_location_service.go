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

	"jinya.de/goji/database"
)

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func PopulateDatabase() error {
	databaseUrl := "https://download.db-ip.com/free/dbip-city-lite-{YEAR}-{MONTH}.csv.gz"
	conf, err := database.FindJinyaConfiguration("ip_database_url", "jinya")
	if err == nil && conf != nil {
		databaseUrl = conf.Value
	}

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
	_, err = dbMap.Exec("DELETE FROM ip_address")
	if err != nil {
		return err
	}

	batchSize := 1000
	counter := 0

	trans, err := dbMap.Begin()
	if err != nil {
		return err
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			trans.Rollback()
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
		ipTo := net.ParseIP(ipToStr)

		var addressType string
		var ipFromBytes, ipToBytes []byte

		if ipFrom.To4() != nil {
			addressType = "ipv4"
			ipFromBytes = ipFrom.To4()
			ipToBytes = ipTo.To4()
		} else {
			addressType = "ipv6"
			ipFromBytes = ipFrom.To16()
			ipToBytes = ipTo.To16()
		}

		addr := database.IpAddress{
			AddressType: addressType,
			IpStart:     ipFromBytes,
			IpEnd:       ipToBytes,
			Country:     countryCode,
			City:        city,
		}

		err = trans.Insert(&addr)
		if err != nil {
			trans.Rollback()
			return err
		}

		counter++
		if counter%batchSize == 0 {
			err = trans.Commit()
			if err != nil {
				return err
			}
			trans, err = dbMap.Begin()
			if err != nil {
				return err
			}
		}
	}

	return trans.Commit()
}

func LocateIp(ipStr string) (Location, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return Location{Country: "-1", City: "-1"}, nil
	}

	var addressType string
	var ipBytes []byte
	if ip.To4() != nil {
		addressType = "ipv4"
		ipBytes = ip.To4()
	} else {
		addressType = "ipv6"
		ipBytes = ip.To16()
	}

	addr, err := database.SelectOne[database.IpAddress](`
SELECT country, city
FROM ip_address
WHERE address_type = ? AND ip_start <= ? AND ? <= ip_end
LIMIT 1
`, addressType, ipBytes, ipBytes)

	if err != nil {
		return Location{Country: "-1", City: "-1"}, nil
	}

	return Location{
		Country: addr.Country,
		City:    addr.City,
	}, nil
}
