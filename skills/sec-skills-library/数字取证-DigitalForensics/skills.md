---
name: digital-forensics
description: Comprehensive digital forensics framework covering disk imaging and evidence acquisition, memory forensics analysis (Volatility), Windows digital forensics analysis, Linux digital forensics analysis, and browser and email forensics. Use when conducting digital forensics investigations, incident response, or evidence collection.
---

# 数字取证技能

## 概述

本技能提供系统化的数字取证框架，涵盖磁盘镜像与证据获取、内存取证分析(Volatility)、Windows数字取证分析、Linux数字取证分析和浏览器与邮件取证等核心能力。

**核心能力:**
- 磁盘镜像 - 磁盘镜像与证据获取
- 内存取证 - 内存取证分析(Volatility)
- Windows取证 - Windows数字取证分析
- Linux取证 - Linux数字取证分析
- 浏览器取证 - 浏览器与邮件取证

**适用场景:**
- 数字取证调查
- 事件响应取证
- 电子证据收集
- 法律诉讼支持
- 取证分析培训

## 快速开始

### 数字取证流程

1. **现场保护** - 30-60 分钟
   - 保护现场环境
   - 隔离目标系统
   - 记录现场情况

2. **证据收集** - 1-2 小时
   - 收集数字证据
   - 创建磁盘镜像
   - 记录证据链

3. **证据分析** - 2-4 小时
   - 分析磁盘镜像
   - 分析内存数据
   - 分析系统日志

4. **报告编写** - 1-2 小时
   - 编写取证报告
   - 提供分析结论
   - 支持法律诉讼

## 核心技能模块

### 1. 磁盘镜像与证据获取 (Disk Imaging and Evidence Acquisition)

磁盘镜像与证据获取。

**核心工具:**
- FTK Imager - 磁盘镜像工具
- Guymager - 磁盘镜像工具
- dd - 磁盘复制工具
- EnCase - 商业取证工具

**详细文档:** [磁盘镜像与证据获取-DiskImagingEvidenceAcquisition.md](skills/磁盘镜像与证据获取-DiskImagingEvidenceAcquisition.md)

### 2. 内存取证分析(Volatility) (Memory Forensics Analysis - Volatility)

内存取证分析(Volatility)。

**核心工具:**
- Volatility 2/3 - 内存取证工具
- Rekall - 内存取证框架
- LiME - Linux内存提取
- Windows内存分析 - Windows内存工具

**详细文档:** [内存取证分析-Volatility-MemoryForensicsVolatility.md](skills/内存取证分析-Volatility-MemoryForensicsVolatility.md)

### 3. Windows数字取证分析 (Windows Digital Forensics Analysis)

Windows数字取证分析。

**核心工具:**
- KAPE - Windows取证工件收集
- Sleuth Kit - 文件系统分析
- Autopsy - 取证分析平台
- Hayabusa - 事件日志分析

**详细文档:** [Windows数字取证分析-WindowsDigitalForensics.md](skills/Windows数字取证分析-WindowsDigitalForensics.md)

### 4. Linux数字取证分析 (Linux Digital Forensics Analysis)

Linux数字取证分析。

**核心工具:**
- Sleuth Kit - 文件系统分析
- Autopsy - 取证分析平台
- Linux取证工具 - Linux取证工具
- 日志分析工具 - 日志分析工具

**详细文档:** [Linux数字取证分析-LinuxDigitalForensics.md](skills/Linux数字取证分析-LinuxDigitalForensics.md)

### 5. 浏览器与邮件取证 (Browser and Email Forensics)

浏览器与邮件取证。

**核心工具:**
- 浏览器取证工具 - 浏览器取证
- 邮件取证工具 - 邮件取证
- 取证分析工具 - 取证分析
- 自定义工具 - 自定义取证工具

**详细文档:** [浏览器与邮件取证-BrowserEmailForensics.md](skills/浏览器与邮件取证-BrowserEmailForensics.md)

## 工具和资源

### 推荐工具套件

**磁盘镜像工具:**
- [FTK Imager](https://accessdata.com/product-download) - 磁盘镜像工具
- [Guymager](https://guymager.sourceforge.io/) - 磁盘镜像工具
- [dd](https://linux.die.net/man/1/dd) - 磁盘复制工具
- [EnCase](https://www.opentext.com/) - 商业取证工具

**内存取证工具:**
- [Volatility](https://www.volatilityfoundation.org/) - 内存取证工具
- [Rekall](https://github.com/google/rekall) - 内存取证框架
- [LiME](https://github.com/504ensicsLabs/LiME) - Linux内存提取

**取证分析工具:**
- [Sleuth Kit](https://www.sleuthkit.org/) - 文件系统分析
- [Autopsy](https://www.autopsy.com/) - 取证分析平台
- [KAPE](https://www.kroll.com/en/services/cyber-risk/incident-response-litigation-support) - Windows取证工件收集
- [Hayabusa](https://github.com/Yamato-Security/hayabusa) - 事件日志分析

**浏览器和邮件取证:**
- [Browser History Viewer](https://www.magnetforensics.com/) - 浏览器历史查看
- [Mail Viewer](https://www.magnetforensics.com/) - 邮件查看
- [FTK Imager](https://accessdata.com/product-download) - 取证分析

### 学习资源

- [SANS Digital Forensics](https://www.sans.org/course/for508) - 数字取证培训
- [Volatility Foundation](https://www.volatilityfoundation.org/) - Volatility文档
- [Sleuth Kit](https://www.sleuthkit.org/) - Sleuth Kit文档
- [Autopsy](https://www.autopsy.com/) - Autopsy文档

## 最佳实践

### 数字取证原则

1. **合法性** - 确保取证合法
2. **完整性** - 确保证据完整
3. **可追溯性** - 记录取证过程
4. **专业性** - 使用专业工具
5. **合规性** - 遵守法律法规

### 常见误区

- 忽略现场保护
- 没有证据链记录
- 使用不专业工具
- 缺乏取证流程
- 违反法律法规

## 参考文档

### 标准和框架
- [NIST SP 800-86](https://csrc.nist.gov/publications/detail/sp/800-86/final) - 数字取证
- [ISO 27037](https://www.iso.org/standard/44400.html) - 数字证据收集
- [SWGDE Guidelines](https://www.swgde.org/) - 数字取证指南

### 技术指南
- [Volatility Documentation](https://www.volatilityfoundation.org/documentation) - Volatility文档
- [Sleuth Kit Documentation](https://www.sleuthkit.org/sleuthkit/docs/) - Sleuth Kit文档
- [Autopsy Documentation](https://www.autopsy.com/docs/) - Autopsy文档

## 维护与更新

本技能基于最新的数字取证实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 5
**工具数量**: 20+

---

*本技能旨在帮助数字取证专家、事件响应人员和法律调查人员系统性地进行数字取证。*