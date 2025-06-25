export default function VpnUserForm() {
  return (
    <form className="space-y-2">
      <input type="text" placeholder="Username" className="border p-2 w-full" />
      <input type="text" placeholder="Group" className="border p-2 w-full" />
      <button className="bg-blue-600 text-white px-3 py-2">Save</button>
    </form>
  );
}
