"use client";

import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";

const pageTitles: Record<string, string> = {
  "/": "Dashboard",
  "/users": "Manage Users",
  "/barang": "Barang",
  "/tipe-barang": "Tipe Barang",
  "/satuan-barang": "Satuan Barang",
};

export default function Navbar() {
  const pathname = usePathname();
  const { user } = useAuth();
  const title = pageTitles[pathname] || "Dashboard";

  return (
    <header className="h-16 bg-white border-b border-zinc-200 flex items-center justify-between px-6">
      <h2 className="text-lg font-semibold text-zinc-800">{title}</h2>
      <div className="flex items-center gap-3">
        <div className="text-right">
          <div className="text-sm font-medium text-zinc-700">{user?.name}</div>
          <div className="text-xs text-zinc-400 capitalize">{user?.role_name}</div>
        </div>
        <div className="w-9 h-9 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600 font-semibold text-sm">
          {user?.name?.charAt(0)?.toUpperCase() || "U"}
        </div>
      </div>
    </header>
  );
}
