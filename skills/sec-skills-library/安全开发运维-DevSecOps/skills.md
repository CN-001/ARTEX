---
name: devsecops
description: Comprehensive DevSecOps framework covering CI/CD pipeline security audit, IaC security scanning, SAST, DAST, software supply chain security, and security requirements and threat modeling. Use when implementing DevSecOps practices, securing CI/CD pipelines, or integrating security into SDLC.
---

# 安全开发运维技能

## 概述

本技能提供系统化的DevSecOps框架，涵盖CI/CD管道安全审计、IaC安全扫描、SAST静态应用安全测试、DAST动态应用安全测试、软件供应链安全和安全需求与威胁建模等核心能力。

**核心能力:**
- CI/CD安全 - CI/CD管道安全审计
- IaC安全 - 基础设施即代码安全扫描
- SAST - 静态应用安全测试
- DAST - 动态应用安全测试
- 供应链安全 - 软件供应链安全
- 威胁建模 - 安全需求和威胁建模

**适用场景:**
- DevSecOps实施
- CI/CD安全加固
- 安全左移实践
- 软件供应链安全
- 安全开发培训

## 快速开始

### DevSecOps实施流程

1. **安全需求分析** - 1-2 天
   - 分析安全需求
   - 进行威胁建模
   - 制定安全策略

2. **安全工具集成** - 2-3 天
   - 选择安全工具
   - 集成到CI/CD
   - 配置安全检查

3. **安全测试自动化** - 3-5 天
   - 自动化SAST
   - 自动化DAST
   - 自动化依赖检查

4. **安全监控和响应** - 持续进行
   - 监控安全指标
   - 响应安全事件
   - 持续改进流程

## 核心技能模块

### 1. CI/CD管道安全审计 (CI/CD Pipeline Security Audit)

CI/CD管道安全审计。

**核心工具:**
- GitLab CI/CD Security - GitLab安全
- Jenkins Security - Jenkins安全
- GitHub Actions Security - GitHub安全

**详细文档:** [CI-CD管道安全审计-CICDPipelineSecurity.md](skills/CI-CD管道安全审计-CICDPipelineSecurity.md)

### 2. IaC安全扫描 (Infrastructure as Code Security)

基础设施即代码安全扫描。

**核心工具:**
- Checkov - IaC安全扫描
- Terrascan - Terraform安全
- tfsec - Terraform安全检查

**详细文档:** [IaC安全扫描-InfrastructureAsCodeSecurity.md](skills/IaC安全扫描-InfrastructureAsCodeSecurity.md)

### 3. SAST静态应用安全测试 (SAST)

静态应用安全测试。

**核心工具:**
- SonarQube - 代码质量平台
- Semgrep - 语义代码分析
- CodeQL - 代码查询语言

**详细文档:** [SAST静态应用安全测试-SAST.md](skills/SAST静态应用安全测试-SAST.md)

### 4. DAST动态应用安全测试 (DAST)

动态应用安全测试。

**核心工具:**
- OWASP ZAP - Web应用扫描
- Burp Suite - Web安全测试
- Acunetix - 商业扫描器

**详细文档:** [DAST动态应用安全测试-DAST.md](skills/DAST动态应用安全测试-DAST.md)

### 5. 软件供应链安全 (Software Supply Chain Security)

软件供应链安全。

**核心工具:**
- Snyk - 依赖安全扫描
- Dependabot - 依赖更新
- Trivy - 容器安全扫描
- SBOM - 软件物料清单

**详细文档:** [软件供应链安全-SoftwareSupplyChainSecurity.md](skills/软件供应链安全-SoftwareSupplyChainSecurity.md)

### 6. 安全需求与威胁建模 (Security Requirements and Threat Modeling)

安全需求和威胁建模。

**核心工具:**
- STRIDE - 威胁建模方法
- PASTA - 威胁建模方法
- Threat Dragon - 威胁建模工具

**详细文档:** [安全需求与威胁建模-ThreatModeling.md](skills/安全需求与威胁建模-ThreatModeling.md)

## 工具和资源

### 推荐工具套件

**CI/CD安全:**
- [GitLab CI/CD Security](https://docs.gitlab.com/ee/user/application_security/) - GitLab安全
- [Jenkins Security](https://www.jenkins.io/security/) - Jenkins安全
- [GitHub Actions Security](https://docs.github.com/en/actions/security-guides) - GitHub安全

**IaC安全:**
- [Checkov](https://www.checkov.io/) - IaC安全扫描
- [Terrascan](https://runterrascan.io/) - Terraform安全
- [tfsec](https://github.com/aquasecurity/tfsec) - Terraform安全检查

**SAST工具:**
- [SonarQube](https://www.sonarqube.org/) - 代码质量平台
- [Semgrep](https://semgrep.dev/) - 语义代码分析
- [CodeQL](https://codeql.github.com/) - 代码查询语言

**DAST工具:**
- [OWASP ZAP](https://www.zaproxy.org/) - Web应用扫描
- [Burp Suite](https://portswigger.net/burp) - Web安全测试
- [Acunetix](https://www.acunetix.com/) - 商业扫描器

**供应链安全:**
- [Snyk](https://snyk.io/) - 依赖安全扫描
- [Dependabot](https://github.com/dependabot) - 依赖更新
- [Trivy](https://aquasecurity.github.io/trivy/) - 容器安全扫描
- [SBOM](https://www.ntia.gov/SBOM) - 软件物料清单

### 学习资源

- [OWASP DevSecOps Guideline](https://owasp.org/www-project-devsecops-guideline/)
- [CIS Software Supply Chain Security](https://www.cisecurity.org/insights/blog/cis-build-software-supply-chain-security)
- [NIST SSDF (SP 800-218)](https://csrc.nist.gov/publications/detail/sp/800-218/final)
- [SLSA Framework](https://slsa.dev/)
- [CNCF Software Supply Chain Best Practices](https://www.cncf.io/)

## 最佳实践

### DevSecOps原则

1. **安全左移** - 尽早集成安全
2. **自动化** - 自动化安全测试
3. **持续监控** - 持续监控安全
4. **文化驱动** - 建立安全文化
5. **持续改进** - 持续改进流程

### 常见误区

- 安全测试太晚
- 过度依赖自动化
- 缺乏安全培训
- 没有威胁建模
- 忽略供应链安全

## 参考文档

### 标准和框架
- [OWASP DevSecOps Guideline](https://owasp.org/www-project-devsecops-guideline/)
- [NIST SSDF (SP 800-218)](https://csrc.nist.gov/publications/detail/sp/800-218/final)
- [SLSA Framework](https://slsa.dev/)
- [ISO 27034](https://www.iso.org/standard/44557.html)

### 技术指南
- [GitLab Security Documentation](https://docs.gitlab.com/ee/user/application_security/)
- [SonarQube Documentation](https://docs.sonarqube.org/)
- [OWASP ZAP Documentation](https://www.zaproxy.org/docs/)

## 维护与更新

本技能基于最新的DevSecOps实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 6
**工具数量**: 20+

---

*本技能旨在帮助DevSecOps工程师、安全开发人员和DevOps工程师系统性地实施DevSecOps。*