package libwebsocketd

import (
	"context"
	elastic "gopkg.in/olivere/elastic.v5"
)

type EsHandler struct {
	client *elastic.Client
	log *LogScope
	config *Config
}

func NewESHandler(config *Config,log *LogScope) (esh *EsHandler) {
	client, err := elastic.NewClient(elastic.SetURL(config.EsUrl))
	if err != nil {
		panic(err)
	}
	return &EsHandler{client:client,log:log,config:config}
}

func (esh *EsHandler) index(command CommandInfo) {
	// Add a document to the index
	_, err := esh.client.Index().
		Index("opsagent").
		Type("command").
		BodyJson(command).
		Do(context.TODO())
	if err != nil {
		// Handle error
		esh.log.Error("server", "index to es error : %s", err)
	}
}
