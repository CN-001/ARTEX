---
name: supply-chain-security
description: Comprehensive supply chain security framework covering SBOM generation and verification, software dependency and open source compliance audit, code signing and supply chain integrity, container image security scanning, third-party vendor risk assessment, and supply chain attack detection and response. Use when implementing supply chain security, conducting vendor risk assessments, or managing software dependencies.
---

# 供应链安全技能

## 概述

本技能提供系统化的供应链安全框架，涵盖SBOM生成与验证、软件依赖与开源合规审计、代码签名与供应链完整性、容器镜像安全扫描、第三方供应商风险评估和供应链攻击检测与响应等核心能力。

**核心能力:**
- SBOM - SBOM生成与验证
- 依赖审计 - 软件依赖与开源合规审计
- 代码签名 - 代码签名与供应链完整性
- 容器扫描 - 容器镜像安全扫描
- 供应商评估 - 第三方供应商风险评估
- 攻击检测 - 供应链攻击检测与响应

**适用场景:**
- 供应链安全实施
- 软件依赖管理
- 供应商风险评估
- 容器安全扫描
- 供应链攻击响应

## 快速开始

### 供应链安全实施流程

1. **资产发现** - 1-2 天
   - 发现供应链资产
   - 分析依赖关系
   - 识别关键组件

2. **风险评估** - 2-3 天
   - 评估供应链风险
   - 分析供应商风险
   - 检查依赖漏洞

3. **安全加固** - 3-5 天
   - 生成SBOM
   - 实施代码签名
   - 扫描容器镜像

4. **持续监控** - 持续进行
   - 监控供应链威胁
   - 响应供应链攻击
   - 更新安全策略

## 核心技能模块

### 1. SBOM生成与验证 (SBOM Generation and Verification)

SBOM生成与验证。

**核心工具:**
- Syft - SBOM生成工具
- CycloneDX - SBOM标准
- SPDX - SBOM标准
- Trivy - 容器安全扫描

**详细文档:** [SBOM生成与验证-SBOMGeneration.md](skills/SBOM生成与验证-SBOMGeneration.md)

### 2. 软件依赖与开源合规审计 (Software Dependency and Open Source Compliance Audit)

软件依赖与开源合规审计。

**核心工具:**
- FOSSA - 开源合规平台
- Black Duck - 开源合规工具
- Snyk - 依赖安全扫描
- OWASP DC - 依赖检查

**详细文档:** [软件依赖与开源合规审计-DependencyLicenseCompliance.md](skills/软件依赖与开源合规审计-DependencyLicenseCompliance.md)

### 3. 代码签名与供应链完整性 (Code Signing and Supply Chain Integrity)

代码签名与供应链完整性。

**核心工具:**
- cosign - 容器签名
- sigstore - 代码签名
- GnuPG - PGP签名
- in-toto - 供应链框架

**详细文档:** [代码签名与供应链完整性-CodeSigningIntegrity.md](skills/代码签名与供应链完整性-CodeSigningIntegrity.md)

### 4. 容器镜像安全扫描 (Container Image Security Scanning)

容器镜像安全扫描。

**核心工具:**
- Trivy - 容器安全扫描
- Clair - 容器漏洞扫描
- Grype - 容器漏洞扫描
- Anchore - 容器安全平台

**详细文档:** [容器镜像安全扫描-ContainerImageScan.md](skills/容器镜像安全扫描-ContainerImageScan.md)

### 5. 第三方供应商风险评估 (Third-party Vendor Risk Assessment)

第三方供应商风险评估。

**核心工具:**
- BitSight - 安全评级平台
- SecurityScorecard - 安全评分
- OneTrust - 合规管理平台
- 自定义工具 - 自定义评估工具

**详细文档:** [第三方供应商风险评估-ThirdPartyVendorRisk.md](skills/第三方供应商风险评估-ThirdPartyVendorRisk.md)

### 6. 供应链攻击检测与响应 (Supply Chain Attack Detection and Response)

供应链攻击检测与响应。

**核心工具:**
- GuardDog - 供应链安全
- Falco - 运行时安全
- OPA - 策略引擎
- Dependency-Check - 依赖检查

**详细文档:** [供应链攻击检测与响应-SupplyChainAttackResponse.md](skills/供应链攻击检测与响应-SupplyChainAttackResponse.md)

## 工具和资源

### 推荐工具套件

**SBOM工具:**
- [Syft](https://github.com/anchore/syft) - SBOM生成
- [CycloneDX](https://cyclonedx.org/) - SBOM标准
- [SPDX](https://spdx.dev/) - SBOM标准
- [Trivy](https://aquasecurity.github.io/trivy/) - 容器安全扫描

**依赖管理工具:**
- [FOSSA](https://fossa.com/) - 开源合规平台
- [Black Duck](https://www.synopsys.com/software-integrity/security-testing/software-composition-analysis) - 开源合规工具
- [Snyk](https://snyk.io/) - 依赖安全扫描
- [OWASP Dependency-Check](https://owasp.org/www-project-dependency-check/) - 依赖检查

**代码签名工具:**
- [cosign](https://github.com/sigstore/cosign) - 容器签名
- [sigstore](https://www.sigstore.dev/) - 代码签名
- [GnuPG](https://gnupg.org/) - PGP签名
- [in-toto](https://in-toto.io/) - 供应链框架

**容器安全工具:**
- [Trivy](https://aquasecurity.github.io/trivy/) - 容器安全扫描
- [Clair](https://github.com/quay/clair) - 容器漏洞扫描
- [Grype](https://github.com/anchore/grype) - 容器漏洞扫描
- [Anchore](https://anchore.com/) - 容器安全平台

### 学习资源

- [NIST SP 800-161](https://csrc.nist.gov/publications/detail/sp/800-161/final) - 供应链风险管理
- [CISA Supply Chain](https://www.cisa.gov/supply-chain-risk-management) - CISA供应链
- [SLSA Framework](https://slsa.dev/) - SLSA框架
- [OpenSSF Best Practices](https://openssf.org/) - OpenSSF最佳实践
- [NTIA SBOM](https://www.ntia.gov/sbom) - NTIA SBOM

## 最佳实践

### 供应链安全原则

1. **可见性** - 了解供应链资产
2. **最小化** - 最小化供应链依赖
3. **验证** - 验证供应链完整性
4. **监控** - 持续监控供应链
5. **响应** - 快速响应供应链攻击

### 常见误区

- 忽略供应链风险
- 没有生成SBOM
- 缺乏代码签名
- 没有供应商评估
- 忽略持续监控

## 参考文档

### 标准和框架
- [NIST SP 800-161](https://csrc.nist.gov/publications/detail/sp/800-161/final) - 供应链风险管理
- [SLSA Framework](https://slsa.dev/) - SLSA框架
- [ISO 28001](https://www.iso.org/standard/45654.html) - 供应链安全
- [NTIA SBOM](https://www.ntia.gov/sbom) - NTIA SBOM

### 技术指南
- [OpenSSF Best Practices](https://openssf.org/) - OpenSSF最佳实践
- [CISA Supply Chain](https://www.cisa.gov/supply-chain-risk-management) - CISA供应链
- [CycloneDX Documentation](https://cyclonedx.org/) - CycloneDX文档

## 维护与更新

本技能基于最新的供应链安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 6
**工具数量**: 25+

---

*本技能旨在帮助安全研究人员、供应链安全专家和DevSecOps工程师系统性地实施供应链安全。*