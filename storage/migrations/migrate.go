package migrations

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/modules/Error"
	"strings"
	"time"
)

type Schema struct {
	Level     int       `gorm:"primary key"`
	Success   bool      `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	TookMs    int       `gorm:"not null"`
}

type Migration interface {
	Name() string
	Level() int
	Run(db *gorm.DB) error
}

type migration struct {
	f     func(db *gorm.DB) error
	name  string
	level int
}

func (m migration) Name() string {
	return m.name
}

func (m migration) Level() int {
	return m.level
}

func (m migration) Run(db *gorm.DB) error {
	return m.f(db)
}

var Migrations = []Migration{
	migration{level: 1, name: "initial schema", f: initialSchema},
}

// Run migrations, if target = -1, run all migrations, otherwise migrate to given level
func RunMigrations(db *gorm.DB, target int) error {
	db.AutoMigrate(&Schema{})
	var migrated bool = false

	tx := db.Begin()
	defer EndTransaction(tx, &migrated)
	failed := &[]Schema{}

	err := tx.Where("success = false").Find(failed).Error
	if err != nil {

	}
	if len(*failed) > 0 {
		var levels []string
		for _, v := range *failed {
			levels = append(levels, string(v.Level))
		}
		text := strings.Join(levels, ", ")
		return &Error.Error{Code: Error.Einternal, Err: errors.New(
			fmt.Sprintf("previous migrations have failed. Unable to continue. Migrations: %s", text))}
	}

	last := &Schema{}
	err = tx.Order("level DESC").First(&last).Error
	if err != nil {
		if err.Error() != "record not found" {
			return Error.Wrap(&err, "failed to run migrations")
		}
		last.Level = 0
	}

	latest := Migrations[len(Migrations)-1].Level()

	if last.Level == latest {
		logrus.Debug("no new migrations")
		return nil
	}

	if target == -1 {
		target = len(Migrations)
	}

	err = migrateDo(tx, last.Level, target)
	if err == nil {
		migrated = true
		return nil
	}
	return Error.Wrap(&err, "migrations failed")
}

func EndTransaction(tx *gorm.DB, success *bool) {
	if *success {
		tx.Commit()
	} else {
		tx.Rollback()
	}
}

func migrateDo(tx *gorm.DB, current int, target int) error {
	logrus.Warnf("running migrations: %d -> %d", current, target)

	for _, v := range Migrations[current:target] {
		start := time.Now()

		err := v.Run(tx)
		if err != nil {
			return Error.Wrap(&err, fmt.Sprintf(
				"Tried to migrate from %d to %d, but failed with %d, '%s'", current, target, v.Level(), v.Name()))
		}
		took := time.Since(start)

		s := &Schema{
			Level:     v.Level(),
			Success:   true,
			Timestamp: time.Now(),
			TookMs:    int(took.Nanoseconds() / 1000000),
		}
		err = tx.Create(&s).Error
		if err != nil {
			return err
		}
	}
	logrus.Warn("migrations ok")
	return nil
}
