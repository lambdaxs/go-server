package psql_client

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "strconv"
    "time"
)

type PsqlConfig struct {
    DSN            string        //"host=myhost port=5432 user=gorm password=mypassword dbname=gorm sslmode=disable"
    MaxIdle        int           //5
    MaxLifeTime    time.Duration //60s
    MaxActive      int           //50
    ConnectTimeout int           //2s
    Log            bool          //false
}

func (p *PsqlConfig) Connect() (db *gorm.DB, err error) {
    if p.MaxIdle == 0 {
        p.MaxIdle = 5
    }
    if p.MaxLifeTime == 0 {
        p.MaxLifeTime = time.Second * 60
    }
    if p.MaxActive == 0 {
        p.MaxActive = 50
    }

    params := map[string]string{}
    if p.ConnectTimeout == 0 {
        p.ConnectTimeout = 2
        params["connect_timeout"] = strconv.Itoa(p.ConnectTimeout)
    }

    //update dsn
    p.DSN = updateDSNQuery(p.DSN, params)

    db, err = gorm.Open("postgres", p.DSN)
    if err != nil {
        return
    }
    db.DB().SetMaxIdleConns(p.MaxIdle)
    db.DB().SetConnMaxLifetime(p.MaxLifeTime)
    db.DB().SetMaxOpenConns(p.MaxActive)
    db.LogMode(p.Log)
    return
}

func updateDSNQuery(dsn string, kv map[string]string) (string) {
    for k, v := range kv {
        dsn = dsn + fmt.Sprintf(" %s=%s", k, v)
    }
    return dsn
}
