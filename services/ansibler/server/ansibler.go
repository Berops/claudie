package main

import (
	"fmt"

	"github.com/Berops/platform/proto/pb"
	"github.com/Berops/platform/utils"
)

const (
	baseDirectory         = "services/ansibler/server"
	inventoryFile         = "inventory.ini"
	nodesInventoryFileTpl = "all-node-inventory.goini"
	outputDirectory       = "clusters"
	privateKeyExt         = "pem"
)

type NodepoolInfo struct {
	Nodepools  []*pb.NodePool
	PrivateKey string
	ID         string
}

type AllNodesInventoryData struct {
	NodepoolInfos []*NodepoolInfo
}

type LbInventoryData struct {
	K8sNodepools []*pb.NodePool
	LBClusters   []*pb.LBcluster
}

func generateInventoryFile(inventoryTemplate, directory string, data interface{}) error {
	templateLoader := utils.TemplateLoader{Directory: utils.AnsiblerTemplates}
	tpl, err := templateLoader.LoadTemplate(inventoryTemplate)
	if err != nil {
		return fmt.Errorf("error while loading template %s : %v", inventoryTemplate, err)
	}
	template := utils.Templates{Directory: directory}
	err = template.Generate(tpl, inventoryFile, data)
	if err != nil {
		return fmt.Errorf("error while generating from template %s : %v", inventoryTemplate, err)
	}
	return nil
}
