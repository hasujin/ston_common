package log2paul

import (
	"os"

	logging "github.com/hasujin/ston_common/go-logging2"
)

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Level int

const (
	CRITICAL Level = iota // Lv 6
	ERROR                 // Lv 5
	WARNING               // Lv 4
	NOTICE                // Lv 3
	INFO                  // Lv 2
	DEBUG                 // Lv 1
)

type Password string

type LOG *logging.Logger

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

type NopBackend struct {
}

func LEVEL_STR(level string) (level_ Level) {

	switch level {
	case "DEBUG":
		level_ = DEBUG
	case "INFO":
		level_ = INFO
	case "NOTICE":
		level_ = NOTICE
	case "WARNING":
		level_ = WARNING
	case "ERROR":
		level_ = ERROR
	case "CRITICAL":
		level_ = CRITICAL
	default:
		level_ = DEBUG
	}
	return
}

func SetupLog(modulename string, level Level, location string) *logging.Logger {
	/*
	   	t := time.Now().UTC()    // Get UTC
	       t = t.In(time.FixedZone("KST", 9*60*60)) // KST adding 9hr
	*/
	var log *logging.Logger
	//log = logging.MustGetLogger(modulename)

	if location == "" {
		log = logging.MustGetLogger(modulename)
	} else {
		log = logging.MustGetLoggerPaul(modulename, location)
	}

	//format := logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05 UTC} [%{module}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`, )
	format := logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05} [%{module}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

	//log.Noticef("FORMAT : %v\n", format)
	//log.Noticef("RECORD : %v\n", logging.Record)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	//logging.SetBackend(backendFormatter).SetLevel(logging.Level(-1), "")
	logging.SetBackend(backendFormatter).SetLevel(logging.Level(level), "")

	/*
		log.Debugf("debug %s", Password("secret"))
		log.Infof("info")
		log.Noticef("notice")
		log.Warningf("warning")
		log.Errorf("err")
		log.Criticalf("crit")
	*/
	return log
}

func GetLogger(module string) *logging.Logger {
	return &logging.Logger{Module: module}
}

func MustGetLogger(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}

func Reset(log *logging.Logger, modulename string, level Level, location string) *logging.Logger {
	/*
		t := time.Now().UTC()    // Get UTC
		t = t.In(time.FixedZone("KST", 9*60*60)) // KST adding 9hr
	*/
	if log == nil {
		log.Critical("Cannot RESET (nil) of the given log")
		return nil
	}
	//var log *logging.Logger
	//log = logging.MustGetLogger(modulename)

	log = logging.MustGetLoggerPaul(modulename, location)

	format := logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05} [%{module}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

	log.Noticef("RESET FORMAT\n")
	//log.Noticef("RECORD : %v\n", logging.Record)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	//logging.SetBackend(backendFormatter).SetLevel(logging.Level(-1), "")
	logging.SetBackend(backendFormatter).SetLevel(logging.Level(level), "")

	//logging.ChangeFormat()
	/*
		log.Debugf("debug %s", Password("secret"))
		log.Infof("info")
		log.Noticef("notice")
		log.Warningf("warning")
		log.Errorf("err")
		log.Criticalf("crit")
	*/
	return log
}

func SetupLog_(modulename string) *logging.Logger {

	// For demo purposes, create two backend for os.Stderr.
	format := logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05 MST} [%{module}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)

	//backend1Leveled.SetLevel(logging.ERROR, "")
	backend1Leveled.SetLevel(logging.Level(-1), "")

	log := logging.MustGetLogger(modulename)

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

	//log.Debugf("LoggLevel %s\n", logging.GetLevel(modulename))
	//log.Debugf("LoggLevel %s\n", logging.GetLevel(backend1Leveled.GetLevel(modulename)))
	log.Debugf("debug %s", Password("secret"))
	log.Infof("info")
	log.Noticef("notice")
	log.Warningf("warning")
	log.Errorf("err")
	log.Criticalf("crit")

	return log

}

func IsEnabledForLogLevel(log *logging.Logger, logLevel string) bool {
	lvl, _ := logging.LogLevel(logLevel)
	return log.IsEnabledFor(lvl)
}
