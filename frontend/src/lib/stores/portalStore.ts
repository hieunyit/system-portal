import { create } from 'zustand';
import type { PortalUser } from '@/types/portal';

interface PortalState {
  users: PortalUser[];
  setUsers: (u: PortalUser[]) => void;
}

export const usePortalStore = create<PortalState>((set) => ({
  users: [],
  setUsers: (u) => set({ users: u }),
}));
