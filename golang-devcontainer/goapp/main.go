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

	//for i := 1; i < 500; i++ {
	selectByYearQuery(session, cluster.Consistency, logger)
	selectByArtistQuery(session, cluster.Consistency, logger)
	insertQuery(session, cluster.Consistency, 50000, "Anas's Song", "anas", "Song Album", 2022, logger)
	selectQueryWhere(session, cluster.Consistency, "Anas", 2022, logger)
	deleteQuery(session, cluster.Consistency, "Anas", 2022, logger)
	selectQueryWhere(session, cluster.Consistency, "Anas", 2022, logger)
	//}
}

func deleteQuery(session gocqlx.Session, consistency gocql.Consistency, artist string, yearReleased int, logger *zap.Logger) {
	logger.Info("Deleting " + artist + "......")
	r := Record{
		YearReleased: yearReleased,
		Artist:       artist,
	}
	err := session.Query(songsByYearTable.Delete()).BindStruct(r).ExecRelease()
	if err != nil {
		downgradedCl, consistencyError := downgradConsistency(consistency, err, logger)

		if consistencyError != nil {
			logger.Warn("Error ", zap.Error(consistencyError))
			return
		}

		logger.Error("delete", zap.Error(err))
		deleteQuery(session, downgradedCl, artist, yearReleased, logger)
	}
}

func insertQuery(session gocqlx.Session, consistency gocql.Consistency, songId int, title string, artist string, album string, yearReleased int, logger *zap.Logger) {
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

	query := session.Query(songsByYearTable.Insert())
	query.SetConsistency(consistency)

	err := query.BindStruct(r).ExecRelease()
	if err != nil {
		downgradedCl, consistencyError := downgradConsistency(consistency, err, logger)

		if consistencyError != nil {
			logger.Warn("Error ", zap.Error(consistencyError))
			return
		}

		logger.Error("insert", zap.Error(err))
		insertQuery(session, downgradedCl, songId, title, artist, album, yearReleased, logger)
	}
}

func selectByYearQuery(session gocqlx.Session, consistency gocql.Consistency, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record

	query := session.Query(songsByYearTable.SelectAll())
	query.SetConsistency(consistency)

	err := query.SelectRelease(&rs)
	if err != nil {
		downgradedCl, consistencyError := downgradConsistency(consistency, err, logger)

		if consistencyError != nil {
			logger.Warn("Error ", zap.Error(consistencyError))
			return
		}

		logger.Warn("select songs_by_year", zap.Error(err))
		selectByYearQuery(session, downgradedCl, logger)
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

func selectQueryWhere(session gocqlx.Session, consistency gocql.Consistency, artist string, yearReleased int, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record
	r := Record{
		YearReleased: yearReleased,
		Artist:       artist,
	}

	query := session.Query(songsByYearTable.Select())
	query.SetConsistency(consistency)

	err := query.BindStruct(r).SelectRelease(&rs)
	if err != nil {
		downgradedCl, consistencyError := downgradConsistency(consistency, err, logger)

		if consistencyError != nil {
			logger.Warn("Error ", zap.Error(consistencyError))
			return
		}

		logger.Warn("select songs_by_year", zap.Error(err))
		selectQueryWhere(session, downgradedCl, artist, yearReleased, logger)
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

func selectByArtistQuery(session gocqlx.Session, consistency gocql.Consistency, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record

	query := session.Query(songsByArtistTable.SelectAll())
	query.SetConsistency(consistency)

	err := query.SelectRelease(&rs)
	if err != nil {
		downgradedCl, consistencyError := downgradConsistency(consistency, err, logger)

		if consistencyError != nil {
			logger.Warn("Error ", zap.Error(consistencyError))
			return
		}

		logger.Warn("select songs_by_artist", zap.Error(err))
		selectByArtistQuery(session, downgradedCl, logger)
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

func downgradConsistency(current gocql.Consistency, original error, logger *zap.Logger) (gocql.Consistency, error) {
	if current == gocql.LocalQuorum {
		logger.Info("\t" + "Consistency downgraded to LocalOne")
		return gocql.LocalOne, nil
	}

	logger.Warn("\t" + "Can't downgrade more the consistency")
	return current, original
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
