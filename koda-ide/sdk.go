package main

import (
	"koda/api"
)

const studioVersion = "0.4.0"

func (a *App) bootstrapSDK() {
	api.EnsureSDKFromExecutable()
}

// CheckSDK returns a lightweight health report for the welcome screen.
func (a *App) CheckSDK() api.SDKStatus {
	a.bootstrapSDK()
	return api.CheckSDK(studioVersion)
}

// ListProjectTemplates returns hello, game, graphics, … for the new-project wizard.
func (a *App) ListProjectTemplates() []string {
	return api.ListProjectTemplates()
}
