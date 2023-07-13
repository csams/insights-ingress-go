package logging

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	lc "github.com/redhatinsights/platform-go-middlewares/logging/cloudwatch"
	"github.com/sirupsen/logrus"
)

var (
    // Log is the global `*logrus.Logger` instance
    Log *logrus.Logger
)

// Setup initializes the Entitlements API logger
func Setup(c CompletedConfig) {
	Log = &logrus.Logger{
		Out:          os.Stdout,
		Level:        c.LogLevel,
		Formatter:    NewCloudwatchFormatter(c.Common.Hostname),
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: true,
	}

	if c.AwsAccessKeyId != "" {
		cred := credentials.NewStaticCredentials(c.AwsAccessKeyId, c.AwsSecretAccessKey, "")
		awsconf := aws.NewConfig().WithRegion(c.AwsRegion).WithCredentials(cred)
		hook, err := lc.NewBatchingHook(c.LogGroup, c.Common.Hostname, awsconf, 10*time.Second)
		if err != nil {
			Log.Info(err)
		}
		Log.Hooks.Add(hook)
	}
}

// CustomCloudwatch adds hostname and app name
type CustomCloudwatch struct {
	Hostname string
}

// Marshaller is an interface any type can implement to change its output in our production logs.
type Marshaller interface {
	MarshalLog() map[string]interface{}
}

// NewCloudwatchFormatter creates a new log formatter
func NewCloudwatchFormatter(hostname string) *CustomCloudwatch {
	return &CustomCloudwatch{
		Hostname: hostname,
	}
}

// Format is the log formatter for the entry
func (f *CustomCloudwatch) Format(entry *logrus.Entry) ([]byte, error) {
	data := map[string]interface{}{
		"@timestamp":  time.Now().Format("2006-01-02T15:04:05.999Z"),
		"@version":    1,
		"message":     entry.Message,
		"levelname":   entry.Level.String(),
		"source_host": f.Hostname,
		"app":         "ingress",
		"caller":      entry.Caller.Func.Name(),
	}

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		case Marshaller:
			data[k] = v.MarshalLog()
		default:
			data[k] = v
		}
	}

	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add newline to make stdout readable
	j = append(j, 0x0a)

	b := &bytes.Buffer{}
	b.Write(j)

	return b.Bytes(), nil
}
