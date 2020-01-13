package usage

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	versionutil "github.com/jfrog/jfrog-client-go/utils/version"
	"github.com/pkg/errors"
	"net/http"
)

const minArtifactoryVersion = "6.9.0"

func SendReportUsage(productId, commandName string, serviceManager *artifactory.ArtifactoryServicesManager) error {
	config := serviceManager.GetConfig()
	if config == nil {
		return errorutils.NewError("Expected full config, but no configuration exists.")
	}
	rtDetails := config.GetArtDetails()
	if rtDetails == nil {
		return errorutils.NewError("Artifactory details not configured.")
	}
	url, err := utils.BuildArtifactoryUrl(rtDetails.GetUrl(), "api/system/usage", make(map[string]string))
	if err != nil {
		return err
	}
	clientDetails := rtDetails.CreateHttpClientDetails()
	// Check Artifactory version
	artifactoryVersion, err := rtDetails.GetVersion()
	if err != nil {
		return err
	}
	if !isVersionCompatible(artifactoryVersion) {
		log.Debug(fmt.Sprintf("Expected Artifactory version %s or above, got %s", minArtifactoryVersion, artifactoryVersion))
		return nil
	}

	bodyContent, err := reportUsageToJson(productId, commandName)
	if err != nil {
		return err
	}
	utils.AddHeader("Content-Type", "application/json", &clientDetails.Headers)
	resp, body, err := serviceManager.Client().SendPost(url, bodyContent, &clientDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Debug("Usage info sent successfully.")
	return nil
}

// Returns an error if the Artifactory version is not compatible
func isVersionCompatible(artifactoryVersion string) bool {
	// API exists from Artifactory version 6.9.0 and above:
	version := versionutil.NewVersion(artifactoryVersion)
	return version.AtLeast(minArtifactoryVersion)
}

func reportUsageToJson(productId, commandName string) ([]byte, error) {
	featureInfo := feature{FeatureId: commandName}
	params := reportUsageParams{ProductId: productId, Features: []feature{featureInfo}}
	bodyContent, err := json.Marshal(params)
	return bodyContent, errorutils.WrapError(err)
}

type reportUsageParams struct {
	ProductId string    `json:"productId"`
	Features  []feature `json:"features,omitempty"`
}

type feature struct {
	FeatureId string `json:"featureId,omitempty"`
}
