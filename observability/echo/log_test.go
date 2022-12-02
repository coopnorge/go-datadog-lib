package echo

import (
    "bytes"
    "testing"

    "github.com/labstack/gommon/log"
)

func TestWrappedEchoLogger(t *testing.T) {
    l := NewWrappedEchoLogger()
    jTest := log.JSON{"name": "value"}

    b := new(bytes.Buffer)
    l.SetOutput(b)
    l.SetLevel(log.DEBUG)

    l.Debug("debug")
    l.Debugf("debug%s", "f")
    l.Debugj(jTest)
    l.Print("print")
    l.Printf("print%s", "f")
    l.Printj(jTest)
    l.Info("info")
    l.Infof("info%s", "f")
    l.Infoj(jTest)
    l.Warn("warn")
    l.Warnf("warn%s", "f")
    l.Warnj(jTest)
    l.Error("error")
    l.Errorf("error%s", "f")
    l.Errorj(jTest)
}
