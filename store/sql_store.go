// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	dbsql "database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	sqltrace "log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	l4g "github.com/alecthomas/log4go"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"
)

const (
	INDEX_TYPE_FULL_TEXT = "full_text"
	INDEX_TYPE_DEFAULT   = "default"
	MAX_DB_CONN_LIFETIME = 60
)

const (
	EXIT_CREATE_TABLE                = 100
	EXIT_DB_OPEN                     = 101
	EXIT_PING                        = 102
	EXIT_NO_DRIVER                   = 103
	EXIT_TABLE_EXISTS                = 104
	EXIT_TABLE_EXISTS_MYSQL          = 105
	EXIT_COLUMN_EXISTS               = 106

	EXIT_DOES_COLUMN_EXISTS_MYSQL    = 108
	EXIT_DOES_COLUMN_EXISTS_MISSING  = 109

	EXIT_CREATE_COLUMN_MYSQL         = 111
	EXIT_CREATE_COLUMN_MISSING       = 112
	EXIT_REMOVE_COLUMN               = 113
	EXIT_RENAME_COLUMN               = 114
	EXIT_MAX_COLUMN                  = 115
	EXIT_ALTER_COLUMN                = 116

	EXIT_CREATE_INDEX_MYSQL          = 118
	EXIT_CREATE_INDEX_FULL_MYSQL     = 119
	EXIT_CREATE_INDEX_MISSING        = 120

	EXIT_REMOVE_INDEX_MYSQL          = 122
	EXIT_REMOVE_INDEX_MISSING        = 123
)

type SqlStore struct {
	master         *gorp.DbMap
	replicas       []*gorp.DbMap
	searchReplicas []*gorp.DbMap
	system         SystemStore
	SchemaVersion  string
	rrCounter      int64
	srCounter      int64
}

func initConnection() *SqlStore {
	sqlStore := &SqlStore{
		rrCounter: 0,
		srCounter: 0,
	}

	sqlStore.master = setupConnection("master", utils.Cfg.SqlSettings.DriverName,
		utils.Cfg.SqlSettings.DataSource, utils.Cfg.SqlSettings.MaxIdleConns,
		utils.Cfg.SqlSettings.MaxOpenConns, utils.Cfg.SqlSettings.Trace)

	if len(utils.Cfg.SqlSettings.DataSourceReplicas) == 0 {
		sqlStore.replicas = make([]*gorp.DbMap, 1)
		sqlStore.replicas[0] = sqlStore.master
	} else {
		sqlStore.replicas = make([]*gorp.DbMap, len(utils.Cfg.SqlSettings.DataSourceReplicas))
		for i, replica := range utils.Cfg.SqlSettings.DataSourceReplicas {
			sqlStore.replicas[i] = setupConnection(fmt.Sprintf("replica-%v", i), utils.Cfg.SqlSettings.DriverName, replica,
				utils.Cfg.SqlSettings.MaxIdleConns, utils.Cfg.SqlSettings.MaxOpenConns,
				utils.Cfg.SqlSettings.Trace)
		}
	}

	if len(utils.Cfg.SqlSettings.DataSourceSearchReplicas) == 0 {
		sqlStore.searchReplicas = sqlStore.replicas
	} else {
		sqlStore.searchReplicas = make([]*gorp.DbMap, len(utils.Cfg.SqlSettings.DataSourceSearchReplicas))
		for i, replica := range utils.Cfg.SqlSettings.DataSourceSearchReplicas {
			sqlStore.searchReplicas[i] = setupConnection(fmt.Sprintf("search-replica-%v", i), utils.Cfg.SqlSettings.DriverName, replica,
				utils.Cfg.SqlSettings.MaxIdleConns, utils.Cfg.SqlSettings.MaxOpenConns,
				utils.Cfg.SqlSettings.Trace)
		}
	}

	sqlStore.SchemaVersion = sqlStore.GetCurrentSchemaVersion()
	return sqlStore
}

