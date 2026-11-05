# 双模型策略（敏感数据合规）

## 背景

ARTEX 的 agent 默认用云端 API（DeepSeek）推理。渗透测试中，**目标系统信息、客户资产、凭据、
漏洞细节**都属于敏感数据——按《数据安全法》与《个人信息保护法》，客户数据不应随意出境/交给
第三方 API。

## 策略：按项目敏感度选择模型

| 项目类型 | 模型选择 | 理由 |
|---|---|---|
| **敏感项目**（客户真实资产、生产系统、含个人信息） | **本地模型** `ollama-artex_qwen3` | 数据不出本机 |
| 一般项目（SRC 公开资产、靶场、CTF） | 云端 `deepseek-v4-flash`（默认） | 效率/质量更好 |
| 需要更强推理的敏感项目 | 本地 `artex_qwen3_long`（32K 上下文） | 兼顾上下文与合规 |

## 怎么切（三处，任选）

### ① 任务级（推荐）—— 只影响这个任务

```powershell
# 切本地模型（敏感项目）
docker exec artex-main-postgres-1 psql -U artex -d artex -c "UPDATE tasks SET llm_profile_id=2 WHERE id=<任务ID>"

# 切回默认云端
docker exec artex-main-postgres-1 psql -U artex -d artex -c "UPDATE tasks SET llm_profile_id=NULL WHERE id=<任务ID>"
```

或直接跑现成脚本：`E:\ARTEX-main\data\_ops\dual_model.ps1 task <任务ID> local`

### ② Agent 级 —— 影响该 agent 的所有任务

```powershell
# 例：把 worker(4) 切本地（只影响执行 agent）
docker exec artex-main-postgres-1 psql -U artex -d artex -c "UPDATE agents SET llm_profile_id=2 WHERE id=4"
```

### ③ 界面设置（ARTEX 后台）

任务/Agent 设置里有模型配置项，选 `ollama-artex_qwen3` 即可。

## 本地模型说明

- 模型：Qwen3-8B Q4（ARTEX 定制版，支持工具调用）
- 上下文：16K（`artex_qwen3`）/ 32K（`artex_qwen3_long`）
- 速度：约 7-10 tokens/s（8GB 显存，部分 CPU 卸载）；简单侦察/执行够用，复杂规划偏慢
- 前置：Ollama 服务运行中（`E:\Ollama\ollama.exe serve`），模型已就绪

## 合规红线

1. 客户真实数据 → **必须**本地模型
2. 本地模型产出的报告 → 人工复核后再交付
3. 云端模型会话中不要粘贴客户凭据/数据库导出
4. 报告归档加密；原始证据不上传任何云端服务
