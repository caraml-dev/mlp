package delete

import (
	"strings"

	"github.com/gojek/mlp/api/pkg/client/mlflow"
	"github.com/gojek/mlp/api/pkg/gcs"
)

type deleteClient struct {
	Client     mlflow.Mlflow
	GcsPackage gcs.GcsPackage
}

type DeletePackage interface {
	DeleteExperiment(trackingURL string, idExperiment string, deleteArtifact bool)
	DeleteRun(trackingURL string, idRun string, delArtifact bool)
}

func NewDeleteClient(mlfclient mlflow.Mlflow, gcspkg gcs.GcsPackage) *deleteClient {
	return &deleteClient{
		Client:     mlfclient,
		GcsPackage: gcspkg,
	}
}

func (dc *deleteClient) DeleteExperiment(idExperiment string) error {

	err := dc.Client.DeleteExperiment(idExperiment)
	if err != nil {
		return err
	}

	relatedRunId, err := dc.Client.SearchRunForExperiment(idExperiment)
	if err != nil {
		return err
	}

	var deletedRunId []string
	var failDeletedRunId []string
	for _, run := range relatedRunId.RunsData {
		err = dc.DeleteRun(run.Info.RunId, false)
		if err != nil {
			failDeletedRunId = append(failDeletedRunId, run.Info.RunId)
			// return err
		} else {
			deletedRunId = append(deletedRunId, run.Info.RunId)
		}
	}

	if len(relatedRunId.RunsData) > 0 {
		path := relatedRunId.RunsData[0].Info.ArtifactURI[5:]
		splitPath := strings.SplitN(path, "/", 4)
		folderPath := strings.Join(splitPath[0:3], "/")
		// deleting folder
		err = dc.GcsPackage.DeleteArtifact(folderPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dc *deleteClient) DeleteRun(idRun string, delArtifact bool) error {

	err := dc.Client.DeleteRun(idRun)
	if err != nil {
		return err
	}
	if delArtifact {
		runDetail, err := dc.Client.SearchRunData(idRun)
		if err != nil {
			return err
		}

		err = dc.GcsPackage.DeleteArtifact(runDetail.RunData.Info.ArtifactURI[5:])
		if err != nil {
			return err
		}

	}
	return nil
}
