---
name: mobile-security
description: Comprehensive mobile security framework covering Android security testing and iOS security testing. Use when conducting mobile application security assessments, penetration testing, or mobile security research.
---

# 移动安全技能

## 概述

本技能提供系统化的移动安全框架，涵盖Android安全测试和iOS安全测试等核心能力。

**核心能力:**
- Android安全测试 - Android应用安全测试
- iOS安全测试 - iOS应用安全测试
- 静态分析 - 应用静态分析
- 动态分析 - 应用动态分析
- 流量分析 - 网络流量分析
- 数据存储分析 - 数据存储安全

**适用场景:**
- 移动应用安全测试
- 移动渗透测试
- 移动安全评估
- 移动漏洞研究
- 移动安全培训

## 快速开始

### Android安全测试流程

1. **信息收集** - 10-20 分钟
   - 收集应用信息
   - 分析应用权限
   - 识别应用组件

2. **静态分析** - 20-40 分钟
   - 反编译应用
   - 分析代码逻辑
   - 检查安全配置

3. **动态分析** - 30-60 分钟
   - Hook应用函数
   - 监控应用行为
   - 测试安全漏洞

### iOS安全测试流程

1. **信息收集** - 10-20 分钟
   - 收集应用信息
   - 分析应用权限
   - 识别应用组件

2. **静态分析** - 20-40 分钟
   - 反编译应用
   - 分析代码逻辑
   - 检查安全配置

3. **动态分析** - 30-60 分钟
   - Hook应用函数
   - 监控应用行为
   - 测试安全漏洞

## 核心技能模块

### 1. Android安全测试 (Android Security Testing)

对Android应用进行安全测试。

**核心工具:**
- drozer - Android安全测试框架
- Frida - 动态插桩工具
- APKTool - APK反编译工具
- MobSF - 移动安全框架

**详细文档:** [Android安全测试-AndroidSecurityTest.md](skills/Android安全测试-AndroidSecurityTest.md)

### 2. iOS安全测试 (iOS Security Testing)

对iOS应用进行安全测试。

**核心工具:**
- Frida - 动态插桩工具
- Objection - 运行时移动探索
- class-dump - Objective-C类信息
- Ghidra - 逆向工程工具

**详细文档:** [iOS安全测试-iOSSecurityTest.md](skills/iOS安全测试-iOSSecurityTest.md)

## 工具和资源

### 推荐工具套件

**Android工具:**
- [drozer](https://github.com/WithSecureLabs/drozer) - Android安全测试框架
- [Frida](https://frida.re/) - 动态插桩工具
- [APKTool](https://ibotpeaches.github.io/Apktool/) - APK反编译工具
- [MobSF](https://github.com/MobSF/Mobile-Security-Framework-MobSF) - 移动安全框架

**iOS工具:**
- [Frida](https://frida.re/) - 动态插桩工具
- [Objection](https://github.com/sensepost/objection) - 运行时移动探索
- [class-dump](https://github.com/nygard/class-dump) - Objective-C类信息
- [Ghidra](https://ghidra-sre.org/) - 逆向工程工具

**通用工具:**
- [Burp Suite](https://portswigger.net/burp) - Web安全测试
- [Wireshark](https://www.wireshark.org/) - 网络协议分析
- [Charles Proxy](https://www.charlesproxy.com/) - HTTP代理
- [mitmproxy](https://mitmproxy.org/) - HTTP代理

### 学习资源

- [OWASP Mobile Security Testing Guide](https://owasp.org/www-project-mobile-security-testing-guide/)
- [Mobile Security Framework (MobSF)](https://github.com/MobSF/Mobile-Security-Framework-MobSF)
- [Frida Documentation](https://frida.re/docs/)
- [NIST SP 800-163](https://csrc.nist.gov/publications/detail/sp/800-163/rev-1/final)

## 最佳实践

### 移动安全测试原则

1. **授权性** - 确保获得测试授权
2. **完整性** - 覆盖所有测试场景
3. **准确性** - 确保测试结果准确
4. **可追溯性** - 记录所有测试活动
5. **合规性** - 遵守法律法规

### 常见误区

- 过度依赖自动化工具
- 忽略业务逻辑测试
- 没有测试数据存储
- 缺乏流量分析
- 违反道德和法律

## 参考文档

### 标准和框架
- [OWASP Mobile Security Testing Guide](https://owasp.org/www-project-mobile-security-testing-guide/)
- [OWASP Mobile Top 10](https://owasp.org/www-project-mobile-top-10/)
- [NIST SP 800-163](https://csrc.nist.gov/publications/detail/sp/800-163/rev-1/final)

### 技术指南
- [MobSF Documentation](https://github.com/MobSF/Mobile-Security-Framework-MobSF)
- [Frida Documentation](https://frida.re/docs/)
- [drozer Documentation](https://github.com/WithSecureLabs/drozer)

## 维护与更新

本技能基于最新的移动安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 2
**工具数量**: 15+

---

*本技能旨在帮助安全研究人员、渗透测试人员和移动安全工程师系统性地进行移动安全测试。*