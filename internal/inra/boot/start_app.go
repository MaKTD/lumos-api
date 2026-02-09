package boot

type App interface {
	Init() error
	Run() error
	Shutdown()
}

func StartApp(app App) error {
	defer app.Shutdown()

	err := app.Init()
	if err != nil {
		return err
	}

	return app.Run()
}
