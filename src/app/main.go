package app

import (
	"../goVirtualHost"
	"../param"
	"../serverErrHandler"
	"../serverLog"
	"../vhostHandler"
	"os"
)

type App struct {
	vhostSvc      *goVirtualHost.Service
	vhostHandlers []*vhostHandler.VhostHandler
}

func (app *App) Open() {
	errors := app.vhostSvc.Open()
	for _, err := range errors {
		serverErrHandler.CheckError(err)
	}
}

func (app *App) Close() {
	for _, vhHandler := range app.vhostHandlers {
		vhHandler.Close()
	}

	app.vhostSvc.Close()
}

func (app *App) ReOpenLog() {
	for _, vhhandler := range app.vhostHandlers {
		vhhandler.ReOpenLog()
	}
}

func NewApp(params []*param.Param) *App {
	vhSvc := goVirtualHost.NewService()
	vhHandlers := make([]*vhostHandler.VhostHandler, 0, len(params))

	for _, p := range params {
		// logger
		logger := serverLog.NewLogger(p.AccessLog, p.ErrorLog)
		errors := logger.Open()
		serverErrHandler.CheckFatal(errors...)

		// ErrHandler
		errHandler := serverErrHandler.NewErrHandler(logger)

		// ServeMux
		vhHandler := vhostHandler.NewHandler(p, logger, errHandler)
		vhHandlers = append(vhHandlers, vhHandler)

		// init vhost
		listens := p.Listens
		if len(listens) == 0 && len(p.ListensPlain) == 0 && len(p.ListensTLS) == 0 {
			if p.Certificate == nil {
				listens = []string{":80"}
			} else {
				listens = []string{":443"}
			}
		}

		errors = vhSvc.Add(&goVirtualHost.HostInfo{
			Listens:      listens,
			ListensPlain: p.ListensPlain,
			ListensTLS:   p.ListensTLS,
			Cert:         p.Certificate,
			HostNames:    p.HostNames,
			Handler:      vhHandler.Handler,
		})
		if len(errors) > 0 {
			serverErrHandler.CheckFatal(errors...)
			logger.LogErrors(errors...)
			os.Exit(1)
		}
	}

	return &App{
		vhostSvc:      vhSvc,
		vhostHandlers: vhHandlers,
	}
}
