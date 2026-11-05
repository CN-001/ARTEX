---
name: iot-security
description: Comprehensive IoT security framework covering firmware reverse engineering, BLE/Zigbee/Z-Wave wireless security testing, IoT communication protocol security, embedded device hardware security testing, IoT platform and cloud security, and smart home and connected vehicle security. Use when conducting IoT security assessments, embedded device penetration testing, or IoT security operations.
---

# 物联网安全技能

## 概述

本技能提供系统化的物联网安全框架，涵盖固件逆向与分析、BLE/Zigbee/Z-Wave无线安全测试、物联网通信协议安全、嵌入式设备硬件安全测试、物联网平台与云安全和智能家居与车联网安全等核心能力。

**核心能力:**
- 固件逆向 - 固件逆向与分析
- 无线安全 - BLE/Zigbee/Z-Wave安全测试
- 通信安全 - 物联网通信协议安全
- 硬件安全 - 嵌入式设备硬件安全测试
- 平台安全 - 物联网平台与云安全
- 智能安全 - 智能家居与车联网安全

**适用场景:**
- IoT安全评估
- 嵌入式设备渗透测试
- IoT安全运营
- IoT安全合规
- IoT安全培训

## 快速开始

### IoT安全评估流程

1. **设备发现** - 1-2 天
   - 发现IoT设备
   - 分析设备架构
   - 识别攻击面

2. **固件分析** - 2-3 天
   - 提取固件
   - 分析固件结构
   - 逆向固件代码

3. **协议分析** - 1-2 天
   - 分析通信协议
   - 测试协议安全
   - 发现协议漏洞

4. **硬件测试** - 2-3 天
   - 测试硬件接口
   - 分析硬件安全
   - 发现硬件漏洞

## 核心技能模块

### 1. 固件逆向与分析 (Firmware Reverse Engineering)

固件逆向与分析。

**核心工具:**
- Binwalk - 固件分析工具
- Firmadyne - 固件模拟框架
- Ghidra - 逆向工程工具
- QEMU - 模拟器

**详细文档:** [固件逆向与分析-FirmwareReverseEngineering.md](skills/固件逆向与分析-FirmwareReverseEngineering.md)

### 2. BLE/Zigbee/Z-Wave无线安全测试 (BLE/Zigbee/Z-Wave Wireless Security Testing)

BLE/Zigbee/Z-Wave无线安全测试。

**核心工具:**
- BLEAH - BLE安全测试
- GATTool - GATT工具
- Z3sec - Zigbee安全
- HackRF - 软件无线电

**详细文档:** [BLE-Zigbee-Z-Wave无线安全测试-WirelessProtocolSecurity.md](skills/BLE-Zigbee-Z-Wave无线安全测试-WirelessProtocolSecurity.md)

### 3. 物联网通信协议安全 (IoT Communication Protocol Security)

物联网通信协议安全。

**核心工具:**
- MQTT Pwn - MQTT安全测试
- MQTTSA - MQTT安全分析
- Wireshark - 网络协议分析

**详细文档:** [物联网通信协议安全-IoTCommunicationSecurity.md](skills/物联网通信协议安全-IoTCommunicationSecurity.md)

### 4. 嵌入式设备硬件安全测试 (Embedded Device Hardware Security Testing)

嵌入式设备硬件安全测试。

**核心工具:**
- JTAGulator - JTAG调试
- OpenOCD - 开源片上调试器
- Saleae Logic - 逻辑分析仪

**详细文档:** [嵌入式设备硬件安全测试-EmbeddedHardwareSecurity.md](skills/嵌入式设备硬件安全测试-EmbeddedHardwareSecurity.md)

### 5. 物联网平台与云安全 (IoT Platform and Cloud Security)

物联网平台与云安全。

**核心工具:**
- AWS IoT Auditor - AWS IoT审计
- Azurite - Azure IoT测试
- 云安全工具 - 云安全平台

**详细文档:** [物联网平台与云安全-IoTPlatformCloudSecurity.md](skills/物联网平台与云安全-IoTPlatformCloudSecurity.md)

### 6. 智能家居与车联网安全 (Smart Home and Connected Vehicle Security)

智能家居与车联网安全。

**核心工具:**
- Home Assistant Audit - Home Assistant审计
- CANtact - CAN总线工具
- 自定义工具 - 自定义测试工具

**详细文档:** [智能家居与车联网安全-SmartHomeConnectedVehicleSecurity.md](skills/智能家居与车联网安全-SmartHomeConnectedVehicleSecurity.md)

## 工具和资源

### 推荐工具套件

**固件分析工具:**
- [Binwalk](https://github.com/ReFirmLabs/binwalk) - 固件分析
- [Firmadyne](https://github.com/pr0v3rbs/FirmAE) - 固件模拟
- [Ghidra](https://ghidra-sre.org/) - 逆向工程
- [QEMU](https://www.qemu.org/) - 模拟器

**无线安全工具:**
- [BLEAH](https://github.com/virtualabs/btlejack) - BLE安全测试
- [GATTool](https://github.com/IoT-Research-Lab/GATTool) - GATT工具
- [Z3sec](https://github.com/zigbeefw/Z3sec) - Zigbee安全
- [HackRF](https://greatscottgadgets.com/hackrf/) - 软件无线电

**协议分析工具:**
- [MQTT Pwn](https://github.com/akira1g/mqtt-pwn) - MQTT安全测试
- [Wireshark](https://www.wireshark.org/) - 网络协议分析
- [Scapy](https://scapy.net/) - 网络包操作

**硬件测试工具:**
- [JTAGulator](https://www.grandideastudio.com/portfolio/jtagulator/) - JTAG调试
- [OpenOCD](http://openocd.org/) - 开源片上调试器
- [Saleae Logic](https://www.saleae.com/) - 逻辑分析仪

### 学习资源

- [OWASP IoT Testing Guide](https://owasp.org/www-project-iot-security-testing-guide/)
- [NIST SP 800-183](https://csrc.nist.gov/publications/detail/sp/800-183/final)
- [Matter and Thread Security](https://csa-iot.org/)
- [IEC 62443-4-2](https://www.iec.ch/) - IoT组件安全
- [Pwn2Own IoT](https://www.zerodayinitiative.com/) - IoT安全竞赛

## 最佳实践

### IoT安全原则

1. **安全设计** - 从设计阶段考虑安全
2. **最小权限** - 遵循最小权限原则
3. **加密通信** - 使用加密通信
4. **安全更新** - 提供安全更新机制
5. **隐私保护** - 保护用户隐私

### 常见误区

- 忽略固件安全
- 使用弱密码
- 没有加密通信
- 缺乏安全更新
- 忽略隐私保护

## 参考文档

### 标准和框架
- [OWASP IoT Testing Guide](https://owasp.org/www-project-iot-security-testing-guide/)
- [NIST SP 800-183](https://csrc.nist.gov/publications/detail/sp/800-183/final)
- [IEC 62443-4-2](https://www.iec.ch/) - IoT组件安全
- [ISO 27001](https://www.iso.org/isoiec-27001-information-security.html)

### 技术指南
- [Matter and Thread Security](https://csa-iot.org/)
- [OWASP IoT Security](https://owasp.org/www-project-iot-security/)
- [NIST IoT Security](https://www.nist.gov/iot)

## 维护与更新

本技能基于最新的IoT安全实践和安全研究持续更新。

**版本**: v1.0
**最后更新**: 2026-06-17
**技能数量**: 6
**工具数量**: 20+

---

*本技能旨在帮助安全研究人员、IoT安全工程师和嵌入式安全专家系统性地进行IoT安全评估。*