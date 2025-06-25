import { NextRequest, NextResponse } from 'next/server';
import axios from 'axios';

export async function POST(req: NextRequest) {
  const data = await req.json();
  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/refresh`,
      data,
    );
    const body = res.data;
    if (body?.success?.data) {
      const { accessToken, refreshToken } = body.success.data;
      return NextResponse.json({
        access_token: accessToken,
        refresh_token: refreshToken,
      });
    }
    return NextResponse.json(body);
  } catch (err) {
    console.error(err);
    return new NextResponse('Failed to refresh token', { status: 500 });
  }
}
