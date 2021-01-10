package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/nmalensek/go-user-form/fileusermodel"
	"github.com/nmalensek/go-user-form/model"
)

const (
	connFlag = "conn"
	fileDb   = "file"
)

var connString = flag.String(connFlag, "", "The database connection string (absolute file path if using a file as a database).")
var dbType = flag.String("db", "", fmt.Sprintf("The type of database to use, options follow:\n %v", printDatabaseTypes()))

var validPath = regexp.MustCompile("^/(users)/([a-zA-Z0-9]*)$")

type flagOption struct {
	Name        string
	Description string
}

func (f *flagOption) String() string {
	return fmt.Sprintf("%v:\t%v", f.Name, f.Description)
}

var databaseTypes = []flagOption{
	{Name: fileDb, Description: fmt.Sprintf("Use a JSON file as a pseudo-database (provide the absolute filepath as the \"%v\" flag).", connFlag)},
}

func printDatabaseTypes() string {
	b := strings.Builder{}
	for _, v := range databaseTypes {
		fmt.Fprint(&b, v.String())
	}
	return b.String()
}

//Env contains all environment variables that the app needs to run (database info, loggers, etc.)
type Env struct {
	Datastore model.UserDataStore
	ErrorLog  *log.Logger
}

//Start initializes all environment dependencies for use in the application.
func Start() (*Env, error) {
	db, err := initDb()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	env := Env{Datastore: db}

	fileLog, err := initLogger()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	env.ErrorLog = fileLog

	return &env, nil
}

//initDb constructs the database connection depending on the type specified in the command line.
func initDb() (model.UserDataStore, error) {
	switch *dbType {
	case fileDb:
		return registerFileDb()
	case "":
		return nil, errors.New("initDb: database type not specified")
	default:
		return nil, errors.New("initDb: unable to parse specified database type")
	}
}

func initLogger() (*log.Logger, error) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile), nil
}

//registerFileDb determines the filepath permissions given in the userFilePath argument, and if the file has the correct permissions the path is stored for future "database" uses.
func registerFileDb() (model.UserDataStore, error) {
	if connString == nil || *connString == "" {
		return nil, fmt.Errorf("registerFileDb: using a file as a database but no file path was provided through the %v flag", connFlag)
	}
	//TODO: ensure read/write permissions on file. should this be logged?
	return &fileusermodel.FileUserModel{Filepath: *connString}, nil
}
