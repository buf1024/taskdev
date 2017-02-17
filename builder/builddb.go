package main

import (
	"database/sql"
	"sync"

	"fmt"

	"time"

	"taskdev/util"

	_ "github.com/mattn/go-sqlite3"
)

const (
	buildStateStop         int32 = 0
	buildStateBuildFailed  int32 = util.KEError
	buildStateBuildSuccess int32 = util.KESuccess
)

type buildDb struct {
	db   *sql.DB
	lock sync.Locker
	init bool
}

type buildState struct {
	groupname string
	name      string
	state     int32
	outbin    string
	log       string
	uptime    time.Time
}

var sqlStr string = `
create table task (
    id            integer primary key autoincrement,
    groupname     varchar(64) not null,
    name          varchar(64) not null,
    state         integer,
    outbin        varchar(256),
    log           varchar(256),
    uptime        timestamp not null default (datetime('now','localtime'))
);
`

func NewBuildDb(file string) (*buildDb, error) {
	b := &buildDb{}
	var err error
	b.db, err = sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	_, err = b.db.Query("select count(1) from task")
	if err != nil {
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
	b.init = true

	return b, nil

}
func BuildDbDestroy(b *buildDb) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.init = false
	return b.db.Close()
}
func (b *buildDb) insert(st *buildState) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if !b.init {
		return fmt.Errorf("not init")
	}

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

func (b *buildDb) gets(groupname string, name string) ([]*buildState, error) {
	if !b.init {
		return nil, fmt.Errorf("not init")
	}

	stm, err := b.db.Prepare("select groupname, name, state, outbin, log, uptime from task where groupname = ? and name = ? order by uptime desc")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	row, err := stm.Query(groupname, name)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var states []*buildState
	for row.Next() {
		st := &buildState{}
		err := row.Scan(&st.groupname, &st.name, &st.state, &st.outbin, &st.log, &st.uptime)
		if err != nil {
			return nil, err
		}
		states = append(states, st)
	}

	return states, nil
}
func (b *buildDb) get(groupname string, name string) (*buildState, error) {
	if !b.init {
		return nil, fmt.Errorf("not init")
	}

	stm, err := b.db.Prepare("select groupname, name, state, outbin, log, uptime from task where groupname = ? and name = ? order by uptime desc")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	row, err := stm.Query(groupname, name)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		st := &buildState{}
		err := row.Scan(&st.groupname, &st.name, &st.state, &st.outbin, &st.log, &st.uptime)
		if err != nil {
			return nil, err
		}

		return st, nil
	}

	return nil, nil
}
