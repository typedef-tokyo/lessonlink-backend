package rdb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
)

type (
	IMySQL interface {
		GetConn() *sql.DB
	}
	MySQL struct {
		Conn *sql.DB
	}

	MySQLConnectorConfig struct {
		DBAddress  string
		DBUser     string
		DBPassword string
		DBName     string
		Option     *MySQLConnectorConfigOption
	}
)

// ////////////////////////
func NewConfig(env configs.EnvConfig) *MySQLConnectorConfig {
	return &MySQLConnectorConfig{
		DBAddress:  env.DbAddress,
		DBUser:     env.DbUser,
		DBPassword: env.DbPassword,
		DBName:     env.DbName,
		Option:     NewMySQLConnectorConfigOption(),
	}
}

//////////////////////////

func NewMySQL(config *MySQLConnectorConfig) IMySQL {
	result := &MySQL{}
	err := result.establishConnection(config.DBUser, config.DBPassword, config.DBAddress, config.DBName, config.Option)
	if err != nil {
		log.Print("establishConnection err")
		log.Fatalln(err)
	}
	return result
}

func (m *MySQL) GetConn() *sql.DB {
	return m.Conn
}

func (m *MySQL) establishConnection(user, password, address, name string, option *MySQLConnectorConfigOption) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
		user,
		password,
		address,
		name,
	)

	m.Conn, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	m.Conn.SetMaxOpenConns(option.MaxOpenConnection)
	m.Conn.SetMaxIdleConns(option.MaxIdleConnection)
	m.Conn.SetConnMaxLifetime(*option.ConnectionMaxLifeTime)
	return m.Conn.Ping()
}
