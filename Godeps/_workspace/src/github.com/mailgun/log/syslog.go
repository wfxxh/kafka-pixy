package log

import (
	"fmt"
	"io"
	"log/syslog"
)

// sysLogger logs messages to rsyslog MAIL_LOG facility.
type sysLogger struct {
	sev Severity

	infoW  io.Writer
	warnW  io.Writer
	errorW io.Writer
}

func NewSysLogger(conf Config) (Logger, error) {
	infoW, err := syslog.New(syslog.LOG_MAIL|syslog.LOG_INFO, appname)
	if err != nil {
		return nil, err
	}

	warnW, err := syslog.New(syslog.LOG_MAIL|syslog.LOG_WARNING, appname)
	if err != nil {
		return nil, err
	}

	errorW, err := syslog.New(syslog.LOG_MAIL|syslog.LOG_ERR, appname)
	if err != nil {
		return nil, err
	}

	sev, err := severityFromString(conf.Severity)
	if err != nil {
		return nil, err
	}

	return &sysLogger{sev, infoW, warnW, errorW}, nil
}

func (l *sysLogger) Writer(sev Severity) io.Writer {
	// is this logger configured to log at the provided severity?
	if sev >= l.sev {
		// return an appropriate writer
		switch sev {
		case SeverityInfo:
			return l.infoW
		case SeverityWarning:
			return l.warnW
		default:
			return l.errorW
		}
	}
	return nil
}

func (l *sysLogger) FormatMessage(sev Severity, caller *CallerInfo, format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", sev, fmt.Sprintf(format, args...))
}
