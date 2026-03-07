// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const assetPackageURL = "https://d1.awsstatic.com/onedam/marketing-channels/website/aws/en_US/architecture/approved/architecture-icons/Asset-Package_07312025.49d3aab7f9e6131e51ade8f7c6c8b961ee7d3bb1.zip"
const assetPackageCacheFile = "aws-asset-package-07312025.zip"

// dacTypeToSVGPath mapeia tipos DAC para o caminho do SVG dentro do Asset Package zip.
var dacTypeToSVGPath = map[string]string{
	// ── Compute ───────────────────────────────────────────────────────────
	"AWS::EC2::Instance":        "Architecture-Service-Icons_07312025/Arch_Compute/64/Arch_Amazon-EC2_64.svg",
	"AWS::Lambda::Function":     "Architecture-Service-Icons_07312025/Arch_Compute/64/Arch_AWS-Lambda_64.svg",
	"AWS::ECS::Cluster":         "Architecture-Service-Icons_07312025/Arch_Containers/64/Arch_Amazon-Elastic-Container-Service_64.svg",
	"AWS::ECS::Service":         "Architecture-Service-Icons_07312025/Arch_Containers/64/Arch_Amazon-Elastic-Container-Service_64.svg",
	"AWS::ECS::TaskDefinition":  "Architecture-Service-Icons_07312025/Arch_Containers/64/Arch_Amazon-Elastic-Container-Service_64.svg",
	"AWS::ECR::Repository":      "Architecture-Service-Icons_07312025/Arch_Containers/64/Arch_Amazon-Elastic-Container-Registry_64.svg",
	"AWS::EKS::Cluster":         "Architecture-Service-Icons_07312025/Arch_Containers/64/Arch_Amazon-Elastic-Kubernetes-Service_64.svg",

	// ── Storage ───────────────────────────────────────────────────────────
	"AWS::S3::Bucket": "Architecture-Service-Icons_07312025/Arch_Storage/64/Arch_Amazon-Simple-Storage-Service_64.svg",

	// ── Database ──────────────────────────────────────────────────────────
	"AWS::RDS::DBInstance":           "Architecture-Service-Icons_07312025/Arch_Database/64/Arch_Amazon-RDS_64.svg",
	"AWS::RDS::DBCluster":            "Architecture-Service-Icons_07312025/Arch_Database/64/Arch_Amazon-RDS_64.svg",
	"AWS::DynamoDB::Table":           "Architecture-Service-Icons_07312025/Arch_Database/64/Arch_Amazon-DynamoDB_64.svg",
	"AWS::ElastiCache::CacheCluster": "Architecture-Service-Icons_07312025/Arch_Database/64/Arch_Amazon-ElastiCache_64.svg",
	"AWS::ElastiCache::ReplicationGroup": "Architecture-Service-Icons_07312025/Arch_Database/64/Arch_Amazon-ElastiCache_64.svg",

	// ── Networking & CDN ──────────────────────────────────────────────────
	"AWS::ElasticLoadBalancingV2::LoadBalancer": "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Elastic-Load-Balancing_64.svg",
	"AWS::ElasticLoadBalancing::LoadBalancer":   "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Elastic-Load-Balancing_64.svg",
	"AWS::CloudFront::Distribution":             "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Amazon-CloudFront_64.svg",
	"AWS::CloudFront":                           "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Amazon-CloudFront_64.svg",
	"AWS::ApiGateway::RestApi":                  "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Amazon-API-Gateway_64.svg",
	"AWS::ApiGateway":                           "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Amazon-API-Gateway_64.svg",
	"AWS::ApiGatewayV2::Api":                    "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Amazon-API-Gateway_64.svg",
	"AWS::EC2::InternetGateway":                 "Resource-Icons_07312025/Res_Networking-Content-Delivery/Res_Amazon-VPC_Internet-Gateway_48.svg",
	"AWS::EC2::NatGateway":                      "Resource-Icons_07312025/Res_Networking-Content-Delivery/Res_Amazon-VPC_NAT-Gateway_48.svg",
	"AWS::Route53::HostedZone":                  "Architecture-Service-Icons_07312025/Arch_Networking-Content-Delivery/64/Arch_Amazon-Route-53_64.svg",

	// ── App Integration / Messaging ───────────────────────────────────────
	"AWS::SNS::Topic":  "Architecture-Service-Icons_07312025/Arch_App-Integration/64/Arch_Amazon-Simple-Notification-Service_64.svg",
	"AWS::SQS::Queue":  "Architecture-Service-Icons_07312025/Arch_App-Integration/64/Arch_Amazon-Simple-Queue-Service_64.svg",
	"AWS::SFN::StateMachine": "Architecture-Service-Icons_07312025/Arch_App-Integration/64/Arch_AWS-Step-Functions_64.svg",

	// ── General ───────────────────────────────────────────────────────────
	"AWS::Diagram::Resource": "Resource-Icons_07312025/Res_General-Icons/Res_48_Light/Res_User_48_Light.svg",
}

