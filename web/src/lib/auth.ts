export interface CurrentUser {
  id: string;
  name: string;
  username: string;
  email: string;
  avatar: string;
  role: string;
}

const TOKEN_KEY = "artex_token";

// 纯 Cookie 认证：会话 JWT 只存在于服务端下发的 HttpOnly cookie 中，JS 读不到。
// 本模块只保留「会话探针」与「本地清理」两件事；不再读写 localStorage 里的 token
// （clearToken 仍会顺手清掉历史遗留的 artex_token，避免旧版本残留被误用）。
export const auth = {
  // checkSession 通过 /api/auth/me 判断当前 cookie 会话是否有效。
  // 401 视为未登录；网络错误也按未登录处理（宁可让用户重登，也不放行受保护 UI）。
  async checkSession(): Promise<CurrentUser | null> {
    if (typeof window === "undefined") return null;
    if (process.env.NEXT_PUBLIC_MOCK === "1") {
      return { id: "1", name: "ARTEX", username: "artex", email: "", avatar: "", role: "operator" };
    }
    try {
      const r = await fetch("/api/auth/me", { credentials: "same-origin" });
      if (!r.ok) return null;
      const u = (await r.json()) as { username?: string; name?: string; role?: string };
      const name = u.name || u.username || "ARTEX";
      return { id: "1", name, username: u.username || "artex", email: "", avatar: "", role: u.role || "admin" };
    } catch {
      return null;
    }
  },

  // 清掉历史遗留的 localStorage token（旧版本用于 Authorization 头；纯 Cookie 后不再使用）。
  clearToken(): void {
    if (typeof window === "undefined") return;
    localStorage.removeItem(TOKEN_KEY);
  },
};
