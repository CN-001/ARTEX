---
name: reconnaissance
description: Comprehensive information gathering and reconnaissance framework covering passive OSINT, active scanning, DNS enumeration, subdomain discovery, and tech stack identification. Use when conducting penetration testing, security assessments, or threat intelligence gathering.
---

# 信息搜集技能

## 概述

本技能提供系统化的信息搜集框架，涵盖被动信息搜集、主动探测、DNS枚举、子域名发现和技术栈识别等核心能力。

**核心能力:**
- 被动信息搜集（OSINT）- 不直接接触目标系统
- 主动信息搜集 - 直接探测目标系统
- DNS枚举 - 发现域名解析信息
- 子域名探测 - 发现隐藏的子域名
- 网络空间搜索引擎 - 利用专业搜索引擎
- 社会工程学信息 - 收集人员和组织信息
- 目标技术栈识别 - 识别目标使用的技术

**适用场景:**
- 渗透测试前期信息搜集
- 安全评估和威胁建模
- 红队演练和蓝队防御
- 竞争情报收集
- 应急响应溯源

## 快速开始

### 被动信息搜集流程

1. **Google Hacking** - 5-10 分钟
   - 使用 Google Dork 搜索敏感信息
   - 发现泄露的配置文件和文档
   - 识别公开的管理后台

2. **网络空间搜索引擎** - 10-20 分钟
   - 使用 Shodan、Fofa、ZoomEye 搜索
   - 发现开放的端口和服务
   - 识别目标的基础设施

3. **社交媒体和人员信息** - 15-30 分钟
   - 使用 theHarvester 收集邮箱和社交媒体
   - 使用 Sherlock 查找用户名
   - 构建目标组织的人员画像

### 主动信息搜集流程

1. **端口扫描** - 10-30 分钟
   - 使用 Nmap 进行端口和服务发现
   - 使用 Masscan 进行快速扫描
   - 识别开放的服务和版本

2. **DNS枚举** - 5-15 分钟
   - 使用 dnsenum 进行 DNS 枚举
   - 使用 dig 和 nslookup 进行 DNS 查询
   - 发现 DNS 服务器和记录

3. **子域名探测** - 15-45 分钟
   - 使用 Subfinder 进行被动子域名发现
   - 使用 Amass 进行深度子域名枚举
   - 使用 Sublist3r 进行快速扫描

4. **技术栈识别** - 10-20 分钟
   - 使用 WhatWeb 识别 Web 技术
   - 使用 Wappalyzer 分析浏览器指纹
   - 使用 Nmap NSE 进行服务指纹识别

## 核心技能模块

### 1. 被动信息搜集 (Passive Reconnaissance)

不直接与目标系统交互，通过公开渠道收集目标信息。

**核心工具:**
- Google Dork - 高级搜索语法
- Shodan - 物联网搜索引擎
- Maltego - 可视化情报分析
- theHarvester - 邮箱和社交媒体收集

**详细文档:** [被动信息搜集-PassiveRecon.md](skills/被动信息搜集-PassiveRecon.md)

### 2. 主动信息搜集 (Active Reconnaissance)

直接与目标系统交互，探测系统信息。

**核心工具:**
- Nmap - 网络映射和端口扫描
- Masscan - 高速端口扫描器
- Netcat - 网络调试工具

**详细文档:** [主动信息搜集-ActiveRecon.md](skills/主动信息搜集-ActiveRecon.md)

### 3. DNS枚举 (DNS Enumeration)

发现目标域名的 DNS 解析信息。

**核心工具:**
- dnsenum - DNS 枚举工具
- dig - DNS 查询工具
- nslookup - DNS 查询工具
- DNSRecon - DNS 侦察框架

**详细文档:** [DNS枚举-DNSEnumeration.md](skills/DNS枚举-DNSEnumeration.md)

### 4. 子域名探测 (Subdomain Discovery)

发现目标的所有子域名。

**核心工具:**
- Subfinder - 被动子域名发现
- Amass - 深度子域名枚举
- Sublist3r - 快速子域名扫描
- Assetfinder - 子域名发现工具