var (
	assetZipOnce  sync.Once
	assetZipBytes []byte
	assetZipErr   error
)

// cacheDir retorna o diretório de cache do awsdac.
func cacheDir() string {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	return filepath.Join(base, "awsdac")
}

// loadAssetPackage baixa e cacheia o Asset Package zip na primeira chamada.
func loadAssetPackage() ([]byte, error) {
	assetZipOnce.Do(func() {
		cachePath := filepath.Join(cacheDir(), assetPackageCacheFile)

		// Usa cache local se já existe
		if data, err := os.ReadFile(cachePath); err == nil {
			log.Infof("drawio: using cached asset package: %s", cachePath)
			assetZipBytes = data
			return
		}

		log.Infof("drawio: downloading AWS Asset Package (~14MB)...")
		resp, err := http.Get(assetPackageURL)
		if err != nil {
			assetZipErr = fmt.Errorf("failed to download asset package: %w", err)
			return
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			assetZipErr = fmt.Errorf("failed to read asset package: %w", err)
			return
		}

		// Salva no cache
		if err := os.MkdirAll(cacheDir(), 0755); err == nil {
			if err := os.WriteFile(cachePath, data, 0644); err != nil {
				log.Warnf("drawio: could not cache asset package: %v", err)
			}
		}

		assetZipBytes = data
	})
	return assetZipBytes, assetZipErr
}

// extractSVGFromZip extrai o conteúdo de um arquivo SVG do zip.
func extractSVGFromZip(zipData []byte, svgPath string) (string, error) {
	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return "", err
	}
	for _, f := range r.File {
		if f.Name == svgPath {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			data, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}
	}
	return "", fmt.Errorf("SVG not found in zip: %s", svgPath)
}

// svgToDataURI converte conteúdo SVG para data URI usando URL-encoding (sem base64).
// Evita o problema do ';' no style do draw.io que ocorre com 'data:image/svg+xml;base64,...'
func svgToDataURI(svgContent string) string {
	// Encoding mínimo para SVG em data URI: substitui chars especiais
	// Troca aspas duplas por simples para evitar encoding (SVG aceita ambas)
	s := strings.ReplaceAll(svgContent, `"`, `'`)
	// Encoding de caracteres obrigatórios
	s = strings.ReplaceAll(s, "%", "%25") // deve ser o primeiro!
	s = strings.ReplaceAll(s, "#", "%23")
	s = strings.ReplaceAll(s, "<", "%3C")
	s = strings.ReplaceAll(s, ">", "%3E")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", " ")
	// Colapsa múltiplos espaços
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return "data:image/svg+xml," + s
}

// svgToBase64DataURI converte SVG para data URI base64.
func svgToBase64DataURI(svgContent string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(svgContent))
	return "data:image/svg+xml;base64," + b64
}

// resolveSVGPath resolve o caminho SVG no zip para um tipo DAC.
// Retorna "" se não houver mapeamento.
func resolveSVGPath(dacType string) string {
	if path, ok := dacTypeToSVGPath[dacType]; ok {
		return path
	}
	// fallback para tipo de serviço (ex: AWS::ApiGateway::X → AWS::ApiGateway)
	parts := strings.SplitN(dacType, "::", 3)
	if len(parts) >= 2 {
		if path, ok := dacTypeToSVGPath[strings.Join(parts[:2], "::")]; ok {
			return path
		}
	}
	return ""
}

// getAWSIconSVGContent extrai e retorna o conteúdo SVG bruto do Asset Package.
func getAWSIconSVGContent(dacType string) string {
	svgPath := resolveSVGPath(dacType)
	if svgPath == "" {
		return ""
	}
	zipData, err := loadAssetPackage()
	if err != nil {
		log.Warnf("assets: package unavailable: %v", err)
		return ""
	}
	svgContent, err := extractSVGFromZip(zipData, svgPath)
	if err != nil {
		log.Warnf("assets: SVG not found for %s: %v", dacType, err)
		return ""
	}
	return svgContent
}

// GetAWSIconSVG retorna o conteúdo SVG bruto do Asset Package para um tipo DAC.
// Usado pelo pipeline PNG para renderizar o ícone via oksvg.
func GetAWSIconSVG(dacType string) string {
	return getAWSIconSVGContent(dacType)
}

// GetAWSIconDataURI retorna o data URI do SVG oficial AWS para uso no draw.io.
// Usa URL-encoding para evitar o ';' no style do draw.io (data:image/svg+xml;base64,...).
func GetAWSIconDataURI(dacType string) string {
	svgContent := getAWSIconSVGContent(dacType)
	if svgContent == "" {
		return ""
	}
	return svgToDataURI(svgContent)
}
