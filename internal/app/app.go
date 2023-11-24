package app

import (
	"fmt"
	"go-Tsv/internal/api/endpoint"
	"go-Tsv/internal/config"
	"go-Tsv/internal/database"
	"go-Tsv/internal/database/databaseStruct"
	filereaderwriter "go-Tsv/internal/pkg/fileReaderWriter"
	"os"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
)

type App struct {
	db  *database.Database
	cfg *config.Config
	e   *endpoint.Endpoint

	echo *echo.Echo
}

func New(cfg *config.Config) *App {
	app := &App{
		cfg:  cfg,
		echo: echo.New(),
	}

	app.db = database.New(cfg)

	app.e = endpoint.New(cfg, app.db)

	app.echo.GET("/getData", app.e.HandlerGetData)

	//Создаём таблицы (если не созданы)
	app.db.Migrate(cfg)
	return app
}

func (a *App) Run() {
	a.cfg.InfoLog.Println("Server running")

	err := a.echo.Start(fmt.Sprintf(":%d", a.cfg.ApiPort))
	if err != nil {
		a.echo.Logger.Fatal(err)
	}

}

// Функция ищет не загруженный файл и вызывает функцию загрузки файла в postgres
func (a *App) MonitFiles() {

	reader := filereaderwriter.New()
	for {
		//Выбираем уже загруженные файлы
		res, err := a.db.SelectReadFiles()
		if err != nil {
			a.cfg.ErrorLog.Println("func app.MonitFiles.SelectReadFiles: ", err)
		}
		//Смотрим файлы в дирректории
		files, err := os.ReadDir(a.cfg.DirInput)
		if err != nil {
			a.cfg.ErrorLog.Println("func app.MonitFiles.ReadDir: ", err)
		}
		for _, file := range files {
			if slices.Index(res, file.Name()) == -1 {

				a.cfg.InfoLog.Println("Load file: ", file.Name())

				//Читаем файл
				err := reader.ReadFile(file.Name(), a.db, a.cfg)
				if err != nil {
					a.cfg.ErrorLog.Println("func app.MonitFiles.ReadFile: ", err)

					//Если ошибка то записываем в файл и в базу
					f, createErr := os.Create(fmt.Sprintf("%s/%s.txt", a.cfg.DirOutput, file.Name()))
					if createErr != nil {
						a.cfg.ErrorLog.Println(createErr)
					}

					parseError := databaseStruct.Errors{
						ErrorMessage: err.Error(),
						Filename:     file.Name(),
					}

					a.db.InsertError(parseError)

					f.Write([]byte(err.Error()))
					continue
				}

				//записываем ответ
				err = reader.CreateFile(file.Name(), a.db, a.cfg)
				if err != nil {
					a.cfg.ErrorLog.Println("func app.MonitFiles.CreateFile: ", err)

					//Если ошибка то записываем в файл и в базу
					f, createErr := os.Create(fmt.Sprintf("%s/%s.txt", a.cfg.DirOutput, file.Name()))
					if createErr != nil {
						a.cfg.ErrorLog.Println(createErr)
					}

					parseError := databaseStruct.Errors{
						ErrorMessage: err.Error(),
						Filename:     file.Name(),
					}

					a.db.InsertError(parseError)

					f.Write([]byte(err.Error()))
				}
			}
		}
		time.Sleep(a.cfg.TimeSleep)
	}
}
