package storages

import (
	"database/sql"
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	_ "github.com/marcboeker/go-duckdb"
	"log"
)

type Database struct {
	db *sql.DB
}

var GlobalDB *Database

var createTableQueries = []string{
	`CREATE TABLE security_info (
		code VARCHAR NOT NULL PRIMARY KEY,
		name VARCHAR NOT NULL,
		market INT NOT NULL,
		daily_limit_rate REAL NOT NULL ,
	);`,
}

func InitDatabase() {
	db := connect()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	storage := newDatabaseStorage(db)
	storage.createAllTables()

	GlobalDB = storage
}

func newDatabaseStorage(db *sql.DB) *Database {
	return &Database{db: db}
}

func (s *Database) createAllTables() {
	for _, q := range createTableQueries {
		_, err := s.db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
}

func (d *Database) insertSecurityInfo(code string, name string, market int, dailyLimitRate float32) {
	_, err := d.db.Exec("INSERT INTO security_info (code, name, market, daily_limit_rate) VALUES (?,?,?,?)", code, name, market, dailyLimitRate)
	if err != nil {
		log.Println(err)
	}
}

func (d *Database) GetStockCodeList(removeST bool) []string {
	queryString := "SELECT security_code FROM security_info"
	if removeST {
		queryString = "SELECT * FROM security_info WHERE name NOT LIKE '%ST%' AND name NOT LIKE '%退%' AND name NOT LIKE '%摘牌%'"
	}
	rows, err := d.db.Query(queryString)
	defer rows.Close()
	if err != nil {
		return nil
	}

	var stockCodes []string
	for rows.Next() {
		var (
			code string
		)
		err := rows.Scan(&code)
		if err != nil {
			return nil
		}
		stockCodes = append(stockCodes, code)
	}

	return stockCodes
}

func (d *Database) UpdateSecurityInfo() {
	addSecurityCode := func(codeBegin, codeEnd int, format string, callback func(string, string, int, float32)) {
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf(format, i)
			securityInfo, ok := securities.CheckoutSecurityInfo(fc)
			if !ok {
				continue
			}
			marketId, _, code := exchange.DetectMarket(fc)
			callback(code, securityInfo.Name, int(marketId), float32(exchange.MarketLimit(fc)))
		}
	}

	// 上海证券交易所
	// 上海
	// sh600000-sh609999
	addSecurityCode(600000, 609999, "sh%d", d.insertSecurityInfo)
	//
	// 科创板
	// sh688000-sh688999
	addSecurityCode(688000, 689999, "sh%d", d.insertSecurityInfo)

	// 深圳证券交易所
	// 深圳主板: sz000000-sz000999
	addSecurityCode(0, 999, "sz000%03d", d.insertSecurityInfo)

	// 中小板
	// 中小板: sz001000-sz009999
	addSecurityCode(1000, 9999, "sz00%04d", d.insertSecurityInfo)

	// 创业板
	// 创业板: sz300000-sz300999
	addSecurityCode(300000, 309999, "sz%06d", d.insertSecurityInfo)
}

func connect() *sql.DB {
	quant1XConfig, found := config.LoadConfig()
	if !found {
		log.Fatal("Failed to load config file")
	}

	db, err := sql.Open("duckdb", quant1XConfig.BaseDir+"/quant1x.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