func NewSqlStore() Store {

	sqlStore := initConnection()




	err := sqlStore.master.CreateTablesIfNotExists()
	if err != nil {
		l4g.Critical(utils.T("store.sql.creating_tables.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_CREATE_TABLE)
	}






	return sqlStore
}

func setupConnection(con_type string, driver string, dataSource string, maxIdle int, maxOpen int, trace bool) *gorp.DbMap {

	db, err := dbsql.Open(driver, dataSource)
	if err != nil {
		l4g.Critical(utils.T("store.sql.open_conn.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_DB_OPEN)
	}

	l4g.Info(utils.T("store.sql.pinging.info"), con_type)
	err = db.Ping()
	if err != nil {
		l4g.Critical(utils.T("store.sql.ping.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_PING)
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)
	db.SetConnMaxLifetime(time.Duration(MAX_DB_CONN_LIFETIME) * time.Minute)

	var dbmap *gorp.DbMap

	if driver == "sqlite3" {
		dbmap = &gorp.DbMap{Db: db, TypeConverter: mattermConverter{}, Dialect: gorp.SqliteDialect{}}
	} else if driver == model.DATABASE_DRIVER_MYSQL {
		dbmap = &gorp.DbMap{Db: db, TypeConverter: mattermConverter{}, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
	}  else {
		l4g.Critical(utils.T("store.sql.dialect_driver.critical"))
		time.Sleep(time.Second)
		os.Exit(EXIT_NO_DRIVER)
	}

	if trace {
		dbmap.TraceOn("", sqltrace.New(os.Stdout, "sql-trace:", sqltrace.Lmicroseconds))
	}

	return dbmap
}

func (ss *SqlStore) TotalMasterDbConnections() int {
	return ss.GetMaster().Db.Stats().OpenConnections
}

func (ss *SqlStore) TotalReadDbConnections() int {

	if len(utils.Cfg.SqlSettings.DataSourceReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.replicas {
		count = count + db.Db.Stats().OpenConnections
	}

	return count
}

func (ss *SqlStore) TotalSearchDbConnections() int {
	if len(utils.Cfg.SqlSettings.DataSourceSearchReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.searchReplicas {
		count = count + db.Db.Stats().OpenConnections
	}

	return count
}

func (ss *SqlStore) GetCurrentSchemaVersion() string {
	version, _ := ss.GetMaster().SelectStr("SELECT Value FROM Systems WHERE Name='Version'")
	return version
}

func (ss *SqlStore) MarkSystemRanUnitTests() {
	if result := <-ss.System().Get(); result.Err == nil {
		props := result.Data.(model.StringMap)
		unitTests := props[model.SYSTEM_RAN_UNIT_TESTS]
		if len(unitTests) == 0 {
			systemTests := &model.System{Name: model.SYSTEM_RAN_UNIT_TESTS, Value: "1"}
			<-ss.System().Save(systemTests)
		}
	}
}

func (ss *SqlStore) DoesTableExist(tableName string) bool {
 if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {

		count, err := ss.GetMaster().SelectInt(
			`SELECT
		    COUNT(0) AS table_exists
			FROM
			    information_schema.TABLES
			WHERE
			    TABLE_SCHEMA = DATABASE()
			        AND TABLE_NAME = ?
		    `,
			tableName,
		)

		if err != nil {
			l4g.Critical(utils.T("store.sql.table_exists.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_TABLE_EXISTS_MYSQL)
		}

		return count > 0

	} else {
		l4g.Critical(utils.T("store.sql.column_exists_missing_driver.critical"))
		time.Sleep(time.Second)
		os.Exit(EXIT_COLUMN_EXISTS)
		return false
	}
}

func (ss *SqlStore) DoesColumnExist(tableName string, columnName string) bool {
 if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {

		count, err := ss.GetMaster().SelectInt(
			`SELECT
		    COUNT(0) AS column_exists
		FROM
		    information_schema.COLUMNS
		WHERE
		    TABLE_SCHEMA = DATABASE()
		        AND TABLE_NAME = ?
		        AND COLUMN_NAME = ?`,
			tableName,
			columnName,
		)

		if err != nil {
			l4g.Critical(utils.T("store.sql.column_exists.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_DOES_COLUMN_EXISTS_MYSQL)
		}

		return count > 0

	} else {
		l4g.Critical(utils.T("store.sql.column_exists_missing_driver.critical"))
		time.Sleep(time.Second)
		os.Exit(EXIT_DOES_COLUMN_EXISTS_MISSING)
		return false
	}
}

func (ss *SqlStore) CreateColumnIfNotExists(tableName string, columnName string, mySqlColType string, postgresColType string, defaultValue string) bool {

	if ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
		_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " ADD " + columnName + " " + mySqlColType + " DEFAULT '" + defaultValue + "'")
		if err != nil {
			l4g.Critical(utils.T("store.sql.create_column.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_CREATE_COLUMN_MYSQL)
		}

		return true

	} else {
		l4g.Critical(utils.T("store.sql.create_column_missing_driver.critical"))
		time.Sleep(time.Second)
		os.Exit(EXIT_CREATE_COLUMN_MISSING)
		return false
	}
}

func (ss *SqlStore) RemoveColumnIfExists(tableName string, columnName string) bool {

	if !ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " DROP COLUMN " + columnName)
	if err != nil {
		l4g.Critical(utils.T("store.sql.drop_column.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_REMOVE_COLUMN)
	}

	return true
}

func (ss *SqlStore) RenameColumnIfExists(tableName string, oldColumnName string, newColumnName string, colType string) bool {
	if !ss.DoesColumnExist(tableName, oldColumnName) {
		return false
	}

	var err error
	if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
		_, err = ss.GetMaster().Exec("ALTER TABLE " + tableName + " CHANGE " + oldColumnName + " " + newColumnName + " " + colType)
	}

	if err != nil {
		l4g.Critical(utils.T("store.sql.rename_column.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_RENAME_COLUMN)
	}

	return true
}

func (ss *SqlStore) GetMaxLengthOfColumnIfExists(tableName string, columnName string) string {
	if !ss.DoesColumnExist(tableName, columnName) {
		return ""
	}

	var result string
	var err error
	if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
		result, err = ss.GetMaster().SelectStr("SELECT CHARACTER_MAXIMUM_LENGTH FROM information_schema.columns WHERE table_name = '" + tableName + "' AND COLUMN_NAME = '" + columnName + "'")
	}

	if err != nil {
		l4g.Critical(utils.T("store.sql.maxlength_column.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_MAX_COLUMN)
	}

	return result
}

func (ss *SqlStore) AlterColumnTypeIfExists(tableName string, columnName string, mySqlColType string, postgresColType string) bool {
	if !ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	var err error
	if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
		_, err = ss.GetMaster().Exec("ALTER TABLE " + tableName + " MODIFY " + columnName + " " + mySqlColType)
	}

	if err != nil {
		l4g.Critical(utils.T("store.sql.alter_column_type.critical"), err)
		time.Sleep(time.Second)
		os.Exit(EXIT_ALTER_COLUMN)
	}

	return true
}

func (ss *SqlStore) CreateUniqueIndexIfNotExists(indexName string, tableName string, columnName string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, columnName, INDEX_TYPE_DEFAULT, true)
}

func (ss *SqlStore) CreateIndexIfNotExists(indexName string, tableName string, columnName string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, columnName, INDEX_TYPE_DEFAULT, false)
}

func (ss *SqlStore) CreateFullTextIndexIfNotExists(indexName string, tableName string, columnName string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, columnName, INDEX_TYPE_FULL_TEXT, false)
}

func (ss *SqlStore) createIndexIfNotExists(indexName string, tableName string, columnName string, indexType string, unique bool) bool {

	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}

	if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {

		count, err := ss.GetMaster().SelectInt("SELECT COUNT(0) AS index_exists FROM information_schema.statistics WHERE TABLE_SCHEMA = DATABASE() and table_name = ? AND index_name = ?", tableName, indexName)
		if err != nil {
			l4g.Critical(utils.T("store.sql.check_index.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_CREATE_INDEX_MYSQL)
		}

		if count > 0 {
			return false
		}

		fullTextIndex := ""
		if indexType == INDEX_TYPE_FULL_TEXT {
			fullTextIndex = " FULLTEXT "
		}

		_, err = ss.GetMaster().Exec("CREATE  " + uniqueStr + fullTextIndex + " INDEX " + indexName + " ON " + tableName + " (" + columnName + ")")
		if err != nil {
			l4g.Critical(utils.T("store.sql.create_index.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_CREATE_INDEX_FULL_MYSQL)
		}
	} else {
		l4g.Critical(utils.T("store.sql.create_index_missing_driver.critical"))
		time.Sleep(time.Second)
		os.Exit(EXIT_CREATE_INDEX_MISSING)
	}

	return true
}

func (ss *SqlStore) RemoveIndexIfExists(indexName string, tableName string) bool {

if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {

		count, err := ss.GetMaster().SelectInt("SELECT COUNT(0) AS index_exists FROM information_schema.statistics WHERE TABLE_SCHEMA = DATABASE() and table_name = ? AND index_name = ?", tableName, indexName)
		if err != nil {
			l4g.Critical(utils.T("store.sql.check_index.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_REMOVE_INDEX_MYSQL)
		}

		if count <= 0 {
			return false
		}

		_, err = ss.GetMaster().Exec("DROP INDEX " + indexName + " ON " + tableName)
		if err != nil {
			l4g.Critical(utils.T("store.sql.remove_index.critical"), err)
			time.Sleep(time.Second)
			os.Exit(EXIT_REMOVE_INDEX_MYSQL)
		}
	} else {
		l4g.Critical(utils.T("store.sql.create_index_missing_driver.critical"))
		time.Sleep(time.Second)
		os.Exit(EXIT_REMOVE_INDEX_MISSING)
	}

	return true
}

func IsUniqueConstraintError(err string, indexName []string) bool {
	unique := strings.Contains(err, "unique constraint") || strings.Contains(err, "Duplicate entry")
	field := false
	for _, contain := range indexName {
		if strings.Contains(err, contain) {
			field = true
			break
		}
	}

	return unique && field
}

func (ss *SqlStore) GetMaster() *gorp.DbMap {
	return ss.master
}

func (ss *SqlStore) GetSearchReplica() *gorp.DbMap {
	rrNum := atomic.AddInt64(&ss.srCounter, 1) % int64(len(ss.searchReplicas))
	return ss.searchReplicas[rrNum]
}

func (ss *SqlStore) GetReplica() *gorp.DbMap {
	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.replicas))
	return ss.replicas[rrNum]
}

func (ss *SqlStore) GetAllConns() []*gorp.DbMap {
	all := make([]*gorp.DbMap, len(ss.replicas)+1)
	copy(all, ss.replicas)
	all[len(ss.replicas)] = ss.master
	return all
}

func (ss *SqlStore) Close() {
	l4g.Info(utils.T("store.sql.closing.info"))
	ss.master.Db.Close()
	for _, replica := range ss.replicas {
		replica.Db.Close()
	}
}




func (ss *SqlStore) System() SystemStore {
	return ss.system
}


func (ss *SqlStore) DropAllTables() {
	ss.master.TruncateTables()
}

type mattermConverter struct{}

func (me mattermConverter) ToDb(val interface{}) (interface{}, error) {

	switch t := val.(type) {
	case model.StringMap:
		return model.MapToJson(t), nil
	case model.StringArray:
		return model.ArrayToJson(t), nil
	case model.EncryptStringMap:
		return encrypt([]byte(utils.Cfg.SqlSettings.AtRestEncryptKey), model.MapToJson(t))
	case model.StringInterface:
		return model.StringInterfaceToJson(t), nil
	}

	return val, nil
}

func (me mattermConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *model.StringMap:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(utils.T("store.sql.convert_string_map"))
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringArray:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(utils.T("store.sql.convert_string_array"))
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.EncryptStringMap:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(utils.T("store.sql.convert_encrypt_string_map"))
			}

			ue, err := decrypt([]byte(utils.Cfg.SqlSettings.AtRestEncryptKey), *s)
			if err != nil {
				return err
			}

			b := []byte(ue)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringInterface:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(utils.T("store.sql.convert_string_interface"))
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	}

	return gorp.CustomScanner{}, false
}



func encrypt(key []byte, text string) (string, error) {

	if text == "" || text == "{}" {
		return "", nil
	}

	plaintext := []byte(text)
	skey := sha512.Sum512(key)
	ekey, akey := skey[:32], skey[32:]

	block, err := aes.NewCipher(ekey)
	if err != nil {
		return "", err
	}

	macfn := hmac.New(sha256.New, akey)
	ciphertext := make([]byte, aes.BlockSize+macfn.Size()+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize+macfn.Size():], plaintext)
	macfn.Write(ciphertext[aes.BlockSize+macfn.Size():])
	mac := macfn.Sum(nil)
	copy(ciphertext[aes.BlockSize:aes.BlockSize+macfn.Size()], mac)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, cryptoText string) (string, error) {

	if cryptoText == "" || cryptoText == "{}" {
		return "{}", nil
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	skey := sha512.Sum512(key)
	ekey, akey := skey[:32], skey[32:]
	macfn := hmac.New(sha256.New, akey)
	if len(ciphertext) < aes.BlockSize+macfn.Size() {
		return "", errors.New(utils.T("store.sql.short_ciphertext"))
	}

	macfn.Write(ciphertext[aes.BlockSize+macfn.Size():])
	expectedMac := macfn.Sum(nil)
	mac := ciphertext[aes.BlockSize : aes.BlockSize+macfn.Size()]
	if hmac.Equal(expectedMac, mac) != true {
		return "", errors.New(utils.T("store.sql.incorrect_mac"))
	}

	block, err := aes.NewCipher(ekey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New(utils.T("store.sql.too_short_ciphertext"))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize+macfn.Size():]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
