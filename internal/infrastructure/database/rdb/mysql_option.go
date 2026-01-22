package rdb

import "time"

type MySQLConnectorConfigOption struct {
	MaxIdleConnection     int
	MaxOpenConnection     int
	ConnectionMaxLifeTime *time.Duration
}

const (
	_defaultMaxIdleConnection     = 10
	_defaultMaxOpenConnection     = 20
	_defaultConnectionMaxLifeTime = 10 * time.Second
)

func NewMySQLConnectorConfigOption() *MySQLConnectorConfigOption {
	maxIdleConnection := _defaultMaxIdleConnection
	connectionMaxLifeTime := _defaultConnectionMaxLifeTime
	maxOpenConnection := _defaultMaxOpenConnection

	return &MySQLConnectorConfigOption{MaxIdleConnection: maxIdleConnection, MaxOpenConnection: maxOpenConnection, ConnectionMaxLifeTime: &connectionMaxLifeTime}
}
