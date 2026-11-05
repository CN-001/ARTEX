---
name: api-security
description: Comprehensive API security framework covering OWASP API security testing, API authentication and authorization security, and GraphQL and microservice API security. Use when conducting API security assessments, securing REST/GraphQL APIs, or implementing API security controls.
---

# API安全技能

## 概述

本技能提供系统化的API安全框架，涵盖OWASP API安全测试、API认证与授权安全和GraphQL与微服务API安全等核心能力。

**核心能力:**
- OWASP API测试 - OWASP API安全测试
- API认证授权 - API认证与授权安全
- GraphQL安全 - GraphQL与微服务API安全

**适用场景:**
- API安全评估
- API安全测试
- API安全加固
- 微服务安全
- GraphQL安全

## 快速开始

### API安全评估流程

1. **API发现** - 30-60 分钟
   - 发现API端点
   - 分析API文档
   - 识别API类型

2. **安全测试** - 1-2 小时
   - 测试认证机制
   - 测试授权机制
   - 测试输入验证

3. **漏洞分析** - 1-2 小时
   - 分析API漏洞
   - 评估漏洞风险
   - 确定修复优先级

4. **安全加固** - 2-4 小时
   - 实施安全控制
   - 加强API安全
   - 验证加固效果

## 核心技能模块

### 1. OWASP API安全测试 (OWASP API Security Testing)

OWASP API安全测试。

**核心内容:**
- OWASP API Top 10
- API漏洞识别
- API安全测试方法
- API安全最佳实践

**详细文档:** [OWASP API安全测试-OWASPAPISecurityTesting.md](skills/OWASP API安全测试-OWASPAPISecurityTesting.md)

### 2. API认证与授权安全 (API Authentication and Authorization Security)

API认证与授权安全。

**核心内容:**
- API认证机制
- API授权机制
- JWT安全
- OAuth 2.0安全

**详细文档:** [API认证与授权安全-APIAuthAuthorizationSecurity.md](skills/API认证与授权安全-APIAuthAuthorizationSecurity.md)

### 3. GraphQL与微服务API安全 (GraphQL and Microservice API Security)

GraphQL与微服务API安全。

**核心内容:**
- GraphQL安全
- 微服务API安全
- 服务网格安全
- API网关安全

**详细文档:** [GraphQL与微服务API安全-GraphQLMicroserviceAPISecurity.md](skills/GraphQL与微服务API安全-GraphQLMicroserviceAPISecurity.md)

## 工具和资源

### 推荐工具套件

**API测试工具:**
- [Burp Suite](https://portswigger.net/burp) - API测试
- [Postman](https://www.postman.com/) - API测试
- [OWASP ZAP](https://www.zaproxy.org/) - 自动化扫描
- [SoapUI](https://www.soapui.org/) - API测试

**JWT工具:**
- [jwt.io](https://jwt.io/) - JWT调试
- [jose](https://github.com/panva/jose) - JWT库
- [JWT Debugger](https://jwt.io/#debugger) - JWT调试器

**服务网格工具:**
- [Istio](https://istio.io/) - 服务网格
- [Envoy](https://www.envoyproxy.io/) - 代理
- [Linkerd](https://linkerd.io/) - 服务网格

**身份代理:**
- [Authentik](https://goauthentik.io/) - 身份代理
- [Keycloak](https://www.keycloak.org/) - 身份管理
- [Okta](https://www.okta.com/) - 身份管理

### 学习资源

- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/) - OWASP API安全
- [OWASP API Security Testing](https://owasp.org/www-project-api-security/) - OWASP API测试
- [GraphQL Security](https://graphql.org/learn/best-practices/) - GraphQL安全
- [Microservices Security](https://microservices.io/) - 微服务安全

## 最佳实践

### API安全原则

1. **认证优先** - 确保API认证
2. **授权控制** - 实施授权控制
3. **输入验证** - 验证所有输入
4. **速率限制** - 实施速率限制
5. **监控日志** - 监控和日志记录

### 常见误区

- 忽略认证授权
- 没有输入验证
- 缺乏速率限制
- 没有监控日志
- 忽略API文档

## 参考文档

### 标准和框架
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/) - OWASP API安全
- [OAuth 2.0](https://oauth.net/2/) - OAuth 2.0
- [OpenAPI Specification](https://swagger.io/specification/) - OpenAPI规范

### 技术指南
- [OWASP API Security Testing](https://owasp.org/www-project-api-security/) - OWASP API测试
- [GraphQL Security](https://graphql.org/learn/best-practices/) - GraphQL安全
- [Microservices Security](https://microservices.io/) - 微服务安全

## 维护与更新

本技能基于最新的API安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 3
**工具数量**: 15+

---

*本技能旨在帮助API开发人员、安全测试人员和安全架构师系统性地进行API安全评估。*