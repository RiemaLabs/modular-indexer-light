package configs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/RiemaLabs/modular-indexer-committee/checkpoint"

	"github.com/RiemaLabs/modular-indexer-light/internal/constant"
	"github.com/RiemaLabs/modular-indexer-light/internal/logs"
	"github.com/RiemaLabs/modular-indexer-light/internal/utils"
	"github.com/RiemaLabs/modular-indexer-light/internal/wallet"
)

type (
	Config struct {
		ListenAddr string `json:"listenAddr"`

		CommitteeIndexers CommitteeIndexers `json:"committeeIndexers"`

		Verification Verification `json:"verification"`

		Report Report `json:"report"`
	}

	CommitteeIndexers struct {
		S3 []SourceS3 `json:"s3"`
		DA []SourceDA `json:"da"`
	}

	Verification struct {
		BitcoinRPC        string `json:"bitcoinRPC"`
		MinimalCheckpoint int    `json:"minimalCheckpoint"`
		MetaProtocol      string `json:"metaProtocol"`
	}

	Report struct {
		Name        string `json:"name"`
		Network     string `json:"network"`
		NamespaceID string `json:"namespaceID"`
		GasCoupon   string `json:"gasCoupon"`
		Timeout     int    `json:"timeout"`

		// PrivateKey loaded from files.
		PrivateKey string `json:"-"`
	}
)

type (
	DenyList struct {
		Evidence *Evidence `json:"evidence"`
		SourceS3 *SourceS3 `json:"sourceS3"`
		SourceDA *SourceDA `json:"sourceDa"`
	}

	Evidence struct {
		Height            uint   `json:"height"`
		Hash              string `json:"hash"`
		CorrectCommitment string `json:"correctCommitment"`
		FraudCommitment   string `json:"fraudCommitment"`
	}
)

type (
	SourceS3 struct {
		Region string `json:"region"`
		Bucket string `json:"bucket"`
		Name   string `json:"name"`
	}

	SourceDA struct {
		Network     string `json:"network"`
		NamespaceID string `json:"namespaceID"`
		Name        string `json:"name"`
	}
)

func (s *SourceS3) Equal(rhs *SourceS3) bool {
	return s.Region == rhs.Region && s.Bucket == rhs.Bucket && s.Name == rhs.Name
}

func (s *SourceDA) Equal(rhs *SourceDA) bool {
	return s.Network == rhs.Network && s.NamespaceID == rhs.NamespaceID && s.Name == rhs.Name
}

type CheckpointExport struct {
	Checkpoint *checkpoint.Checkpoint `json:"checkPoint"`
	SourceS3   *SourceS3              `json:"sourceS3,omitempty"`
	SourceDA   *SourceDA              `json:"sourceDa,omitempty"`
}

var C *Config

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func ReadDenyList(path string) ([]*DenyList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer func() { _ = f.Close() }()

	var ret []*DenyList
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println()
		data := []byte(strings.TrimSpace(scanner.Text()))
		var item DenyList
		if err := json.Unmarshal(data, &item); err != nil {
			logs.Warn.Printf("Parse deny list error: loaded=%d, err=%v", len(ret), err)
			return ret, nil
		}
		ret = append(ret, &item)
	}
	if err := scanner.Err(); err != nil {
		logs.Warn.Printf("Scan deny list error: loaded=%d, err=%v", len(ret), err)
	}

	return ret, nil
}

func (r *Report) LoadPrivate(path string) error {
	if data, err := os.ReadFile(path); err == nil {
		r.PrivateKey = string(data)
		return nil
	}

	logs.Info.Printf("Failed to read the private key from the local directory, generating a new one...")
	pwd := constant.DefaultPassword

	wall := wallet.NewWallet(&pwd)
	if !wall.GenerateBip39Seed(&pwd, &pwd) {
		return errors.New("failed to generate BIP39 seed")
	}
	account := wall.GenerateAccount(&pwd)

	key := utils.EcdsaToPrivateStr(account.PrivateKey(&pwd))
	if err := os.WriteFile(path, []byte(key), 0644); err != nil {
		return fmt.Errorf("write private key to file error: %v", err)
	}

	logs.Info.Printf("Store your private file %q carefully and don't share it!", path)
	r.PrivateKey = key

	return nil
}

func Init(configPath, denyListPath string) error {
	c, err := ReadConfig(configPath)
	if err != nil {
		return err
	}

	denials, err := ReadDenyList(denyListPath)
	if err != nil {
		return err
	}
	for _, b := range denials {
		if d := b.SourceDA; d != nil {
			c.CommitteeIndexers.DA = slices.DeleteFunc(c.CommitteeIndexers.DA, func(s SourceDA) bool { return s.Equal(d) })
		}
		if d := b.SourceS3; d != nil {
			c.CommitteeIndexers.S3 = slices.DeleteFunc(c.CommitteeIndexers.S3, func(s SourceS3) bool { return s.Equal(d) })
		}
	}

	C = c
	return nil
}

func AppendDenyList(path string, item *DenyList) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open error: %v", err)
	}
	defer func() { _ = f.Close() }()

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	_, err = f.WriteString(string(data) + "\n")
	return err
}
