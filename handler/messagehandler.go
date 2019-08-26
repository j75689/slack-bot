package handler

// SlackMessageHandler process Slack message
type SlackMessageHandler struct {
	Processer map[StageType]StageRunner
}

// Do ...
func (obj *SlackMessageHandler) Do() (string, error) {
	return "", nil
}

// DryRun ...
func (obj *SlackMessageHandler) DryRun() (string, error) {
	return "", nil
}

func newSlackMessageHandler() *SlackMessageHandler {
	return &SlackMessageHandler{
		Processer: map[StageType]StageRunner{
			renderTag: &RenderProcesser{},
			actionTag: &ActionProcesser{},
		},
	}
}
