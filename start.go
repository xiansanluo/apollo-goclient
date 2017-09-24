package agollo

//start apollo
func Start() {
	initApp()
	initNotify()
	//first sync
	notifySyncConfigServices()

	//start auto refresh config
	go StartRefreshConfig(&AutoRefreshConfigComponent{})

	//start long poll sync config
	go StartRefreshConfig(&NotifyConfigComponent{})
}
