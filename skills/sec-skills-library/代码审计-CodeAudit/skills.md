---
name: code-audit
description: Comprehensive code audit framework covering PHP, Java, JavaScript, Python, C, C++, Rust, and Go code auditing, as well as AI Agent code auditing. Use when conducting security code reviews, vulnerability assessments, or secure code development.
---

# 代码审计技能

## 概述

本技能提供系统化的代码审计框架，涵盖PHP、Java、JavaScript、Python、C、C++、Rust和Go代码审计，以及AI Agent代码审计等核心能力。

**核心能力:**
- PHP代码审计 - PHP源代码安全审计
- Java代码审计 - Java源代码安全审计
- JavaScript代码审计 - JavaScript源代码安全审计
- Python代码审计 - Python源代码安全审计
- C代码审计 - C源代码安全审计
- C++代码审计 - C++源代码安全审计
- Rust代码审计 - Rust源代码安全审计
- Go代码审计 - Go源代码安全审计
- AI Agent代码审计 - AI代码安全审计

**适用场景:**
- 安全代码审计
- 漏洞发现和修复
- 安全开发实践
- 合规性评估
- 代码质量提升

## 快速开始

### 代码审计流程

1. **代码收集** - 10-20 分钟
   - 收集源代码
   - 分析代码结构
   - 识别关键模块

2. **静态分析** - 20-40 分钟
   - 使用静态分析工具
   - 检查常见漏洞
   - 分析代码质量

3. **手动审计** - 30-60 分钟
   - 手动审查代码
   - 识别业务逻辑漏洞
   - 评估安全风险

## 核心技能模块

### 1. PHP代码审计 (PHP Code Audit)

对PHP源代码进行安全审计。

**核心工具:**
- RIPS - PHP静态分析
- phpcs - PHP代码规范
- 手动审查 - 人工代码审计

**详细文档:** [PHP代码审计-PHPCodeAudit.md](skills/PHP代码审计-PHPCodeAudit.md)

### 2. Java代码审计 (Java Code Audit)

对Java源代码进行安全审计。

**核心工具:**
- FindSecBugs - Java安全插件
- SpotBugs - Java静态分析
- Fortify - 商业代码审计

**详细文档:** [Java代码审计-JavaCodeAudit.md](skills/Java代码审计-JavaCodeAudit.md)

### 3. JavaScript代码审计 (JavaScript Code Audit)

对JavaScript源代码进行安全审计。

**核心工具:**
- ESLint - JavaScript代码检查
- RetireJS - 依赖漏洞检测
- Semgrep - 语义代码分析

**详细文档:** [JavaScript代码审计-JSCodeAudit.md](skills/JavaScript代码审计-JSCodeAudit.md)

### 4. Python代码审计 (Python Code Audit)

对Python源代码进行安全审计。

**核心工具:**
- Bandit - Python安全检查
- Semgrep - 语义代码分析
- Pyre - Python类型检查

**详细文档:** [Python代码审计-PythonCodeAudit.md](skills/Python代码审计-PythonCodeAudit.md)

### 5. C代码审计 (C Code Audit)

对C源代码进行安全审计。

**核心工具:**
- Flawfinder - C安全检查
- Cppcheck - C/C++静态分析
- ASan - AddressSanitizer
- Valgrind - 内存调试

**详细文档:** [C代码审计-CCodeAudit.md](skills/C代码审计-CCodeAudit.md)

### 6. C++代码审计 (C++ Code Audit)

对C++源代码进行安全审计。

**核心工具:**
- Cppcheck - C/C++静态分析
- Clang-Tidy - Clang静态分析
- PVS-Studio - 商业代码审计

**详细文档:** [C++代码审计-CPPCodeAudit.md](skills/C++代码审计-CPPCodeAudit.md)

### 7. Rust代码审计 (Rust Code Audit)

对Rust源代码进行安全审计。

**核心工具:**
- Clippy - Rust代码检查
- cargo-audit - 依赖审计
- cargo-geiger - 不安全代码检测
- Miri - Rust解释器

**详细文档:** [Rust代码审计-RustCodeAudit.md](skills/Rust代码审计-RustCodeAudit.md)

