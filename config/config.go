package config

import (
	"errors"
	"flag"
	"fmt"
	"regexp"

	"github.com/nmalensek/go-user-form/model"
)

var userFilePath = flag.String("ufile", "", "The absolute path for the file to use as a pseudo-database")
var dbType = flag.String("db", "", "The type of database to use. If db=file, the ufile parameter must also be specified")

var validPath = regexp.MustCompile("^/(users)/([a-zA-Z0-9]*)$")

//Env contains all environment variables that the app needs to run (database info, loggers, etc.)
type Env struct {
	Datastore model.UserDataStore
}

//Start initializes all environment dependencies for use in the application.
func Start() (*Env, error) {
	db, err := initDb()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	env := Env{Datastore: db}

	return &env, nil
}

//initDb constructs the database connection depending on the type specified in the command line.
func initDb() (model.UserDataStore, error) {
	switch *dbType {
	case "file":
		return registerFileDb()
	case "":
		return nil, errors.New("initDb: database type not specified")
	default:
		return nil, errors.New("initDb: unable to parse specified database type")
	}
}

//registerFileDb determines the filepath permissions given in the userFilePath argument, and if
//the file has the correct permissions the path is stored for future "database" uses.
func registerFileDb() (model.UserDataStore, error) {
	if userFilePath == nil || *userFilePath == "" {
		return nil, errors.New("registerFileDb: using a file as a database but no file path was provided")
	}
	//TODO: ensure read/write permissions on file. should this be logged?
	return &model.FileUserModel{Filepath: userFilePath}, nil
}
