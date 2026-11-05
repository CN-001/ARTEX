---
name: persistence
description: Comprehensive persistence framework covering persistence techniques, boot and logon autostart, account persistence, Office application persistence, and Bootkit and firmware persistence. Use when conducting penetration testing, security assessments, or persistence research.
---

# 持久化技能

## 概述

本技能提供系统化的持久化框架，涵盖持久化技术、启动项与登录自动执行、账户持久化、Office应用程序持久化和Bootkit与固件持久化等核心能力。

**核心能力:**
- 持久化控制 - 建立持久化后门
- 启动项 - 利用启动项持久化
- 登录自动执行 - 登录时自动执行
- 账户持久化 - 创建隐藏账户
- Office持久化 - 利用Office应用
- Bootkit - 固件级持久化
- 固件持久化 - 硬件级持久化

**适用场景:**
- 渗透测试持久化
- 安全评估和红队演练
- 漏洞研究和PoC开发
- 安全培训和演练
- 应急响应和取证

## 快速开始

### 持久化流程

1. **持久化点识别** - 10-20 分钟
   - 识别系统持久化点
   - 分析启动流程
   - 发现自动执行机制

2. **持久化建立** - 15-30 分钟
   - 选择持久化技术
   - 配置持久化机制
   - 测试持久化效果

3. **持久化维护** - 20-40 分钟
   - 监控持久化状态
   - 更新持久化机制
   - 确保持续访问

## 核心技能模块

### 1. 持久化 (Persistence)

建立系统持久化后门。

**核心工具:**
- MSF Venom - Metasploit后门生成
- Cobalt Strike - 商业C2框架
- 计划任务 - 系统计划任务
- 自启动项 - 系统自启动

**详细文档:** [持久化-Persistence.md](skills/持久化-Persistence.md)

### 2. 启动项与登录自动执行 (Boot and Logon Autostart)

利用启动项和登录自动执行机制。

**核心工具:**
- Autoruns - 启动项管理工具
- schtasks - Windows计划任务
- systemd - Linux系统服务
- launchd - macOS启动服务

**详细文档:** [启动项与登录自动执行-BootLogonAutostart.md](skills/启动项与登录自动执行-BootLogonAutostart.md)

### 3. 账户持久化 (Account Persistence)

创建隐藏账户和凭证持久化。

**核心工具:**
- net user - Windows用户管理
- Mimikatz - 凭证提取
- AWS CLI - AWS凭证管理
- Azure AD - Azure凭证管理

**详细文档:** [账户持久化-AccountPersistence.md](skills/账户持久化-AccountPersistence.md)

### 4. Office应用程序持久化 (Office Application Persistence)

利用Office应用程序进行持久化。

**核心工具:**
- EvilClippy - VBA宏混淆工具
- macro_pack - 宏生成工具
- VBA - Visual Basic for Applications
- Office模板 - Office模板利用

**详细文档:** [Office应用程序持久化-OfficePersistence.md](skills/Office应用程序持久化-OfficePersistence.md)

### 5. Bootkit与固件持久化 (Bootkit and Firmware Persistence)

在固件级别建立持久化。

**核心工具:**
- Chipsec - 硬件安全分析
- Flashrom - 固件读写工具
- UEFITool - UEFI固件分析
- Bootkit - 引导级后门

**详细文档:** [Bootkit与固件持久化-BootkitFirmwarePersistence.md](skills/Bootkit与固件持久化-BootkitFirmwarePersistence.md)

## 工具和资源

### 推荐工具套件

**持久化工具:**
- [Cobalt Strike](https://www.cobaltstrike.com/) - 商业C2框架
- [Metasploit](https://www.metasploit.com/) - 开源C2框架
- [Empire](https://github.com/EmpireProject/Empire) - PowerShell C2
- [Sliver](https://github.com/BishopFox/sliver) - 现代C2框架

**启动项工具:**
- [Autoruns](https://learn.microsoft.com/en-us/sysinternals/downloads/autoruns) - 启动项管理
- [schtasks](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/schtasks) - 计划任务
- [systemd](https://www.freedesktop.org/software/systemd/) - Linux系统服务
- [launchd](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html) - macOS启动服务

**Office工具:**
- [EvilClippy](https://github.com/outflanknl/EvilClippy) - VBA宏混淆
- [macro_pack](https://github.com/sevagas/macro_pack) - 宏生成工具
- [Phishing](https://github.com/trustedsec/phishing) - 钓鱼工具包

**固件工具:**
- [Chipsec](https://github.com/chipsec/chipsec) - 硬件安全分析
- [Flashrom](https://www.flashrom.org/) - 固件读写
- [UEFITool](https://github.com/LongSoft/UEFITool) - UEFI固件分析

### 学习资源

- [MITRE ATT&CK - Persistence](https://attack.mitre.org/tactics/TA0003/)
- [HackTricks - Persistence](https://book.hacktricks.xyz/windows-hardening/stealth-persistence)
- [SharPersist](https://github.com/fireeye/SharPersist) - Windows持久化工具包
- [PayloadsAllTheThings - Persistence](https://github.com/swisskyrepo/PayloadsAllTheThings)

## 最佳实践

### 持久化原则

1. **授权性** - 确保获得持久化授权
2. **隐蔽性** - 避免被检测和清除
3. **稳定性** - 确保持久化稳定运行
4. **可维护性** - 便于更新和维护
5. **合规性** - 遵守法律法规

### 常见误区

- 过度依赖单一持久化方法
- 忽略检测规避
- 没有测试持久化效果
- 缺乏备份机制
- 违反道德和法律

## 参考文档

### 标准和框架
- [PTES - Penetration Testing Execution Standard](https://www.pentest-standard.org/)
- [MITRE ATT&CK - Persistence](https://attack.mitre.org/tactics/TA0003/)
- [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)

### 技术指南
- [HackTricks - Windows Persistence](https://book.hacktricks.xyz/windows-hardening/stealth-persistence)
- [HackTricks - Linux Persistence](https://book.hacktricks.xyz/linux-hardening/privilege-escalation#linux-persistence)
- [SharPersist Documentation](https://github.com/fireeye/SharPersist)

## 维护与更新

本技能基于最新的持久化实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 5
**工具数量**: 20+

---

*本技能旨在帮助安全研究人员、渗透测试人员和红队成员系统性地建立持久化。*