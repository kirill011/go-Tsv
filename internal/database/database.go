package database

import (
	"fmt"
	"go-Tsv/internal/config"
	"go-Tsv/internal/database/databaseStruct"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Db *gorm.DB
}

func New(cfg *config.Config) *Database {
	dbLogger := logger.New(
		log.New(os.Stdout, "INFO\t DATABASE\t", log.Ldate|log.Ltime|log.Lmicroseconds),
		logger.Config{
			SlowThreshold:             1 * time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=%s TimeZone=%s", cfg.Host, cfg.User, cfg.Password, cfg.Port, cfg.Sslmode, cfg.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		cfg.ErrorLog.Fatal("func database.New: ", err)
	}

	dbs := fmt.Sprintf("%s dbname=%s", dsn, cfg.Dbname)
	count := 0
	db.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", cfg.Dbname).Scan(&count)
	if count == 0 {
		sql := fmt.Sprintf("CREATE DATABASE %s;", cfg.Dbname)
		db.Exec(sql)
	}

	db, err = gorm.Open(postgres.Open(dbs), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		cfg.ErrorLog.Fatal("func database.New: ", err)
	}

	return &Database{db}
}

func (base *Database) Migrate(cfg *config.Config) {
	err := base.Db.AutoMigrate(&databaseStruct.LoadFiles{}, &databaseStruct.Messages{}, &databaseStruct.Errors{})
	if err != nil {
		cfg.ErrorLog.Fatal("func database.Migrate: ", err)
	}

	cfg.InfoLog.Println("Migration success")
}

func (base *Database) SelectReadFiles() ([]string, error) {
	files := []databaseStruct.LoadFiles{}
	result := base.Db.Model(&databaseStruct.LoadFiles{}).Where("file_saved = true").Scan(&files)
	if result.Error != nil {
		return nil, result.Error
	}

	filesName := []string{}

	for _, val := range files {
		filesName = append(filesName, val.Filename)
	}
	return filesName, nil
}

func (base *Database) InsertMessages(file *databaseStruct.LoadFiles) error {

	result := base.Db.Create(file)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

type GuidFileid struct {
	Unit_guid   string
	LoadFilesId int
}

func (base *Database) SelectGuid(file string) ([]GuidFileid, error) {
	res := []GuidFileid{}

	result := base.Db.Model(&databaseStruct.LoadFiles{}).Distinct("messages.unit_guid, messages.load_files_id").Joins("Inner join messages on messages.load_files_id = load_files.id").Where("load_files.filename = ?", file).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}

func (base *Database) SelectMessages(fileId int, guid string) ([]databaseStruct.Messages, error) {
	mes := []databaseStruct.Messages{}

	result := base.Db.Model(&databaseStruct.LoadFiles{}).Select("messages.n, messages.mqtt, messages.invid, messages.unit_guid, messages.msg_id, messages.text, messages.context, messages.class, messages.level, messages.area, messages.addr, messages.block, messages.type, messages.bit, messages.invert_bit").Joins("Inner join messages on messages.load_files_id = load_files.id").Where("load_files_id = ? and unit_guid = ?", fileId, guid).Scan(&mes)
	if result.Error != nil {
		return nil, result.Error
	}

	return mes, nil
}

func (base *Database) InsertError(dbError databaseStruct.Errors) error {
	result := base.Db.Model(&databaseStruct.Errors{}).Create(&dbError)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (base *Database) GetData(page int, pageSize int) ([]databaseStruct.Messages, error) {

	var messages []databaseStruct.Messages

	offset := 0

	if pageSize != -1 && page != -1 {
		offset = pageSize * (page - 1)
	}
	limit := pageSize

	//если conditions содержит нулевые поля ("" для string, 0 для int), то такие поля не будут использоваться в Where
	//если limit = -1, то он будет игнорироваться
	result := base.Db.Model(&databaseStruct.Messages{}).Offset(offset).Limit(limit).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}
