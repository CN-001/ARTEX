---
name: incident-response
description: Comprehensive incident response framework covering incident triage, log collection and analysis, memory forensics, disk forensics and data recovery, network traffic analysis, malware sample analysis, containment and eradication, ransomware incident response, cloud incident response, threat hunting, SOAR automation, and lessons learned reporting. Use when conducting incident response, digital forensics, or security operations.
---

# 应急响应技能

## 概述

本技能提供系统化的应急响应框架，涵盖事件分类与优先级评估、日志收集与分析、内存取证分析、磁盘取证与数据恢复、网络流量分析、恶意样本分析、事件遏制与清除、勒索软件应急响应、云环境应急响应、威胁狩猎、SOAR自动化响应和事件复盘与报告等核心能力。

**核心能力:**
- 事件分类 - 事件分类和优先级评估
- 日志分析 - 日志收集和分析
- 内存取证 - 内存取证分析
- 磁盘取证 - 磁盘取证和数据恢复
- 流量分析 - 网络流量分析
- 恶意样本 - 恶意样本分析
- 事件遏制 - 事件遏制和清除
- 勒索软件 - 勒索软件应急响应
- 云环境 - 云环境应急响应
- 威胁狩猎 - 威胁狩猎技术
- SOAR - 自动化响应
- 事件复盘 - 事件复盘和报告

**适用场景:**
- 安全事件响应
- 数字取证分析
- 安全运营中心
- 威胁狩猎
- 应急响应培训

## 快速开始

### 应急响应流程

1. **事件检测** - 10-30 分钟
   - 检测安全事件
   - 评估事件严重性
   - 确定响应优先级

2. **事件分析** - 30-60 分钟
   - 收集事件证据
   - 分析事件影响
   - 确定攻击来源

3. **事件遏制** - 15-30 分钟
   - 遏制事件扩散
   - 隔离受影响系统
   - 保护关键资产

4. **事件清除** - 30-60 分钟
   - 清除恶意软件
   - 修复安全漏洞
   - 恢复系统正常

5. **事件恢复** - 1-2 小时
   - 恢复系统功能
   - 验证系统安全
   - 恢复业务运营

6. **事件复盘** - 1-2 天
   - 分析事件根因
   - 总结经验教训
   - 改进响应流程

## 核心技能模块

### 1. 事件分类与优先级评估 (Incident Triage)

事件分类和优先级评估。

**核心工具:**
- 事件管理平台 - 事件管理
- SIRP工具 - 安全事件响应平台
- 分类标准 - 事件分类标准

**详细文档:** [事件分类与优先级评估-IncidentTriage.md](skills/事件分类与优先级评估-IncidentTriage.md)

### 2. 日志收集与分析 (Log Collection and Analysis)

日志收集和分析。

**核心工具:**
- ELK Stack - 日志分析平台
- Splunk - 日志分析平台
- Wazuh - 安全监控平台

**详细文档:** [日志收集与分析-LogCollectionAnalysis.md](skills/日志收集与分析-LogCollectionAnalysis.md)

### 3. 内存取证分析 (Memory Forensics)

内存取证分析。

**核心工具:**
- Volatility - 内存取证工具
- Rekall - 内存取证框架
- LiME - Linux内存提取

**详细文档:** [内存取证分析-MemoryForensics.md](skills/内存取证分析-MemoryForensics.md)

### 4. 磁盘取证与数据恢复 (Disk Forensics and Data Recovery)

磁盘取证和数据恢复。

**核心工具:**
- FTK Imager - 磁盘镜像工具
- Autopsy - 取证分析平台
- dd - 磁盘复制工具

**详细文档:** [磁盘取证与数据恢复-DiskForensicsDataRecovery.md](skills/磁盘取证与数据恢复-DiskForensicsDataRecovery.md)

### 5. 网络流量分析 (Network Traffic Analysis)

网络流量分析。

**核心工具:**
- Wireshark - 网络协议分析
- Zeek - 网络安全监控
- Suricata - 入侵检测系统

**详细文档:** [网络流量分析-NetworkTrafficAnalysis.md](skills/网络流量分析-NetworkTrafficAnalysis.md)

### 6. 恶意样本分析 (Malware Sample Analysis)

恶意样本分析。

**核心工具:**
- Cuckoo Sandbox - 恶意软件沙箱
- Ghidra - 逆向工程工具
- YARA - 恶意软件识别

**详细文档:** [恶意样本分析-MalwareSampleAnalysis.md](skills/恶意样本分析-MalwareSampleAnalysis.md)

### 7. 事件遏制与清除 (Containment and Eradication)

事件遏制和清除。

**核心工具:**
- Firewall ACL - 防火墙规则
- EDR隔离 - 端点检测和响应
- GPO - 组策略管理

**详细文档:** [事件遏制与清除-ContainmentEradication.md](skills/事件遏制与清除-ContainmentEradication.md)

### 8. 勒索软件应急响应 (Ransomware Incident Response)

