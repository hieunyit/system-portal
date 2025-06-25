interface VpnGroupDetailProps { params: { groupName: string } }

export default async function VpnGroupDetail({ params }: VpnGroupDetailProps) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/openvpn/groups/${params.groupName}`);
  const group = await res.json();
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold">{group.group_name}</h1>
      <p>Auth Method: {group.auth_method}</p>
    </div>
  );
}
