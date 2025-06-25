import { useState } from 'react';

export function useLocalStorage(key: string, initial: string) {
  const [value, setValue] = useState(() => localStorage.getItem(key) || initial);
  const set = (v: string) => {
    setValue(v);
    localStorage.setItem(key, v);
  };
  return [value, set] as const;
}
