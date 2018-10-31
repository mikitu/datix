package db

import (
	"bufio"
	"encoding/csv"
	"github.com/hashicorp/go-memdb"
	"github.com/kataras/iris/core/errors"
	"github.com/mikitu/datix/model"
	"github.com/mikitu/datix/util"
	"io"
	"os"
	"strconv"
)

type Db struct {
	db *memdb.MemDB
}

func New() *Db {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": &memdb.TableSchema{
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	util.FailOnError(err, "Cannot create database")
	return &Db{db:db}
}

func (d *Db) Import() {
	// Create a write transaction
	txn := d.db.Txn(true)
	csvFile, err := os.Open("./import/user_data.csv")
	util.FailOnError(err, "Cannot open csv file")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			util.FailOnError(err, "Cannot read csv line")
		}
		// Insert a new person
		id, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			// go to next line if first column is not an Id
			continue
		}
		u := &model.User{
			Id:        id,
			FirstName: line[1],
			LastName:  line[2],
			Email:     line[3],
			Gender:    line[4],
			IpAddress: line[5],
		}
		err = txn.Insert("user", u)
		util.FailOnError(err, "Cannot insert record")
	}
	// Commit the transaction
	txn.Commit()
}

func (d Db) FindById(id uint64) (model.User, error) {
	txn := d.db.Txn(false)
	defer txn.Abort()
	user := &model.User{}
	raw, err := txn.First("user", "id", id)
	if raw == nil {
		return *user, errors.New("User not found")
	}
	user = raw.(*model.User)
	return *user, err
}

func (d Db) FindAll() ([]model.User, error) {
	txn := d.db.Txn(false)
	defer txn.Abort()
	users := model.Users{}
	raw, err := txn.Get("user", "id")
	for user := raw.Next(); user != nil; user = raw.Next() {
		users = append(users, *user.(*model.User))
	}
	return users, err
}

