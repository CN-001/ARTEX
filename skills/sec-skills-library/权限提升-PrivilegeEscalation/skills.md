---
name: privilege-escalation
description: Comprehensive privilege escalation framework covering Linux privilege escalation, Windows privilege escalation, kernel and service misconfiguration exploitation, and credential theft. Use when conducting penetration testing, security assessments, or privilege escalation research.
---

# 权限提升技能

## 概述

本技能提供系统化的权限提升框架，涵盖Linux提权、Windows提权、内核漏洞与服务配置错误提权和凭证窃取与利用等核心能力。

**核心能力:**
- Linux提权 - Linux系统权限提升
- Windows提权 - Windows系统权限提升
- 内核漏洞提权 - 利用内核漏洞提权
- 服务配置错误提权 - 利用服务配置错误提权
- 凭证窃取 - 窃取系统凭证
- 哈希传递 - 利用哈希进行认证

**适用场景:**
- 渗透测试权限提升
- 安全评估和红队演练
- 漏洞研究和PoC开发
- 安全培训和演练
- 应急响应和取证

## 快速开始

### Linux提权流程

1. **信息收集** - 10-20 分钟
   - 收集系统信息
   - 检查用户权限
   - 识别潜在提权点

2. **漏洞识别** - 15-30 分钟
   - 检查内核版本
   - 查找已知漏洞
   - 分析配置错误

3. **提权利用** - 20-40 分钟
   - 利用已知漏洞
   - 利用配置错误
   - 获取root权限

### Windows提权流程

1. **信息收集** - 10-20 分钟
   - 收集系统信息
   - 检查用户权限
   - 识别潜在提权点

2. **漏洞识别** - 15-30 分钟
   - 检查系统版本
   - 查找已知漏洞
   - 分析配置错误

3. **提权利用** - 20-40 分钟
   - 利用已知漏洞
   - 利用配置错误
   - 获取管理员权限

## 核心技能模块

### 1. Linux提权 (Linux Privilege Escalation)

在Linux系统中提升权限。

**核心工具:**
- LinPEAS - Linux提权枚举脚本
- Linux Exploit Suggester - 漏洞建议工具
- GTFOBins - 二进制文件利用
- SUID/GUID - 特殊权限利用

**详细文档:** [Linux提权-LinuxPrivEsc.md](skills/Linux提权-LinuxPrivEsc.md)

### 2. Windows提权 (Windows Privilege Escalation)

在Windows系统中提升权限。

**核心工具:**
- WinPEAS - Windows提权枚举脚本
- PowerUp - PowerShell提权工具
- PrivescCheck - 权限检查工具
- Watson - 漏洞检测工具

**详细文档:** [Windows提权-WindowsPrivEsc.md](skills/Windows提权-WindowsPrivEsc.md)

### 3. 内核漏洞与服务配置错误提权 (Kernel and Service Misconfiguration Privilege Escalation)

利用内核漏洞和服务配置错误进行提权。

**核心工具:**
- 漏洞利用代码 - 内核漏洞利用
- MSF模块 - Metasploit模块
- 服务配置检查 - 服务配置分析
- 计划任务利用 - 计划任务提权

**详细文档:** [内核漏洞与服务配置错误提权-KernelServicePrivEsc.md](skills/内核漏洞与服务配置错误提权-KernelServicePrivEsc.md)

### 4. 凭证窃取与利用 (Credential Theft and Utilization)

窃取和利用系统凭证。

**核心工具:**
- Mimikatz - Windows凭证提取
- Lazagne - 凭证提取工具
- LaZagne - 密码恢复工具
- 哈希传递 - 凭证传递攻击

**详细文档:** [凭证窃取与利用-CredentialTheft.md](skills/凭证窃取与利用-CredentialTheft.md)

## 工具和资源

### 推荐工具套件

**Linux提权:**
- [LinPEAS](https://github.com/carlospolop/linPEAS) - Linux提权枚举
- [Linux Exploit Suggester](https://github.com/mzet-/linux-exploit-suggester) - 漏洞建议
- [GTFOBins](https://gtfobins.github.io/) - 二进制文件利用
- [LinEnum](https://github.com/rebootuser/LinEnum) - Linux枚举脚本

**Windows提权:**
- [WinPEAS](https://github.com/carlospolop/WinPEAS) - Windows提权枚举
- [PowerUp](https://github.com/PowerShellMafia/PowerSploit/tree/master/Privesc) - PowerShell提权
- [PrivescCheck](https://github.com/itm4n/PrivescCheck) - 权限检查
- [Watson](https://github.com/rasta-mouse/Watson) - 漏洞检测

**凭证窃取:**
- [Mimikatz](https://github.com/gentilkiwi/mimikatz) - Windows凭证提取
- [LaZagne](https://github.com/AlessandroZ/LaZagne) - 密码恢复
- [Hashcat](https://hashcat.net/hashcat/) - 密码破解
- [John the Ripper](https://www.openwall.com/john/) - 密码破解

### 学习资源

- [GTFOBins](https://gtfobins.github.io/) - Unix二进制利用
- [LOLBAS](https://lolbas-project.github.io/) - Windows二进制利用
- [HackTricks - Privilege Escalation](https://book.hacktricks.xyz/)
- [MITRE ATT&CK - Privilege Escalation](https://attack.mitre.org/tactics/TA0004/)

## 最佳实践

### 权限提升原则

1. **授权性** - 确保获得提权授权
2. **最小影响** - 避免对系统造成破坏
3. **可追溯性** - 记录所有提权活动
4. **安全性** - 保护提权工具和环境
5. **合规性** - 遵守法律法规

### 常见误区

- 过度依赖自动化工具
- 忽略配置错误
- 没有验证提权结果
- 缺乏回退计划
- 违反道德和法律

## 参考文档

### 标准和框架
- [PTES - Penetration Testing Execution Standard](https://www.pentest-standard.org/)
- [MITRE ATT&CK - Privilege Escalation](https://attack.mitre.org/tactics/TA0004/)
- [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)

### 技术指南
- [GTFOBins Documentation](https://gtfobins.github.io/)
- [LOLBAS Documentation](https://lolbas-project.github.io/)
- [HackTricks - Linux Privilege Escalation](https://book.hacktricks.xyz/linux-hardening/privilege-escalation)

## 维护与更新

本技能基于最新的权限提升实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 4
**工具数量**: 20+

---

*本技能旨在帮助安全研究人员、渗透测试人员和红队成员系统性地进行权限提升。*