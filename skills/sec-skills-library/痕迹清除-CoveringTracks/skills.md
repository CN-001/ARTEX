---
name: covering-tracks
description: Comprehensive anti-forensics framework covering track clearing and anti-forensics, process injection and code injection, code obfuscation and anti-analysis, and AMSI bypass and EDR evasion. Use when conducting penetration testing, security assessments, or anti-forensics research.
---

# 痕迹清除技能

## 概述

本技能提供系统化的痕迹清除框架，涵盖痕迹清除与反取证、进程注入与代码注入、代码混淆与反分析和AMSI绕过与EDR规避等核心能力。

**核心能力:**
- 痕迹清除 - 清除操作痕迹
- 反取证 - 防止取证分析
- 进程注入 - 注入恶意代码
- 代码注入 - 注入代码到进程
- 代码混淆 - 混淆恶意代码
- 反分析 - 对抗分析工具
- AMSI绕过 - 绕过反恶意软件扫描
- EDR规避 - 规避端点检测

**适用场景:**
- 渗透测试痕迹清除
- 安全评估和红队演练
- 漏洞研究和PoC开发
- 安全培训和演练
- 应急响应和取证

## 快速开始

### 痕迹清除流程

1. **日志清除** - 10-20 分钟
   - 识别相关日志
   - 清除敏感记录
   - 修改时间戳

2. **文件隐藏** - 15-30 分钟
   - 隐藏恶意文件
   - 修改文件属性
   - 使用隐藏技术

3. **反取证** - 20-40 分钟
   - 清除内存痕迹
   - 清除网络痕迹
   - 清除注册表痕迹

## 核心技能模块

### 1. 痕迹清除 (Covering Tracks)

清除渗透过程中的操作痕迹。

**核心工具:**
- 日志清除 - 系统日志清除
- 文件隐藏 - 文件隐藏技术
- 时间戳修改 - 修改文件时间
- 反取证技术 - 反取证方法

**详细文档:** [痕迹清除-CoveringTracks.md](skills/痕迹清除-CoveringTracks.md)

### 2. 进程注入与代码注入 (Process and Code Injection)

注入恶意代码到合法进程。

**核心工具:**
- Cobalt Strike - 商业C2框架
- Meterpreter - Metasploit后门
- PowerSploit - PowerShell工具集
- 自定义注入 - 自定义注入代码

**详细文档:** [进程注入与代码注入-ProcessInjection.md](skills/进程注入与代码注入-ProcessInjection.md)

### 3. 代码混淆与反分析 (Code Obfuscation and Anti-Analysis)

混淆恶意代码，对抗分析工具。

**核心工具:**
- ConfuserEx - .NET混淆器
- Obfuscator-LLVM - LLVM混淆器
- UPX - 可执行文件压缩
- VMP - 虚拟化保护

**详细文档:** [代码混淆与反分析-ObfuscationAntiAnalysis.md](skills/代码混淆与反分析-ObfuscationAntiAnalysis.md)

### 4. AMSI绕过与EDR规避 (AMSI Bypass and EDR Evasion)

绕过反恶意软件扫描和端点检测。

**核心工具:**
- SysWhispers - 系统调用混淆
- AMSI.fail - AMSI绕过工具
- SharpSploit - .NET安全工具
- 自定义规避 - 自定义规避技术

**详细文档:** [AMSI绕过与EDR规避-AMSIByPassEDREvasion.md](skills/AMSI绕过与EDR规避-AMSIByPassEDREvasion.md)

## 工具和资源

### 推荐工具套件

**痕迹清除:**
- [Timestomp](https://learn.microsoft.com/en-us/sysinternals/downloads/timestomp) - 时间戳修改
- [SDelete](https://learn.microsoft.com/en-us/sysinternals/downloads/sdelete) - 安全删除
- [CCleaner](https://www.ccleaner.com/) - 系统清理
- [BleachBit](https://www.bleachbit.org/) - 磁盘清理

**进程注入:**
- [Cobalt Strike](https://www.cobaltstrike.com/) - 商业C2框架
- [Metasploit](https://www.metasploit.com/) - 开源C2框架
- [PowerSploit](https://github.com/PowerShellMafia/PowerSploit) - PowerShell工具
- [SharpHound](https://github.com/BloodHoundAD/SharpHound) - .NET注入工具

**代码混淆:**
- [ConfuserEx](https://github.com/mkaring/ConfuserEx) - .NET混淆器
- [Obfuscator-LLVM](https://github.com/obfuscator-llvm/obfuscator) - LLVM混淆器
- [UPX](https://upx.github.io/) - 可执行文件压缩
- [VMP](https://vmpsoft.com/) - 虚拟化保护

**EDR规避:**
- [SysWhispers](https://github.com/jthuraisamy/SysWhispers) - 系统调用混淆
- [AMSI.fail](https://amsi.fail/) - AMSI绕过
- [SharpSploit](https://github.com/cobbr/SharpSploit) - .NET安全工具
- [Havoc](https://github.com/HavocFramework/Havoc) - 现代C2框架

### 学习资源

- [MITRE ATT&CK - Defense Evasion](https://attack.mitre.org/tactics/TA0005/)
- [Anti-Forensic Techniques](https://forensicswiki.xyz/)
- [Timestomp Documentation](https://learn.microsoft.com/en-us/sysinternals/downloads/timestomp)
- [DFIR Review](https://dfirreview.com/)

## 最佳实践

### 痕迹清除原则

1. **授权性** - 确保获得清除授权
2. **完整性** - 避免破坏系统完整性
3. **隐蔽性** - 避免被检测和发现
4. **可追溯性** - 记录清除活动
5. **合规性** - 遵守法律法规

### 常见误区

- 过度依赖自动化工具
- 忽略检测规避
- 没有评估清除效果
- 缺乏备份机制
- 违反道德和法律

## 参考文档

### 标准和框架
- [PTES - Penetration Testing Execution Standard](https://www.pentest-standard.org/)
- [MITRE ATT&CK - Defense Evasion](https://attack.mitre.org/tactics/TA0005/)
- [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)

### 技术指南
- [HackTricks - Anti-Forensics](https://book.hacktricks.xyz/windows-hardening/stealth-methods)
- [SysWhispers Documentation](https://github.com/jthuraisamy/SysWhispers)
- [AMSI.fail Documentation](https://amsi.fail/)

## 维护与更新

本技能基于最新的痕迹清除实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 4
**工具数量**: 20+

---

*本技能旨在帮助安全研究人员、渗透测试人员和红队成员系统性地进行痕迹清除。*