---
name: os-security
description: Comprehensive OS security framework covering Windows security hardening and baseline checking, Windows attack and lateral movement techniques, Linux security hardening and baseline checking, Linux attack and persistence techniques, macOS security assessment and hardening, and domestic OS security hardening. Use when conducting OS security assessments, hardening systems, or OS penetration testing.
---

# 操作系统安全技能

## 概述

本技能提供系统化的操作系统安全框架，涵盖Windows安全加固与基线检查、Windows攻击与横向移动技术、Linux安全加固与基线检查、Linux攻击与权限维持技术、macOS安全评估与加固和国产操作系统安全加固等核心能力。

**核心能力:**
- Windows加固 - Windows安全加固与基线检查
- Windows攻击 - Windows攻击与横向移动技术
- Linux加固 - Linux安全加固与基线检查
- Linux攻击 - Linux攻击与权限维持技术
- macOS安全 - macOS安全评估与加固
- 国产OS - 国产操作系统安全加固

**适用场景:**
- 操作系统安全评估
- 系统安全加固
- 操作系统渗透测试
- 系统安全运营
- 操作系统安全培训

## 快速开始

### OS安全评估流程

1. **系统发现** - 30-60 分钟
   - 发现操作系统
   - 分析系统版本
   - 识别系统配置

2. **安全评估** - 1-2 小时
   - 评估系统安全
   - 检查安全配置
   - 识别安全漏洞

3. **安全加固** - 2-4 小时
   - 应用安全基线
   - 修复安全漏洞
   - 加强安全配置

4. **持续监控** - 持续进行
   - 监控系统安全
   - 检测安全事件
   - 响应安全威胁

## 核心技能模块

### 1. Windows安全加固与基线检查 (Windows Security Hardening and Baseline Checking)

Windows安全加固与基线检查。

**核心工具:**
- CIS Benchmarks - CIS基准
- GPO - 组策略管理
- LGPO - 本地组策略
- Microsoft Security Compliance Toolkit - 安全合规工具包

**详细文档:** [Windows安全加固与基线检查-WindowsHardeningBaseline.md](skills/Windows安全加固与基线检查-WindowsHardeningBaseline.md)

### 2. Windows攻击与横向移动技术 (Windows Attack and Lateral Movement Techniques)

Windows攻击与横向移动技术。

**核心工具:**
- Impacket - Python网络工具
- Mimikatz - 凭证提取
- Rubeus - Kerberos操作
- BloodHound - AD攻击路径分析
- CrackMapExec - 网络评估工具

**详细文档:** [Windows攻击与横向移动技术-WindowsAttackLateralMovement.md](skills/Windows攻击与横向移动技术-WindowsAttackLateralMovement.md)

### 3. Linux安全加固与基线检查 (Linux Security Hardening and Baseline Checking)

Linux安全加固与基线检查。

**核心工具:**
- Lynis - 安全审计工具
- OpenSCAP - 安全内容自动化
- auditd - 审计系统
- SELinux - 安全增强Linux
- AppArmor - 应用程序armor

**详细文档:** [Linux安全加固与基线检查-LinuxHardeningBaseline.md](skills/Linux安全加固与基线检查-LinuxHardeningBaseline.md)

### 4. Linux攻击与权限维持技术 (Linux Attack and Persistence Techniques)

Linux攻击与权限维持技术。

**核心工具:**
- LinPEAS - Linux提权枚举
- GTFO Bins - 二进制文件利用
- systemctl - 系统服务管理
- Chisel - HTTP隧道
- Diamorphine - Linux rootkit

**详细文档:** [Linux攻击与权限维持技术-LinuxAttackPersistence.md](skills/Linux攻击与权限维持技术-LinuxAttackPersistence.md)

### 5. macOS安全评估与加固 (macOS Security Assessment and Hardening)

macOS安全评估与加固。

**核心工具:**
- SIP - 系统完整性保护
- TCC - 透明同意与控制
- MDM - 移动设备管理
- osquery - 操作系统查询
- Santa - 二进制白名单
- Goal - 端点保护
- BlockBlock - 端点保护

