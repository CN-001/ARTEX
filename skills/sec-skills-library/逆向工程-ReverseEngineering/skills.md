---
name: reverse-engineering
description: Comprehensive reverse engineering framework covering static reverse analysis, dynamic debug analysis, and malware analysis. Use when conducting binary analysis, malware research, or vulnerability discovery.
---

# 逆向工程技能

## 概述

本技能提供系统化的逆向工程框架，涵盖静态逆向分析、动态调试分析和恶意软件分析等核心能力。

**核心能力:**
- 静态逆向分析 - 静态分析二进制程序
- 动态调试分析 - 动态调试程序行为
- 恶意软件分析 - 分析恶意软件
- 反汇编 - 反汇编二进制代码
- 反编译 - 反编译二进制代码
- 调试 - 调试程序执行
- 行为分析 - 分析程序行为

**适用场景:**
- 恶意软件分析
- 漏洞研究和发现
- 软件安全评估
- 知识产权保护
- 安全培训和教育

## 快速开始

### 静态逆向分析流程

1. **文件识别** - 5-10 分钟
   - 识别文件类型
   - 分析文件结构
   - 收集文件信息

2. **反汇编** - 15-30 分钟
   - 反汇编二进制代码
   - 分析汇编代码
   - 识别关键函数

3. **反编译** - 20-40 分钟
   - 反编译二进制代码
   - 分析高级代码
   - 理解程序逻辑

### 动态调试分析流程

1. **环境准备** - 10-20 分钟
   - 准备调试环境
   - 配置调试工具
   - 加载目标程序

2. **断点设置** - 15-30 分钟
   - 设置关键断点
   - 配置断点条件
   - 准备监控变量

3. **调试执行** - 30-60 分钟
   - 单步执行程序
   - 监控程序状态
   - 分析程序行为

## 核心技能模块

### 1. 静态逆向分析 (Static Reverse Analysis)

静态分析二进制程序。

**核心工具:**
- IDA Pro - 专业反汇编工具
- Ghidra - 开源逆向工程平台
- Radare2 - 开源逆向框架
- Binary Ninja - 现代逆向工具

**详细文档:** [静态逆向分析-StaticReverseAnalysis.md](skills/静态逆向分析-StaticReverseAnalysis.md)

### 2. 动态调试分析 (Dynamic Debug Analysis)

动态调试程序行为。

**核心工具:**
- x64dbg - Windows调试器
- OllyDbg - Windows调试器
- GDB - GNU调试器
- WinDbg - Windows调试器

**详细文档:** [动态调试分析-DynamicDebugAnalysis.md](skills/动态调试分析-DynamicDebugAnalysis.md)

### 3. 恶意软件分析 (Malware Analysis)

分析恶意软件。

**核心工具:**
- Cuckoo Sandbox - 恶意软件沙箱
- CAPE - 恶意软件分析平台
- 动态调试 - 动态分析工具
- 静态分析 - 静态分析工具

**详细文档:** [恶意软件分析-MalwareAnalysis.md](skills/恶意软件分析-MalwareAnalysis.md)

## 工具和资源

### 推荐工具套件

**静态分析工具:**
- [IDA Pro](https://hex-rays.com/ida-pro/) - 专业反汇编工具
- [Ghidra](https://ghidra-sre.org/) - 开源逆向工程平台
- [Radare2](https://rada.re/n/) - 开源逆向框架
- [Binary Ninja](https://binary.ninja/) - 现代逆向工具

**动态调试工具:**
- [x64dbg](https://x64dbg.com/) - Windows调试器
- [OllyDbg](http://www.ollydbg.de/) - Windows调试器
- [GDB](https://www.gnu.org/software/gdb/) - GNU调试器
- [WinDbg](https://learn.microsoft.com/en-us/windows-hardware/drivers/debugger/) - Windows调试器

**恶意软件分析工具:**
- [Cuckoo Sandbox](https://cuckoosandbox.org/) - 恶意软件沙箱
- [CAPE](https://github.com/kevoreilly/CAPEv2) - 恶意软件分析平台
- [VirusTotal](https://www.virustotal.com/) - 病毒检测平台
- [Hybrid Analysis](https://www.hybrid-analysis.com/) - 恶意软件分析

**辅助工具:**
- [PEiD](https://www.aldeid.com/wiki/PEiD) - PE文件分析
- [Resource Hacker](http://www.angusj.com/resourcehacker/) - 资源编辑器
- [Process Monitor](https://learn.microsoft.com/en-us/sysinternals/downloads/procmon) - 进程监控
- [Wireshark](https://www.wireshark.org/) - 网络协议分析

### 学习资源

- [Reverse Engineering Reddit](https://www.reddit.com/r/ReverseEngineering/)
- [Malware Analysis Reddit](https://www.reddit.com/r/Malware/)
- [Practical Malware Analysis](https://www.nostarch.com/malware)
- [The IDA Pro Book](https://www.nostarch.com/idapro2.htm)

## 最佳实践

### 逆向工程原则

1. **合法性** - 确保获得逆向授权
2. **道德性** - 遵守道德准则
3. **安全性** - 在隔离环境中分析
4. **可追溯性** - 记录分析过程
5. **合规性** - 遵守法律法规

### 常见误区

- 过度依赖自动化工具
- 忽略动态分析
- 没有理解程序逻辑
- 缺乏文档记录
- 违反道德和法律

## 参考文档

### 标准和框架
- [NIST SP 800-183](https://csrc.nist.gov/publications/detail/sp/800-183/final)
- [OWASP Malware Analysis](https://owasp.org/www-community/attacks/Malware)
- [SANS Malware Analysis](https://www.sans.org/course/for610)

### 技术指南
- [Ghidra Documentation](https://ghidra-sre.org/)
- [IDA Pro Documentation](https://hex-rays.com/products/ida/support/)
- [x64dbg Documentation](https://x64dbg.com/)

## 维护与更新

本技能基于最新的逆向工程实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 3
**工具数量**: 20+

---

*本技能旨在帮助安全研究人员、逆向工程师和恶意软件分析师系统性地进行逆向工程。*