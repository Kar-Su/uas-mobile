"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import api from "@/lib/api";
import DataTable from "@/components/DataTable";
import ConfirmModal from "@/components/ConfirmModal";
import {
  HiOutlineSearch,
  HiOutlineSortAscending,
  HiOutlineSortDescending,
} from "react-icons/hi";

interface TipeOption {
  id: number;
  name: string;
}

interface SatuanOption {
  id: number;
  satuan: string | null;
}

interface Barang {
  kode: string;
  name: string;
  tipe: { id: number; name: string };
  satuan: { id: number; satuan: string | null; keterangan: string | null };
  quantity: number;
  created_at: string;
  updated_at: string;
}

export default function BarangPage() {
  const [data, setData] = useState<Barang[]>([]);
  const [initialLoading, setInitialLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingKode, setEditingKode] = useState<string | null>(null);
  const [tipeOptions, setTipeOptions] = useState<TipeOption[]>([]);
  const [satuanOptions, setSatuanOptions] = useState<SatuanOption[]>([]);
  const [form, setForm] = useState({
    kode: "",
    name: "",
    tipe_id: 0,
    satuan_id: 0,
    quantity: "",
  });
  const [error, setError] = useState("");
  const [stayOpen, setStayOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Barang | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState("");
  const [tipeFilter, setTipeFilter] = useState(0);
  const [sortDir, setSortDir] = useState<"asc" | "desc" | "">("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);

  const toggleSort = () => {
    if (sortDir === "") setSortDir("asc");
    else if (sortDir === "asc") setSortDir("desc");
    else setSortDir("");
  };

  // Debounce Effect
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearch(search);
    }, 500);

    // Cleanup function: batalkan timeout jika user mengetik lagi sebelum 500ms
    return () => {
      clearTimeout(handler);
    };
  }, [search]);

  const fetchData = useCallback(async () => {
    try {
      const params = new URLSearchParams({
        page: page.toString(),
      });

      if (debouncedSearch.trim())
        params.append("search", debouncedSearch.trim());
      if (tipeFilter > 0) params.append("tipe", tipeFilter.toString());
      if (sortDir !== "") params.append("qty_order", sortDir);

      const res = await api.get(`/barang?${params.toString()}`);

      setData(res.data.data || []);
      setTotalPages(res.data.meta?.total_pages || 1);
      setTotalItems(res.data.meta?.total_items || 0);
    } catch {
      // ignore
    } finally {
      setInitialLoading(false);
    }
  }, [page, debouncedSearch, tipeFilter, sortDir]); // <-- Gunakan debouncedSearch di sini

  const fetchOptions = useCallback(async () => {
    try {
      const [tipeRes, satuanRes] = await Promise.all([
        api.get("/tipe-barang"),
        api.get("/satuan-barang"),
      ]);
      setTipeOptions(tipeRes.data.data || []);
      setSatuanOptions(satuanRes.data.data || []);
    } catch {
      // ignore
    }
  }, []);

  // const filteredData = useMemo(() => {
  //   let result = data;
  //   if (search) {
  //     const q = search.toLowerCase();
  //     result = result.filter(
  //       (item) =>
  //         item.name.toLowerCase().includes(q) ||
  //         item.kode.toLowerCase().includes(q),
  //     );
  //   }
  //   if (tipeFilter > 0) {
  //     result = result.filter((item) => item.tipe.id === tipeFilter);
  //   }
  //   if (sortDir === "asc") {
  //     result = [...result].sort((a, b) => a.quantity - b.quantity);
  //   } else if (sortDir === "desc") {
  //     result = [...result].sort((a, b) => b.quantity - a.quantity);
  //   }
  //   return result;
  // }, [data, search, tipeFilter, sortDir]);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, tipeFilter, sortDir]); // <-- Ubah search menjadi debouncedSearch

  useEffect(() => {
    fetchData();
    fetchOptions();

    const es = new EventSource("/api/sse", { withCredentials: true });
    es.addEventListener("barang", () => {
      fetchData();
    });
    es.onerror = () => {
      es.close();
    };
    return () => es.close();
  }, [fetchData, fetchOptions]);

  const openCreate = () => {
    setEditingKode(null);
    setForm({ kode: "", name: "", tipe_id: 0, satuan_id: 0, quantity: "" });
    setError("");
    setShowModal(true);
  };

  const openEdit = (item: Barang) => {
    setEditingKode(item.kode);
    setForm({
      kode: item.kode,
      name: item.name,
      tipe_id: item.tipe.id,
      satuan_id: item.satuan.id,
      quantity: String(item.quantity),
    });
    setError("");
    setShowModal(true);
  };

  const handleDelete = (item: Barang) => {
    setDeleteTarget(item);
    setDeleteError("");
  };

  const confirmDelete = async () => {
    if (!deleteTarget) return;
    setDeleting(true);
    setDeleteError("");
    try {
      await api.delete(`/barang/${deleteTarget.kode}`);
      setDeleteTarget(null);
      fetchData();
    } catch (err: unknown) {
      setDeleteError(
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error || "Gagal menghapus",
      );
    } finally {
      setDeleting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      if (editingKode) {
        const payload: Record<string, unknown> = {
          name: form.name,
          tipe_id: form.tipe_id,
          satuan_id: form.satuan_id,
        };
        if (form.quantity !== "") {
          payload.quantity = Number(form.quantity);
        }
        await api.put(`/barang/${editingKode}`, payload);
        setShowModal(false);
      } else {
        await api.post("/barang", {
          kode: form.kode,
          name: form.name,
          tipe_id: form.tipe_id,
          satuan_id: form.satuan_id,
          quantity: form.quantity === "" ? 0 : Number(form.quantity),
        });
        if (stayOpen) {
          setForm({
            kode: "",
            name: "",
            tipe_id: 0,
            satuan_id: 0,
            quantity: "",
          });
          setError("");
        } else {
          setShowModal(false);
        }
      }
      fetchData();
    } catch (err: unknown) {
      setError(
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error || "Operation failed",
      );
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-zinc-500">View all inventory items</p>
        <button
          onClick={openCreate}
          className="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors"
        >
          + Create Barang
        </button>
      </div>

      <div className="flex gap-3 mb-4">
        <div className="relative flex-1">
          <HiOutlineSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 text-lg" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by name or code..."
            className="w-full pl-10 pr-4 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <select
          value={tipeFilter}
          onChange={(e) => setTipeFilter(Number(e.target.value))}
          className="px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option value={0}>All Tipe</option>
          {tipeOptions.map((t) => (
            <option key={t.id} value={t.id}>
              {t.name}
            </option>
          ))}
        </select>
        <button
          onClick={toggleSort}
          className={`flex items-center gap-1 px-3 py-2 border rounded-lg text-sm transition-colors ${
            sortDir !== ""
              ? "bg-indigo-50 border-indigo-200 text-indigo-700"
              : "border-zinc-300 text-zinc-600 hover:bg-zinc-50"
          }`}
          title="Sort by quantity"
        >
          {sortDir === "asc" && <HiOutlineSortAscending className="text-lg" />}
          {sortDir === "desc" && (
            <HiOutlineSortDescending className="text-lg" />
          )}
          {sortDir === "" && (
            <HiOutlineSortDescending className="text-lg opacity-40" />
          )}{" "}
          {/* Tampilan default/non-filter */}
          Qty
        </button>
      </div>

      <DataTable
        columns={[
          { key: "kode", label: "Kode" },
          { key: "name", label: "Name" },
          {
            key: "tipe",
            label: "Tipe",
            render: (row) => (row.tipe as { name: string }).name,
          },
          {
            key: "satuan",
            label: "Satuan",
            render: (row) =>
              (row.satuan as { satuan: string | null }).satuan || "-",
          },
          { key: "quantity", label: "Quantity" },
        ]}
        data={data}
        loading={initialLoading}
        onEdit={openEdit}
        onDelete={handleDelete}
      />

      {!initialLoading && (
        <div className="flex items-center justify-between mt-4 text-sm text-zinc-500">
          <span>
            {totalItems} item{totalItems !== 1 ? "s" : ""}
          </span>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="px-3 py-1.5 border border-zinc-300 rounded-lg hover:bg-zinc-50 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              Prev
            </button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
              <button
                key={p}
                onClick={() => setPage(p)}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium ${
                  p === page
                    ? "bg-indigo-600 text-white"
                    : "border border-zinc-300 hover:bg-zinc-50"
                }`}
              >
                {p}
              </button>
            ))}
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="px-3 py-1.5 border border-zinc-300 rounded-lg hover:bg-zinc-50 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      )}

      <ConfirmModal
        open={!!deleteTarget}
        title="Delete Barang"
        message={`Hapus "${deleteTarget?.name}"?`}
        onConfirm={confirmDelete}
        onCancel={() => {
          setDeleteTarget(null);
          setDeleteError("");
        }}
        loading={deleting}
        error={deleteError}
      />

      {showModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl shadow-xl p-6 w-full max-w-md mx-4">
            <h3 className="text-lg font-semibold mb-4">
              {editingKode ? "Edit Barang" : "Create Barang"}
            </h3>
            {error && (
              <div className="bg-red-50 text-red-600 text-sm rounded-lg p-3 mb-4">
                {error}
              </div>
            )}
            <form onSubmit={handleSubmit} className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Kode
                </label>
                <input
                  value={form.kode}
                  onChange={(e) => setForm({ ...form, kode: e.target.value })}
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required={!editingKode}
                  disabled={!!editingKode}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Name
                </label>
                <input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Tipe
                </label>
                <select
                  value={form.tipe_id}
                  onChange={(e) =>
                    setForm({ ...form, tipe_id: Number(e.target.value) })
                  }
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                >
                  <option value={0} disabled>
                    Select tipe
                  </option>
                  {tipeOptions.map((t) => (
                    <option key={t.id} value={t.id}>
                      {t.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Satuan
                </label>
                <select
                  value={form.satuan_id}
                  onChange={(e) =>
                    setForm({ ...form, satuan_id: Number(e.target.value) })
                  }
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                >
                  <option value={0} disabled>
                    Select satuan
                  </option>
                  {satuanOptions.map((s) => (
                    <option key={s.id} value={s.id}>
                      {s.satuan || "-"}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Quantity
                </label>
                <input
                  type="number"
                  min="0"
                  value={form.quantity}
                  onChange={(e) =>
                    setForm({ ...form, quantity: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="0"
                />
              </div>
              <div className="flex items-center justify-between pt-2">
                {!editingKode && (
                  <label className="flex items-center gap-2 text-sm text-zinc-500 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={stayOpen}
                      onChange={(e) => setStayOpen(e.target.checked)}
                      className="rounded border-zinc-300 text-indigo-600 focus:ring-indigo-500"
                    />
                    Tambah lagi
                  </label>
                )}
                <div className="flex gap-3 ml-auto">
                  <button
                    type="button"
                    onClick={() => setShowModal(false)}
                    className="px-4 py-2 text-sm text-zinc-600 hover:text-zinc-800"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700"
                  >
                    {editingKode ? "Update" : "Create"}
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
