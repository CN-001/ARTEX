---
name: reporting
description: Comprehensive reporting framework covering penetration testing report writing, vulnerability rating and CVSS, and security report templates in Markdown, HTML, and Word formats. Use when conducting penetration testing, security assessments, or report generation.
---

# 报告撰写技能

## 概述

本技能提供系统化的报告撰写框架，涵盖渗透测试报告编写、漏洞评级与CVSS评分和安全报告模板（Markdown、HTML、Word）等核心能力。

**核心能力:**
- 报告编写 - 渗透测试报告编写方法论
- 漏洞评级 - CVSS 3.1/4.0评分
- Markdown报告 - Markdown格式报告
- HTML报告 - 交互式HTML报告
- Word报告 - Word/PDF格式报告
- 统一数据模型 - JSON Schema数据模型
- AI Agent接口 - 原生AI Agent支持

**适用场景:**
- 渗透测试报告撰写
- 安全评估报告生成
- 漏洞评级和分类
- 合规报告生成
- 批量报告生成

## 快速开始

### 报告撰写流程

1. **数据收集** - 30-60 分钟
   - 收集测试数据
   - 整理漏洞信息
   - 准备证据材料

2. **报告结构** - 15-30 分钟
   - 确定报告结构
   - 编写执行摘要
   - 组织报告内容

3. **报告生成** - 30-60 分钟
   - 选择报告格式
   - 生成报告内容
   - 审核报告质量

### 漏洞评级流程

1. **漏洞分析** - 10-20 分钟
   - 分析漏洞特征
   - 评估漏洞影响
   - 确定利用难度

2. **CVSS评分** - 15-30 分钟
   - 选择CVSS版本
   - 计算基础分数
   - 评估环境和时间分数

3. **风险评级** - 10-20 分钟
   - 确定风险等级
   - 提供修复建议
   - 设置修复优先级

## 核心技能模块

### 1. 报告编写 (Penetration Testing Report Writing)

渗透测试报告编写方法论。

**核心内容:**
- 报告结构设计
- 执行摘要编写
- 技术细节描述
- 修复建议提供
- 质量标准制定

**详细文档:** [报告编写-PentestReport.md](skills/报告编写-PentestReport.md)

### 2. 漏洞评级与CVSS (Vulnerability Rating and CVSS)

漏洞评级和CVSS评分。

**核心内容:**
- CVSS 3.1/4.0评分
- DREAD评级方法
- 风险等级确定
- 修复优先级设置
- 评级标准制定

**详细文档:** [漏洞评级与CVSS-VulnRatingCVSS.md](skills/漏洞评级与CVSS-VulnRatingCVSS.md)

### 3. 安全报告模板-Markdown (Security Report Template - Markdown)

Markdown格式安全报告模板。

**核心特性:**
- 完整模板结构
- Mermaid图表支持
- 5种模板变体
- Python生成器
- AI Agent接口

**详细文档:** [安全报告模板-Markdown.md](skills/安全报告模板-Markdown.md)

### 4. 安全报告模板-HTML (Security Report Template - HTML)

交互式HTML安全报告模板。

**核心特性:**
- Chart.js动态图表
- 暗色/亮色模式
- 侧边栏导航
- 打印/PDF导出
- Python生成器

**详细文档:** [安全报告模板-HTML.md](skills/安全报告模板-HTML.md)

### 5. 安全报告模板-Word (Security Report Template - Word)

Word/PDF格式安全报告模板。

**核心特性:**
- 完整封面+目录
- 风险等级颜色标注
- 页眉页脚+页码
- python-docx生成器
- AI Agent接口

**详细文档:** [安全报告模板-Word.md](skills/安全报告模板-Word.md)

## 工具和资源

### 推荐工具套件

**Markdown工具:**
- [Typora](https://typora.io/) - Markdown编辑器
- [VS Code](https://code.visualstudio.com/) - 代码编辑器
- [Mermaid](https://mermaid.js.org/) - 图表生成
- [Pandoc](https://pandoc.org/) - 文档转换

**HTML工具:**
- [Chart.js](https://www.chartjs.org/) - 图表库
- [Bootstrap](https://getbootstrap.com/) - CSS框架
- [jQuery](https://jquery.com/) - JavaScript库
- [wkhtmltopdf](https://wkhtmltopdf.org/) - HTML转PDF

**Word工具:**
- [python-docx](https://python-docx.readthedocs.io/) - Word生成库
- [Microsoft Word](https://www.microsoft.com/word) - Word编辑器
- [LibreOffice](https://www.libreoffice.org/) - 开源办公套件
- [ReportLab](https://www.reportlab.com/) - PDF生成库

**CVSS工具:**
- [CVSS Calculator](https://www.first.org/cvss/calculator/3.1) - CVSS计算器
- [CVSS Python](https://github.com/RedHatProductSecurity/cvss) - CVSS Python库
- [DREAD Calculator](https://dread.readthedocs.io/) - DREAD计算器

### 学习资源

- [PTES Reporting Guidelines](http://www.pentest-standard.org/index.php/Reporting)
- [SANS Report Writing](https://www.sans.org/white-papers/36057/)
- [OWASP Testing Guide - Reporting](https://owasp.org/www-project-web-security-testing-guide/latest/4-Reporting)
- [CVSS 4.0 Specification](https://www.first.org/cvss/v4-0/)

## 最佳实践

### 报告撰写原则

1. **准确性** - 确保报告内容准确
2. **完整性** - 覆盖所有测试内容
3. **可读性** - 便于客户理解
4. **专业性** - 体现专业水准
5. **及时性** - 按时提交报告

### 常见误区

- 报告结构不清晰
- 技术细节过多
- 缺乏修复建议
- 漏洞评级不准确
- 报告质量不高

## 参考文档

### 标准和框架
- [PTES - Penetration Testing Execution Standard](https://www.pentest-standard.org/)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)
- [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)
- [PCI DSS v4.0](https://www.pcisecuritystandards.org/document_library/)

### 技术指南
- [CVSS 4.0 Specification](https://www.first.org/cvss/v4-0/)
- [python-docx Documentation](https://python-docx.readthedocs.io/)
- [Chart.js Documentation](https://www.chartjs.org/docs/latest/)
- [Mermaid Documentation](https://mermaid.js.org/intro/)

## 维护与更新

本技能基于最新的报告撰写实践和安全研究持续更新。

**版本**: v2.0
**最后更新**: 2026-06-17
**技能数量**: 5
**模板数量**: 3
**工具数量**: 15+

---

*本技能旨在帮助安全研究人员、渗透测试人员和报告撰写人员系统性地生成高质量安全报告。*