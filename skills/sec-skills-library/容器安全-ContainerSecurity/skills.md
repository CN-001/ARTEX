---
name: container-security
description: Comprehensive container security framework covering container image security and vulnerability scanning, Kubernetes RBAC and security policies, container runtime security (Falco), and container escape detection and defense. Use when securing containerized applications, hardening Kubernetes clusters, or implementing container security controls.
---

# 容器安全技能

## 概述

本技能提供系统化的容器安全框架，涵盖容器镜像安全与漏洞扫描、Kubernetes RBAC与安全策略、容器运行时安全(Falco)和容器逃逸检测与防御等核心能力。

**核心能力:**
- 镜像安全 - 容器镜像安全与漏洞扫描
- K8s安全 - Kubernetes RBAC与安全策略
- 运行时安全 - 容器运行时安全(Falco)
- 逃逸防御 - 容器逃逸检测与防御

**适用场景:**
- 容器安全评估
- Kubernetes安全加固
- 容器运行时安全
- 容器逃逸防护
- 容器安全运营

## 快速开始

### 容器安全评估流程

1. **镜像扫描** - 30-60 分钟
   - 扫描容器镜像
   - 识别安全漏洞
   - 评估漏洞风险

2. **集群加固** - 2-4 小时
   - 加固Kubernetes集群
   - 配置RBAC策略
   - 实施网络策略

3. **运行时监控** - 持续进行
   - 监控容器运行
   - 检测异常行为
   - 响应安全事件

4. **逃逸防护** - 持续进行
   - 检测逃逸行为
   - 防御容器逃逸
   - 改进安全策略

## 核心技能模块

### 1. 容器镜像安全与漏洞扫描 (Container Image Security and Vulnerability Scanning)

容器镜像安全与漏洞扫描。

**核心工具:**
- Trivy - 容器安全扫描
- Grype - 容器漏洞扫描
- Clair - 容器漏洞扫描
- Cosign - 镜像签名

**详细文档:** [容器镜像安全与漏洞扫描-ContainerImageSecurityScanning.md](skills/容器镜像安全与漏洞扫描-ContainerImageSecurityScanning.md)

### 2. Kubernetes RBAC与安全策略 (Kubernetes RBAC and Security Policies)

Kubernetes RBAC与安全策略。

**核心工具:**
- kubectl - Kubernetes命令行
- kube-bench - 集群安全审计
- kube-hunter - 集群安全测试
- OPA - 策略引擎

**详细文档:** [Kubernetes RBAC与安全策略-KubernetesRBACSecurityPolicy.md](skills/Kubernetes RBAC与安全策略-KubernetesRBACSecurityPolicy.md)

### 3. 容器运行时安全(Falco) (Container Runtime Security - Falco)

容器运行时安全(Falco)。

**核心工具:**
- Falco - 运行时安全
- tracee - 运行时安全
- eBPF - 内核监控
- 自定义工具 - 自定义监控工具

**详细文档:** [容器运行时安全-Falco-ContainerRuntimeSecurityFalco.md](skills/容器运行时安全-Falco-ContainerRuntimeSecurityFalco.md)

### 4. 容器逃逸检测与防御 (Container Escape Detection and Defense)

容器逃逸检测与防御。

**核心工具:**
- Falco - 运行时安全
- tracee - 运行时安全
- 自定义工具 - 自定义检测工具
- 安全策略 - 安全策略配置

**详细文档:** [容器逃逸检测与防御-ContainerEscapeDetectionDefense.md](skills/容器逃逸检测与防御-ContainerEscapeDetectionDefense.md)

## 工具和资源

### 推荐工具套件

**镜像扫描工具:**
- [Trivy](https://aquasecurity.github.io/trivy/) - 容器安全扫描
- [Grype](https://github.com/anchore/grype) - 容器漏洞扫描
- [Clair](https://github.com/quay/clair) - 容器漏洞扫描
- [Cosign](https://github.com/sigstore/cosign) - 镜像签名

**运行时安全工具:**
- [Falco](https://falco.org/) - 运行时安全
- [tracee](https://github.com/aquasecurity/tracee) - 运行时安全
- [eBPF](https://ebpf.io/) - 内核监控
- [Cilium](https://cilium.io/) - 网络安全

**Kubernetes安全工具:**
- [kube-bench](https://github.com/aquasecurity/kube-bench) - 集群安全审计
- [kube-hunter](https://github.com/aquasecurity/kube-hunter) - 集群安全测试
- [OPA Gatekeeper](https://github.com/open-policy-agent/gatekeeper) - 策略引擎
- [Kyverno](https://kyverno.io/) - 策略引擎

**策略引擎:**
- [OPA](https://www.openpolicyagent.org/) - 策略引擎
- [Kyverno](https://kyverno.io/) - 策略引擎
- [Gatekeeper](https://github.com/open-policy-agent/gatekeeper) - 策略引擎

### 学习资源

- [Kubernetes Security](https://kubernetes.io/docs/concepts/security/) - Kubernetes安全
- [Falco Documentation](https://falco.org/docs/) - Falco文档
- [Container Security Best Practices](https://www.cisecurity.org/benchmark/kubernetes) - 容器安全最佳实践
- [NIST Container Security](https://www.nist.gov/) - NIST容器安全

## 最佳实践

### 容器安全原则

1. **镜像安全** - 确保镜像安全
2. **最小权限** - 遵循最小权限
3. **运行时监控** - 监控运行时行为
4. **网络隔离** - 实施网络隔离
5. **持续扫描** - 持续扫描漏洞

### 常见误区

- 忽略镜像安全
- 没有最小权限
- 缺乏运行时监控
- 没有网络隔离
- 缺乏持续扫描

## 参考文档

### 标准和框架
- [CIS Kubernetes Benchmark](https://www.cisecurity.org/benchmark/kubernetes) - CIS Kubernetes基准
- [NIST SP 800-190](https://csrc.nist.gov/publications/detail/sp/800-190/final) - 容器安全
- [ISO 27001](https://www.iso.org/isoiec-27001-information-security.html) - 信息安全

### 技术指南
- [Kubernetes Security](https://kubernetes.io/docs/concepts/security/) - Kubernetes安全
- [Falco Documentation](https://falco.org/docs/) - Falco文档
- [Container Security Best Practices](https://www.cisecurity.org/benchmark/kubernetes) - 容器安全最佳实践

## 维护与更新

本技能基于最新的容器安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 4
**工具数量**: 20+

---

*本技能旨在帮助容器开发人员、Kubernetes管理员和安全工程师系统性地进行容器安全评估。*