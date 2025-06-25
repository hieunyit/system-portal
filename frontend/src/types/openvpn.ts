export interface VpnUser {
  username: string;
  group: string;
}

export interface VpnGroup {
  group_name: string;
  auth_method: string;
}

export interface Session {
  username: string;
  ip: string;
  connected_since?: string;
}
