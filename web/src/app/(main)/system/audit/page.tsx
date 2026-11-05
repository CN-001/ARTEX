"use client";

import * as React from "react";
import { Download, RefreshCw, ShieldCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";
import type { SecurityAuditEntry } from "@/lib/types";
import { cn } from "@/lib/utils";

// 安全审计页：看「谁、从哪、对平台做了什么、结果如何」。
// 数据来自 /api/security/audit（写操作自动记录 + 登录/改密/门禁拦截等安全事件）。
function fmtTime(ts: string) {
  const d = new Date(ts);
  if (Number.isNaN(d.getTime())) return ts;
  const p = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}

// 动作 → 展示色：失败/拦截用红色，认证事件用琥珀色，普通写操作用默认色。
function actionTone(action: string, status: number) {
  if (status >= 400) return "text-red-600 dark:text-red-400";
  if (action.startsWith("auth.")) return "text-amber-600 dark:text-amber-400";
  if (action.startsWith("scope.")) return "text-red-600 dark:text-red-400";
  return "text-foreground/80";
}

function toCSV(rows: SecurityAuditEntry[]) {
  const head = ["id", "ts", "actor", "ip", "method", "path", "action", "target", "status", "detail"];
  const esc = (v: unknown) => {
    const s = String(v ?? "");
    return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
  };
  return [head.join(","), ...rows.map((r) => head.map((k) => esc((r as unknown as Record<string, unknown>)[k])).join(","))].join(
    "\n",
  );
}

export default function SecurityAuditPage() {
  const [rows, setRows] = React.useState<SecurityAuditEntry[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");
  const [filter, setFilter] = React.useState("");
  const [limit, setLimit] = React.useState(200);
  const [hasMore, setHasMore] = React.useState(false);

  const load = React.useCallback(
    async (beforeId?: number) => {
      setLoading(true);
      setError("");
      try {
        const res = await api.securityAudit({ limit, beforeId });
        const entries = res.entries ?? [];
        if (beforeId) {
          setRows((old) => [...old, ...entries]);
        } else {
          setRows(entries);
        }
        setHasMore(entries.length >= limit);
      } catch (e) {
        setError(e instanceof Error ? e.message : "读取审计日志失败");
      } finally {
        setLoading(false);
      }
    },
    [limit],
  );

  React.useEffect(() => {
    void load();
  }, [load]);

  const filtered = React.useMemo(() => {
    const q = filter.trim().toLowerCase();
    if (!q) return rows;
    return rows.filter((r) =>
      [r.action, r.actor, r.ip, r.path, r.target, r.detail, String(r.status)].some((v) =>
        String(v ?? "").toLowerCase().includes(q),
      ),
    );
  }, [rows, filter]);

  function exportCSV() {
    const blob = new Blob([toCSV(filtered)], { type: "text/csv;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `artex-security-audit-${new Date().toISOString().slice(0, 19).replace(/[:T]/g, "-")}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="flex h-full flex-col gap-4 p-4 md:p-6">
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-2">
          <ShieldCheck className="size-5 text-muted-foreground" />
          <h1 className="font-semibold text-lg">安全审计</h1>
          <span className="text-muted-foreground text-xs">写操作 / 登录 / 门禁拦截留痕</span>
        </div>
        <div className="ml-auto flex flex-wrap items-center gap-2">
          <Input
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="过滤：动作 / 用户 / IP / 路径…"
            className="h-8 w-64"
          />
          <select
            className="h-8 rounded-md border bg-background px-2 text-sm"
            value={limit}
            onChange={(e) => setLimit(Number(e.target.value))}
          >
            <option value={100}>100 条</option>
            <option value={200}>200 条</option>
            <option value={500}>500 条</option>
          </select>
          <Button variant="outline" size="sm" onClick={() => void load()}>
            <RefreshCw className={cn("size-4", loading && "animate-spin")} />
            刷新
          </Button>
          <Button variant="outline" size="sm" onClick={exportCSV} disabled={filtered.length === 0}>
            <Download className="size-4" />
            导出 CSV
          </Button>
        </div>
      </div>

      {error ? <div className="rounded-md border border-red-500/40 bg-red-500/5 p-3 text-red-600 text-sm">{error}</div> : null}

      <div className="min-h-0 flex-1 overflow-auto rounded-md border">
        <table className="w-full border-collapse text-sm">
          <thead className="sticky top-0 bg-muted/80 backdrop-blur">
            <tr className="text-muted-foreground">
              <th className="px-3 py-2 text-left font-medium">时间</th>
              <th className="px-3 py-2 text-left font-medium">动作</th>
              <th className="px-3 py-2 text-left font-medium">用户</th>
              <th className="px-3 py-2 text-left font-medium">来源 IP</th>
              <th className="px-3 py-2 text-left font-medium">路径</th>
              <th className="px-3 py-2 text-left font-medium">目标</th>
              <th className="px-3 py-2 text-right font-medium">状态</th>
              <th className="px-3 py-2 text-left font-medium">明细</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((r) => (
              <tr key={r.id} className="border-t hover:bg-muted/40">
                <td className="whitespace-nowrap px-3 py-1.5 tabular-nums">{fmtTime(r.ts)}</td>
                <td className={cn("whitespace-nowrap px-3 py-1.5 font-mono text-xs", actionTone(r.action, r.status))}>
                  {r.action}
                </td>
                <td className="px-3 py-1.5">{r.actor || "-"}</td>
                <td className="px-3 py-1.5 font-mono text-xs">{r.ip || "-"}</td>
                <td className="max-w-[280px] truncate px-3 py-1.5 font-mono text-xs" title={r.path}>
                  {r.method} {r.path}
                </td>
                <td className="max-w-[240px] truncate px-3 py-1.5" title={r.target}>
                  {r.target || "-"}
                </td>
                <td
                  className={cn(
                    "px-3 py-1.5 text-right tabular-nums",
                    r.status >= 400 ? "text-red-600 dark:text-red-400" : "text-emerald-600 dark:text-emerald-400",
                  )}
                >
                  {r.status}
                </td>
                <td className="max-w-[320px] truncate px-3 py-1.5 text-muted-foreground text-xs" title={r.detail}>
                  {r.detail || "-"}
                </td>
              </tr>
            ))}
            {filtered.length === 0 && !loading ? (
              <tr>
                <td colSpan={8} className="px-3 py-10 text-center text-muted-foreground text-sm">
                  暂无审计记录
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>

      <div className="flex items-center justify-between text-muted-foreground text-xs">
        <span>共 {filtered.length} 条（保留 180 天；审计写入为异步，失败不影响业务）</span>
        {hasMore ? (
          <Button
            variant="ghost"
            size="sm"
            disabled={loading || rows.length === 0}
            onClick={() => void load(rows[rows.length - 1]?.id)}
          >
            加载更早记录
          </Button>
        ) : null}
      </div>
    </div>
  );
}
