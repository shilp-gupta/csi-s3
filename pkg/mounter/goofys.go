package mounter

import (
	"fmt"
	"path"

	"context"

	"github.com/ctrox/csi-s3/pkg/s3"
	goofysApi "github.com/kahing/goofys/api"
	goofysApiCommon "github.com/kahing/goofys/api/common"
)

const (
	goofysCmd     = "goofys"
	defaultRegion = "us-east-1"
)

// Implements Mounter
type goofysMounter struct {
	meta            *s3.FSMeta
	endpoint        string
	region          string
	accessKeyID     string
	secretAccessKey string
}

func newGoofysMounter(meta *s3.FSMeta, cfg *s3.Config) (Mounter, error) {
	region := cfg.Region
	// if endpoint is set we need a default region
	if region == "" && cfg.Endpoint != "" {
		region = defaultRegion
	}
	return &goofysMounter{
		meta:            meta,
		endpoint:        cfg.Endpoint,
		region:          region,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: cfg.SecretAccessKey,
	}, nil
}

func (goofys *goofysMounter) Stage(stageTarget string) error {
	return nil
}

func (goofys *goofysMounter) Unstage(stageTarget string) error {
	return nil
}

func (goofys *goofysMounter) Mount(source string, target string) error {
	goofysCfg := &goofysApiCommon.FlagStorage{
		MountOptions: map[string]string{
			"allow_other": "",
		},
		MountPoint: target,
		DirMode: 0o777,
		FileMode: 0o666,
		Endpoint: goofys.endpoint,
		Backend: &goofysApiCommon.S3Config{
			Region: goofys.region,
			AccessKey: goofys.accessKeyID,
			SecretKey: goofys.secretAccessKey,
			StorageClass: "STANDARD", // Need to figure out how to expose this as a parameter
		},
		StatCacheTTL: 0, // Need to figure out how to expose this as a parameter
		TypeCacheTTL: 0, // Need to figure out how to expose this as a parameter
		DebugFuse: false, // Need to figure out how to expose this as a parameter
		DebugS3: false, // Need to figure out how to expose this as a parameter
	}

	fullPath := fmt.Sprintf("%s:%s", goofys.meta.BucketName, path.Join(goofys.meta.Prefix, goofys.meta.FSPath))

	_, _, err := goofysApi.Mount(context.Background(), fullPath, goofysCfg)

	if err != nil {
		return fmt.Errorf("Error mounting via goofys: %s", err)
	}
	return nil
}
