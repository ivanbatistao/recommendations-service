package main

type Application struct {
	config Config
	recommendationService *recommendation.Service
	httpServer *HTTPServer
}

func (app *Application) mount() http.Handler {
	return nil
}

func (app *Application) run(h http.Handler) error {
	return nil
}