勒索软件应急响应。

**核心工具:**
- EDR工具 - 端点检测和响应
- 解密工具 - 勒索软件解密
- 备份恢复 - 备份和恢复

**详细文档:** [勒索软件应急响应-RansomwareIncidentResponse.md](skills/勒索软件应急响应-RansomwareIncidentResponse.md)

### 9. 云环境应急响应 (Cloud Incident Response)

云环境应急响应。

**核心工具:**
- AWS GuardDuty - AWS威胁检测
- Azure Sentinel - Azure安全监控
- 云安全工具 - 云安全平台

**详细文档:** [云环境应急响应-CloudIncidentResponse.md](skills/云环境应急响应-CloudIncidentResponse.md)

### 10. 威胁狩猎 (Threat Hunting)

威胁狩猎技术。

**核心工具:**
- KQL - Kusto查询语言
- YARA - 恶意软件识别
- Sigma - 日志检测规则
- MITRE ATT&CK - 威胁框架

**详细文档:** [威胁狩猎-ThreatHunting.md](skills/威胁狩猎-ThreatHunting.md)

### 11. SOAR自动化响应 (SOAR Automation)

SOAR自动化响应。

**核心工具:**
- Splunk SOAR - 安全编排自动化
- Shuffle - 工作流自动化
- Tines - 安全自动化平台

**详细文档:** [SOAR自动化响应-SOARAutomation.md](skills/SOAR自动化响应-SOARAutomation.md)

### 12. 事件复盘与报告 (Lessons Learned and Reporting)

事件复盘和报告。

**核心工具:**
- 复盘模板 - 事件复盘模板
- Root Cause分析 - 根因分析工具
- 报告工具 - 报告生成工具

**详细文档:** [事件复盘与报告-LessonsLearnedReporting.md](skills/事件复盘与报告-LessonsLearnedReporting.md)

### 13. AI安全应急响应 (AI Security Incident Response)

AI安全应急响应。

**核心工具:**
- AI IR Framework - AI应急响应框架
- MITRE ATLAS - AI威胁框架
- TheHive - 安全事件平台

**详细文档:** [AI安全应急响应-AISecurityIncidentResponse.md](skills/AI安全应急响应-AISecurityIncidentResponse.md)

## 工具和资源

### 推荐工具套件

**取证工具:**
- [Volatility](https://www.volatilityfoundation.org/) - 内存取证
- [Autopsy](https://www.autopsy.com/) - 取证分析
- [FTK Imager](https://www.accessdata.com/products/ftk-imager) - 磁盘镜像
- [Wireshark](https://www.wireshark.org/) - 网络分析

**恶意软件分析:**
- [Cuckoo Sandbox](https://cuckoosandbox.org/) - 恶意软件沙箱
- [Ghidra](https://ghidra-sre.org/) - 逆向工程
- [YARA](https://virustotal.github.io/yara/) - 恶意软件识别

**日志分析:**
- [ELK Stack](https://www.elastic.co/) - 日志分析平台
- [Splunk](https://www.splunk.com/) - 日志分析平台
- [Wazuh](https://wazuh.com/) - 安全监控平台

**SOAR工具:**
- [Splunk SOAR](https://www.splunk.com/en_us/products/splunk-soar.html) - 安全编排
- [Shuffle](https://shuffle.dev/) - 工作流自动化
- [Tines](https://www.tines.com/) - 安全自动化

### 学习资源

- [NIST SP 800-61 Rev 2](https://csrc.nist.gov/publications/detail/sp/800-61/rev-2/final)
- [SANS Incident Handling](https://www.sans.org/course/for508)
- [ISO/IEC 27035](https://www.iso.org/standard/78973.html)
- [FIRST](https://www.first.org/)
- [MITRE ATT&CK](https://attack.mitre.org/)

## 最佳实践

### 应急响应原则

1. **及时性** - 快速响应事件
2. **准确性** - 准确分析事件
3. **完整性** - 完整记录过程
4. **协调性** - 协调各方资源
5. **合规性** - 遵守法律法规

### 常见误区

- 响应时间过长
- 缺乏事件分类
- 没有遏制措施
- 缺乏取证意识
- 没有复盘总结

## 参考文档

### 标准和框架
- [NIST SP 800-61 Rev 2](https://csrc.nist.gov/publications/detail/sp/800-61/rev-2/final)
- [SANS PICERL](https://www.sans.org/white-papers/33901/)
- [ISO/IEC 27035](https://www.iso.org/standard/78973.html)
- [FIRST](https://www.first.org/)

### 技术指南
- [Volatility Documentation](https://www.volatilityfoundation.org/documentation)
- [Wireshark Documentation](https://www.wireshark.org/docs/)
- [Splunk Documentation](https://docs.splunk.com/)

## 维护与更新

本技能基于最新的应急响应实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 13
**工具数量**: 40+

---

*本技能旨在帮助安全响应人员、数字取证专家和安全运营人员系统性地进行应急响应。*