**详细文档:** [子域名探测-SubdomainDiscovery.md](skills/子域名探测-SubdomainDiscovery.md)

### 5. 网络空间搜索引擎 (OSINT Search Engine)

利用专业搜索引擎收集目标信息。

**核心平台:**
- Shodan - 物联网设备搜索引擎
- Fofa - 中国网络空间搜索引擎
- ZoomEye - 网络空间测绘
- Censys - 互联网资产发现

**详细文档:** [网络空间搜索引擎-OSINT-SearchEngine.md](skills/网络空间搜索引擎-OSINT-SearchEngine.md)

### 6. 社会工程学信息 (Social Engineering Information)

收集目标组织和人员的相关信息。

**核心工具:**
- theHarvester - 邮箱和社交媒体收集
- Sherlock - 用户名搜索
- Maltego - 可视化情报分析
- SpiderFoot - OSINT 自动化工具

**详细文档:** [社会工程学信息-SocialEngineeringInfo.md](skills/社会工程学信息-SocialEngineeringInfo.md)

### 7. 目标技术栈识别 (Tech Stack Fingerprint)

识别目标使用的技术栈和框架。

**核心工具:**
- WhatWeb - Web 技术识别
- Wappalyzer - 浏览器扩展
- Nmap NSE - 服务指纹识别
- BuiltWith - 技术栈分析

**详细文档:** [目标技术栈识别-TechStackFingerprint.md](skills/目标技术栈识别-TechStackFingerprint.md)

## 工具和资源

### 推荐工具套件

**被动信息搜集:**
- [OSINT Framework](https://osintframework.com/) - OSINT 工具集合
- [Recon-ng](https://github.com/lanmaster53/recon-ng) - 模块化侦察框架
- [Maltego](https://www.maltego.com/) - 可视化情报分析

**主动信息搜集:**
- [Nmap](https://nmap.org/) - 网络映射器
- [Masscan](https://github.com/robertdavidgraham/masscan) - 高速端口扫描
- [Unicornscan](https://github.com/unicornscan/unicornscan) - 端口扫描器

**DNS 和子域名:**
- [Subfinder](https://github.com/projectdiscovery/subfinder) - 被动子域名发现
- [Amass](https://github.com/OWASP/Amass) - 深度子域名枚举
- [DNSRecon](https://github.com/darkoperator/dnsrecon) - DNS 侦察框架

**技术栈识别:**
- [WhatWeb](https://github.com/urbanadventurer/WhatWeb) - Web 技术识别
- [Wappalyzer](https://www.wappalyzer.com/) - 浏览器扩展
- [BuiltWith](https://builtwith.com/) - 在线技术栈分析

### 学习资源

- [Google Hacking Database (GHDB)](https://www.exploit-db.com/google-hacking-database)
- [OWASP Reconnaissance](https://owasp.org/www-community/Reconnaissance)
- [PTES (Penetration Testing Execution Standard)](https://www.pentest-standard.org/)

## 最佳实践

### 信息搜集原则

1. **合法性** - 确保所有活动获得授权
2. **隐蔽性** - 被动搜集优先，避免被检测
3. **系统性** - 按照流程逐步收集信息
4. **完整性** - 确保信息覆盖所有攻击面
5. **时效性** - 及时更新收集到的信息

### 常见误区

- 过度依赖自动化工具
- 忽略被动信息搜集的价值
- 信息收集不完整
- 没有建立信息管理机制
- 违反法律法规和道德准则

## 参考文档

### 标准和框架
- [PTES - Penetration Testing Execution Standard](https://www.pentest-standard.org/)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)
- [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)

### 技术指南
- [Google Hacking 教程](https://www.exploit-db.com/google-hacking-database)
- [Nmap 官方文档](https://nmap.org/book/man.html)
- [DNS 安全最佳实践](https://www.icann.org/en/system/files/files/dns-security-best-practices-28may14-en.pdf)

## 维护与更新

本技能基于最新的渗透测试实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 7
**工具数量**: 30+

---

*本技能旨在帮助安全研究人员、渗透测试人员和红队成员系统性地进行信息搜集。*