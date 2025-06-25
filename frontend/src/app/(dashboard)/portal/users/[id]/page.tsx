interface UserDetailProps {
  params: { id: string };
}

export default async function UserDetail({ params }: UserDetailProps) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/portal/users/${params.id}`);
  const user = await res.json();
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold">{user.username}</h1>
      <p>Email: {user.email}</p>
      <p>Full Name: {user.full_name}</p>
    </div>
  );
}
