import { useAuthContext } from '../auth/context';

export default function useAuth() {
  return useAuthContext();
}
