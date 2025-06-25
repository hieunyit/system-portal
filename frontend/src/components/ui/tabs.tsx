import { ReactNode, useState } from 'react';

export function Tabs({ tabs }: { tabs: { label: string; content: ReactNode }[] }) {
  const [index, setIndex] = useState(0);
  return (
    <div>
      <div className="flex space-x-2 mb-2">
        {tabs.map((t, i) => (
          <button key={t.label} onClick={() => setIndex(i)} className={i === index ? 'font-bold' : ''}>{t.label}</button>
        ))}
      </div>
      <div>{tabs[index].content}</div>
    </div>
  );
}
