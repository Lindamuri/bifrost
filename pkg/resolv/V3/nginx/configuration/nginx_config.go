package configuration

import (
	"bytes"
	"encoding/json"
	"github.com/tremendouscan/bifrost/internal/pkg/code"
	"github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	utilsV3 "github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
	"strings"
	"sync"
	"time"
)

type NginxConfig interface {
	Main() local.MainContext
	UpdateFromJsonBytes(data []byte) error
	UpdatedTimestamp() time.Time
	TextLines() []string
	Json() []byte
	Dump() map[string]*bytes.Buffer
}

type nginxConfig struct {
	mainContext local.MainContext
	timestamp   time.Time
	rwLocker    *sync.RWMutex
}

func (n *nginxConfig) Main() local.MainContext {
	n.rwLocker.RLock()
	defer n.rwLocker.RUnlock()
	return n.mainContext
}

func (n *nginxConfig) UpdateFromJsonBytes(data []byte) error {
	m, err := local.JsonLoader(data).Load()
	if err != nil {
		return err
	}
	// write lock operation in renewMainContext()
	return n.renewMainContext(m)
}

func (n *nginxConfig) UpdatedTimestamp() time.Time {
	n.rwLocker.RLock()
	defer n.rwLocker.RUnlock()
	return n.timestamp
}

func (n *nginxConfig) TextLines() []string {
	n.rwLocker.RLock()
	defer n.rwLocker.RUnlock()
	lines, _ := n.mainContext.ConfigLines(false)
	return lines
}

func (n *nginxConfig) Json() []byte {
	n.rwLocker.RLock()
	defer n.rwLocker.RUnlock()
	data, err := json.Marshal(n.mainContext)
	if err != nil {
		return nil
	}
	return data
}

func (n *nginxConfig) Dump() map[string]*bytes.Buffer {
	n.rwLocker.RLock()
	defer n.rwLocker.RUnlock()
	return dumpMainContext(n.mainContext)
}

func (n *nginxConfig) renewMainContext(m local.MainContext) error {
	oldFP := utilsV3.NewConfigFingerprinterWithTimestamp(n.Dump(), time.Time{})
	newFP := utilsV3.NewConfigFingerprinterWithTimestamp(dumpMainContext(m), time.Time{})
	n.rwLocker.Lock()
	defer n.rwLocker.Unlock()
	if !oldFP.Diff(newFP.Fingerprints()) {
		return errors.WithCode(code.ErrSameConfigFingerprint, "same config fingerprint")
	}
	n.mainContext = m
	n.timestamp = time.Now()
	return nil
}

func NewNginxConfigFromJsonBytes(data []byte) (NginxConfig, error) {
	m, err := local.JsonLoader(data).Load()
	if err != nil {
		return nil, err
	}
	return newNginxConfig(m)
}

func NewNginxConfigFromFS(filepath string) (NginxConfig, error) {
	m, t, err := loadMainContextFromFS(filepath)
	if err != nil {
		logV1.Warnf("load nginx config failed: %w", err)
		return nil, err
	}
	return newNginxConfigWithTimestamp(m, t)
}

func loadMainContextFromFS(filepath string) (local.MainContext, time.Time, error) {
	timestamp := time.UnixMicro(0)
	m, err := local.FileLoader(filepath).Load()
	if err != nil {
		return nil, timestamp, err
	}
	for _, config := range m.ListConfigs() {
		tt, err := utilsV3.FileModifyTime(config.FullPath())
		if err != nil {
			return nil, timestamp, err
		}
		if tt.After(timestamp) {
			timestamp = *tt
		}
	}
	return m, timestamp, nil
}

func dumpMainContext(m local.MainContext) map[string]*bytes.Buffer {
	if m == nil {
		return nil
	}
	dumps := make(map[string]*bytes.Buffer)
	for _, config := range m.ListConfigs() {
		lines, err := config.ConfigLines(true)
		if err != nil {
			return nil
		}
		buff := bytes.NewBuffer([]byte{})
		for _, line := range lines {
			buff.WriteString(line + "\n")
		}
		dumps[strings.TrimSpace(config.FullPath())] = buff
	}
	return dumps
}

func newNginxConfigWithTimestamp(m local.MainContext, timestamp time.Time) (NginxConfig, error) {
	if m == nil {
		return nil, errors.WithCode(code.ErrV3InvalidContext, "new nginx config with a nil main context")
	}
	return &nginxConfig{
		mainContext: m,
		rwLocker:    new(sync.RWMutex),
		timestamp:   timestamp,
	}, nil
}

func newNginxConfig(m local.MainContext) (NginxConfig, error) {
	return newNginxConfigWithTimestamp(m, time.Now())
}
