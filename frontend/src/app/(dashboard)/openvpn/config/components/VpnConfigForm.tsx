export default function VpnConfigForm() {
  return (
    <form className="space-y-2">
      <input type="text" placeholder="Key" className="border p-2 w-full" />
      <input type="text" placeholder="Value" className="border p-2 w-full" />
      <button className="bg-blue-600 text-white px-3 py-2">Save</button>
    </form>
  );
}
