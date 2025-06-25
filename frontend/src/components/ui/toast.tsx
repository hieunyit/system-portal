export function Toast({ message }: { message: string }) {
  return <div className="fixed bottom-2 right-2 bg-black text-white px-3 py-2">{message}</div>;
}
