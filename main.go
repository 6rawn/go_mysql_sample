package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	masterDB *sqlx.DB
	slaveDB1 *sqlx.DB
	slaveDB2 *sqlx.DB
	db       = map[int]*sqlx.DB{}
)

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	UserID    int       `db:"user_id"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
}

func (c *Comment) String() string {
	return "Comment{ID: " + strconv.Itoa(c.ID) + ", PostID: " + strconv.Itoa(c.PostID) + ", UserID: " + strconv.Itoa(c.UserID) + ", Comment: " + c.Comment + ", CreatedAt: " + c.CreatedAt.String() + "}"
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\nError: %s", msg, err)
	}
}

func setupDB() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	masterPort := os.Getenv("DB_MASTER_PORT")
	if masterPort == "" {
		masterPort = "3306"
	}
	_, err := strconv.Atoi(masterPort)
	failOnError(err, "Failed to read DB port number from an environment variable DB_MASTER_PORT.")

	slave1Port := os.Getenv("DB_SLAVE1_PORT")
	if slave1Port == "" {
		slave1Port = "3307"
	}
	_, err = strconv.Atoi(slave1Port)
	failOnError(err, "Failed to read DB port number from an environment variable DB_SLAVE1_PORT.")

	slave2Port := os.Getenv("DB_SLAVE2_PORT")
	if slave2Port == "" {
		slave2Port = "3308"
	}
	_, err = strconv.Atoi(slave2Port)
	failOnError(err, "Failed to read DB port number from an environment variable DB_SLAVE2_PORT.")

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASSWORD")

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "isuconp"
	}

	masterDsn := mysql.Config{
		DBName:               dbname,
		User:                 user,
		Passwd:               password,
		AllowNativePasswords: true,
		Addr:                 host + ":" + masterPort,
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  time.Local,
		InterpolateParams:    true,
	}

	slave1Dsn := mysql.Config{
		DBName:               dbname,
		User:                 user,
		Passwd:               password,
		AllowNativePasswords: true,
		Addr:                 host + ":" + slave1Port,
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  time.Local,
		InterpolateParams:    true,
	}

	slave2Dsn := mysql.Config{
		DBName:               dbname,
		User:                 user,
		Passwd:               password,
		AllowNativePasswords: true,
		Addr:                 host + ":" + slave2Port,
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  time.Local,
		InterpolateParams:    true,
	}

	masterDB, err = sqlx.Open("mysql", masterDsn.FormatDSN())
	failOnError(err, "Failed to connect to DB master.")

	slaveDB1, err = sqlx.Open("mysql", slave1Dsn.FormatDSN())
	failOnError(err, "Failed to connect to DB slave1.")

	slaveDB2, err = sqlx.Open("mysql", slave2Dsn.FormatDSN())
	failOnError(err, "Failed to connect to DB slave2.")

	db[0] = masterDB
	db[1] = slaveDB1
	db[2] = slaveDB2
}

func init() {
	log.Print("Initializing...")
	setupDB()
	log.Print("Initialized.")
}

func main() {
	defer masterDB.Close()
	defer slaveDB1.Close()
	defer slaveDB2.Close()

	comment32 := Comment{}
	commentId := 32
	err := db[commentId%len(db)].Get(&comment32, "SELECT id, post_id, user_id, comment, created_at FROM comments WHERE id = ?", commentId)
	failOnError(err, "Failed to get comment.")
	log.Print(comment32.String())

	comment33 := Comment{}
	commentId = 33
	err = db[commentId%len(db)].Get(&comment33, "SELECT id, post_id, user_id, comment, created_at FROM comments WHERE id = ?", commentId)
	failOnError(err, "Failed to get comment.")
	log.Print(comment33.String())

	comment34 := Comment{}
	commentId = 34
	err = db[commentId%len(db)].Get(&comment34, "SELECT id, post_id, user_id, comment, created_at FROM comments WHERE id = ?", commentId)
	failOnError(err, "Failed to get comment.")
	log.Print(comment34.String())
}
