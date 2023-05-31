package fs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"go.keploy.io/server/pkg"
	"go.keploy.io/server/pkg/models"
	"go.keploy.io/server/pkg/persistence"
)

type testReport struct {
	// Map to hold all the test results.
	// It is keyed on run ID.
	results map[string][]models.TestResult

	// Mutex to guard shared access of results map.
	m sync.Mutex

	// The native filesystem for test reports.
	native persistence.Filesystem
}

func NewTestReportFS(native persistence.Filesystem) *testReport {
	return &testReport{
		results: map[string][]models.TestResult{},
		m:       sync.Mutex{},
		native:  native,
	}
}

func (tr *testReport) SetResult(runId string, test models.TestResult) {
	// TODO: send runId to the historyConfig
	tr.m.Lock()
	defer tr.m.Unlock()

	tr.results[runId] = append(tr.results[runId], test)
	results, _ := tr.results[runId]
	results = append(results, test)
	tr.results[runId] = results
}

func (tr *testReport) GetResults(runId string) ([]models.TestResult, error) {
	tr.m.Lock()
	defer tr.m.Unlock()

	results, ok := tr.results[runId]
	if !ok {
		return nil, fmt.Errorf("found no test results for test report with id: %s", runId)
	}
	return results, nil
}

func (tr *testReport) Read(ctx context.Context, path, name string) (models.TestReport, error) {
	if !pkg.IsValidPath(path) {
		return models.TestReport{},
			fmt.Errorf("file path should be absolute. got test report path: %s "+
				"and its name: %s", pkg.SanitiseInput(path), pkg.SanitiseInput(name))
	}
	if strings.Contains(name, "/") || !pkg.IsValidPath(name) {
		return models.TestReport{},
			errors.New("invalid name for test-report. It should not include any slashes")
	}
	file, err := tr.native.OpenFile(filepath.Join(path, name+".yaml"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return models.TestReport{}, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var doc models.TestReport
	err = decoder.Decode(&doc)
	if err != nil {
		return models.TestReport{},
			fmt.Errorf("failed to decode the yaml file documents. error: %v", err.Error())
	}
	return doc, nil
}

func (tr *testReport) Write(ctx context.Context, path string, doc models.TestReport) error {
	if strings.Contains(doc.Name, "/") || !pkg.IsValidPath(doc.Name) {
		return errors.New("invalid name for test-report. It should not include any slashes")
	}

	_, err := tr.native.CreateYamlFile(path, doc.Name)
	if err != nil {
		return err
	}

	var data []byte
	d, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document to yaml. error: %s", err.Error())
	}
	data = append(data, d...)

	err = tr.native.WriteFile(filepath.Join(path, doc.Name+".yaml"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write test report in yaml file. error: %s", err.Error())
	}
	return nil
}
