---
name: ransomware-defense
description: Comprehensive ransomware defense framework covering ransomware attack chain analysis and detection, ransomware incident response and recovery, and anti-ransomware hardening and backup strategies. Use when defending against ransomware attacks, implementing ransomware prevention controls, or conducting ransomware incident response.
---

# 勒索软件防御技能

## 概述

本技能提供系统化的勒索软件防御框架，涵盖勒索软件攻击链分析与检测、勒索软件应急响应与恢复和反勒索软件加固与备份策略等核心能力。

**核心能力:**
- 攻击链分析 - 勒索软件攻击链分析与检测
- 应急响应 - 勒索软件应急响应与恢复
- 防御加固 - 反勒索软件加固与备份策略

**适用场景:**
- 勒索软件防御
- 勒索软件检测
- 勒索软件应急响应
- 勒索软件恢复
- 勒索软件预防

## 快速开始

### 勒索软件防御流程

1. **预防加固** - 2-4 小时
   - 实施安全加固
   - 配置备份策略
   - 部署检测工具

2. **检测分析** - 持续进行
   - 监控可疑活动
   - 分析攻击行为
   - 检测勒索软件

3. **应急响应** - 1-2 小时
   - 隔离受感染系统
   - 分析攻击影响
   - 制定恢复计划

4. **恢复验证** - 2-4 小时
   - 恢复系统数据
   - 验证系统完整性
   - 改进防御策略

## 核心技能模块

### 1. 勒索软件攻击链分析与检测 (Ransomware Attack Chain Analysis and Detection)

勒索软件攻击链分析与检测。

**核心内容:**
- 勒索软件攻击链
- 攻击行为分析
- 检测规则配置
- 威胁情报应用

**详细文档:** [勒索软件攻击链分析与检测-RansomwareAttackChainAnalysisDetection.md](skills/勒索软件攻击链分析与检测-RansomwareAttackChainAnalysisDetection.md)

### 2. 勒索软件应急响应与恢复 (Ransomware Incident Response and Recovery)

勒索软件应急响应与恢复。

**核心内容:**
- 应急响应流程
- 系统隔离策略
- 数据恢复方法
- 恢复验证流程

**详细文档:** [勒索软件应急响应与恢复-RansomwareIncidentResponseRecovery.md](skills/勒索软件应急响应与恢复-RansomwareIncidentResponseRecovery.md)

### 3. 反勒索软件加固与备份策略 (Anti-Ransomware Hardening and Backup Strategies)

反勒索软件加固与备份策略。

**核心内容:**
- 系统安全加固
- 备份策略制定
- 不可变备份配置
- 防御策略优化

**详细文档:** [反勒索软件加固与备份策略-AntiRansomwareHardeningBackup.md](skills/反勒索软件加固与备份策略-AntiRansomwareHardeningBackup.md)

## 工具和资源

### 推荐工具套件

**检测工具:**
- [YARA](https://yara.readthedocs.io/) - 检测规则
- [Sigma](https://github.com/SigmaHQ/sigma) - 检测规则
- [Sysmon](https://docs.microsoft.com/en-us/sysinternals/downloads/sysmon) - 系统监控
- [Velociraptor](https://www.velocidex.com/) - 取证收集

**备份工具:**
- [Veeam](https://www.veeam.com/) - 备份恢复
- [Rubrik](https://www.rubrik.com/) - 云备份
- [Cohesity](https://www.cohesity.com/) - 数据管理
- [Druva](https://www.druva.com/) - 云备份

**取证工具:**
- [KAPE](https://www.kroll.com/) - 取证收集
- [Velociraptor](https://www.velocidex.com/) - 取证收集
- [Autopsy](https://www.autopsy.com/) - 取证分析
- [FTK Imager](https://accessdata.com/) - 取证分析

**识别工具:**
- [ID Ransomware](https://id-ransomware.malwarehunterteam.com/) - 勒索软件识别
- [NoMoreRansom](https://www.nomoreransom.org/) - 解密工具
- [Emsisoft Decryptor](https://www.emsisoft.com/ransomware-decryption-tools) - 解密工具

### 学习资源

- [NoMoreRansom](https://www.nomoreransom.org/) - 勒索软件防御
- [CISA Ransomware](https://www.cisa.gov/ransomware) - CISA勒索软件
- [NIST Ransomware](https://www.nist.gov/) - NIST勒索软件
- [SANS Ransomware](https://www.sans.org/white-papers/) - SANS勒索软件

## 最佳实践

### 勒索软件防御原则

1. **预防优先** - 预防优于检测
2. **备份关键** - 备份是关键
3. **快速响应** - 快速响应攻击
4. **持续改进** - 持续改进防御
5. **员工培训** - 培训员工意识

### 常见误区

- 没有备份策略
- 备份不可靠
- 响应时间过长
- 缺乏员工培训
- 没有持续改进

## 参考文档

### 标准和框架
- [NIST SP 800-53](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final) - 安全控制
- [CISA Ransomware](https://www.cisa.gov/ransomware) - CISA勒索软件
- [ISO 27001](https://www.iso.org/isoiec-27001-information-security.html) - 信息安全

### 技术指南
- [NoMoreRansom](https://www.nomoreransom.org/) - 勒索软件防御
- [CISA Ransomware Guide](https://www.cisa.gov/ransomware) - CISA勒索软件指南
- [NIST Ransomware Guide](https://www.nist.gov/) - NIST勒索软件指南

## 维护与更新

本技能基于最新的勒索软件防御实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 3
**工具数量**: 15+

---

*本技能旨在帮助安全运营人员、事件响应人员和系统管理员系统性地进行勒索软件防御。*