package alert

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"text/template"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/log"
)

type Alert interface {
	Create(ctx context.Context) error
	Delete(ctx context.Context) error
	Get(ctx context.Context) error
	List(ctx context.Context) error
	Update(ctx context.Context) error
}

type alert struct {
	cfg         *config.AlertConfig
	atomicToken *atomic.Value
}

func New(cfg *config.AlertConfig, atomicToken *atomic.Value) Alert {
	return &alert{
		cfg:         cfg,
		atomicToken: atomicToken,
	}
}

type AlertData struct {
	ID               int64
	Name             string
	TeamName         string
	ApplicationName  string
	ReceiverTeamName string
	Expression       string
	Period           string
	Severity         string
	Summary          string
	Description      string
	DashboardURL     string
	PlaybookURL      string
	Status           string
}

func (a *alert) Create(ctx context.Context) error {
	alertData := AlertData{
		Name:             "arief-harshil-test-alert",
		TeamName:         "harshil_test",
		ApplicationName:  "test_harshil_application",
		ReceiverTeamName: "harshil_test",
		Expression:       `sum(rate(container_cpu_usage_seconds_total{container_name=\"harshil\"}[5m])) by (pod_name) > 0.5`,
		Period:           "5m",
		Severity:         "critical",
		Summary:          "harshil is using too much CPU",
		Description:      "harshil is using too much CPU",
		DashboardURL:     "https://grafana.caraml.io/d/000000001/",
		PlaybookURL:      "this one is empty",
		Status:           "enabled",
	}

	print(a.cfg.CreateBodyTemplate)
	log.Infof("a.cfg.CreateBodyTemplate %+v", a.cfg.CreateBodyTemplate)

	tmpl, err := template.New("tmpl").Parse(a.cfg.CreateBodyTemplate)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, &alertData)
	if err != nil {
		return err
	}

	print(b.String())

	req, err := http.NewRequest("POST", a.cfg.Host+a.cfg.CreateEndpoint, bytes.NewBuffer(b.Bytes()))
	if err != nil {
		log.Errorf("error fetching alerts: %s", err)
		return err
	}

	token := a.atomicToken.Load().(string)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Host", a.cfg.Host)
	req.Header.Set("Content-Type", a.cfg.CreateContentType)

	if a.cfg.CreateAdditionalHeaders != nil {
		for k, v := range a.cfg.CreateAdditionalHeaders {
			req.Header.Set(k, v)
		}
	}

	log.Infof("req: %+v", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error creating alerts: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading alerts: %s", err)
		return err
	}
	log.Infof("alert.Create response body: %s", string(body))
	return nil
}

func (a *alert) Delete(ctx context.Context) error {
	return nil
}

func (a *alert) Get(ctx context.Context) error {
	return nil
}

func (a *alert) List(ctx context.Context) error {
	req, err := http.NewRequest("GET", a.cfg.Host+a.cfg.ListEndpoint+"?team_name=harshil_test&application_name=test_harshil_application", nil)
	if err != nil {
		log.Errorf("error fetching alerts: %s", err)
		return err
	}

	token := a.atomicToken.Load().(string)
	log.Infof("url: %s", a.cfg.Host+a.cfg.ListEndpoint)
	log.Infof("token %s", token)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Host", a.cfg.Host)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error fetching alerts: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading alerts: %s", err)
		return err
	}
	log.Infof("alert.List %s", string(body))
	return nil
}

func (a *alert) Update(ctx context.Context) error {
	alertData := AlertData{
		ID:               1646,
		Name:             "arief-harshil-test-alert",
		TeamName:         "harshil_test",
		ApplicationName:  "test_harshil_application",
		ReceiverTeamName: "harshil_test",
		Expression:       `sum(rate(container_cpu_usage_seconds_total{container_name=\"harshil\"}[5m])) by (pod_name) > 0.5`,
		Period:           "5m",
		Severity:         "critical",
		Summary:          "harshil is using too much CPU",
		Description:      "harshil is using too much CPU",
		DashboardURL:     "https://grafana.caraml.io/d/000000001/",
		PlaybookURL:      "this one is empty",
		Status:           "enabled",
	}

	print(a.cfg.UpdateBodyTemplate)
	log.Infof("a.cfg.UpdateBodyTemplate %+v", a.cfg.UpdateBodyTemplate)

	tmpl, err := template.New("tmpl").Parse(a.cfg.UpdateBodyTemplate)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, &alertData)
	if err != nil {
		return err
	}

	print(b.String())

	req, err := http.NewRequest(a.cfg.UpdateMethod, fmt.Sprintf("%s%s/%d", a.cfg.Host, a.cfg.UpdateEndpoint, alertData.ID), bytes.NewBuffer(b.Bytes()))
	if err != nil {
		log.Errorf("error fetching alerts: %s", err)
		return err
	}

	token := a.atomicToken.Load().(string)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Host", a.cfg.Host)
	req.Header.Set("Content-Type", a.cfg.UpdateContentType)

	if a.cfg.UpdateAdditionalHeaders != nil {
		for k, v := range a.cfg.UpdateAdditionalHeaders {
			req.Header.Set(k, v)
		}
	}

	log.Infof("req: %+v", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error creating alerts: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading alerts: %s", err)
		return err
	}
	log.Infof("alert.Update response body: %s", string(body))
	return nil
}
