"use client";

import { useEffect, useState } from "react";

import { auth, type CurrentUser } from "@/lib/auth";

const FALLBACK: CurrentUser = { id: "1", name: "ARTEX", username: "artex", email: "", avatar: "", role: "admin" };

// 纯 Cookie 认证后前端无法解析 JWT，当前用户由 /api/auth/me 返回。
export function useCurrentUser(): CurrentUser {
  const [user, setUser] = useState<CurrentUser>(FALLBACK);

  useEffect(() => {
    let cancelled = false;
    void auth.checkSession().then((u) => {
      if (!cancelled && u) setUser(u);
    });
    return () => {
      cancelled = true;
    };
  }, []);

  return user;
}
