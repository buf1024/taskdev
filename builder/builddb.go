package main

import (
	"database/sql"
	"sync"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	buildStateStop         = 0
	buildStateBuilding     = 1
	buildStateBuildFailed  = 2
	buildStateBuildSuccess = 3
)

type buildDb struct {
	db   *sql.DB
	lock sync.Locker
}

type buildState struct {
	groupname string
	name      string
	state     int
	outbin    string
	log       string
}

var sqlStr string = ""

//create table task (
//     id            integer primary key autoincrement,
//     groupname     varchar(64),
//     name          varchar(64),
//     state         integer,
//     outbin        varchar(256),
//     log           varchar(256)
// );
// `

func NewBuildDb(file string) (*buildDb, error) {
	b := &buildDb{}
	var err error
	b.db, err = sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	_, err = b.db.Query("select count(1) from task")
	if err == nil {
		stm, err := b.db.Prepare(sqlStr)
		if err != nil {
			b.db.Close()
			return nil, err
		}
		_, err = stm.Exec()
		if err != nil {
			b.db.Close()
			return nil, err
		}
		stm.Close()
	}
	b.lock = &sync.Mutex{}
	return b, nil

}

func (b *buildDb) insert(st *buildState) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	stm, err := b.db.Prepare("insert into task(groupname, name, state, outbin, log) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	rst, err := stm.Exec(st.groupname, st.name, st.state, st.outbin, st.log)
	if err != nil {
		return err
	}
	id, err := rst.LastInsertId()
	if err != nil {
		return err
	}
	stm.Close()

	fmt.Printf("last insertid = %d\n", id)

	return nil
}

func (b *buildDb) get(groupname string, name string) (*buildState, error) {
	row, err := b.db.Query("select groupname, name, state, outbin, log from task where groupname = ? and name = ?")
	if err != nil {
		return nil, err
	}

	if row.Next() {
		st := &buildState{}
		err := row.Scan(&st.groupname, &st.name, &st.state, &st.outbin, &st.log)
		if err != nil {
			row.Close()
			return nil, err
		}
		row.Close()
		return st, nil
	}
	return nil, nil
}
