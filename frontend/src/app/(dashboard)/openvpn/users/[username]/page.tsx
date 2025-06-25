interface VpnUserDetailProps { params: { username: string } }

export default async function VpnUserDetail({ params }: VpnUserDetailProps) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/openvpn/users/${params.username}`);
  const user = await res.json();
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold">{user.username}</h1>
      <p>Group: {user.group}</p>
    </div>
  );
}
