package core

import (
	"go.oease.dev/goe/contracts"
	"go.oease.dev/goe/modules/msearch"
)

type Container struct {
	config      contracts.Config
	mongo       contracts.MongoDB
	mailer      contracts.Mailer
	meilisearch contracts.Meilisearch
	logger      contracts.Logger
	queue       contracts.Queue
	appConfig   *GoeConfig
}

func NewContainer(config contracts.Config, logger contracts.Logger, appConfig *GoeConfig) *Container {
	return &Container{
		config:    config,
		logger:    logger,
		appConfig: appConfig,
	}
}

func (c *Container) InitMongo() {
	// Initialize MongoDB
	mdb, err := NewGoeMongoDB(c.appConfig, c.logger)
	if err != nil {
		c.logger.Panic("Failed to initialize MongoDB: ", err)
		return
	} else {
		c.mongo = mdb
	}
}

func (c *Container) InitMeilisearch() {
	if c.appConfig.Features.MeilisearchEnabled {
		if c.appConfig.Meilisearch.ApiKey == "" {
			c.logger.Panic("meilisearch api key is required")
			return
		}
		if c.appConfig.Meilisearch.Endpoint == "" {
			c.logger.Panic("meilisearch endpoint is required")
			return
		}
		ms := msearch.NewMSearch(c.appConfig.Meilisearch.Endpoint, c.appConfig.Meilisearch.ApiKey, c.logger)
		if ms == nil {
			c.logger.Panic("Failed to initialize Meilisearch")
			return
		}
		c.meilisearch = ms
		if c.appConfig.Features.SearchDBSyncEnabled {
			err := c.mongo.(*GoeMongoDB).SetMeilisearch(ms)
			if err != nil {
				c.logger.Panic("Failed to bind Meilisearch to MongoDB: ", err)
				return
			}
		}
	}
}

func (c *Container) GetConfig() contracts.Config {
	return c.config
}

func (c *Container) GetMongo() contracts.MongoDB {
	return c.mongo
}

func (c *Container) GetMailer() contracts.Mailer {
	return c.mailer
}

func (c *Container) GetMeilisearch() contracts.Meilisearch {
	return c.meilisearch
}

func (c *Container) GetLogger() contracts.Logger {
	return c.logger
}

func (c *Container) GetQueue() contracts.Queue {
	return c.queue
}
