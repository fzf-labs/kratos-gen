package utils

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MySQL    = "mysql"
	Postgres = "postgres"
)

// NewDB 创建数据库连接
func NewDB(db, dsn string) (*gorm.DB, error) {
	switch db {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	case "postgres":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	default:
		return nil, errors.New(" db not support")
	}
}

// GetPartitionTableToChildTables 获取分区表到子表的映射
func GetPartitionTableToChildTables(db *gorm.DB) (resp map[string][]string, err error) {
	switch db.Dialector.Name() {
	case Postgres:
		resp, err = getPGPartitionTableToChildTables(db)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case MySQL:
		resp, err = getMySQLPartitionTableToChildTables(db)
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, nil
	}
}

// getPartitionTableToChildTable 获取PG分区表到子表的映射
func getPGPartitionTableToChildTables(db *gorm.DB) (map[string][]string, error) {
	resp := make(map[string][]string)
	type tmp struct {
		PartitionedTable string `gorm:"column:partitioned_table" json:"partitioned_table"`
		ChildTables      string `gorm:"column:child_tables" json:"child_tables"`
	}
	result := make([]tmp, 0)
	sql := `SELECT p.relname AS partitioned_table,array_to_string(array_agg(c.relname),',')AS child_tables FROM pg_catalog.pg_class c JOIN pg_catalog.pg_inherits i ON c.oid=i.inhrelid JOIN pg_catalog.pg_class p ON p.oid=i.inhparent GROUP BY p.relname;`
	err := db.Raw(sql).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		resp[v.PartitionedTable] = append(resp[v.PartitionedTable], strings.Split(v.ChildTables, ",")...)
	}
	return resp, nil
}

// getMySQLPartitionTableToChildTable 获取MySQL分区表到子表的映射
func getMySQLPartitionTableToChildTables(db *gorm.DB) (map[string][]string, error) {
	resp := make(map[string][]string)
	type tmp struct {
		TableName                string `gorm:"column:TABLE_NAME"`
		PartitionMethod          string `gorm:"column:PARTITION_METHOD"`
		SubPartitionMethod       string `gorm:"column:SUBPARTITION_METHOD"`
		PartitionOrdinalPosition string `gorm:"column:PARTITION_ORDINAL_POSITION"`
	}
	result := make([]tmp, 0)
	sql := fmt.Sprintf(`SELECT TABLE_NAME,PARTITION_METHOD,SUBPARTITION_METHOD,PARTITION_ORDINAL_POSITION FROM INFORMATION_SCHEMA.PARTITIONS WHERE PARTITION_NAME IS NOT NULL AND TABLE_SCHEMA='%s' ORDER BY TABLE_NAME`, db.Name())
	err := db.Raw(sql).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	// mysql 特殊处理成这种形式
	for _, v := range result {
		resp[v.TableName] = []string{
			v.TableName,
		}
	}
	return resp, nil
}
