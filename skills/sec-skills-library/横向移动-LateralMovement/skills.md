---
name: lateral-movement
description: Comprehensive lateral movement framework covering lateral movement techniques, internal proxy and tunneling, and PsExec and WMI remote execution. Use when conducting penetration testing, security assessments, or lateral movement research.
---

# 横向移动技能

## 概述

本技能提供系统化的横向移动框架，涵盖横向移动技术、内网代理与隧道和PsExec与WMI远程执行等核心能力。

**核心能力:**
- 横向移动 - 在内网中移动
- 内网代理 - 建立内网代理
- 隧道技术 - 绕过网络限制
- PsExec - 远程命令执行
- WMI - Windows管理接口
- WinRM - Windows远程管理
- Impacket - Python网络工具

**适用场景:**
- 渗透测试横向移动
- 安全评估和红队演练
- 漏洞研究和PoC开发
- 安全培训和演练
- 应急响应和取证

## 快速开始

### 横向移动流程

1. **网络发现** - 10-20 分钟
   - 发现内网主机
   - 识别开放端口
   - 绘制网络拓扑

2. **凭证收集** - 15-30 分钟
   - 收集本地凭证
   - 发现域凭证
   - 识别重用凭证

3. **横向移动** - 20-40 分钟
   - 选择移动技术
   - 执行远程命令
   - 建立持久化

### 内网代理流程

1. **代理建立** - 10-20 分钟
   - 选择代理工具
   - 配置代理节点
   - 建立代理链

2. **隧道建立** - 15-30 分钟
   - 选择隧道协议
   - 配置隧道参数
   - 测试隧道连接

3. **流量转发** - 20-40 分钟
   - 配置流量转发
   - 测试访问能力
   - 优化代理性能

## 核心技能模块

### 1. 横向移动 (Lateral Movement)

在内网中进行横向移动。

**核心工具:**
- BloodHound - AD攻击路径分析
- CrackMapExec - 网络评估工具
- Impacket - Python网络工具
- 自定义脚本 - 横向移动脚本

**详细文档:** [横向移动-LateralMovement.md](skills/横向移动-LateralMovement.md)

### 2. 内网代理与隧道 (Internal Proxy and Tunneling)

建立内网代理和隧道。

**核心工具:**
- Chisel - HTTP隧道
- FRP - 反向代理
- Stowaway - 多级代理
- Ngrok - 内网穿透

**详细文档:** [内网代理与隧道-InternalProxyTunnel.md](skills/内网代理与隧道-InternalProxyTunnel.md)

### 3. PsExec与WMI远程执行 (PsExec and WMI Remote Execution)

使用PsExec和WMI进行远程执行。

**核心工具:**
- PsExec - 远程命令执行
- WMI - Windows管理接口
- WinRM - Windows远程管理
- Impacket - Python网络工具

**详细文档:** [PsExec与WMI远程执行-PsExecWMI.md](skills/PsExec与WMI远程执行-PsExecWMI.md)

## 工具和资源

### 推荐工具套件

**横向移动:**
- [BloodHound](https://github.com/BloodHoundAD/BloodHound) - AD攻击路径分析
- [CrackMapExec](https://github.com/Porchetta-Industries/CrackMapExec) - 网络评估工具
- [Impacket](https://github.com/fortra/impacket) - Python网络工具
- [Responder](https://github.com/SpiderLabs/Responder) - LLMNR/NBT-NS欺骗

**内网代理:**
- [Chisel](https://github.com/jpillora/chisel) - HTTP隧道
- [FRP](https://github.com/fatedier/frp) - 反向代理
- [Stowaway](https://github.com/ph4ntonn/Stowaway) - 多级代理
- [Ngrok](https://ngrok.com/) - 内网穿透

**远程执行:**
- [PsExec](https://learn.microsoft.com/en-us/sysinternals/downloads/psexec) - 远程命令执行
- [WinRM](https://learn.microsoft.com/en-us/windows/win32/winrm/portal) - Windows远程管理
- [WMI](https://learn.microsoft.com/en-us/windows/win32/wmisdk/wmi-start-page) - Windows管理接口
- [Impacket](https://github.com/fortra/impacket) - Python网络工具

### 学习资源

- [MITRE ATT&CK - Lateral Movement](https://attack.mitre.org/tactics/TA0008/)
- [CrackMapExec Documentation](https://github.com/Porchetta-Industries/CrackMapExec)
- [Impacket Documentation](https://github.com/fortra/impacket)
- [BloodHound Documentation](https://bloodhound.readthedocs.io/)

## 最佳实践

### 横向移动原则

1. **授权性** - 确保获得横向移动授权
2. **隐蔽性** - 避免被检测和阻断
3. **渐进性** - 逐步扩展访问范围
4. **可追溯性** - 记录所有移动活动
5. **合规性** - 遵守法律法规

### 常见误区

- 过度依赖自动化工具
- 忽略检测规避
- 没有评估移动风险
- 缺乏回退计划
- 违反道德和法律

## 参考文档

### 标准和框架
- [PTES - Penetration Testing Execution Standard](https://www.pentest-standard.org/)
- [MITRE ATT&CK - Lateral Movement](https://attack.mitre.org/tactics/TA0008/)
- [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)

### 技术指南
- [BloodHound Documentation](https://bloodhound.readthedocs.io/)
- [Impacket Examples](https://github.com/fortra/impacket/tree/master/examples)
- [CrackMapExec Wiki](https://github.com/Porchetta-Industries/CrackMapExec/wiki)

## 维护与更新

本技能基于最新的横向移动实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 3
**工具数量**: 15+

---

*本技能旨在帮助安全研究人员、渗透测试人员和红队成员系统性地进行横向移动。*