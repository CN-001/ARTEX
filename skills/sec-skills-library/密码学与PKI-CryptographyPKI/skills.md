---
name: cryptography-pki
description: Comprehensive cryptography and PKI framework covering TLS/SSL security configuration and audit, PKI architecture and certificate management, and encryption algorithms and key management. Use when implementing cryptographic controls, managing PKI infrastructure, or securing communications.
---

# 密码学与PKI技能

## 概述

本技能提供系统化的密码学与PKI框架，涵盖TLS/SSL安全配置与审计、PKI架构与证书管理和加密算法与密钥管理等核心能力。

**核心能力:**
- TLS配置 - TLS/SSL安全配置与审计
- PKI管理 - PKI架构与证书管理
- 密钥管理 - 加密算法与密钥管理

**适用场景:**
- 密码学实施
- PKI基础设施管理
- TLS/SSL配置
- 密钥管理
- 加密算法应用

## 快速开始

### 密码学与PKI实施流程

1. **需求分析** - 1-2 小时
   - 分析加密需求
   - 确定安全要求
   - 选择加密算法

2. **架构设计** - 2-4 小时
   - 设计PKI架构
   - 规划证书体系
   - 制定密钥策略

3. **实施部署** - 4-8 小时
   - 部署PKI基础设施
   - 配置TLS/SSL
   - 实施密钥管理

4. **运维管理** - 持续进行
   - 监控证书状态
   - 管理密钥生命周期
   - 审计安全配置

## 核心技能模块

### 1. TLS/SSL安全配置与审计 (TLS/SSL Security Configuration and Audit)

TLS/SSL安全配置与审计。

**核心工具:**
- OpenSSL - 加密工具集
- certbot - 证书自动化
- acme.sh - 证书自动化
- SSL Labs - SSL测试

**详细文档:** [TLS-SSL安全配置与审计-TLSSSLSecurityConfigurationAudit.md](skills/TLS-SSL安全配置与审计-TLSSSLSecurityConfigurationAudit.md)

### 2. PKI架构与证书安全管理 (PKI Architecture and Certificate Management)

PKI架构与证书安全管理。

**核心工具:**
- CFSSL - 证书签发工具
- OpenSSL - 加密工具集
- HashiCorp Vault - 密钥管理
- HSM - 硬件安全模块

**详细文档:** [PKI架构与证书安全管理-PKIArchitectureCertificateManagement.md](skills/PKI架构与证书安全管理-PKIArchitectureCertificateManagement.md)

### 3. 加密算法与密钥管理 (Encryption Algorithms and Key Management)

加密算法与密钥管理。

**核心工具:**
- OpenSSL - 加密工具集
- HashiCorp Vault - 密钥管理
- HSM - 硬件安全模块
- 自定义工具 - 自定义加密工具

**详细文档:** [加密算法与密钥管理-EncryptionAlgorithmsKeyManagement.md](skills/加密算法与密钥管理-EncryptionAlgorithmsKeyManagement.md)

## 工具和资源

### 推荐工具套件

**加密工具:**
- [OpenSSL](https://www.openssl.org/) - 加密工具集
- [GnuPG](https://gnupg.org/) - PGP加密
- [LibreSSL](https://www.libressl.org/) - 加密库
- [Bouncy Castle](https://www.bouncycastle.org/) - 加密库

**证书工具:**
- [certbot](https://certbot.eff.org/) - 证书自动化
- [acme.sh](https://acme.sh/) - 证书自动化
- [CFSSL](https://github.com/cloudflare/cfssl) - 证书签发
- [Let's Encrypt](https://letsencrypt.org/) - 免费证书

**密钥管理:**
- [HashiCorp Vault](https://www.vaultproject.io/) - 密钥管理
- [AWS KMS](https://aws.amazon.com/kms/) - 密钥管理
- [Azure Key Vault](https://azure.microsoft.com/key-vault/) - 密钥管理
- [Google Cloud KMS](https://cloud.google.com/kms) - 密钥管理

**硬件安全模块:**
- [YubiHSM](https://www.yubico.com/products/yubihsm/) - 硬件安全模块
- [CloudHSM](https://aws.amazon.com/cloudhsm/) - 云硬件安全模块
- [Thales Luna](https://cpl.thalesgroup.com/) - 硬件安全模块

### 学习资源

- [NIST Cryptography](https://csrc.nist.gov/projects/cryptographic-standards-and-guidelines) - NIST密码学
- [OpenSSL Documentation](https://www.openssl.org/docs/) - OpenSSL文档
- [PKI Best Practices](https://www.cisecurity.org/benchmark/pki) - PKI最佳实践
- [TLS Configuration](https://ssl-config.mozilla.org/) - TLS配置

## 最佳实践

### 密码学原则

1. **算法选择** - 选择安全算法
2. **密钥管理** - 安全管理密钥
3. **证书管理** - 有效管理证书
4. **配置审计** - 定期审计配置
5. **持续更新** - 持续更新算法

### 常见误区

- 使用弱加密算法
- 密钥管理不当
- 证书过期未更新
- 配置不安全
- 缺乏定期审计

## 参考文档

### 标准和框架
- [NIST SP 800-57](https://csrc.nist.gov/publications/detail/sp/800-57/part-1/rev-5/final) - 密钥管理
- [NIST SP 800-52](https://csrc.nist.gov/publications/detail/sp/800-52/rev-2/final) - TLS安全
- [ISO 27001](https://www.iso.org/isoiec-27001-information-security.html) - 信息安全

### 技术指南
- [OpenSSL Documentation](https://www.openssl.org/docs/) - OpenSSL文档
- [PKI Best Practices](https://www.cisecurity.org/benchmark/pki) - PKI最佳实践
- [TLS Configuration](https://ssl-config.mozilla.org/) - TLS配置

## 维护与更新

本技能基于最新的密码学与PKI实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 3
**工具数量**: 15+

---

*本技能旨在帮助密码学专家、PKI管理员和安全架构师系统性地实施密码学与PKI。*