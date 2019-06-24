package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/modules/gateway"
	"github.com/tryffel/market/modules/logger"
	"github.com/tryffel/market/storage"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

type Service struct {
	Config     *config.Config
	lock       sync.RWMutex
	logRequest *os.File
	logSql     *os.File
	logFile    *os.File
	running    bool
	Store      *storage.Store
	Api        modules.Tasker
}

func NewService(config *config.Config) (*Service, error) {
	service := &Service{
		Config: config,
	}

	logFormat := &prefixed.TextFormatter{
		FullTimestamp:  true,
		QuoteCharacter: "'",
	}
	logFormat.ForceFormatting = true
	logrus.SetFormatter(logFormat)

	var sqlLogger *logrus.Logger

	err := service.openLogFiles()
	if err != nil {
		err = Error.Wrap(&err, "Could not open all log files")
		e := service.closeLogFiles()
		if e != nil {
			err = Error.Wrap(&err, e.Error())
		}
		return service, err
	}

	if service.Config.Logging.LogSql {
		sqlLogger = logrus.New()
		sqlLogger.SetOutput(service.logSql)
		sqlLogger.SetFormatter(logFormat)
		sqlLogger.SetLevel(logrus.InfoLevel)
	}

	sqlLog := &logger.SqlLogger{
		Logger: sqlLogger,
	}

	service.Store, err = storage.NewStore(config, sqlLog)
	if err != nil {
		err = Error.Wrap(&err, "failed to initialize database")
		logrus.Fatal(err)
	}

	logrus.SetOutput(service.logFile)
	logrus.AddHook(&logger.StdLogger{})

	service.Api, err = gateway.NewApi(config, service.Store)
	if err != nil {
		err = Error.Wrap(&err, "Failed to create api server")
		logrus.Fatal(err)
	}

	return service, nil
}

// Start Start service
func (s *Service) Start() {
	s.lock.Lock()
	if !s.running {
		s.running = true
		s.lock.Unlock()

		logrus.Info("")
		logrus.Warn("Starting service.")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		signal.Notify(c, os.Interrupt, syscall.SIGINT)
		go func() {
			<-c
			s.Stop()
		}()

		err := s.Api.Start()
		if err != nil {
			err = Error.Wrap(&err, "failed to start api server")
			logrus.Error(err)
		}

	} else {
	}

}

// Stop Stop service
func (s *Service) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.running {
		s.running = false
		logrus.Warn("Stopping service")
		s.Api.Stop()
		s.Store.Close()
	}
}

func (s *Service) openLogFiles() error {

	mode := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	perm := os.FileMode(0760)

	conf := &s.Config.Logging
	file, err := os.OpenFile(filepath.Join(conf.Directory, config.LogMainFile), mode, perm)
	if err != nil {
		return err
	}
	s.logFile = file

	if conf.LogSql {
		file, err = os.OpenFile(filepath.Join(conf.Directory, config.LogSqlFile), mode, perm)
		if err != nil {
			return err
		}
		s.logSql = file
	}
	return nil
}

func (s *Service) closeLogFiles() error {
	var e error
	err := s.logFile.Close()
	if err != nil {
		e = Error.Wrap(&err, "failed to close main log file")
	}
	if s.logSql != nil {
		err = s.logSql.Close()
		if err != nil {
			e = Error.Wrap(&e, "failed to close sql log file")
		}
	}
	return e
}
