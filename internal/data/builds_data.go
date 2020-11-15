package data

import (
	"deploy-lander/internal/structs"
	"github.com/pkg/errors"
	"time"
)

type BuildData struct {
	version  string
	created  time.Time
	deployed *time.Time
}

type BuildsData struct {
	data map[string]*BuildData
}

func NewBuildsData() *BuildsData {
	return &BuildsData{make(map[string]*BuildData)}
}

func (d *BuildsData) RegisterBuild(version structs.BuildVersioner) error {

	d.data[version.Value()] = &BuildData{
		version:  version.Value(),
		created:  time.Now(),
		deployed: nil,
	}

	return nil
}

func (d *BuildsData) IsRegistered(version structs.BuildVersioner) (found bool) {
	_, found = d.data[version.Value()]
	return
}

func (d *BuildsData) MarkDeployed(version structs.BuildVersioner) error {
	if !d.IsRegistered(version) {
		return errors.Errorf("version %v not found", version.Value())
	}

	now := time.Now()
	d.data[version.Value()].deployed = &now

	return nil
}
