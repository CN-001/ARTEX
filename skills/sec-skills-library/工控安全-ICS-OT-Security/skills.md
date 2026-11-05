---
name: ics-ot-security
description: Comprehensive ICS/OT security framework covering SCADA system security assessment, PLC/RTU security testing, ICS network protocol security, IEC 62443 compliance audit, ICS incident response, and industrial firewall and network segmentation. Use when conducting industrial control system security assessments, ICS penetration testing, or OT security operations.
---

# 工控安全技能

## 概述

本技能提供系统化的工控安全框架，涵盖SCADA系统安全评估、PLC与RTU安全测试、工控网络协议安全、工控安全合规审计（IEC 62443）、工控安全应急预案和工业防火墙与网络分段等核心能力。

**核心能力:**
- SCADA安全 - SCADA系统安全评估
- PLC/RTU安全 - PLC与RTU安全测试
- 协议安全 - 工控网络协议安全
- 合规审计 - IEC 62443合规审计
- 应急响应 - 工控安全应急预案
- 网络分段 - 工业防火墙和网络分段

**适用场景:**
- 工控安全评估
- ICS渗透测试
- OT安全运营
- 工控安全合规
- 工控安全培训

## 快速开始

### 工控安全评估流程

1. **资产发现** - 1-2 天
   - 发现工控资产
   - 分析工控架构
   - 识别关键设备

2. **安全评估** - 2-3 天
   - 评估SCADA系统
   - 测试PLC/RTU设备
   - 分析网络协议

3. **合规审计** - 1-2 天
   - 审计IEC 62443合规性
   - 检查安全配置
   - 评估安全策略

4. **安全加固** - 2-3 天
   - 修复安全漏洞
   - 加强安全配置
   - 改进安全策略

## 核心技能模块

### 1. SCADA系统安全评估 (SCADA System Security Assessment)

SCADA系统安全评估。

**核心工具:**
- PLCScan - PLC扫描工具
- Modbus/TCP Scanner - Modbus扫描
- Shodan - 工控设备搜索

**详细文档:** [SCADA系统安全评估-SCADASecurityAssessment.md](skills/SCADA系统安全评估-SCADASecurityAssessment.md)

### 2. PLC与RTU安全测试 (PLC/RTU Security Testing)

PLC与RTU安全测试。

**核心工具:**
- MSPT - Modbus安全测试工具
- ISF - ICS利用框架
- 自定义工具 - 自定义测试工具

**详细文档:** [PLC与RTU安全测试-PLC-RTU-SecurityTesting.md](skills/PLC与RTU安全测试-PLC-RTU-SecurityTesting.md)

### 3. 工控网络协议安全 (ICS Network Protocol Security)

工控网络协议安全。

**核心工具:**
- Wireshark - 网络协议分析
- Scapy - 网络包操作
- Zabbix - 网络监控

**详细文档:** [工控网络协议安全-ICS-NetworkProtocolSecurity.md](skills/工控网络协议安全-ICS-NetworkProtocolSecurity.md)

### 4. 工控安全合规审计 (IEC 62443 Audit)

工控安全合规审计。

**核心工具:**
- IEC 62443 Checklist - 合规检查清单
- Compliance Scanner - 合规扫描工具
- 审计工具 - 审计工具

**详细文档:** [工控安全合规审计-IEC62443-Audit.md](skills/工控安全合规审计-IEC62443-Audit.md)

### 5. 工控安全应急预案 (ICS Incident Response)

工控安全应急预案。

**核心工具:**
- ICS-CERT Tools - ICS应急工具
- Dragos Platform - 工控安全平台
- 应急预案 - 应急预案模板

**详细文档:** [工控安全应急预案-ICS-IncidentResponse.md](skills/工控安全应急预案-ICS-IncidentResponse.md)

### 6. 工业防火墙与网络分段 (Industrial Firewall and Network Segmentation)

工业防火墙和网络分段。

**核心工具:**
- Tofino Security - 工业防火墙
- Claroty - 工控安全平台
- Nozomi - 工控安全平台

**详细文档:** [工业防火墙与网络分段-Industrial-Firewall-Segmentation.md](skills/工业防火墙与网络分段-Industrial-Firewall-Segmentation.md)

## 工具和资源

### 推荐工具套件

**工控扫描工具:**
- [PLCScan](https://github.com/SCADAFence/PLCScan) - PLC扫描工具
- [Modbus Scanner](https://github.com/sourceperl/Modbus) - Modbus扫描
- [Shodan](https://www.shodan.io/) - 工控设备搜索

**工控安全平台:**
- [Dragos](https://dragos.com/) - 工控安全平台
- [Claroty](https://www.claroty.com/) - 工控安全平台
- [Nozomi Networks](https://www.nozominetworks.com/) - 工控安全平台

**工控防火墙:**
- [Tofino Security](https://www.tofinosecurity.com/) - 工业防火墙
- [Moxa](https://www.moxa.com/) - 工业网络设备
- [Cisco ICS](https://www.cisco.com/c/en/us/solutions/industrial/) - 思科工控解决方案

**协议分析工具:**
- [Wireshark](https://www.wireshark.org/) - 网络协议分析
- [Scapy](https://scapy.net/) - 网络包操作
- [Zabbix](https://www.zabbix.com/) - 网络监控

### 学习资源

- [IEC 62443 Series](https://www.iec.ch/standards/62443)
- [NIST SP 800-82 Rev.3](https://csrc.nist.gov/publications/detail/sp/800-82/rev-3/final)
- [ICS-CERT Advisories](https://www.cisa.gov/ics)
- [MITRE ATT&CK for ICS](https://attack.mitre.org/techniques/ics/)
- [SANS ICS Security](https://www.sans.org/ics-security/)

## 最佳实践

### 工控安全原则

1. **可用性优先** - 确保系统可用性
2. **安全分区** - 实施网络分区
3. **最小权限** - 遵循最小权限原则
4. **持续监控** - 持续监控系统
5. **应急准备** - 准备应急响应

### 常见误区

- 忽略可用性要求
- 直接连接互联网
- 没有网络分区
- 缺乏应急准备
- 忽略合规要求

## 参考文档

### 标准和框架
- [IEC 62443 Series](https://www.iec.ch/standards/62443)
- [NIST SP 800-82 Rev.3](https://csrc.nist.gov/publications/detail/sp/800-82/rev-3/final)
- [ISA/IEC 62443](https://www.isa.org/standards-and-publications/isa-62443)
- [NERC CIP](https://www.nerc.com/pa/Stand/Pages/StandardActivities.aspx)

### 技术指南
- [ICS-CERT Advisories](https://www.cisa.gov/ics)
- [MITRE ATT&CK for ICS](https://attack.mitre.org/techniques/ics/)
- [SANS ICS Security Resources](https://www.sans.org/ics-security/)

## 维护与更新

本技能基于最新的工控安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 6
**工具数量**: 20+

---

*本技能旨在帮助安全研究人员、工控安全工程师和OT安全专家系统性地进行工控安全评估。*