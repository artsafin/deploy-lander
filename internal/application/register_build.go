package application

import (
	"deploy-lander/internal/structs"
	"github.com/pkg/errors"
)

func (app *App) RegisterBuild(alias string, version structs.BuildVersioner) error {
	var project structs.Project
	var found bool
	if project, found = app.projects[alias]; !found {
		return ErrUnknownAlias
	}

	if !version.IsTag(project.TagsRe) && !version.IsBranch(project.BranchesRe) {
		return errors.Errorf("supplied version %v doesn't match project %v", version.Value(), alias)
	}

	return app.builds.RegisterBuild(version)
}
