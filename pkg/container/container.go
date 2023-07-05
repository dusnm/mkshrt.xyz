package container

import (
	"database/sql"
	"embed"
	"github.com/dusnm/mkshrt.xyz/pkg/config"
	"github.com/dusnm/mkshrt.xyz/pkg/database"
	"github.com/dusnm/mkshrt.xyz/pkg/repositories/mapping"
	"log"
)

type (
	Container struct {
		Views  embed.FS
		Assets embed.FS
		cfg    *config.Config
		db     *sql.DB

		mappingRepo mapping.Interface
	}
)

func New(
	views embed.FS,
	assets embed.FS,
) *Container {
	return &Container{
		Views:  views,
		Assets: assets,
	}
}

func (c *Container) GetConfig() *config.Config {
	if c.cfg == nil {
		cfg, err := config.New()
		if err != nil {
			log.Fatal(err)
		}

		c.cfg = cfg
	}

	return c.cfg
}

func (c *Container) GetDB() *sql.DB {
	if c.db == nil {
		db, err := database.New(c.GetConfig())
		if err != nil {
			log.Fatal(err)
		}

		c.db = db
	}

	return c.db
}

func (c *Container) GetMappingRepository() mapping.Interface {
	if c.mappingRepo == nil {
		c.mappingRepo = mapping.New(c.GetDB())
	}

	return c.mappingRepo
}

func (c *Container) Close() error {
	if err := c.db.Close(); err != nil {
		return err
	}

	return nil
}