### 8. Go代码审计 (Go Code Audit)

对Go源代码进行安全审计。

**核心工具:**
- gosec - Go安全检查
- staticcheck - Go静态分析
- govulncheck - Go漏洞检查

**详细文档:** [Go代码审计-GoCodeAudit.md](skills/Go代码审计-GoCodeAudit.md)

### 9. AI Agent代码审计 (AI Agent Code Audit)

对AI应用代码进行安全审计。

**核心工具:**
- Semgrep - 语义代码分析
- Bandit - Python安全检查
- CodeQL - 代码查询语言

**详细文档:** [AI Agent代码审计-AIAgentCodeAudit.md](skills/AI%20Agent代码审计-AIAgentCodeAudit.md)

## 工具和资源

### 推荐工具套件

**PHP工具:**
- [RIPS](https://www.ripstech.com/) - PHP静态分析
- [PHP_CodeSniffer](https://github.com/squizlabs/PHP_CodeSniffer) - PHP代码规范
- [SonarQube](https://www.sonarqube.org/) - 代码质量平台

**Java工具:**
- [FindSecBugs](https://find-sec-bugs.github.io/) - Java安全插件
- [SpotBugs](https://spotbugs.github.io/) - Java静态分析
- [Fortify](https://www.microfocus.com/en-us/products/fortify-static-code-analyzer) - 商业代码审计

**JavaScript工具:**
- [ESLint](https://eslint.org/) - JavaScript代码检查
- [Retire.js](https://retirejs.github.io/retire.js/) - 依赖漏洞检测
- [Semgrep](https://semgrep.dev/) - 语义代码分析

**Python工具:**
- [Bandit](https://github.com/PyCQA/bandit) - Python安全检查
- [Semgrep](https://semgrep.dev/) - 语义代码分析
- [Pyre](https://pyre-check.org/) - Python类型检查

**C/C++工具:**
- [Flawfinder](https://github.com/david-a-wheeler/flawfinder) - C安全检查
- [Cppcheck](https://cppcheck.sourceforge.io/) - C/C++静态分析
- [Clang-Tidy](https://clang.llvm.org/extra/clang-tidy/) - Clang静态分析

**Rust工具:**
- [Clippy](https://github.com/rust-lang/rust-clippy) - Rust代码检查
- [cargo-audit](https://github.com/RustSec/cargo-audit) - 依赖审计
- [cargo-geiger](https://github.com/rust-secure-code/cargo-geiger) - 不安全代码检测

**Go工具:**
- [gosec](https://github.com/securego/gosec) - Go安全检查
- [staticcheck](https://staticcheck.io/) - Go静态分析
- [govulncheck](https://golang.org/x/vuln/cmd/govulncheck) - Go漏洞检查

### 学习资源

- [OWASP Code Review Guide](https://owasp.org/www-project-code-review-guide/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
- [Semgrep Rules](https://semgrep.dev/explore)
- [CodeQL Queries](https://codeql.github.com/codeql-query-help/)

## 最佳实践

### 代码审计原则

1. **全面性** - 覆盖所有代码路径
2. **准确性** - 确保审计结果准确
3. **可追溯性** - 记录审计过程
4. **及时性** - 及时发现和修复漏洞
5. **合规性** - 遵守安全编码规范

### 常见误区

- 过度依赖自动化工具
- 忽略业务逻辑审计
- 没有测试漏洞利用
- 缺乏代码质量评估
- 忽略依赖安全

## 参考文档

### 标准和框架
- [OWASP Code Review Guide](https://owasp.org/www-project-code-review-guide/)
- [CWE/SANS Top 25](https://cwe.mitre.org/top25/)
- [NIST SP 800-53](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final)

### 技术指南
- [Semgrep Documentation](https://semgrep.dev/docs/)
- [CodeQL Documentation](https://codeql.github.com/docs/)
- [SonarQube Documentation](https://docs.sonarqube.org/)

## 维护与更新

本技能基于最新的代码审计实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 9
**工具数量**: 30+

---

*本技能旨在帮助安全研究人员、代码审计人员和开发人员系统性地进行代码安全审计。*