package integration

import (
	bssms "bssms/bssms"
	"bssms/internal/installer"
	"bssms/internal/proxy"
	provisioner "bssms/provisioner"
	"context"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	ProxyHost = "localhost"
)

func ProxyAddress(proxyPort string) string { return fmt.Sprintf("%s:%s", ProxyHost, proxyPort) }

func GetTestDataDir() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("path discovery")
	}
	return filepath.Join(filepath.Dir(thisFile), "testdata")
}

func GetIhs(dataFile string) ([]provisioner.InstallHost, error) {
	return provisioner.LoadYamlInstallHosts(filepath.Join(GetTestDataDir(), dataFile))
}

func RunTestProxy(proxyPort string) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	config := bssms.ProxyConfig{
		BaseConfig: bssms.BaseConfig{Ctx: ctx},
		UnsafeTls:  true,
		ListenAddr: fmt.Sprintf(":%s", proxyPort),
		Host:       ProxyHost,
	}
	var bgErr error
	go func() {
		err := proxy.RunProxy(&config)
		if err != nil {
			bgErr = err
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if bgErr != nil {
		cancel()
		return nil, bgErr
	}
	return cancel, nil
}

func RunTestProvisioner(dataFile string, proxyPort string) (context.CancelFunc, error) {
	ihs, err := GetIhs(dataFile)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	var bgErr error
	go func() {
		pihs := provisioner.PInstallHosts(ihs)
		err := provisioner.RunIhs(
			&bssms.ProvisionerConfig{
				BaseConfig:   bssms.BaseConfig{Ctx: ctx},
				UnsafeTls:    true,
				ProxyAddress: ProxyAddress(proxyPort),
			},
			pihs,
		)
		if err != nil {
			bgErr = err
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if bgErr != nil {
		cancel()
		return nil, bgErr
	}
	return cancel, nil
}

func RunTestTofu(td string, index int, proxyPort string) (context.CancelFunc, error) {
	ihs, err := provisioner.LoadInstallHosts(td)
	if err != nil {
		return nil, err
	}
	if index >= len(ihs) {
		return nil, fmt.Errorf("index out of range")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	var bgErr error
	go func() {
		err := provisioner.TofuRun(
			&bssms.ProvisionerConfig{
				BaseConfig:   bssms.BaseConfig{Ctx: ctx},
				UnsafeTls:    true,
				ProxyAddress: ProxyAddress(proxyPort),
			},
			td,
			ihs[index],
		)
		if err != nil {
			bgErr = err
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if bgErr != nil {
		cancel()
		return nil, bgErr
	}
	return cancel, nil
}

func RunTestInstaller(t *testing.T, dataFile string, index int, proxyPort string) (context.CancelFunc, error) {
	ihs, err := GetIhs(dataFile)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	var bgErr error
	go func() {
		err := installer.Run(&bssms.InstallerConfig{
			BaseConfig:   bssms.BaseConfig{Ctx: ctx},
			UnsafeTls:    true,
			ProxyAddress: ProxyAddress(proxyPort),
			PrivateKey:   ihs[index].PrivateKey,
			JsonSecf:     path.Join(t.TempDir(), fmt.Sprintf("secrets-%s.json", ihs[index].Name)),
			Installable: bssms.Installable{
				ServerUuid: ihs[index].ServerUuid,
				MacAddress: ihs[index].MacExtAddress,
				IPAddress:  ihs[index].IPExtAddress,
			},
		})
		if err != nil {
			bgErr = err
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if bgErr != nil {
		cancel()
		return nil, bgErr
	}
	return cancel, nil
}
