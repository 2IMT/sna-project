package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
    db *sql.DB
    scoreQuery string
    scoreInsert string
    scoreUpdate string
    scoreExists string
    leaderboardQuery string
}

type Score struct {
    Id int64
    Score int64
}

func NewDatabase() (Database, error) {
    var result Database
    var err error

    connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s "+
        "sslmode=disable",
        Env.DbUser,
        Env.DbPass,
        Env.DbName,
        Env.DbHost,
        Env.DbPort)
    result.db, err = sql.Open("postgres", connStr)
    if err != nil {
        return Database{}, fmt.Errorf("failed to open database: %s", err)
    }

    err = result.db.Ping()
    if err != nil {
        return Database{}, fmt.Errorf("failed to ping database: %s", err)
    }

    err = loadScripts(&result)
    if err != nil {
        return Database{}, err
    }

    return result, nil
}

func (db *Database) QueryScore(id int64) (int64, error) {
    var result int64

    row := db.db.QueryRow(db.scoreQuery, id)
    err := row.Scan(&result)
    if err != nil {
        return 0, fmt.Errorf("failed to query: %s", err)
    }

    return result, nil
}

func (db *Database) UpdateScore(id int64, score int64) error {
    _, err := db.db.Exec(db.scoreUpdate, id, score)
    if err != nil {
        return fmt.Errorf("failed to update: %s", err)
    }

    return nil
}

func (db *Database) InsertScore(id int64, score int64) error {
    _, err := db.db.Exec(db.scoreInsert, id, score)
    if err != nil {
        return fmt.Errorf("failed to insert: %s", err)
    }

    return nil
}

func (db *Database) ScoreExists(id int64) (bool, error) {
    row := db.db.QueryRow(db.scoreExists, id)
    var exists bool
    err := row.Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check if score exists: %s", err)
    }

    return exists, nil
}

func (db *Database) QueryLeaderboard() ([]Score, error) {
    scores := make([]Score, 10)
    rows, err := db.db.Query(db.leaderboardQuery)
    if err != nil {
        return nil, fmt.Errorf("failed to query leaderboard: %s", err)
    }

    for rows.Next() {
        var score Score
        rows.Scan(&score.Id, &score.Score)
        scores = append(scores, score)
    }

    return scores, nil
}

func readSql(path string) (string, error) {
    str, err := ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("failed to read %s: %s", path, err)
    }

    return str, nil
}

func loadScripts(db *Database) error {
    var err error

    db.scoreQuery, err = readSql("sql/query_score.sql")
    if err != nil {
        return err
    }

    db.scoreUpdate, err = readSql("sql/update_score.sql")
    if err != nil {
        return err
    }

    db.scoreInsert, err = readSql("sql/insert_score.sql")
    if err != nil {
        return err
    }

    db.scoreExists, err = readSql("sql/exists_score.sql")
    if err != nil {
        return err
    }

    db.leaderboardQuery, err = readSql("sql/query_leaderboard.sql")
    if err != nil {
        return err
    }

    return nil
}
