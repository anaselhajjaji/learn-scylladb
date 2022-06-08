package main

import (
	"goapp/internal/log"
	"goapp/internal/scylla"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"go.uber.org/zap"
)

var stmts = createStatements()

func main() {
	logger := log.CreateLogger("info")

	cluster := scylla.CreateCluster(gocql.Quorum, "songs", "scylla-node1", "scylla-node2", "scylla-node3")
	session, err := gocql.NewSession(*cluster)
	if err != nil {
		logger.Fatal("unable to connect to scylla", zap.Error(err))
	}
	defer session.Close()

	selectQuery(session, logger)
	//insertQuery(session, "Mike", "Tyson", "12345 Foo Lane", "http://www.facebook.com/mtyson", logger)
	//insertQuery(session, "Alex", "Jones", "56789 Hickory St", "http://www.facebook.com/ajones", logger)
	//selectQuery(session, logger)
	//deleteQuery(session, "Mike", "Tyson", logger)
	//selectQuery(session, logger)
	//deleteQuery(session, "Alex", "Jones", logger)
	//selectQuery(session, logger)
}

/*func deleteQuery(session *gocql.Session, firstName string, lastName string, logger *zap.Logger) {
	logger.Info("Deleting " + firstName + "......")
	r := Record{
		FirstName: firstName,
		LastName:  lastName,
	}
	err := gocqlx.Query(session.Query(stmts.del.stmt), stmts.del.names).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error("delete catalog.mutant_data", zap.Error(err))
	}
}*/

/*func insertQuery(session *gocql.Session, firstName, lastName, address, pictureLocation string, logger *zap.Logger) {
	logger.Info("Inserting " + firstName + "......")
	r := Record{
		FirstName:       firstName,
		LastName:        lastName,
		Address:         address,
		PictureLocation: pictureLocation,
	}
	err := gocqlx.Query(session.Query(stmts.ins.stmt), stmts.ins.names).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error("insert catalog.mutant_data", zap.Error(err))
	}
}*/

func selectQuery(session *gocql.Session, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record
	err := gocqlx.Query(session.Query(stmts.sel.stmt), stmts.sel.names).SelectRelease(&rs)
	if err != nil {
		logger.Warn("select songs_by_year", zap.Error(err))
		return
	}
	for _, r := range rs {
		logger.Info("\t" + strconv.Itoa(r.YearReleased) + " " + r.Artist + ", " + r.Title + ", " + r.Album)
	}
}

func createStatements() *statements {
	m := table.Metadata{
		Name:    "songs_by_year",
		Columns: []string{"song_id", "title", "artist", "album", "year_released", "duration", "tempo", "loudness"},
		PartKey: []string{"year_released"},
		SortKey: []string{"artist"},
	}
	tbl := table.New(m)

	deleteStmt, deleteNames := tbl.Delete()
	insertStmt, insertNames := tbl.Insert()
	// Normally a select statement such as this would use `tbl.Select()` to select by
	// primary key but now we just want to display all the records...
	selectStmt, selectNames := qb.Select(m.Name).Columns(m.Columns...).ToCql()
	return &statements{
		del: query{
			stmt:  deleteStmt,
			names: deleteNames,
		},
		ins: query{
			stmt:  insertStmt,
			names: insertNames,
		},
		sel: query{
			stmt:  selectStmt,
			names: selectNames,
		},
	}
}

type query struct {
	stmt  string
	names []string
}

type statements struct {
	del query
	ins query
	sel query
}

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
