"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import api from "@/lib/api";
import DataTable from "@/components/DataTable";
import ConfirmModal from "@/components/ConfirmModal";
import { HiOutlineSearch } from "react-icons/hi";

interface User {
  id: string;
  email: string;
  name: string;
  role_name: string;
}

interface UserForm {
  name: string;
  email: string;
  password: string;
  role_name: string;
}

const emptyForm: UserForm = {
  name: "",
  email: "",
  password: "",
  role_name: "user",
};

export default function UsersPage() {
  const router = useRouter();
  const { user } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [initialLoading, setInitialLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<UserForm>(emptyForm);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [stayOpen, setStayOpen] = useState(false);
  const [roleFilter, setRoleFilter] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<User | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);

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

  const fetchUsers = useCallback(async () => {
    try {
      const params = new URLSearchParams({
        page: page.toString(),
      });

      if (debouncedSearch.trim())
        params.append("search", debouncedSearch.trim());
      if (roleFilter) params.append("role", roleFilter.toString());

      const res = await api.get(`/users?${params.toString()}`);
      setUsers(res.data.data || []);
      setTotalPages(res.data.meta?.total_pages || 1);
      setTotalItems(res.data.meta?.total_items || 0);
    } catch {
      // ignore
    } finally {
      setInitialLoading(false);
    }
  }, [page, debouncedSearch, roleFilter]);

  // const filteredUsers = useMemo(() => {
  //   let result = users;

  //   if (search) {
  //     const q = search.toLowerCase();
  //     result = result.filter(
  //       (u) =>
  //         u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q),
  //     );
  //   }
  //   if (roleFilter) {
  //     result = result.filter((u) => u.role_name === roleFilter);
  //   }
  //   return result;
  // }, [users, search, roleFilter]);

  useEffect(() => {
    if (user && user.role_name !== "super-admin") {
      router.push("/");
      return;
    }
    fetchUsers();
    const es = new EventSource("/api/sse", { withCredentials: true });
    es.addEventListener("user", () => { fetchUsers(); });
    es.onerror = () => { es.close(); };
    return () => es.close();
  }, [user, router, fetchUsers]);

  useEffect(() => {
    setPage(1);
  }, [search, roleFilter]);

  const openCreate = () => {
    setEditingId(null);
    setForm(emptyForm);
    setError("");
    setShowModal(true);
  };

  const openEdit = (item: User) => {
    setEditingId(item.id);
    setForm({
      name: item.name,
      email: item.email,
      password: "",
      role_name: item.role_name,
    });
    setError("");
    setShowModal(true);
  };

  const handleDelete = (item: User) => {
    setDeleteTarget(item);
    setDeleteError("");
  };

  const confirmDelete = async () => {
    if (!deleteTarget) return;
    setDeleting(true);
    setDeleteError("");
    try {
      await api.delete(`/users/${deleteTarget.id}`);
      setDeleteTarget(null);
      fetchUsers();
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
      if (editingId) {
        const payload: Record<string, unknown> = {
          name: form.name,
          email: form.email,
          role_name: form.role_name,
        };
        if (form.password) payload.password = form.password;
        await api.put(`/users/${editingId}`, payload);
        setShowModal(false);
      } else {
        await api.post("/users", form);
        if (stayOpen) {
          setForm(emptyForm);
          setError("");
        } else {
          setShowModal(false);
        }
      }
      fetchUsers();
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
        <p className="text-zinc-500">Manage all users in the system</p>
        <button
          onClick={openCreate}
          className="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors"
        >
          + Create User
        </button>
      </div>

      <div className="flex gap-3 mb-4">
        <div className="relative flex-1">
          <HiOutlineSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 text-lg" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by name or email..."
            className="w-full pl-10 pr-4 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <select
          value={roleFilter}
          onChange={(e) => setRoleFilter(e.target.value)}
          className="px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option value="">All Roles</option>
          <option value="super-admin">Super Admin</option>
          <option value="admin-gudang">Admin Gudang</option>
          <option value="user">User</option>
        </select>
      </div>

      <DataTable
        columns={[
          { key: "name", label: "Name" },
          { key: "email", label: "Email" },
          { key: "role_name", label: "Role" },
        ]}
        data={users}
        loading={initialLoading}
        onEdit={openEdit}
        onDelete={handleDelete}
      />

      {!initialLoading && (
        <div className="flex items-center justify-between mt-4 text-sm text-zinc-500">
          <span>
            {totalItems} user{totalItems !== 1 ? "s" : ""}
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
        title="Delete User"
        message={`Hapus user "${deleteTarget?.name}"?`}
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
              {editingId ? "Edit User" : "Create User"}
            </h3>
            {error && (
              <div className="bg-red-50 text-red-600 text-sm rounded-lg p-3 mb-4">
                {error}
              </div>
            )}
            <form onSubmit={handleSubmit} className="space-y-3">
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
                  Email
                </label>
                <input
                  type="email"
                  value={form.email}
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Password {editingId && "(leave empty to keep)"}
                </label>
                <input
                  type="password"
                  value={form.password}
                  onChange={(e) =>
                    setForm({ ...form, password: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  {...(!editingId && { required: true })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-zinc-700 mb-1">
                  Role
                </label>
                <select
                  value={form.role_name}
                  onChange={(e) =>
                    setForm({ ...form, role_name: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-zinc-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="super-admin">Super Admin</option>
                  <option value="admin-gudang">Admin Gudang</option>
                  <option value="user">User</option>
                </select>
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
