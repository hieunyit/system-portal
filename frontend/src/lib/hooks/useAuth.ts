import { useContext } from 'react';
import { AuthContext } from '../auth/context';

export default function useAuth() {
  const user = useContext(AuthContext);
  return { user, loading: false };
}
