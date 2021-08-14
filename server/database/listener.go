package database

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/bebo-dot-dev/go-postgres-sockets/server/socket"
    "github.com/lib/pq"
    "log"
    "os"
    "time"
)

const (
    dbport = 5432
    dbname = "notifications"
)

type PostgresDbListener struct {
    socketHub *socket.Hub
}

func NewPostgresDbListener(hub *socket.Hub) *PostgresDbListener {
    return &PostgresDbListener{
        socketHub: hub,
    }
}

func (l *PostgresDbListener) getDbListener() *pq.Listener {
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        dbport,
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        dbname)

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
            l.doLog(true, "Postgres database listener state change: ", err)
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
            l.doLog(false, "DB listener received data from channel [", n.Channel, "]")
            var prettyJSON bytes.Buffer
            err := json.Indent(&prettyJSON, []byte(n.Extra), "", "    ")
            if err != nil {
                l.doLog(true, "DB listener error processing JSON: ", err)
            }
            log.Println(string(prettyJSON.Bytes()))
            l.socketHub.Broadcast <- prettyJSON.Bytes()
        case <-time.After(90 * time.Second):
            l.doLog(false, "DB listener received no notification events for 90 seconds, pinging DB connection")
            go func() {
                err := dbl.Ping()
                if err != nil {
                    l.doLog(true, err)
                }
            }()
        }
    }
}

func (l *PostgresDbListener) doLog(isError bool, v ...interface{}) {
    if isError {
        log.SetOutput(os.Stderr)
    } else {
        log.SetOutput(os.Stdout)
    }
    log.Println(v)
}

func (l *PostgresDbListener) Listen() {
    l.doLog(false, "Starting database notifications listener")
    dbl := l.getDbListener()
    go func() {
        for {
            l.waitForNotification(dbl)
        }
    }()
}