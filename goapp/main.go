package main

import (
	"goapp/internal/log"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	"go.uber.org/zap"
)

func main() {
	logger := log.CreateLogger("info")

	cluster := gocql.NewCluster("scylla-node1", "scylla-node2", "scylla-node3")
	cluster.Keyspace = "songs"
	cluster.Consistency = gocql.LocalQuorum
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("DC1"))

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		logger.Fatal("unable to connect to scylla", zap.Error(err))
	}
	defer session.Close()

	//for i := 1; i < 20; i++ {
	selectByYearQuery(session, logger)
	selectByArtistQuery(session, logger)
	insertQuery(session, 50000, "Anas's Song", "anas", "Song Album", 2022, logger)
	selectQueryWhere(session, "Anas", 2022, logger)
	deleteQuery(session, "Anas", 2022, logger)
	selectQueryWhere(session, "Anas", 2022, logger)
	//}
}

func deleteQuery(session gocqlx.Session, artist string, yearReleased int, logger *zap.Logger) {
	logger.Info("Deleting " + artist + "......")
	r := Record{
		YearReleased: yearReleased,
		Artist:       artist,
	}
	err := session.Query(songsByYearTable.Delete()).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error("delete", zap.Error(err))
	}
}

func insertQuery(session gocqlx.Session, songId int, title string, artist string, album string, yearReleased int, logger *zap.Logger) {
	logger.Info("Inserting " + title + "......")
	r := Record{
		SongId:       songId,
		Title:        title,
		Artist:       artist,
		Album:        album,
		YearReleased: yearReleased,
		Duration:     0,
		Tempo:        0,
		Loudness:     0,
	}
	err := session.Query(songsByYearTable.Insert()).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error("insert", zap.Error(err))
	}
}

func selectByYearQuery(session gocqlx.Session, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record
	err := session.Query(songsByYearTable.SelectAll()).SelectRelease(&rs)
	if err != nil {
		logger.Warn("select songs_by_year", zap.Error(err))
		return
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

func selectQueryWhere(session gocqlx.Session, artist string, yearReleased int, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record
	r := Record{
		YearReleased: yearReleased,
		Artist:       artist,
	}

	err := session.Query(songsByYearTable.Select()).BindStruct(r).SelectRelease(&rs)
	if err != nil {
		logger.Warn("select songs_by_year", zap.Error(err))
		return
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

func selectByArtistQuery(session gocqlx.Session, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record
	err := session.Query(songsByArtistTable.SelectAll()).SelectRelease(&rs)
	if err != nil {
		logger.Warn("select songs_by_artist", zap.Error(err))
		return
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

// Songs By Artist
var songsByArtistMetadata = table.Metadata{
	Name:    "songs_by_artist",
	Columns: []string{"song_id", "title", "artist", "album", "year_released", "duration", "tempo", "loudness"},
	PartKey: []string{"artist"},
	SortKey: []string{"year_released"},
}
var songsByArtistTable = table.New(songsByArtistMetadata)

// Songs By Year
var songsByYearMetadata = table.Metadata{
	Name:    "songs_by_year",
	Columns: []string{"song_id", "title", "artist", "album", "year_released", "duration", "tempo", "loudness"},
	PartKey: []string{"year_released"},
	SortKey: []string{"artist"},
}
var songsByYearTable = table.New(songsByYearMetadata)

type Record struct {
	SongId       int     `db:"song_id"`
	Title        string  `db:"title"`
	Artist       string  `db:"artist"`
	Album        string  `db:"album"`
	YearReleased int     `db:"year_released"`
	Duration     float64 `db:"duration"`
	Tempo        float64 `db:"tempo"`
	Loudness     float64 `db:"loudness"`
}
