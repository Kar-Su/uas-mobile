"use client";

import { useEffect, useState, useCallback } from "react";
import api from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import {
  HiOutlineCube,
  HiOutlineTag,
  HiOutlineScale,
  HiOutlineUsers,
} from "react-icons/hi";

interface DashboardData {
  totalBarang: number;
  totalTipe: number;
  totalSatuan: number;
  totalUsers: number;
}


export default function DashboardPage() {
  const { user } = useAuth();
  const isSuperAdmin = user?.role_name === "super-admin";
  const [data, setData] = useState<DashboardData | null>(null);

  const fetchData = useCallback(() => {
    const promises = [
      api
        .get("/barang?page=1")
        .catch(() => ({ data: { data: [], meta: { total_items: 0 } } })),
      api.get("/tipe-barang").catch(() => ({ data: { data: [] } })),
      api.get("/satuan-barang").catch(() => ({ data: { data: [] } })),
    ];
    if (isSuperAdmin) {
      promises.push(
        api
          .get("/users?page=1")
          .catch(() => ({ data: { data: [], meta: { total_items: 0 } } })),
      );
    }
    Promise.all(promises).then(([barang, tipe, satuan, users]) => {
      setData({
        totalBarang: barang.data.meta?.total_items || 0,
        totalTipe: Array.isArray(tipe.data.data) ? tipe.data.data.length : 0,
        totalSatuan: Array.isArray(satuan.data.data)
          ? satuan.data.data.length
          : 0,
        totalUsers: users?.data?.meta?.total_items || 0,
      });
    });
  }, [isSuperAdmin]);

  useEffect(() => {
    fetchData();
    const es = new EventSource("/api/sse", { withCredentials: true });
    const refresh = () => fetchData();
    es.addEventListener("barang", refresh);
    es.addEventListener("tipe_barang", refresh);
    es.addEventListener("satuan_barang", refresh);
    es.addEventListener("user", refresh);
    es.onerror = () => { es.close(); };
    return () => es.close();
  }, [fetchData]);

  const cards = [
    {
      label: "Total Barang",
      value: data?.totalBarang ?? 0,
      icon: HiOutlineCube,
      color: "bg-blue-500",
    },
    {
      label: "Tipe Barang",
      value: data?.totalTipe ?? 0,
      icon: HiOutlineTag,
      color: "bg-emerald-500",
    },
    {
      label: "Satuan Barang",
      value: data?.totalSatuan ?? 0,
      icon: HiOutlineScale,
      color: "bg-amber-500",
    },
  ];
  if (isSuperAdmin) {
    cards.push({
      label: "Users",
      value: data?.totalUsers ?? 0,
      icon: HiOutlineUsers,
      color: "bg-purple-500",
    });
  }

  return (
    <div>
      <p className="text-zinc-500 mb-6">Overview of your inventory system</p>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {cards.map((card) => (
          <div
            key={card.label}
            className="bg-white rounded-xl shadow-sm border border-zinc-200 p-6"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-zinc-500">{card.label}</p>
                <p className="text-3xl font-bold mt-1">{card.value}</p>
              </div>
              <div
                className={`w-12 h-12 rounded-lg ${card.color} flex items-center justify-center`}
              >
                <card.icon className="text-white text-xl" />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
