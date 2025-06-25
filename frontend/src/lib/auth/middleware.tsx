import { NextRequest, NextResponse } from 'next/server';

export function middleware(req: NextRequest) {
  // placeholder check
  const token = req.cookies.get('access_token');
  if (!token && !req.nextUrl.pathname.startsWith('/login')) {
    return NextResponse.redirect(new URL('/login', req.url));
  }
}
