"use client";

import { useEffect, useState, useCallback } from "react";
import api from "@/lib/api";
import DataTable from "@/components/DataTable";
import ConfirmModal from "@/components/ConfirmModal";

interface SatuanBarang {
  id: number;
  satuan: string | null;
  keterangan: string | null;
}

export default function SatuanBarangPage() {
  const [data, setData] = useState<SatuanBarang[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState({ satuan: "", keterangan: "" });
  const [error, setError] = useState("");
  const [stayOpen, setStayOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<SatuanBarang | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState("");

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.get("/satuan-barang");
      setData(res.data.data || []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const es = new EventSource("/api/sse", { withCredentials: true });
    es.addEventListener("satuan_barang", () => { fetchData(); });
    es.onerror = () => { es.close(); };
    return () => es.close();
  }, [fetchData]);

  const openCreate = () => {
    setEditingId(null);
    setForm({ satuan: "", keterangan: "" });
    setError("");
    setShowModal(true);
  };

  const openEdit = (item: SatuanBarang) => {
    setEditingId(item.id);
    setForm({ satuan: item.satuan || "", keterangan: item.keterangan || "" });
    setError("");
    setShowModal(true);
  };

  const handleDelete = (item: SatuanBarang) => {
    setDeleteTarget(item);
    setDeleteError("");
  };

  const confirmDelete = async () => {
    if (!deleteTarget) return;
    setDeleting(true);
    setDeleteError("");
    try {
      await api.delete(`/satuan-barang/${deleteTarget.id}`);
      setDeleteTarget(null);
      fetchData();
    } catch (err: unknown) {
      setDeleteError((err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal menghapus");
    } finally {
      setDeleting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      if (editingId) {
        await api.put(`/satuan-barang/${editingId}`, form);
        setShowModal(false);
      } else {
        await api.post("/satuan-barang", form);
        if (stayOpen) {
          setForm({ satuan: "", keterangan: "" });
        } else {
          setShowModal(false);
        }
      }
      fetchData();
    } catch (err: unknown) {
      setError((err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Operation failed");
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-zinc-500">View all units of measurement</p>
        <button
          onClick={openCreate}
          className="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors"
        >
          + Create Satuan
        </button>
      </div>

      <DataTable
        columns={[
          { key: "id", label: "ID" },
          { key: "satuan", label: "Satuan", render: (row) => (row.satuan as string) || "-" },
          { key: "keterangan", label: "Keterangan", render: (row) => (row.keterangan as string) || "-" },
        ]}
        data={data}
        loading={loading}
        onEdit={openEdit}
        onDelete={handleDelete}
      />

      <ConfirmModal
        open={!!deleteTarget}
        title="Delete Satuan Barang"
        message={`Hapus "${deleteTarget?.satuan || deleteTarget?.id}"?`}
        onConfirm={confirmDelete}
        onCancel={() => { setDeleteTarget(null); setDeleteError(""); }}
        loading={deleting}
        error={deleteError}
      />

      {showModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl shadow-xl p-6 w-full max-w-md mx-4">
            <h3 className="text-lg font-semibold mb-4">{editingId ? "Edit Satuan Barang" : "Create Satuan Barang"}</h3>
            {error && <div className="bg-red-50 text-red-600 text-sm rounded-lg p-3 mb-4">{error}</div>}
            <form onSubmit={handleSubmit} className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">Satuan</label>
                <input
                  value={form.satuan}
                  onChange={(e) => setForm({ ...form, satuan: e.target.value })}
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="e.g. kg, pcs, liter"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">Keterangan</label>
                <input
                  value={form.keterangan}
                  onChange={(e) => setForm({ ...form, keterangan: e.target.value })}
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Optional description"
                />
              </div>
              <div className="flex items-center justify-between pt-2">
                {!editingId && (
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
                    {editingId ? "Update" : "Create"}
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
