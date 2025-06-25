import { useDebounce } from './useDebounce';
import { useState } from 'react';

export function useSearch(delay = 300) {
  const [query, setQuery] = useState('');
  const debounced = useDebounce(query, delay);
  return { query, setQuery, debounced };
}
