---
name: threat-hunting
description: Comprehensive threat hunting framework covering threat hunting methodology and hypothesis-driven hunting, Sigma rule engineering, MITRE ATT&CK-based threat hunting, and network traffic and log anomaly detection. Use when conducting threat hunting, security operations, or incident response.
---

# 威胁狩猎技能

## 概述

本技能提供系统化的威胁狩猎框架，涵盖威胁狩猎方法论与假设驱动、Sigma规则检测工程、基于ATT&CK的威胁狩猎和网络流量与日志异常检测等核心能力。

**核心能力:**
- 狩猎方法论 - 威胁狩猎方法论与假设驱动
- Sigma规则 - Sigma规则检测工程
- ATT&CK狩猎 - 基于ATT&CK的威胁狩猎
- 异常检测 - 网络流量与日志异常检测

**适用场景:**
- 威胁狩猎
- 安全运营
- 事件响应
- 安全分析
- 威胁情报

## 快速开始

### 威胁狩猎流程

1. **假设生成** - 1-2 小时
   - 分析威胁情报
   - 生成狩猎假设
   - 确定狩猎目标

2. **数据收集** - 2-4 小时
   - 收集相关数据
   - 准备数据环境
   - 验证数据质量

3. **狩猎分析** - 4-8 小时
   - 执行狩猎分析
   - 检测异常行为
   - 验证威胁存在

4. **结果报告** - 1-2 小时
   - 分析狩猎结果
   - 编写狩猎报告
   - 改进狩猎策略

## 核心技能模块

### 1. 威胁狩猎方法论与假设驱动 (Threat Hunting Methodology and Hypothesis-driven Hunting)

威胁狩猎方法论与假设驱动。

**核心内容:**
- 狩猎方法论
- 假设生成技术
- 狩猎策略制定
- 狩猎效果评估

**详细文档:** [威胁狩猎方法论与假设驱动-ThreatHuntingMethodology.md](skills/威胁狩猎方法论与假设驱动-ThreatHuntingMethodology.md)

### 2. Sigma规则检测工程 (Sigma Rule Engineering)

Sigma规则检测工程。

**核心内容:**
- Sigma规则语法
- 规则编写最佳实践
- 规则测试和验证
- 规则管理和维护

**详细文档:** [Sigma规则检测工程-SigmaRuleEngineering.md](skills/Sigma规则检测工程-SigmaRuleEngineering.md)

### 3. 基于ATT&CK的威胁狩猎 (MITRE ATT&CK-based Threat Hunting)

基于ATT&CK的威胁狩猎。

**核心内容:**
- ATT&CK框架应用
- 基于TTP的狩猎
- ATT&CK覆盖度分析
- 狩猎策略优化

**详细文档:** [ATT&CK威胁狩猎-MITREAttackHunting.md](skills/ATT&CK威胁狩猎-MITREAttackHunting.md)

### 4. 网络流量与日志异常检测 (Network Traffic and Log Anomaly Detection)

网络流量与日志异常检测。

**核心内容:**
- 流量异常检测
- 日志异常检测
- 行为基线建立
- 异常行为分析

**详细文档:** [流量与日志异常检测-TrafficLogAnomalyDetection.md](skills/流量与日志异常检测-TrafficLogAnomalyDetection.md)

## 工具和资源

### 推荐工具套件

**狩猎分析工具:**
- [Jupyter Notebooks](https://jupyter.org/) - 狩猎分析环境
- [Elastic Stack](https://www.elastic.co/) - 日志分析与狩猎
- [Splunk](https://www.splunk.com/) - 日志分析与狩猎
- [Sigma](https://github.com/SigmaHQ/sigma) - 通用检测规则

**网络分析工具:**
- [Zeek](https://zeek.org/) - 网络流量分析
- [Suricata](https://suricata.io/) - 入侵检测系统
- [Wireshark](https://www.wireshark.org/) - 网络协议分析

**ATT&CK工具:**
- [MITRE ATT&CK Navigator](https://attack.mitre.org/resources/attack-navigator/) - 覆盖度分析
- [Atomic Red Team](https://atomicredteam.io/) - 原子红队
- [CALDERA](https://caldera.mitre.org/) - 对抗模拟

### 学习资源

- [MITRE ATT&CK](https://attack.mitre.org/) - ATT&CK框架
- [SigmaHQ](https://github.com/SigmaHQ/sigma) - Sigma规则
- [Threat Hunting Methodology](https://www.sans.org/white-papers/) - 狩猎方法论
- [Elastic Security](https://www.elastic.co/guide/en/security/current/index.html) - Elastic安全

## 最佳实践

### 威胁狩猎原则

1. **假设驱动** - 基于假设进行狩猎
2. **数据驱动** - 基于数据进行分析
3. **持续改进** - 持续改进狩猎策略
4. **协作共享** - 协作共享狩猎结果
5. **文档记录** - 记录狩猎过程

### 常见误区

- 没有明确假设
- 过度依赖自动化
- 缺乏数据质量
- 没有持续改进
- 缺乏文档记录

## 参考文档

### 标准和框架
- [MITRE ATT&CK](https://attack.mitre.org/) - ATT&CK框架
- [NIST SP 800-61](https://csrc.nist.gov/publications/detail/sp/800-61/rev-2/final) - 事件响应
- [ISO 27001](https://www.iso.org/isoiec-27001-information-security.html) - 信息安全

### 技术指南
- [Sigma Documentation](https://github.com/SigmaHQ/sigma/wiki) - Sigma文档
- [Elastic Security](https://www.elastic.co/guide/en/security/current/index.html) - Elastic安全
- [Splunk Security](https://docs.splunk.com/Documentation/Security) - Splunk安全

## 维护与更新

本技能基于最新的威胁狩猎实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 4
**工具数量**: 15+

---

*本技能旨在帮助安全分析师、威胁狩猎专家和安全运营人员系统性地进行威胁狩猎。*