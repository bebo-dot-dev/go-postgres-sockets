package database

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/lib/pq"
    "log"
    "os"
    "time"
)

const (
    dbhost = "localhost"
    dbport = 5432
    dbuser = "postgres"
    dbname = "notifications"
)

type PostgresDbListener struct { }

func NewPostgresDbListener() *PostgresDbListener {
    return &PostgresDbListener{}
}

func (l *PostgresDbListener) getDbListener() *pq.Listener {
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        dbhost, dbport, dbuser, os.Getenv("POSTGRES_PASSWORD"), dbname)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    err = db.Close()
    if err != nil {
        panic(err)
    }

    stateChange := func(ev pq.ListenerEventType, err error) {
        if err != nil {
            log.Printf("postgres database listener state change: %v", err.Error())
        }
    }

    listener := pq.NewListener(connStr, 10 * time.Second, time.Minute, stateChange)
    err = listener.Listen("notifications_data_changed")
    if err != nil {
        panic(err)
    }

    return listener
}

func (l *PostgresDbListener) waitForNotification(dbl *pq.Listener) {
    for {
        select {
        case n := <- dbl.Notify:
            log.Println("DB listener received data from channel [", n.Channel, "] :")
            var prettyJSON bytes.Buffer
            err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
            if err != nil {
                log.Println("DB listener error processing JSON: ", err)
            }
            log.Println(string(prettyJSON.Bytes()))
        case <-time.After(90 * time.Second):
            log.Println("DB listener received no notification events for 90 seconds, checking connection")
            go func() {
                err := dbl.Ping()
                if err != nil {
                    log.Println("listener ping error: ", err)
                }
            }()
        }
    }
}

func (l *PostgresDbListener) Listen() {
    log.Println("Starting database notifications listener")
    dbl := l.getDbListener()
    go func() {
        for {
            l.waitForNotification(dbl)
        }
    }()
}