package scheduler

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// OpenDatabase проверяет, существует ли база данных.
// Если да - открывает, если нет - направляет к функции создания бд
func OpenDatabase() (*sql.DB, error) {

	// Получение значения переменной окружения TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")

	// Если переменная окружения пустая, используется стандартный путь
	if dbFile == "" {
		// os.Executable возвращает путь к исполняемому файлу программы
		appPath, err := os.Executable()
		if err != nil {
			// log.Fatal(err)
			fmt.Println("Fatal error:", err)
			os.Exit(1)
		}
		// filepath.Join объединяет несколько путей в один.
		// filepath.Dir возвращает путь к директории, в которой находится файл.
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}
	install := false
	// os.Stat получает информацию о файле.
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить функцию createDatabase
	if install {
		if err := createDatabase(dbFile); err != nil {
			return nil, err
		}
	}
	// Открытие базы данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	// Возвращается объект базы данных и nil как ошибка
	return db, nil
}

// createDatabase - функция, создающая базу данных
func createDatabase(dbFile string) error {
	if _, err := os.Create(dbFile); err != nil {
		return err
	}
	// Проверка существования sql файла
	if _, err := os.Stat("scheduler.sql"); os.IsNotExist(err) {
		// return log.Output(2, "The scheduler.sql file was not found")
		fmt.Printf("The scheduler.sql file was not found\n")
		return err
	}
	// Чтение sql файла
	n, err := os.ReadFile("scheduler.sql")
	if err != nil {
		// return log.Output(2, "Couldn't read the scheduler.sql file: "+err.Error())
		fmt.Printf("Couldn't read the scheduler.sql file: %v\n", err)
		return err
	}
	// Открытие базы данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		// return log.Output(2, "Couldn't connect to the database: "+err.Error())
		fmt.Printf("Couldn't connect to the database: %v\n", err)
		return err
	}
	// Выполнение sql запроса
	m := string(n)
	if _, err = db.Exec(m); err != nil {
		// return log.Output(2, "Sql query execution error: "+err.Error())
		fmt.Printf("Sql query execution error: %v\n", err)
		return err
	}
	return nil
}