**详细文档:** [macOS安全评估与加固-macOSSecurityHardening.md](skills/macOS安全评估与加固-macOSSecurityHardening.md)

### 6. 国产操作系统安全加固 (Domestic OS Security Hardening)

国产操作系统安全加固。

**核心工具:**
- KylinOS - 银河麒麟
- UOS - 统信操作系统
- openEuler - 华为openEuler
- 等保2.0 - 等级保护2.0
- SM2/SM3/SM4 - 国密算法

**详细文档:** [国产操作系统安全加固-DomesticOSSecurity.md](skills/国产操作系统安全加固-DomesticOSSecurity.md)

## 工具和资源

### 推荐工具套件

**Windows工具:**
- [CIS Benchmarks](https://www.cisecurity.org/benchmark/) - CIS基准
- [Microsoft Security Compliance Toolkit](https://www.microsoft.com/en-us/download/details.aspx?id=55319) - 安全合规工具包
- [Impacket](https://github.com/fortra/impacket) - Python网络工具
- [Mimikatz](https://github.com/gentilkiwi/mimikatz) - 凭证提取

**Linux工具:**
- [Lynis](https://cisofy.com/lynis/) - 安全审计工具
- [OpenSCAP](https://www.open-scap.org/) - 安全内容自动化
- [LinPEAS](https://github.com/carlospolop/linPEAS) - Linux提权枚举
- [GTFOBins](https://gtfobins.github.io/) - 二进制文件利用

**macOS工具:**
- [osquery](https://osquery.io/) - 操作系统查询
- [Santa](https://santa.dev/) - 二进制白名单
- [BlockBlock](https://objective-see.com/products/blockblock) - 端点保护
- [Goal](https://github.com/joelsephus/goal) - 端点保护

**国产OS工具:**
- [KylinOS](https://www.kylinos.cn/) - 银河麒麟
- [UOS](https://www.uniontech.com/) - 统信操作系统
- [openEuler](https://www.openeuler.org/) - 华为openEuler

### 学习资源

- [CIS Benchmarks](https://www.cisecurity.org/benchmark/) - CIS基准
- [GB/T 22239-2019](https://openstd.samr.gov.cn/) - 等保2.0
- [MITRE ATT&CK](https://attack.mitre.org/tactics/enterprise/) - ATT&CK框架
- [NIST SP 800-53 Rev 5](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final) - 安全控制
- [Microsoft Security Baselines](https://docs.microsoft.com/en-us/windows/security/threat-protection/windows-security-baselines) - Microsoft安全基线
- [Apple Platform Security](https://support.apple.com/en-us/guide/security/welcome/web) - Apple平台安全

## 最佳实践

### OS安全原则

1. **最小权限** - 遵循最小权限原则
2. **安全基线** - 应用安全基线
3. **持续更新** - 持续更新系统
4. **监控审计** - 监控和审计系统
5. **应急响应** - 准备应急响应

### 常见误区

- 忽略安全基线
- 没有最小权限
- 缺乏系统更新
- 没有监控审计
- 缺乏应急准备

## 参考文档

### 标准和框架
- [CIS Benchmarks](https://www.cisecurity.org/benchmark/) - CIS基准
- [GB/T 22239-2019](https://openstd.samr.gov.cn/) - 等保2.0
- [NIST SP 800-53 Rev 5](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final) - 安全控制
- [ISO 27001](https://www.iso.org/isoiec-27001-information-security.html) - 信息安全

### 技术指南
- [Microsoft Security Baselines](https://docs.microsoft.com/en-us/windows/security/threat-protection/windows-security-baselines) - Microsoft安全基线
- [Apple Platform Security](https://support.apple.com/en-us/guide/security/welcome/web) - Apple平台安全
- [KylinOS Security](https://www.kylinos.cn/support/technical.html) - 银河麒麟安全

## 维护与更新

本技能基于最新的操作系统安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 6
**工具数量**: 30+

---

*本技能旨在帮助安全研究人员、系统管理员和安全工程师系统性地进行操作系统安全评估和加固。*