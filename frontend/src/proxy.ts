import { NextRequest, NextResponse } from "next/server";

const allowedRoles = ["super-admin", "admin-gudang"];
const superAdminRoutes = ["/users"];
const publicRoutes = ["/login"];

export default async function proxy(req: NextRequest) {
  const token = req.cookies.get("access_token")?.value;
  const role = req.cookies.get("role_name")?.value;
  const path = req.nextUrl.pathname;

  const isSuperAdminRoute = superAdminRoutes.some((route) => path === route || path.startsWith(route + "/"));
  const isPublic = publicRoutes.some((route) => path === route || path.startsWith(route + "/"));

  if (!token || !role || !allowedRoles.includes(role)) {
    if (isPublic) return NextResponse.next();
    return NextResponse.redirect(new URL("/login", req.nextUrl));
  }

  if (isSuperAdminRoute && role !== "super-admin") {
    return NextResponse.redirect(new URL("/", req.nextUrl));
  }

  if (isPublic) {
    return NextResponse.redirect(new URL("/", req.nextUrl));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
