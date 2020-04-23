package reader

import (
	"encoding/json"

	"github.com/google/go-github/v31/github"
	"github.com/knqyf263/labeler/config"
	"github.com/knqyf263/labeler/logs"
	"github.com/knqyf263/labeler/remote"
	"github.com/knqyf263/labeler/types"
	yaml "gopkg.in/yaml.v2"
)

// Run executes the write actions against the repo.
func Run(client *github.Client, opt *types.Options) error {
	file := opt.Filename

	// TODO:
	// DryRun should cleanup if missing as well...
	err := config.CreateIfMissing(file)
	if err != nil {
		return err
	}

	lf, err := config.ReadFile(file)
	if err != nil {
		return err
	}

	opt.Repo, err = config.GetRepo(opt, lf)
	if err != nil {
		logs.V(0).Infof("No repo provided")
		return err
	}

	err = opt.ValidateRepo()
	if err != nil {
		logs.V(0).Infof("Failed to parse repo format: owner/name")
		return err
	}

	// Get all remote labels from repo
	labelsRemote, err := remote.GetLabels(client, opt)
	if err != nil {
		return err
	}

	total := len(labelsRemote)

	x, err := json.Marshal(labelsRemote)
	if err != nil {
		logs.V(0).Infof("Failed to marshal labels from remote format")
		return err
	}

	labels := []*types.Label{}

	// TODO:
	// Can we directly unmarshal from labelsRemote?
	err = yaml.Unmarshal(x, &labels)
	if err != nil {
		logs.V(0).Infof("Failed to unmarshal labels to local format")
		return err
	}

	for _, l := range labels {
		logs.V(4).Infof("Fetched '%s' with color '%s'", l.Name, l.Color)
	}

	lf = &types.LabelFile{
		Repo:   opt.Repo,
		Labels: labels,
	}

	if !opt.DryRun {
		err = config.WriteFile(file, lf)
		if err != nil {
			return err
		}
	}

	logs.V(4).Infof("Processed %d labels in total", total)

	return nil
}
