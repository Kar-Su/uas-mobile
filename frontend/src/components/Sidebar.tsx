"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { HiOutlineHome, HiOutlineUsers, HiOutlineCube, HiOutlineTag, HiOutlineScale } from "react-icons/hi";

const superAdminNav = [
  { href: "/", label: "Dashboard", icon: HiOutlineHome },
  { href: "/users", label: "Users", icon: HiOutlineUsers },
  { href: "/barang", label: "Barang", icon: HiOutlineCube },
  { href: "/tipe-barang", label: "Tipe Barang", icon: HiOutlineTag },
  { href: "/satuan-barang", label: "Satuan Barang", icon: HiOutlineScale },
];

const adminGudangNav = [
  { href: "/", label: "Dashboard", icon: HiOutlineHome },
  { href: "/barang", label: "Barang", icon: HiOutlineCube },
  { href: "/tipe-barang", label: "Tipe Barang", icon: HiOutlineTag },
  { href: "/satuan-barang", label: "Satuan Barang", icon: HiOutlineScale },
];

export default function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const navItems = user?.role_name === "super-admin" ? superAdminNav : adminGudangNav;

  return (
    <aside className="w-64 bg-[#1b1b2f] text-white flex flex-col shrink-0">
      <div className="h-16 flex items-center px-6 border-b border-white/10">
        <h1 className="text-lg font-bold">Inventory</h1>
      </div>
      <nav className="flex-1 py-4 space-y-1 px-3">
        {navItems.map((item) => {
          const active = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                active
                  ? "bg-indigo-600 text-white"
                  : "text-zinc-400 hover:bg-[#162447] hover:text-white"
              }`}
            >
              <item.icon className="text-lg" />
              {item.label}
            </Link>
          );
        })}
      </nav>
      <div className="px-4 py-4 border-t border-white/10">
        <div className="text-sm text-zinc-400 truncate">{user?.email}</div>
        <button
          onClick={logout}
          className="mt-2 text-sm text-red-400 hover:text-red-300 transition-colors"
        >
          Logout
        </button>
      </div>
    </aside>
  );
}
