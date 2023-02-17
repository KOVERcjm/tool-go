package fx

import (
	"strings"

	"go.uber.org/fx/fxevent"

	"github.com/kovercjm/tool-go/logger"
)

var _ fxevent.Logger = (*fxLogger)(nil)

type fxLogger struct {
	logger logger.Logger
}

func FxLogger(logger logger.Logger) fxevent.Logger {
	return fxLogger{logger: logger.NoCaller()}
}

func (fl fxLogger) LogEvent(event fxevent.Event) {
	if fl.logger == nil {
		return
	}
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		fl.logger.Info("OnStart hook executing", "callee", e.FunctionName, "caller", e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			fl.logger.Error("OnStart hook failed", "callee", e.FunctionName, "caller", e.CallerName)
		} else {
			fl.logger.Info("OnStart hook executed", "callee", e.FunctionName, "caller", e.CallerName, "runtime", e.Runtime.String())
		}
	case *fxevent.OnStopExecuting:
		fl.logger.Info("OnStop hook executing", "callee", e.FunctionName, "caller", e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			fl.logger.Error("OnStop hook failed", "callee", e.FunctionName, "caller", e.CallerName, "error", e.Err)
		} else {
			fl.logger.Info("OnStop hook executed", "callee", e.FunctionName, "caller", e.CallerName, "runtime", e.Runtime.String())
		}
	case *fxevent.Supplied:
		fl.logger.Info("supplied", "type", e.TypeName, "module", e.ModuleName, "error", e.Err)
	case *fxevent.Provided:
		for _, typeName := range e.OutputTypeNames {
			fl.logger.Info("provided", "constructor", e.ConstructorName, "module", e.ModuleName, "type", typeName)
		}
		if e.Err != nil {
			fl.logger.Error("error encountered while applying options", "module", e.ModuleName, "error", e.Err)
		}
	case *fxevent.Decorated:
		for _, typeName := range e.OutputTypeNames {
			fl.logger.Info("decorated", "decorator", e.DecoratorName, "module", e.ModuleName, "type", typeName)
		}
		if e.Err != nil {
			fl.logger.Error("error encountered while applying options", "module", e.ModuleName, "error", e.Err)
		}
	case *fxevent.Invoking:
		fl.logger.Info("invoking", "function", e.FunctionName, "module", e.ModuleName)
	case *fxevent.Invoked:
		if e.Err != nil {
			fl.logger.Error("invoke failed", "function", e.FunctionName, "module", e.ModuleName, "error", e.Err, "stack", e.Trace)
		}
	case *fxevent.Stopping:
		fl.logger.Info("received signal", "signal", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			fl.logger.Error("stop failed", "error", e.Err)
		}
	case *fxevent.RollingBack:
		fl.logger.Error("start failed, rolling back", "error", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			fl.logger.Error("rollback failed", "error", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			fl.logger.Error("start failed", "error", e.Err)
		} else {
			fl.logger.Info("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			fl.logger.Error("custom logger initialization failed", "error", e.Err)
		} else {
			fl.logger.Info("initialized custom fxevent.Logger", "function", e.ConstructorName)
		}
	}
}
