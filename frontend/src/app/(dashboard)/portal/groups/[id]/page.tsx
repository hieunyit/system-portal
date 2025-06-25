interface GroupDetailProps { params: { id: string } }

export default async function GroupDetail({ params }: GroupDetailProps) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/portal/groups/${params.id}`);
  const group = await res.json();
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold">{group.name}</h1>
      <p>Display Name: {group.display_name}</p>
    </div>
  );
}
