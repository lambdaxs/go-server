package mysql_client

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "net/url"
    "time"
)

type MysqlDB struct {
    DSN string //username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local&readTimeout=3s&writeTime=3s
    MaxIdle int //5
    MaxLifeTime time.Duration //60s
    MaxActive int //50
    DialTimeout time.Duration //2s
    ReadTimeout time.Duration //3s
    WriteTimeout time.Duration //3s
    Log bool //false
}

//创建gorm连接
func (g *MysqlDB)ConnectGORMDB() (db *gorm.DB, err error){
    if g.MaxIdle == 0 {
        g.MaxIdle = 5
    }
    if g.MaxLifeTime == 0 {
        g.MaxLifeTime = time.Second*60
    }
    if g.MaxActive == 0 {
        g.MaxActive = 50
    }
    if g.DialTimeout == 0 {
        g.DialTimeout = time.Second*2
    }
    if g.ReadTimeout == 0 {
        g.ReadTimeout = time.Second*3
    }
    if g.WriteTimeout == 0 {
        g.WriteTimeout = time.Second*3
    }
    //update dsn
    g.DSN,err = updateDSNQuery(g.DSN, map[string]string{
        "timeout": g.DialTimeout.String(),
        "readTimeout": g.ReadTimeout.String(),
        "writeTimeout": g.WriteTimeout.String(),
    })
    if err != nil {
        return
    }
    db,err = gorm.Open("mysql", g.DSN)
    if err != nil {
        return
    }
    db.DB().SetMaxIdleConns( g.MaxIdle)
    db.DB().SetConnMaxLifetime(g.MaxLifeTime)
    db.DB().SetMaxOpenConns(g.MaxActive)
    db.LogMode(g.Log)
    return
}

func updateDSNQuery(dsn string, kv map[string]string) (string,error) {
    u,err := url.Parse(dsn)
    if err != nil {
        return "",err
    }
    q := u.Query()
    for k,v := range kv {
        q.Set(k, v)
    }
    u.RawQuery = q.Encode()
    return u.String(),nil
}

