import { create } from 'zustand';
import type { VpnUser } from '@/types/openvpn';

interface OpenVPNState {
  vpnUsers: VpnUser[];
  setVpnUsers: (u: VpnUser[]) => void;
}

export const useOpenVPNStore = create<OpenVPNState>((set) => ({
  vpnUsers: [],
  setVpnUsers: (u) => set({ vpnUsers: u }),
}));
