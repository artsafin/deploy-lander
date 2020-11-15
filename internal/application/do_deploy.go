package application

import (
	"bytes"
	"deploy-lander/internal/structs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os/exec"
)

var ErrUnknownAlias = errors.New("unkown alias")

func (app *App) DoDeploy(alias string, version structs.BuildVersioner) error {
	var project structs.Project
	var found bool
	if project, found = app.projects[alias]; !found {
		return ErrUnknownAlias
	}

	if !app.builds.IsRegistered(version) {
		return errors.Errorf("version %v is not registered", version.Value())
	}

	gitErr, out := app.updateProjectToVersion(project, version)

	if gitErr != nil {
		app.logger.Error("git error", zap.Error(gitErr), zap.ByteString("output", out.Bytes()))
		return gitErr
	}

	cmdErr, cmdOut := app.execDeployCommand(project)
	if cmdErr != nil {
		app.logger.Error("deploy command error", zap.Error(cmdErr), zap.ByteString("output", cmdOut))
		return cmdErr
	}

	return app.builds.MarkDeployed(version)
}

func (app *App) execDeployCommand(project structs.Project) (Err error, Out []byte) {
	cmd := exec.Command("/bin/sh", "-c", project.DeployCommand)
	cmd.Env = make([]string, 0)

	app.logger.Info("executing command", zap.String("cmd", cmd.String()))

	Out, Err = cmd.CombinedOutput()
	return
}

func (app *App) updateProjectToVersion(project structs.Project, version structs.BuildVersioner) (Err error, Out *bytes.Buffer) {
	Out = &bytes.Buffer{}
	app.logger.Info("Opening repository", zap.String("path", project.Path))

	repo, openerr := git.PlainOpen(project.Path)
	if openerr != nil {
		return openerr, Out
	}

	fetchOptions := git.FetchOptions{
		RemoteName: "origin",
		Tags:       git.AllTags,
		Progress:   Out,
	}
	fetchErr := repo.Fetch(&fetchOptions)
	if fetchErr != nil {
		return fetchErr, Out
	}
	app.logger.Info("Fetched successfully", zap.String("version", version.Value()))

	var ref *plumbing.Reference
	var refErr error

	app.logger.Info("Resolving version", zap.String("version", version.Value()))

	if version.IsTag(project.TagsRe) {
		ref, refErr = repo.Tag(version.Value())
	}

	if version.IsBranch(project.BranchesRe) {
		br, brancherr := repo.Branch(version.Value())
		if brancherr != nil {
			return brancherr, Out
		}
		ref, refErr = repo.Reference(br.Merge, true)
	}

	if ref == nil && refErr == nil {
		refErr = errors.Errorf("couldn't match %v against %v or %v", version.Value(), project.TagsRe, project.BranchesRe)
	}

	if refErr != nil {
		return errors.Errorf("supplied version %v doesn't match project %v: %v", version.Value(), project.Alias, refErr), Out
	}

	app.logger.Info("Resolved ref", zap.String("version", version.Value()), zap.String("hash", ref.String()))

	wt, _ := repo.Worktree()
	resetOptions := git.ResetOptions{
		Commit: ref.Hash(),
		Mode:   git.HardReset,
	}
	reseterr := wt.Reset(&resetOptions)
	if reseterr != nil {
		return reseterr, Out
	}

	return nil, Out
}
